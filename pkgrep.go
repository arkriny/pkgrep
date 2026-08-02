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
	"time"

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
		os.Exit(0)
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

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s QUERY\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

// Checks if repository name is explicitly excluded or not included via flags.
func shouldSkipRepository(name string) bool {
	nameLower := strings.ToLower(name)
	excluded := slices.Contains(flagExclude, nameLower)
	included := len(flagInclude) == 0 || slices.Contains(flagInclude, nameLower)
	return excluded || !included
}
