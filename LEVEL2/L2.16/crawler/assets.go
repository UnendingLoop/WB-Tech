package crawler

import (
	"log"
	"net/url"
	"strings"
)

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
