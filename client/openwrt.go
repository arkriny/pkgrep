package client

import (
	"fmt"
	"net/http"
)

type Openwrt struct {
	HTTPClient *http.Client
}

func (c Openwrt) Name() string {
	return "OpenWrt"
}

func (c Openwrt) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://openwrt.org/packages/pkgdata/%s", query)
	contains, err := getBodyContains(c.HTTPClient, url, "This topic does not exist yet")
	return !contains, err
}
