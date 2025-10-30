// Package crawler contains core-logic of wget
package crawler

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"miniWget/fetcher"
	"miniWget/parser"
	"miniWget/saver"

	"golang.org/x/net/html"
)

// Crawler - конфиг для обхода страниц
type Crawler struct {
	MaxDepth        int              // глубина рекурсии скачивания из os.Args
	Timeout         int              // таймаут соединения при скачивании по ссылкам
	InitURL         string           // введенный пользователем URL из os.Args
	Domain          string           // доменное имя InitURL
	currentPagePath string           // временно хранит путь текущей страницы при обработке HTML
	Queue           []QueueItem      // очередь страниц на скачивание
	Fetcher         *fetcher.Fetcher // Fetcher производит первичное скачивание по InitURL
	Saver           *saver.Saver     // Saver скачивает ассеты текущего URL из очереди
	Downloaded      map[string]bool  // карта для контроля скачанных ассетов/страниц - для избежания повторных скачиваний
	RobotAllows     map[string]bool  // для выполнения требований из robots.txt
	wg              sync.WaitGroup   // для контроля горутин скачивания ассетов
	sync.RWMutex                     // нужно для защиты карты Downloaded
}

type QueueItem struct {
	URL   *url.URL // распарсенный URL
	Depth int      // текущая глубина рекурсии считая от InitURL
}

func (c *Crawler) Crawl() {
	startURL, err := url.Parse(c.InitURL)
	if err != nil {
		log.Printf("invalid start url: %v", err)
		return
	}
	c.Queue = append(c.Queue, QueueItem{URL: startURL, Depth: 0})

	for len(c.Queue) > 0 {
		item := c.Queue[0]
		c.Queue = c.Queue[1:]

		loadedPage, err := c.Fetcher.Fetch(item.URL, c.Timeout)
		if err != nil {
			log.Printf("failed to fetch %s: %v", item.URL.String(), err)
			continue
		}

		links, err := parser.ParseHTML(item.URL, loadedPage.Content)
		if err != nil {
			log.Printf("failed to parse html %s: %v", item.URL.String(), err)
			continue
		}

		c.downloadAssets(links.Assets)
		c.wg.Wait()

		c.updatePagesQueue(links.Pages, item.Depth)

		// Перед заменой ссылок сохраняем текущий путь
		c.currentPagePath = c.makeLocalPath(item.URL)
		localPageData := c.replaceLinks(loadedPage.Content, links)
		c.currentPagePath = ""

		if err := c.saveHTML(item.URL, localPageData); err != nil {
			log.Printf("failed to save html %s: %v", item.URL, err)
		}
	}
}

// replaceLinks — заменяем ссылки на страницы и ассеты
func (c *Crawler) replaceLinks(htmlData []byte, links *parser.ResourceLinks) []byte {
	if c.currentPagePath == "" {
		return htmlData
	}

	// Сначала страницы
	for _, pLink := range links.Pages {
		uPage, err := url.Parse(pLink)
		if err != nil || uPage.Host == "" {
			continue
		}

		if !c.isAllowedPage(uPage) || !c.isRobotAllowed(uPage) {
			continue
		}

		if c.isQueued(normalizeURL(uPage)) || c.isDownloadedURL(uPage) {
			targetLocal := c.makeLocalPath(uPage)
			relPath, err := filepath.Rel(filepath.Dir(c.currentPagePath), targetLocal)
			if err != nil {
				relPath = targetLocal
			}
			relPath = filepath.ToSlash(relPath)
			htmlData = bytes.ReplaceAll(htmlData, []byte(pLink), []byte(relPath))
		}
	}

	// Затем ассеты
	htmlData = replaceAssetLinksInHTML(c, htmlData)

	return htmlData
}

func replaceAssetLinksInHTML(c *Crawler, htmlData []byte) []byte {
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
						if newVal := makeLocalIfDownloaded(c, a.Val); newVal != "" {
							n.Attr[i].Val = newVal
						}

					case "srcset":
						fixed := fixSrcset(c, a.Val)
						if fixed != "" {
							n.Attr[i].Val = fixed
						}
					}
				}

			case "link", "a", "use":
				for i, a := range n.Attr {
					if a.Key == "href" || a.Key == "xlink:href" {
						if newVal := makeLocalIfDownloaded(c, a.Val); newVal != "" {
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

// возвращает локальный путь для скачанных ассетов
func makeLocalIfDownloaded(c *Crawler, link string) string {
	u, err := url.Parse(link)
	if err != nil || u.Host == "" {
		return ""
	}

	if !c.isDownloadedURL(u) {
		return ""
	}

	localPath := c.makeLocalPath(u)
	return c.makeRelativeToMirror(localPath)
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

// помечаем URL как скачанный
func (c *Crawler) markAsDownloadedURL(u *url.URL) {
	key := normalizeURL(u)
	if key == "" {
		return
	}
	c.Lock()
	c.Downloaded[key] = true
	c.Unlock()
}

// isDownloadedURL проверяет по нормализованному URL
func (c *Crawler) isDownloadedURL(u *url.URL) bool {
	key := normalizeURL(u)
	if key == "" {
		return false
	}
	c.RLock()
	defer c.RUnlock()
	return c.Downloaded[key]
}

// проверка разрешенности домена для страниц
func (c *Crawler) isAllowedPage(u *url.URL) bool {
	return u.Host == c.Domain
}

func (c *Crawler) downloadAssets(assetURLs []string) {
	for _, asset := range assetURLs {
		if idx := strings.Index(asset, "#"); idx != -1 {
			asset = asset[:idx] // убираем якорь чтобы корректно скачать svg-спрайты
		}

		u, err := url.Parse(asset)
		if err != nil || u.Host == "" {
			log.Printf("failed to parse asset URL %s: %v", asset, err)
			continue
		}

		if c.isDownloadedURL(u) {
			continue
		}

		c.wg.Add(1)
		go func(u *url.URL) {
			defer c.wg.Done()
			log.Printf("[asset] downloading: %s", u.String())
			result, err := c.Fetcher.Fetch(u, c.Timeout)
			if err != nil {
				log.Printf("failed to fetch asset %s: %v", u.String(), err)
				return
			}
			localPath := c.makeLocalPath(u)
			if _, err := c.Saver.Save(localPath, result.Content); err != nil {
				log.Printf("failed to save asset %s: %v", localPath, err)
			} else {
				c.markAsDownloadedURL(u)
			}
		}(u)
	}
}

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

// сохраняем HTML с правильной структурой и относительными путями
func (c *Crawler) saveHTML(u *url.URL, body []byte) error {
	local := c.makeLocalPath(u)
	log.Printf("[PAGE] downloading: %s", u.String())

	if err := os.MkdirAll(filepath.Dir(local), 0o755); err != nil {
		return err
	}

	fixed := c.makeRelativePaths(local, body)

	if err := os.WriteFile(local, fixed, 0o644); err != nil {
		return err
	}

	c.markAsDownloadedURL(u)

	return nil
}

func (c *Crawler) isRobotAllowed(u *url.URL) bool {
	for path := range c.RobotAllows {
		if strings.HasPrefix(u.Path, path) {
			return false
		}
	}
	return true
}

func (c *Crawler) LoadRobots() {
	u, err := url.Parse(fmt.Sprintf("https://%s/robots.txt", c.Domain))
	if err != nil {
		log.Printf("failed to load robots.txt: %v", err)
		return
	}

	data, err := c.Fetcher.Fetch(u, c.Timeout)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(data.Content))
	allowed := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "User-agent:") {
			allowed = strings.Contains(line, "*")
		} else if allowed && strings.HasPrefix(line, "Disallow:") {
			path := strings.TrimSpace(strings.TrimPrefix(line, "Disallow:"))
			if path != "" {
				c.RobotAllows[path] = false
			}
		}
	}
}

func (c *Crawler) isQueued(rawOrNorm string) bool {
	// rawOrNorm может быть уже нормализованной строкой (normalizeURL)
	for _, item := range c.Queue {
		if normalizeURL(item.URL) == rawOrNorm || item.URL.String() == rawOrNorm {
			return true
		}
	}
	return false
}

// заменяем URL в srcset на локальные пути из _assets.
func fixSrcset(c *Crawler, srcset string) string {
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

// возвращаем нормализованный URL для мапы Downloaded/Queue.
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
