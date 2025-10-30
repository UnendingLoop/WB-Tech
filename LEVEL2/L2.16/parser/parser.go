// Package parser recursively parses provided html-code as a DOM-tree. Result - list of all URLs to htmls and to other stuff to download
package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type ResourceLinks struct {
	Pages  []string // ссылки на другие HTML
	Assets []string // ресурсы: js, css, img
}

func ParseHTML(base *url.URL, data []byte) (*ResourceLinks, error) {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %v", err)
	}

	result := ResourceLinks{}

	var walk func(node *html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "a":
				href := getAttr(node, "href")
				if href != "" && !strings.HasPrefix(href, "#") && !strings.HasPrefix(href, "mailto:") && !strings.HasPrefix(href, "tel:") {
					result.Pages = append(result.Pages, href)
				}

			case "img", "script":
				src := getAttr(node, "src")
				if src != "" {
					result.Assets = append(result.Assets, src)
				}

			case "link":
				rel := getAttr(node, "rel")
				if strings.Contains(rel, "stylesheet") {
					href := getAttr(node, "href")
					result.Assets = append(result.Assets, href)
				}
			case "use":
				href := getAttr(node, "xlink:href")
				if href == "" {
					href = getAttr(node, "href")
				}
				if href != "" {
					if idx := strings.Index(href, "#"); idx != -1 {
						href = href[:idx]
					}
					u, err := url.Parse(href)
					if err == nil {
						abs := base.ResolveReference(u)
						result.Assets = append(result.Assets, abs.String())
					}
				}

			}
		}

		// всегда продолжаем обход потомков
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(doc)

	// избавляемся от дублей ассетов
	uniqueAssets := make(map[string]struct{})
	cleanAssets := make([]string, 0, len(result.Assets))
	for _, a := range result.Assets {
		if _, ok := uniqueAssets[a]; !ok {
			uniqueAssets[a] = struct{}{}
			cleanAssets = append(cleanAssets, a)
		}
	}
	result.Assets = cleanAssets

	// избавляемся от дублей страниц
	uniquePages := make(map[string]struct{})
	cleanPages := make([]string, 0, len(result.Pages))
	for _, a := range result.Pages {
		if _, ok := uniquePages[a]; !ok {
			uniquePages[a] = struct{}{}
			cleanPages = append(cleanPages, a)
		}
	}
	result.Pages = cleanPages

	return &result, nil
}

func getAttr(node *html.Node, key string) string {
	for _, v := range node.Attr {
		vKey := strings.ToLower(v.Key)
		if vKey == key {
			return v.Val
		}
	}
	return ""
}
