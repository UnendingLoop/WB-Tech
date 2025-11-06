package crawler

import (
	"bytes"
	"net/url"
	"path/filepath"

	"miniWget/parser"
)

// replaceLinks - заменяем ссылки на страницы и ассеты
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
	htmlData = c.replaceAssetLinksInHTML(htmlData)

	return htmlData
}

// возвращает локальный путь для скачанных ассетов
func (c *Crawler) makeLocalIfDownloaded(link string) string {
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
