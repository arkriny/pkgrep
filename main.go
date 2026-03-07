package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/arkriny/pkgrep/internal/alpine"
	"github.com/arkriny/pkgrep/internal/archlinux"
	"github.com/arkriny/pkgrep/internal/aur"
	"github.com/arkriny/pkgrep/internal/chocolatey"
	"github.com/arkriny/pkgrep/internal/cratesio"
	"github.com/arkriny/pkgrep/internal/debian"
	"github.com/arkriny/pkgrep/internal/fedora"
	"github.com/arkriny/pkgrep/internal/guix"
	"github.com/arkriny/pkgrep/internal/homebrew"
	"github.com/arkriny/pkgrep/internal/kali"
	"github.com/arkriny/pkgrep/internal/macports"
	"github.com/arkriny/pkgrep/internal/nixpkgs"
	"github.com/arkriny/pkgrep/internal/npm"
	"github.com/arkriny/pkgrep/internal/opensuse"
	"github.com/arkriny/pkgrep/internal/pypi"
	"github.com/arkriny/pkgrep/internal/rubygems"
	"github.com/arkriny/pkgrep/internal/scoop"
	"github.com/arkriny/pkgrep/internal/sisyphus"
	"github.com/arkriny/pkgrep/internal/snapcraft"
	"github.com/arkriny/pkgrep/internal/ubuntu"
	"github.com/arkriny/pkgrep/internal/voidlinux"
)

type Querier interface {
	// Query accepts a search query string and returns a
	// boolean indicating whether the package was found, or an error.
	// The query string should be URL-safe.
	Query(query string) (bool, error)
}

type Repository struct {
	Name    string
	Querier Querier
}

var httpClient = &http.Client{
	// TODO(arkriny): make configurable.
	Timeout: 20 * time.Second,
}

var repos = []Repository{
	{"Alpine", alpine.Client{httpClient}},
	{"Arch", archlinux.Client{httpClient}},
	{"AUR", aur.Client{httpClient}},
	{"Chocolatey", chocolatey.Client{httpClient}},
	{"crates.io", cratesio.Client{httpClient}},
	{"Debian", debian.Client{httpClient}},
	{"Fedora", fedora.Client{httpClient}},
	{"Guix", guix.Client{httpClient}},
	{"Homebrew", homebrew.Client{httpClient}},
	{"Kali", kali.Client{httpClient}},
	{"MacPorts", macports.Client{httpClient}},
	{"Nixpkgs", nixpkgs.Client{httpClient}},
	{"NPM", npm.Client{httpClient}},
	{"openSUSE", opensuse.Client{httpClient}},
	{"PyPI", pypi.Client{httpClient}},
	{"RubyGems", rubygems.Client{httpClient}},
	{"Scoop", scoop.Client{httpClient}},
	{"Sisyphus", sisyphus.Client{httpClient}},
	{"Snapcraft", snapcraft.Client{httpClient}},
	{"Ubuntu", ubuntu.Client{httpClient}},
	{"Void", voidlinux.Client{httpClient}},
}

type Result struct {
	Name  string
	Found bool
}

func main() {
	log.SetPrefix("pkgrep: ")
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s QUERY\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	query := flag.Arg(0)
	query = url.QueryEscape(query)

	var wg sync.WaitGroup
	results := make(chan Result)

	for _, repo := range repos {
		wg.Add(1)
		go func(r Repository) {
			defer wg.Done()

			found, err := r.Querier.Query(query)
			if err != nil {
				log.Println(err)
				return
			}

			results <- Result{
				Name:  r.Name,
				Found: found,
			}
		}(repo)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	foundAny := false
	for result := range results {
		if result.Found {
			fmt.Printf("*")
			foundAny = true
		} else {
			fmt.Print("-")
		}
		fmt.Printf(" %s\n", result.Name)
	}

	if !foundAny {
		os.Exit(1)
	}
}
