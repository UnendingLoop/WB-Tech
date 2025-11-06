// Package crawler contains core-logic of wget
package crawler

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"miniWget/fetcher"
	"miniWget/parser"
	"miniWget/saver"
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

// QueueItem - минимальный элемент в очереди на скачивание
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

func (c *Crawler) isQueued(rawOrNorm string) bool {
	// rawOrNorm может быть уже нормализованной строкой (normalizeURL)
	for _, item := range c.Queue {
		if normalizeURL(item.URL) == rawOrNorm || item.URL.String() == rawOrNorm {
			return true
		}
	}
	return false
}
