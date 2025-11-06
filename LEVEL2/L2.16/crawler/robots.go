package crawler

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strings"
)

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
