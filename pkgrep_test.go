package main

import (
	"encoding/hex"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/arkriny/pkgrep/internal/alpine"
	"github.com/arkriny/pkgrep/internal/aosc"
	"github.com/arkriny/pkgrep/internal/archlinux"
	"github.com/arkriny/pkgrep/internal/aur"
	"github.com/arkriny/pkgrep/internal/chocolatey"
	"github.com/arkriny/pkgrep/internal/clojars"
	"github.com/arkriny/pkgrep/internal/cran"
	"github.com/arkriny/pkgrep/internal/cratesio"
	"github.com/arkriny/pkgrep/internal/debian"
	"github.com/arkriny/pkgrep/internal/dub"
	"github.com/arkriny/pkgrep/internal/fedora"
	"github.com/arkriny/pkgrep/internal/guix"
	"github.com/arkriny/pkgrep/internal/hackage"
	hexclient "github.com/arkriny/pkgrep/internal/hex"
	"github.com/arkriny/pkgrep/internal/homebrew"
	"github.com/arkriny/pkgrep/internal/julia"
	"github.com/arkriny/pkgrep/internal/kali"
	"github.com/arkriny/pkgrep/internal/luarocks"
	"github.com/arkriny/pkgrep/internal/macports"
	"github.com/arkriny/pkgrep/internal/melpa"
	"github.com/arkriny/pkgrep/internal/nixpkgs"
	"github.com/arkriny/pkgrep/internal/npm"
	"github.com/arkriny/pkgrep/internal/nuget"
	"github.com/arkriny/pkgrep/internal/opam"
	"github.com/arkriny/pkgrep/internal/pkggodev"
	"github.com/arkriny/pkgrep/internal/pubdev"
	"github.com/arkriny/pkgrep/internal/pypi"
	"github.com/arkriny/pkgrep/internal/rubygems"
	"github.com/arkriny/pkgrep/internal/scoop"
	"github.com/arkriny/pkgrep/internal/sisyphus"
	"github.com/arkriny/pkgrep/internal/snapcraft"
	"github.com/arkriny/pkgrep/internal/ubuntu"
	"github.com/arkriny/pkgrep/internal/voidlinux"
)

type testcase struct {
	querier Querier
	pkg     string
}

var testUserAgent = "pkgrep-test/0 (https://github.com/arkriny/pkgrep; mailto:arkriny@gmail.com)"

var testHTTPClient = &http.Client{
	Transport: &UserAgentRoundTripper{
		RoundTripper: http.DefaultTransport,
		UserAgent:    testUserAgent,
	},
	Timeout: time.Minute,
}

var tests = []testcase{
	{alpine.Client{testHTTPClient}, "alpine-base"},
	{aosc.Client{testHTTPClient}, "kernel-base"},
	{archlinux.Client{testHTTPClient}, "linux"},
	{aur.Client{testHTTPClient}, "google-chrome"},
	{chocolatey.Client{testHTTPClient}, "go"},
	{clojars.Client{testHTTPClient}, "core.typed"},
	{cran.Client{testHTTPClient}, "ggplot2"},
	{cratesio.Client{testHTTPClient}, "syn"},
	{debian.Client{testHTTPClient}, "linux-base"},
	{dub.Client{testHTTPClient}, "vibe-d"},
	{fedora.Client{testHTTPClient}, "linux-firmware"},
	{guix.Client{testHTTPClient}, "go"},
	{hackage.Client{testHTTPClient}, "ghc"},
	{hexclient.Client{testHTTPClient}, "phoenix"},
	{homebrew.Client{testHTTPClient}, "go"},
	{julia.Client{testHTTPClient}, "plots"},
	{kali.Client{testHTTPClient}, "linux"},
	{luarocks.Client{testHTTPClient}, "lua-cjson"},
	{macports.Client{testHTTPClient}, "go"},
	{melpa.Client{testHTTPClient}, "magit"},
	{nixpkgs.Client{testHTTPClient}, "home-manager"},
	{npm.Client{testHTTPClient}, "npm"},
	{nuget.Client{testHTTPClient}, "Azure.Core"},
	{opam.Client{testHTTPClient}, "ocaml"},
	{pkggodev.Client{testHTTPClient}, "http"},
	{pubdev.Client{testHTTPClient}, "http"},
	{pypi.Client{testHTTPClient}, "pip"},
	{rubygems.Client{testHTTPClient}, "rails"},
	{scoop.Client{testHTTPClient}, "go"},
	{sisyphus.Client{testHTTPClient}, "firmware-linux"},
	{snapcraft.Client{testHTTPClient}, "go"},
	{ubuntu.Client{testHTTPClient}, "linux-firmware"},
	{voidlinux.Client{testHTTPClient}, "linux-firmware"},
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
