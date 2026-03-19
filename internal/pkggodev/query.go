package pkggodev

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "pkg.go.dev"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pkg.go.dev/search?q=%s", query)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var buf strings.Builder
	inAnchor := false
	for scanner.Scan() {
		line := scanner.Text()
		if inAnchor {
			buf.WriteString(line)
			if strings.Contains(line, `class="SearchSnippet-header-path"`) {
				innerText := buf.String()

				if i := strings.Index(innerText, "<span"); i >= 0 {
					innerText = innerText[:i]
				}

				start := strings.Index(innerText, ">")
				if start >= 0 {
					if strings.TrimSpace(innerText[start+1:]) == query {
						return true, nil
					}
				}

				buf.Reset()
				inAnchor = false
			}
		} else {
			if strings.Contains(line, `data-gtmc="search result"`) {
				inAnchor = true
				buf.WriteString(line)
			}
		}
	}
	return false, scanner.Err()
}
