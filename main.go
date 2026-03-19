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
	"github.com/arkriny/pkgrep/internal/opensuse"
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

var httpClient = &http.Client{
	Transport: &UserAgentRoundTripper{
		RoundTripper: http.DefaultTransport,
		UserAgent:    "pkgrep/0 (https://github.com/arkriny/pkgrep; mailto:arkriny@gmail.com)",
	},
	// TODO(arkriny): make configurable.
	Timeout: 20 * time.Second,
}

var queiers = []Querier{
	alpine.Client{httpClient},
	aosc.Client{httpClient},
	archlinux.Client{httpClient},
	aur.Client{httpClient},
	chocolatey.Client{httpClient},
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
	opensuse.Client{httpClient},
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

var flagDryRun = flag.Bool("dry-run", false, "do everything except actually send the requests")
var flagList = flag.Bool("list", false, "list repositories")
var flagVersion = flag.Bool("version", false, "print version")
var flagInclude repoList
var flagExclude repoList

func init() {
	flag.Var(&flagInclude, "include", "search in specified repositories only")
	flag.Var(&flagExclude, "exclude", "skip specified repositories")
}

// Checks if repository name is explicitly excluded or not included via flags.
func shouldSkipRepository(repoName string) bool {
	nameLower := strings.ToLower(repoName)
	excluded := slices.Contains(flagExclude, nameLower)
	included := len(flagInclude) == 0 || slices.Contains(flagInclude, nameLower)
	return excluded || !included
}

func main() {
	log.SetPrefix("pkgrep: ")
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s QUERY\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *flagVersion {
		if build, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(build.Main.Version)
		} else {
			fmt.Println("unknown version")
		}
		return
	}

	if *flagList {
		for _, querier := range queiers {
			fmt.Println(querier.Name())
		}
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		log.Print("missing query")
		flag.Usage()
		os.Exit(2)
	}

	query := flag.Arg(0)
	err := safeURLSegment(query)
	if err != nil {
		log.Fatal("invalid query: ", err)
	}

	type Result struct {
		Name  string
		Found bool
	}
	results := make(chan Result)

	var wg sync.WaitGroup
	for _, querier := range queiers {
		if shouldSkipRepository(querier.Name()) {
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
					log.Println(err)
					return
				}
			}

			results <- Result{
				Name:  q.Name(),
				Found: found,
			}
		}(querier)
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
