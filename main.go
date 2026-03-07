package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/arkriny/pkgrep/internal/alpine"
	"github.com/arkriny/pkgrep/internal/archlinux"
	"github.com/arkriny/pkgrep/internal/aur"
	"github.com/arkriny/pkgrep/internal/cratesio"
	"github.com/arkriny/pkgrep/internal/debain"
	"github.com/arkriny/pkgrep/internal/fedora"
	"github.com/arkriny/pkgrep/internal/guix"
	"github.com/arkriny/pkgrep/internal/homebrew"
	"github.com/arkriny/pkgrep/internal/macports"
	"github.com/arkriny/pkgrep/internal/nixpkgs"
	"github.com/arkriny/pkgrep/internal/npm"
	"github.com/arkriny/pkgrep/internal/opensuse"
	"github.com/arkriny/pkgrep/internal/pypi"
	"github.com/arkriny/pkgrep/internal/rubygems"
	"github.com/arkriny/pkgrep/internal/snapcraft"
	"github.com/arkriny/pkgrep/internal/ubuntu"
	"github.com/arkriny/pkgrep/internal/voidlinux"
)

// QueryFunc is a function that accepts a search query string and returns a
// boolean indicating whether the package was found, or an error.
// The query string should be URL-safe.
type QueryFunc func(query string) (bool, error)

type Repository struct {
	Name string
	Qf   QueryFunc
}

var repos = []Repository{
	{"Alpine", alpine.Query},
	{"Arch", archlinux.Query},
	{"AUR", aur.Query},
	{"crates.io", cratesio.Query},
	{"Debian", debian.Query},
	{"Fedora", fedora.Query},
	{"Guix", guix.Query},
	{"Homebrew", homebrew.Query},
	{"MacPorts", macports.Query},
	{"Nixpkgs", nixpkgs.Query},
	{"NPM", npm.Query},
	{"openSUSE", opensuse.Query},
	{"PyPI", pypi.Query},
	{"RubyGems", rubygems.Query},
	{"Snapcraft", snapcraft.Query},
	{"Ubuntu", ubuntu.Query},
	{"Void", voidlinux.Query},
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

			found, err := r.Qf(query)
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
