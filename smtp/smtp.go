package sfSmtp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"sync"
	
	sfInput "github.com/t8nlab/surface/input"
)

func SmtpSend(input map[string]any) (any, error) {
	return sendInternal(input)
}

func sendInternal(input map[string]any) (any, error) {
	host, _ := sfInput.GetString(input, "host")
	port, _ := sfInput.GetInt(input, "port")
	user, _ := sfInput.GetString(input, "username")
	pass, _ := sfInput.GetString(input, "password")
	from, _ := sfInput.GetString(input, "from")
	to, _ := sfInput.GetString(input, "to")
	cc, _ := sfInput.GetString(input, "cc")
	bcc, _ := sfInput.GetString(input, "bcc")
	subject, _ := sfInput.GetString(input, "subject")
	body, _ := sfInput.GetString(input, "body")
	isSSL := sfInput.GetBool(input, "ssl", false)

	if port == 0 {
		port = 587
	}

	isRaw := sfInput.GetBool(input, "raw", false)

	addr := fmt.Sprintf("%s:%d", host, port)
	var msg []byte
	
	if isRaw {
		msg = []byte(body)
		
		// If to/from are missing, parse them from the headers in the body
		if from == "" || to == "" {
			if m, err := mail.ReadMessage(strings.NewReader(body)); err == nil {
				if from == "" { from = m.Header.Get("From") }
				if to == "" { to = m.Header.Get("To") }
			}
		}
	} else {
		msg = buildMessage(from, to, cc, subject, body)
	}
	
	recipients := getRecipients(to, cc, bcc)
	auth := smtp.PlainAuth("", user, pass, host)

	if port == 465 || isSSL {
		return sendWithSSL(addr, host, auth, from, recipients, msg)
	}

	err := smtp.SendMail(addr, auth, from, recipients, msg)
	if err != nil {
		return nil, err
	}

	return map[string]any{"status": "sent", "to": to}, nil
}

func buildMessage(from, to, cc, subject, body string) []byte {
	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	if cc != "" {
		header["Cc"] = cc
	}
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	return []byte(message)
}

func getRecipients(to, cc, bcc string) []string {
	var r []string
	if to != "" { r = append(r, strings.Split(to, ",")...) }
	if cc != "" { r = append(r, strings.Split(cc, ",")...) }
	if bcc != "" { r = append(r, strings.Split(bcc, ",")...) }
	for i := range r {
		r[i] = strings.TrimSpace(r[i])
	}
	return r
}

func sendWithSSL(addr, host string, auth smtp.Auth, from string, to []string, msg []byte) (any, error) {
	conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: false, ServerName: host})
	if err != nil { return nil, err }
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil { return nil, err }
	defer client.Quit()

	if auth != nil {
		if err = client.Auth(auth); err != nil { return nil, err }
	}

	if err = client.Mail(from); err != nil { return nil, err }
	for _, k := range to {
		if err = client.Rcpt(k); err != nil { return nil, err }
	}

	w, err := client.Data()
	if err != nil { return nil, err }
	_, err = w.Write(msg)
	w.Close()

	return map[string]any{"status": "sent", "ssl": true}, nil
}

// BULK SEND WITH CONNECTION POOLING (REUSE CONNECTIONS)
func SmtpBulkSend(input map[string]any) (any, error) {
	emails, ok := input["emails"].([]any)
	if !ok { return nil, fmt.Errorf("emails must be an array") }

	host, _ := sfInput.GetString(input, "host")
	port, _ := sfInput.GetInt(input, "port")
	user, _ := sfInput.GetString(input, "username")
	pass, _ := sfInput.GetString(input, "password")
	from, _ := sfInput.GetString(input, "from")
	concurrency, _ := sfInput.GetInt(input, "concurrency")
	if concurrency <= 0 { concurrency = 3 }

	results := make([]any, len(emails))
	var wg sync.WaitGroup
	emailChan := make(chan struct {
		index int
		data  map[string]any
	}, len(emails))

	// Start worker pool (each worker maintains ONE persistent connection)
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			addr := fmt.Sprintf("%s:%d", host, port)
			auth := smtp.PlainAuth("", user, pass, host)

			// 1. Establish persistent connection
			client, err := dialClient(addr, host, auth)
			if err != nil {
				return
			}
			defer client.Quit()

			// 2. Process emails over this connection
			for email := range emailChan {
				d := email.data
				to, _ := sfInput.GetString(d, "to")
				cc, _ := sfInput.GetString(d, "cc")
				bcc, _ := sfInput.GetString(d, "bcc")
				subject, _ := sfInput.GetString(d, "subject")
				body, _ := sfInput.GetString(d, "body")

				msg := buildMessage(from, to, cc, subject, body)
				recipients := getRecipients(to, cc, bcc)

				if err := sendOverConnection(client, from, recipients, msg); err != nil {
					results[email.index] = map[string]any{"error": err.Error(), "to": to}
					
					// Reconnect on failure
					client.Quit()
					client, _ = dialClient(addr, host, auth)
					if client == nil { return }
				} else {
					results[email.index] = map[string]any{"status": "sent", "to": to}
				}
			}
		}()
	}

	// Fill email queue
	for i, e := range emails {
		data, _ := e.(map[string]any)
		emailChan <- struct {
			index int
			data  map[string]any
		}{i, data}
	}
	close(emailChan)

	wg.Wait()
	return results, nil
}

func dialClient(addr, host string, auth smtp.Auth) (*smtp.Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil { return nil, err }
	
	client, err := smtp.NewClient(conn, host)
	if err != nil { return nil, err }

	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{InsecureSkipVerify: false, ServerName: host}
		if err := client.StartTLS(config); err != nil {
			return nil, err
		}
	}

	if auth != nil {
		if err = client.Auth(auth); err != nil { return nil, err }
	}
	return client, nil
}

func sendOverConnection(c *smtp.Client, from string, to []string, msg []byte) error {
	if err := c.Mail(from); err != nil { return err }
	for _, k := range to {
		if err := c.Rcpt(k); err != nil { return err }
	}
	w, err := c.Data()
	if err != nil { return err }
	_, err = w.Write(msg)
	if err != nil { return err }
	return w.Close()
}
