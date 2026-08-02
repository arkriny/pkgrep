// Pkgrep queries multiple package repositories by package name.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/arkriny/pkgrep/client"
)

type Querier interface {
	// Query accepts a search query string and returns a
	// boolean indicating whether the package was found, or an error.
	Query(query string) (bool, error)

	// Name returns the name of the querier, usually a repository name.
	Name() string
}

type repoList []string

func (rl *repoList) String() string {
	return fmt.Sprint(*rl)
}

func (rl *repoList) Set(value string) error {
	for p := range strings.SplitSeq(value, ",") {
		*rl = append(*rl, strings.ToLower(p))
	}
	return nil
}

var (
	flagDryRun  = flag.Bool("dry-run", false, "do everything except actually send the requests")
	flagList    = flag.Bool("list", false, "list repositories")
	flagVersion = flag.Bool("version", false, "print version")
	flagTimeout = flag.Int64("timeout", 0, "time limit for a single request in seconds")
	flagInclude repoList
	flagExclude repoList
)

func init() {
	flag.Var(&flagInclude, "include", "search in specified repositories only")
	flag.Var(&flagExclude, "exclude", "skip specified repositories")
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s QUERY\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *flagVersion {
		if build, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(build.Main.Version)
		} else {
			fmt.Println("unknown version")
		}
		return
	}

	httpClient := &http.Client{
		Transport: &UserAgentRoundTripper{
			RoundTripper: http.DefaultTransport,
			UserAgent:    "pkgrep/0 (https://github.com/arkriny/pkgrep; mailto:arkriny@gmail.com)",
		},
		Timeout: time.Duration(*flagTimeout) * time.Second,
	}

	queriers := []Querier{
		client.Alpine{httpClient},
		client.Aosc{httpClient},
		client.Archlinux{httpClient},
		client.Aur{httpClient},
		client.Chocolatey{httpClient},
		client.Clojars{httpClient},
		client.Cran{httpClient},
		client.Cratesio{httpClient},
		client.Debian{httpClient},
		client.Dub{httpClient},
		client.Fedora{httpClient},
		client.Guix{httpClient},
		client.Hackage{httpClient},
		client.Hex{httpClient},
		client.Homebrew{httpClient},
		client.Julia{httpClient},
		client.Kali{httpClient},
		client.Luarocks{httpClient},
		client.Macports{httpClient},
		client.Melpa{httpClient},
		client.Nixpkgs{httpClient},
		client.Npm{httpClient},
		client.Nuget{httpClient},
		client.Opam{httpClient},
		client.Openwrt{httpClient},
		client.Pkggodev{httpClient},
		client.Pubdev{httpClient},
		client.Pypi{httpClient},
		client.Rubygems{httpClient},
		client.Scoop{httpClient},
		client.Sisyphus{httpClient},
		client.Snapcraft{httpClient},
		client.Ubuntu{httpClient},
		client.Voidlinux{httpClient},
	}

	if *flagList {
		for _, querier := range queriers {
			fmt.Println(querier.Name())
		}
		return
	}
	if flag.NArg() != 1 {
		log.Print("missing query")
		usage()
	}

	query := flag.Arg(0)
	foundAny := runQuery(queriers, query)
	if !foundAny {
		os.Exit(1)
	}
}

func runQuery(queriers []Querier, query string) bool {
	if err := safeURLSegment(query); err != nil {
		log.Fatal("invalid query: ", err)
	}

	type Result struct {
		Name  string
		Found bool
	}
	results := make(chan Result)
	var wg sync.WaitGroup
	for _, q := range queriers {
		if shouldSkipRepository(q.Name()) {
			continue
		}
		wg.Add(1)
		go func(q Querier) {
			defer wg.Done()

			found := false
			if !*flagDryRun {
				var err error
				found, err = q.Query(query)
				if err != nil {
					log.Printf("query '%s': %s", q.Name(), err)
					return
				}
			}
			results <- Result{
				Name:  q.Name(),
				Found: found,
			}
		}(q)
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
	return foundAny
}

func shouldSkipRepository(name string) bool {
	nameLower := strings.ToLower(name)
	excluded := slices.Contains(flagExclude, nameLower)
	included := len(flagInclude) == 0 || slices.Contains(flagInclude, nameLower)
	return excluded || !included
}

// safeURLSegment checks whether a string can be safely placed in URL segment.
func safeURLSegment(s string) error {
	for _, r := range s {
		if !unicode.IsLetter(r) &&
			!unicode.IsDigit(r) &&
			r != '-' &&
			r != '.' &&
			r != '_' &&
			r != '~' {
			return fmt.Errorf("disallowed character in URL %q", r)
		}
	}
	return nil
}
