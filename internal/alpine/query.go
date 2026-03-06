package alpine

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

func Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pkgs.alpinelinux.org/packages?branch=edge&name=%s", query)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "No matching packages found...") {
			return false, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}

	return true, nil
}
