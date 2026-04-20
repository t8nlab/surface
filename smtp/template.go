package sfSmtp

import (
	"bytes"
	"html/template"
	sfInput "github.com/t8nlab/surface/input"
)

/**
 * Execute a Go HTML template with data
 */
func SmtpRender(input map[string]any) (any, error) {
	tplStr, err := sfInput.GetString(input, "template")
	if err != nil {
		return nil, err
	}

	data, _ := input["data"]

	// Create and parse template
	tmpl, err := template.New("sf_email").Parse(tplStr)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.String(), nil
}

/**
 * Read and execute a Go HTML template file from disk
 */
func SmtpRenderFile(input map[string]any) (any, error) {
	path, err := sfInput.GetString(input, "path")
	if err != nil {
		return nil, err
	}

	data, _ := input["data"]

	// Parse the file directly in Go
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.String(), nil
}
