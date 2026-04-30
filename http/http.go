package sfHttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Request performs an HTTP request with full control, similar to axios/fetch
func Request(input map[string]any) (any, error) {
	urlStr, _ := input["url"].(string)
	method, _ := input["method"].(string)
	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	// Handle params (query string)
	if params, ok := input["params"].(map[string]any); ok {
		u, err := url.Parse(urlStr)
		if err == nil {
			q := u.Query()
			for k, v := range params {
				q.Set(k, fmt.Sprintf("%v", v))
			}
			u.RawQuery = q.Encode()
			urlStr = u.String()
		}
	}

	var bodyReader io.Reader
	var contentType string

	if data := input["data"]; data != nil {
		switch v := data.(type) {
		case string:
			bodyReader = strings.NewReader(v)
		default:
			jsonData, err := json.Marshal(v)
			if err == nil {
				bodyReader = bytes.NewReader(jsonData)
				contentType = "application/json"
			}
		}
	}

	req, err := http.NewRequest(method, urlStr, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set default content type if we detected JSON
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Handle custom headers
	if headers, ok := input["headers"].(map[string]any); ok {
		for k, v := range headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	// Timeout
	timeoutMs, _ := input["timeout"].(float64)
	if timeoutMs <= 0 {
		timeoutMs = 30000 // default 30s
	}

	client := &http.Client{
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respHeaders := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			respHeaders[k] = v[0]
		}
	}

	// Try to parse JSON if content-type is application/json
	var responseData any = string(respBody)
	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var decoded any
		if err := json.Unmarshal(respBody, &decoded); err == nil {
			responseData = decoded
		}
	}

	return map[string]any{
		"status":     resp.StatusCode,
		"statusText": resp.Status,
		"headers":    respHeaders,
		"data":       responseData,
		"url":        resp.Request.URL.String(),
		"ok":         resp.StatusCode >= 200 && resp.StatusCode < 300,
	}, nil
}

// Get is a shorthand for Request with GET method
func Get(input map[string]any) (any, error) {
	input["method"] = "GET"
	return Request(input)
}

// Post is a shorthand for Request with POST method
func Post(input map[string]any) (any, error) {
	input["method"] = "POST"
	return Request(input)
}
