package luarocks

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://luarocks.org/search?q=%s", query)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Response HTML has the following structure:
	//
	// <h2>Modules</h2>
	// <ul>
	//   ...
	//   <a>_package_</a>
	//   ...
	// </ul>
	// <h2>Users</h2>
	//
	// Scan it until the header and then look for first <a> tag, assuming that
	// the list is sorted and exact match would be first.
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanWords)
	inModules := false
	for scanner.Scan() {
		line := scanner.Text()
		if inModules {
			if strings.Contains(line, "</a>") {
				return strings.Contains(line, fmt.Sprintf(">%s</a>", query)), nil
			}
			if strings.Contains(line, "<h2>Users</h2>") {
				return false, nil
			}
		} else {
			if strings.Contains(line, "<h2>Modules</h2>") {
				inModules = true
			}
		}
	}
	return false, scanner.Err()
}
