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
	"github.com/arkriny/pkgrep/internal/macports"
	"github.com/arkriny/pkgrep/internal/melpa"
	"github.com/arkriny/pkgrep/internal/nixpkgs"
	"github.com/arkriny/pkgrep/internal/npm"
	"github.com/arkriny/pkgrep/internal/nuget"
	"github.com/arkriny/pkgrep/internal/opam"
	"github.com/arkriny/pkgrep/internal/opensuse"
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
	name    string
	querier Querier
	pkg     string
}

var testHTTPClient = &http.Client{
	Timeout: time.Minute,
}

var tests = []testcase{
	{"Alpine", alpine.Client{testHTTPClient}, "alpine-base"},
	{"AOSC", aosc.Client{testHTTPClient}, "kernel-base"},
	{"Arch", archlinux.Client{testHTTPClient}, "linux"},
	{"AUR", aur.Client{testHTTPClient}, "google-chrome"},
	{"Chocolatey", chocolatey.Client{testHTTPClient}, "go"},
	{"CRAN", cran.Client{testHTTPClient}, "ggplot2"},
	{"crates.io", cratesio.Client{testHTTPClient}, "syn"},
	{"Debian", debian.Client{testHTTPClient}, "linux-base"},
	{"DUB", dub.Client{testHTTPClient}, "vibe-d"},
	{"Fedora", fedora.Client{testHTTPClient}, "linux-firmware"},
	{"Guix", guix.Client{testHTTPClient}, "go"},
	{"Hackage", hackage.Client{testHTTPClient}, "ghc"},
	{"Hex", hexclient.Client{testHTTPClient}, "phoenix"},
	{"Homebrew", homebrew.Client{testHTTPClient}, "go"},
	{"Julia", julia.Client{testHTTPClient}, "plots"},
	{"Kali", kali.Client{testHTTPClient}, "linux"},
	{"MacPorts", macports.Client{testHTTPClient}, "go"},
	{"MELPA", melpa.Client{testHTTPClient}, "magit"},
	{"Nixpkgs", nixpkgs.Client{testHTTPClient}, "home-manager"},
	{"NPM", npm.Client{testHTTPClient}, "npm"},
	{"NuGet", nuget.Client{testHTTPClient}, "Azure.Core"},
	{"opam", opam.Client{testHTTPClient}, "ocaml"},
	{"openSUSE", opensuse.Client{testHTTPClient}, "linux-firmware"},
	{"pub.dev", pubdev.Client{testHTTPClient}, "http"},
	{"PyPI", pypi.Client{testHTTPClient}, "pip"},
	{"RubyGems", rubygems.Client{testHTTPClient}, "rails"},
	{"Scoop", scoop.Client{testHTTPClient}, "go"},
	{"Sisyphus", sisyphus.Client{testHTTPClient}, "firmware-linux"},
	{"Snapcraft", snapcraft.Client{testHTTPClient}, "go"},
	{"Ubuntu", ubuntu.Client{testHTTPClient}, "linux-firmware"},
	{"Void", voidlinux.Client{testHTTPClient}, "linux-firmware"},
}

func Test(t *testing.T) {
	t.Parallel()

	b := make([]byte, 10)
	rand.Read(b)
	randomPackage := hex.EncodeToString(b)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			{
				found, err := tt.querier.Query(tt.pkg)
				if err != nil {
					t.Fatalf("%s: %v", tt.name, err)
				}

				if !found {
					t.Errorf("%s: %s not found", tt.name, tt.pkg)
				}
			}
			{
				found, err := tt.querier.Query(randomPackage)
				if err != nil {
					t.Fatalf("%s: %v", tt.name, err)
				}

				if found {
					t.Errorf("%s: %s found", tt.name, randomPackage)
				}
			}
		})
	}
}
