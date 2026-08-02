package main

import (
	"encoding/hex"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/arkriny/pkgrep/client"
)

type testcase struct {
	querier Querier
	pkg     string
}

var testUserAgent = "pkgrep-test"

var testHTTPClient = &http.Client{
	Transport: &UserAgentRoundTripper{
		RoundTripper: http.DefaultTransport,
		UserAgent:    testUserAgent + "/0 " + UserAgentInfo,
	},
	Timeout: time.Minute,
}

var tests = []testcase{
	{client.Alpine{HTTPClient: testHTTPClient}, "alpine-base"},
	{client.Aosc{HTTPClient: testHTTPClient}, "kernel-base"},
	{client.Archlinux{HTTPClient: testHTTPClient}, "linux"},
	{client.Aur{HTTPClient: testHTTPClient}, "google-chrome"},
	{client.Chocolatey{HTTPClient: testHTTPClient}, "go"},
	{client.Clojars{HTTPClient: testHTTPClient}, "core.typed"},
	{client.Cran{HTTPClient: testHTTPClient}, "ggplot2"},
	{client.Cratesio{HTTPClient: testHTTPClient}, "syn"},
	{client.Debian{HTTPClient: testHTTPClient}, "linux-base"},
	{client.Dub{HTTPClient: testHTTPClient}, "vibe-d"},
	{client.Fedora{HTTPClient: testHTTPClient}, "linux-firmware"},
	{client.Guix{HTTPClient: testHTTPClient}, "go"},
	{client.Hackage{HTTPClient: testHTTPClient}, "ghc"},
	{client.Hex{HTTPClient: testHTTPClient}, "phoenix"},
	{client.Homebrew{HTTPClient: testHTTPClient}, "go"},
	{client.Julia{HTTPClient: testHTTPClient}, "plots"},
	{client.Kali{HTTPClient: testHTTPClient}, "linux"},
	{client.Luarocks{HTTPClient: testHTTPClient}, "lua-cjson"},
	{client.Macports{HTTPClient: testHTTPClient}, "go"},
	{client.Melpa{HTTPClient: testHTTPClient}, "magit"},
	{client.Nixpkgs{HTTPClient: testHTTPClient}, "home-manager"},
	{client.Npm{HTTPClient: testHTTPClient}, "npm"},
	{client.Nuget{HTTPClient: testHTTPClient}, "Azure.Core"},
	{client.Opam{HTTPClient: testHTTPClient}, "ocaml"},
	{client.Openwrt{HTTPClient: testHTTPClient}, "opkg"},
	{client.Pkggodev{HTTPClient: testHTTPClient}, "http"},
	{client.Pubdev{HTTPClient: testHTTPClient}, "http"},
	{client.Pypi{HTTPClient: testHTTPClient}, "pip"},
	{client.Rubygems{HTTPClient: testHTTPClient}, "rails"},
	{client.Scoop{HTTPClient: testHTTPClient}, "go"},
	{client.Sisyphus{HTTPClient: testHTTPClient}, "firmware-linux"},
	{client.Snapcraft{HTTPClient: testHTTPClient}, "go"},
	{client.Ubuntu{HTTPClient: testHTTPClient}, "linux-firmware"},
	{client.Voidlinux{HTTPClient: testHTTPClient}, "linux-firmware"},
}

func Test(t *testing.T) {
	t.Parallel()

	b := make([]byte, 10)
	rand.Read(b)
	randomPackage := hex.EncodeToString(b)

	for _, tt := range tests {
		name := tt.querier.Name()
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			{
				found, err := tt.querier.Query(tt.pkg)
				if err != nil {
					t.Fatalf("%s: %v", name, err)
				}

				if !found {
					t.Errorf("%s: %s not found", name, tt.pkg)
				}
			}
			{
				found, err := tt.querier.Query(randomPackage)
				if err != nil {
					t.Fatalf("%s: %v", name, err)
				}

				if found {
					t.Errorf("%s: %s found", name, randomPackage)
				}
			}
		})
	}
}
