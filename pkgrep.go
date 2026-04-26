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
	"github.com/arkriny/pkgrep/internal/hex"
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
	log.SetPrefix("pkgrep: ")
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
		alpine.Client{httpClient},
		aosc.Client{httpClient},
		archlinux.Client{httpClient},
		aur.Client{httpClient},
		chocolatey.Client{httpClient},
		clojars.Client{httpClient},
		cran.Client{httpClient},
		cratesio.Client{httpClient},
		debian.Client{httpClient},
		dub.Client{httpClient},
		fedora.Client{httpClient},
		guix.Client{httpClient},
		hackage.Client{httpClient},
		hex.Client{httpClient},
		homebrew.Client{httpClient},
		julia.Client{httpClient},
		kali.Client{httpClient},
		luarocks.Client{httpClient},
		macports.Client{httpClient},
		melpa.Client{httpClient},
		nixpkgs.Client{httpClient},
		npm.Client{httpClient},
		nuget.Client{httpClient},
		opam.Client{httpClient},
		pkggodev.Client{httpClient},
		pubdev.Client{httpClient},
		pypi.Client{httpClient},
		rubygems.Client{httpClient},
		scoop.Client{httpClient},
		sisyphus.Client{httpClient},
		snapcraft.Client{httpClient},
		ubuntu.Client{httpClient},
		voidlinux.Client{httpClient},
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
