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
	"github.com/arkriny/pkgrep/internal/npm"
	"github.com/arkriny/pkgrep/internal/pypi"
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
	{"NPM", npm.Query},
	{"PyPI", pypi.Query},
	{"Void", voidlinux.Query},
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
	results := make(chan string)

	for _, repo := range repos {
		wg.Add(1)
		go func(r Repository) {
			defer wg.Done()

			found, err := r.Qf(query)
			if err != nil {
				log.Println(err)
				return
			}

			var result string
			if found {
				result = fmt.Sprintf("* %s", r.Name)
			} else {
				result = fmt.Sprintf("- %s", r.Name)
			}
			results <- result
		}(repo)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
