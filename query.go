package main

import (
	"fmt"
	"log"
	"sync"
)

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
					log.Println(err)
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
