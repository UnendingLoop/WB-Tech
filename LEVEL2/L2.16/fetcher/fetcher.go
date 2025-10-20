// Package fetcher requests files using provided URLs and returns their content and type
package fetcher

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

// Fetcher выполняет HTTP-запросы и возвращает содержимое страницы.
type Fetcher struct {
	Client *http.Client
}

// FetchResult — результат скачивания.
type FetchResult struct {
	URL         *url.URL
	Content     []byte
	ContentType string
}

func (f *Fetcher) Fetch(inputURL *url.URL, t int) (*FetchResult, error) {
	client := &http.Client{
		Timeout: time.Duration(t) * time.Second,
	}

	req, _ := http.NewRequest("GET", inputURL.String(), nil)
	req.Header.Set("User-Agent", "GoWget/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to load: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server returned code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp body: %v", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	return &FetchResult{
		URL:         inputURL,
		Content:     makeLinksAbsolute(inputURL, body),
		ContentType: contentType,
	}, nil
}

func makeLinksAbsolute(base *url.URL, htmlData []byte) []byte {
	doc, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return htmlData
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a", "link", "use":
				for i, a := range n.Attr {
					if a.Key == "href" || a.Key == "xlink:href" {
						u, err := url.Parse(a.Val)
						if err == nil && !u.IsAbs() {
							n.Attr[i].Val = base.ResolveReference(u).String()
						}
					}
				}
			case "img", "script":
				for i, a := range n.Attr {
					if a.Key == "src" {
						u, err := url.Parse(a.Val)
						if err == nil && !u.IsAbs() {
							n.Attr[i].Val = base.ResolveReference(u).String()
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		log.Printf("Error during render of loaded page after fixing all URLs to abs-format: %v", err)
		return buf.Bytes()
	}
	return buf.Bytes()
}
