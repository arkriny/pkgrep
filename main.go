package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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
}

type Repository struct {
	Name    string
	Querier Querier
}

var httpClient = &http.Client{
	Transport: &UserAgentRoundTripper{
		RoundTripper: http.DefaultTransport,
		UserAgent:    "pkgrep/0 (https://github.com/arkriny/pkgrep; mailto:arkriny@gmail.com)",
	},
	// TODO(arkriny): make configurable.
	Timeout: 20 * time.Second,
}

var repos = []Repository{
	{"Alpine", alpine.Client{httpClient}},
	{"AOSC", aosc.Client{httpClient}},
	{"Arch", archlinux.Client{httpClient}},
	{"AUR", aur.Client{httpClient}},
	{"Chocolatey", chocolatey.Client{httpClient}},
	{"CRAN", cran.Client{httpClient}},
	{"crates.io", cratesio.Client{httpClient}},
	{"Debian", debian.Client{httpClient}},
	{"DUB", dub.Client{httpClient}},
	{"Fedora", fedora.Client{httpClient}},
	{"Guix", guix.Client{httpClient}},
	{"Hackage", hackage.Client{httpClient}},
	{"Hex", hex.Client{httpClient}},
	{"Homebrew", homebrew.Client{httpClient}},
	{"Julia", julia.Client{httpClient}},
	{"Kali", kali.Client{httpClient}},
	{"LuaRocks", luarocks.Client{httpClient}},
	{"MacPorts", macports.Client{httpClient}},
	{"MELPA", melpa.Client{httpClient}},
	{"Nixpkgs", nixpkgs.Client{httpClient}},
	{"NPM", npm.Client{httpClient}},
	{"NuGet", nuget.Client{httpClient}},
	{"opam", opam.Client{httpClient}},
	{"openSUSE", opensuse.Client{httpClient}},
	{"pkg.go.dev", pkggodev.Client{httpClient}},
	{"pub.dev", pubdev.Client{httpClient}},
	{"PyPI", pypi.Client{httpClient}},
	{"RubyGems", rubygems.Client{httpClient}},
	{"Scoop", scoop.Client{httpClient}},
	{"Sisyphus", sisyphus.Client{httpClient}},
	{"Snapcraft", snapcraft.Client{httpClient}},
	{"Ubuntu", ubuntu.Client{httpClient}},
	{"Void", voidlinux.Client{httpClient}},
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
	included := len(flagInclude) > 0 && slices.Contains(flagInclude, nameLower)
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

	if *flagList {
		for _, repo := range repos {
			fmt.Println(repo.Name)
		}
		os.Exit(0)
	}

	if flag.NArg() != 1 {
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
	for _, repo := range repos {
		if shouldSkipRepository(repo.Name) {
			continue
		}
		wg.Add(1)
		go func(r Repository) {
			defer wg.Done()

			found := false
			if !*flagDryRun {
				var err error
				found, err = r.Querier.Query(query)
				if err != nil {
					log.Println(err)
					return
				}
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
