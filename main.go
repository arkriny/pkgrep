package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/arkriny/pkgrep/internal/archlinux"
	"github.com/arkriny/pkgrep/internal/pypi"
)

// QueryFunc is a function that accepts a search query string and returns a
// boolean indicating whether the package was found, or an error.
type QueryFunc func(query string) (bool, error)

type Repository struct {
	Name string
	Qf   QueryFunc
}

var repos = []Repository{
	{"Arch", archlinux.Query},
	{"PyPI", pypi.Query},
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

	for _, repo := range repos {
		found, err := repo.Qf(query)
		if err != nil {
			log.Fatal(err)
		}

		if found {
			fmt.Print("*")
		} else {
			fmt.Print("-")
		}
		fmt.Printf(" %s\n", repo.Name)
	}
}
