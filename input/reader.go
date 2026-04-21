package sfInput

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

/**
 * Smart Reader: Automatically detects if path is a local file or a URL
 * and returns an io.ReadCloser that can be used directly for streaming.
 */
func GetReader(path string) (io.ReadCloser, error) {
	// Check for Cloud URLs
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return nil, fmt.Errorf("network error: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to fetch cloud resource (HTTP %d)", resp.StatusCode)
		}

		return resp.Body, nil
	}

	// Normal local file
	return os.Open(path)
}
