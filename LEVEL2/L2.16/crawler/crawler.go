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
)

// Crawler - конфиг для обхода страниц
type Crawler struct {
	MaxDepth     int              // глубина рекурсии скачивания из os.Args
	Timeout      int              // таймаут соединения при скачивании по ссылкам
	InitURL      string           // введенный пользователем URL из os.Args
	Domain       string           // доменное имя InitURL
	Queue        []QueueItem      // сама очередь страниц
	Fetcher      *fetcher.Fetcher // Fetcher производит первичное скачивание по InitURL
	Saver        *saver.Saver     // Saver скачивает ассеты текущего URL из очереди
	Downloaded   map[string]bool  // карта для контроля скачанных ассетов/страниц - для избежания повторных скачиваний
	RobotAllows  map[string]bool  // для выполнения требований из robots.txt
	wg           sync.WaitGroup   // для контроля горутин скачивания ассетов
	sync.RWMutex                  // нужно для защиты карты Downloaded
}

type QueueItem struct {
	URL   *url.URL // распарсенный URL
	Depth int      // уровень рекурсии считая от InitURL
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

		// Грузим текущую страницу из очереди
		loadedPage, err := c.Fetcher.Fetch(item.URL, c.Timeout)
		if err != nil {
			log.Printf("failed to fetch %s: %v", item.URL.String(), err)
			continue
		}

		// Парсим текущую страницу - достаем ее ассеты и ссылки на другие страницы
		links, err := parser.ParseHTML(item.URL, loadedPage.Content)
		if err != nil {
			log.Printf("failed to parse html %s: %v", item.URL.String(), err)
			continue
		}

		// Сохраняем ассеты текущей страницы
		c.downloadAssets(links.Assets)
		c.wg.Wait()

		// Добавляем новые страницы в очередь
		c.updatePagesQueue(links.Pages, item.Depth)

		// Подменяем ссылки и сохраняем страницу
		localPageData := c.replaceLinks(loadedPage.Content, links, item.Depth)
		if err := c.saveHTML(item.URL, localPageData); err != nil {
			log.Printf("failed to save html %s: %v", item.URL, err)
		}

	}
}

func (c *Crawler) updatePagesQueue(newPages []string, parentDepth int) {
	if parentDepth+1 > c.MaxDepth {
		return
	}

	for _, page := range newPages {
		nextURL, err := url.Parse(page)
		if err == nil && c.isAllowedPage(nextURL) && c.isRobotAllowed(nextURL) && !c.isDownloaded(page) && !c.isQueued(page) {
			c.Queue = append(c.Queue, QueueItem{URL: nextURL, Depth: parentDepth + 1})
		}
	}
	log.Printf("Updated pages queue is: %v", c.Queue)
}

// Oтметка скачивания:
func (c *Crawler) markAsDownloaded(rawURL string) {
	c.Lock()
	c.Downloaded[rawURL] = true
	c.Unlock()
}

// Проверка скачивания:
func (c *Crawler) isDownloaded(rawURL string) bool {
	c.RLock()
	defer c.RUnlock()
	return c.Downloaded[rawURL]
}

// Для страниц - только строго тот же домен
func (c *Crawler) isAllowedPage(u *url.URL) bool {
	return u.Host == c.Domain
}

// Для ассетов - разрешаем поддомены
func (c *Crawler) isAllowedAsset(u *url.URL) bool {
	return strings.HasSuffix(u.Host, "."+c.Domain) || u.Host == c.Domain
}

// Скачивает и сохраняет ассеты текущей страницы из очереди
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

		if c.isDownloaded(asset) || !c.isAllowedAsset(u) {
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
				c.markAsDownloaded(u.String())
			}
		}(u)
	}
}

// localPath вычисляет путь сохранения с учётом того, что InitURL — корень
func (c *Crawler) makeLocalPath(u *url.URL) string {
	base, _ := url.Parse(c.InitURL)
	rel := strings.TrimPrefix(u.Path, "/")

	// если это самая первая страница — кладём прямо в корень
	if u.String() == base.String() {
		return filepath.Join("mirror", u.Host, "index.html")
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

// Своппинг ссылок на локальные пути
func (c *Crawler) replaceLinks(htmlData []byte, links *parser.ResourceLinks, currentDepth int) []byte {
	for _, pLink := range links.Pages {
		uPage, err := url.Parse(pLink)
		if err != nil || uPage.Host == "" {
			continue
		}

		// Проверяем свой-чужой(домен) и соответствие правилам robots.txt
		if !c.isAllowedPage(uPage) || !c.isRobotAllowed(uPage) {
			continue
		}

		// Проверяем глубину и очередь — чтобы менять ccылки только тех страниц, которые реально скачиваем или уже скачали
		if currentDepth+1 > c.MaxDepth || !c.isAllowedPage(uPage) {
			continue
		}

		// Формируем локальный абсолютный путь
		localPath := c.makeLocalPath(uPage)

		// Преобразуем в "корневой" URL относительно зеркала
		relPath := c.makeRelativeToMirror(localPath)

		htmlData = bytes.ReplaceAll(htmlData, []byte(pLink), []byte(relPath))
	}

	for _, aLink := range links.Assets {
		uAsset, err := url.Parse(aLink)
		if err != nil || uAsset.Host == "" {
			continue
		}

		if !c.isDownloaded(aLink) || !c.isAllowedAsset(uAsset) {
			continue
		}

		localPath := c.makeLocalPath(uAsset)
		relPath := c.makeRelativeToMirror(localPath)

		htmlData = bytes.ReplaceAll(htmlData, []byte(aLink), []byte(relPath))
	}

	return htmlData
}

// превращает абсолютный локальный путь типа mirror/habr.com/assets/* в /assets/*
func (c *Crawler) makeRelativeToMirror(localPath string) string {
	localPath = filepath.ToSlash(localPath)

	// Находим корневой префикс, который нужно обрезать
	prefix := filepath.ToSlash(filepath.Join("mirror", c.Domain)) + "/"

	// Убираем его из начала пути
	if strings.HasPrefix(localPath, prefix) {
		localPath = localPath[len(prefix)-1:] // сохраняем ведущий слеш
	}

	// Гарантируем, что путь начинается с '/'
	if !strings.HasPrefix(localPath, "/") {
		localPath = "/" + localPath
	}

	return localPath
}

// переписывает абсолютные пути внутри HTML в относительные
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

// сохраняет HTML с правильной структурой и относительными путями
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

	c.markAsDownloaded(u.String())

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

func (c *Crawler) isQueued(rawURL string) bool {
	for _, item := range c.Queue {
		if item.URL.String() == rawURL {
			return true
		}
	}
	return false
}
