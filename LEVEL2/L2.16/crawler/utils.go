package crawler

import (
	"bytes"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// вычисляем путь сохранения относительно InitURL
func (c *Crawler) makeLocalPath(u *url.URL) string {
	base, _ := url.Parse(c.InitURL)
	rel := strings.TrimPrefix(u.Path, "/")

	// если это самая первая страница — кладём прямо в корень
	if u.String() == base.String() {
		return filepath.Join("mirror", u.Host, "index.html")
	}

	ext := strings.ToLower(filepath.Ext(rel))
	switch ext {
	case ".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".ico", ".woff", ".woff2", ".ttf":
		filename := filepath.Base(rel)
		return filepath.Join("mirror", c.Domain, "_assets", filename)
	}

	// если путь пустой или заканчивается на /
	if rel == "" || strings.HasSuffix(rel, "/") {
		rel = filepath.Join(rel, "index.html")
	} else {
		// если последний сегмент не содержит точки (расширения)
		baseName := filepath.Base(rel)
		if !strings.Contains(baseName, ".") {
			rel = filepath.Join(rel, "index.html")
		}
	}

	// если сабдомен - положить его внутрь основного домена а не рядом с ним
	subPath := ""
	if u.Host != c.Domain {
		subPath = strings.TrimSuffix(u.Host, "."+c.Domain)
	}
	localDir := filepath.Join("mirror", c.Domain, subPath)
	return filepath.Join(localDir, rel)
}

// превращаем абсолютный локальный путь типа mirror/habr.com/assets/* в /assets/*
func (c *Crawler) makeRelativeToMirror(localPath string) string {
	localPath = filepath.ToSlash(localPath)

	// Находим корневой префикс, который нужно обрезать
	prefix := filepath.ToSlash(filepath.Join("mirror", c.Domain)) + "/"

	// Убираем его из начала пути
	localPath = strings.TrimPrefix(localPath, prefix)

	return localPath
}

// переделываем абсолютные пути внутри HTML в относительные
func (c *Crawler) makeRelativePaths(basePath string, html []byte) []byte {
	depth := strings.Count(filepath.Dir(basePath), string(os.PathSeparator)) - 2 // -2 чтобы не считать mirror/domain
	if depth < 0 {
		depth = 0
	}
	prefix := strings.Repeat("../", depth)
	html = bytes.ReplaceAll(html, []byte(`src="/`), []byte(`src="`+prefix))
	html = bytes.ReplaceAll(html, []byte(`href="/`), []byte(`href="`+prefix))
	return html
}

// возвращает нормализованный URL для мапы Downloaded/Queue.
func normalizeURL(u *url.URL) string {
	if u == nil {
		return ""
	}
	nu := *u // копия
	nu.Fragment = ""
	// Для простоты: если путь == "" => "/"
	if nu.Path == "" {
		nu.Path = "/"
	}
	// Убрать лишний конечный слеш (но не для корня "/")
	if len(nu.Path) > 1 && strings.HasSuffix(nu.Path, "/") {
		nu.Path = strings.TrimRight(nu.Path, "/")
	}
	// Схема и хост в нижний регистр
	nu.Scheme = strings.ToLower(nu.Scheme)
	nu.Host = strings.ToLower(nu.Host)
	// Не включаем User/Password
	nu.User = nil
	return nu.String()
}

func (c *Crawler) replaceAssetLinksInHTML(htmlData []byte) []byte {
	doc, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return htmlData
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "img", "script", "iframe", "video", "audio", "source":
				for i, a := range n.Attr {
					switch a.Key {
					case "src", "poster":
						if newVal := c.makeLocalIfDownloaded(a.Val); newVal != "" {
							n.Attr[i].Val = newVal
						}

					case "srcset":
						fixed := c.fixSrcset(a.Val)
						if fixed != "" {
							n.Attr[i].Val = fixed
						}
					}
				}

			case "link", "a", "use":
				for i, a := range n.Attr {
					if a.Key == "href" || a.Key == "xlink:href" {
						if newVal := c.makeLocalIfDownloaded(a.Val); newVal != "" {
							n.Attr[i].Val = newVal
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
		log.Printf("render HTML after replacing asset links failed: %v", err)
		return htmlData
	}
	return buf.Bytes()
}

// заменяем URL в srcset на локальные пути из _assets.
func (c *Crawler) fixSrcset(srcset string) string {
	const ReplaceAllSrcset = true // если false — оставит сетевые ссылки как fallback

	// Разбиваем по запятым (разные размеры: 480w, 780w и т.п.)
	parts := strings.Split(srcset, ",")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// srcset может быть в формате: "https://...jpg 780w"
		fields := strings.Fields(part)
		if len(fields) == 0 {
			continue
		}

		link := fields[0]
		u, err := url.Parse(strings.TrimSpace(link))
		if err != nil || u.Host == "" {
			continue
		}

		// Если ассет уже скачан — подменяем ссылку на локальный путь
		if c.isDownloadedURL(u) {
			localPath := c.makeLocalPath(u)
			relPath := c.makeRelativeToMirror(localPath)
			fields[0] = relPath
			parts[i] = strings.Join(fields, " ")
			continue
		}

		// Если не скачан — удаляем, если включён ReplaceAllSrcset
		if ReplaceAllSrcset {
			parts[i] = ""
		}
	}

	// Удаляем пустые элементы и собираем обратно
	var result []string
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			result = append(result, p)
		}
	}

	return strings.TrimSpace(strings.Join(result, ", "))
}

// проверка разрешенности домена для страниц
func (c *Crawler) isAllowedPage(u *url.URL) bool {
	return u.Host == c.Domain
}

func (c *Crawler) updatePagesQueue(newPages []string, parentDepth int) {
	if parentDepth+1 > c.MaxDepth {
		return
	}

	for _, page := range newPages {
		nextURL, err := url.Parse(page)
		if err == nil && c.isAllowedPage(nextURL) && c.isRobotAllowed(nextURL) && !c.isDownloadedURL(nextURL) && !c.isQueued(normalizeURL(nextURL)) {
			c.Queue = append(c.Queue, QueueItem{URL: nextURL, Depth: parentDepth + 1})
		}
	}
	// log.Printf("Updated pages queue is: %v", c.Queue)
}
