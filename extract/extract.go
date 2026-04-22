package sfExtract

import (
	"io"
	"net/http"
	"golang.org/x/net/html"
)

// ExtractHTML fetches raw HTML from a URL
func ExtractHTML(input map[string]any) (any, error) {
	url, _ := input["url"].(string)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return string(body), nil
}

// ExtractLinks pulls all hrefs from a URL
func ExtractLinks(input map[string]any) (any, error) {
	url, _ := input["url"].(string)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	var links []string

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()
		if token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
	}

	return links, nil
}

// ExtractMeta pulls SEO/OG meta tags from a URL
func ExtractMeta(input map[string]any) (any, error) {
	url, _ := input["url"].(string)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	meta := make(map[string]string)

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()
		if token.Data == "meta" {
			var name, property, content string
			for _, attr := range token.Attr {
				if attr.Key == "name" {
					name = attr.Val
				}
				if attr.Key == "property" {
					property = attr.Val
				}
				if attr.Key == "content" {
					content = attr.Val
				}
			}

			key := name
			if key == "" { key = property }
			if key != "" && content != "" {
				meta[key] = content
			}
		}
		// Optimization: stop after </head>
		if token.Data == "body" {
			break
		}
	}

	return meta, nil
}
