package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/protsack-stephan/wme-poc/models"
	"github.com/protsack-stephan/wme-poc/pkg/gql"

	"github.com/machinebox/graphql"
	"github.com/protsack-stephan/mediawiki-api-client"
)

const siteURL = "https://en.wikipedia.org"

var defaults = map[int]string{
	0: "Article",
}

type projectQuery struct {
	QueryProject []*models.Project `json:"queryProject"`
}

func main() {
	var workers int
	flag.IntVar(&workers, "w", 5, "Number of workers for concurrency")
	flag.Parse()

	ctx := context.Background()
	client := gql.NewClient()
	res := &projectQuery{[]*models.Project{}}

	err := client.Run(ctx, graphql.NewRequest(`
		query {
			queryProject(filter: {}) {
				identifier
				url
				inLanguage {
					identifier
				}
			}
		}
	`), res)

	if err != nil {
		log.Panicf("project query failed: %v", err)
	}

	if len(res.QueryProject) == 0 {
		log.Panic("first run `projects' job")
	}

	wg := new(sync.WaitGroup)
	queue := make(chan *models.Project, workers)

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			for project := range queue {
				mwiki := mediawiki.NewClient(project.URL)
				namespaces, err := mwiki.Namespaces(ctx)

				if err != nil {
					log.Panicf("can't fetch '%s' namespace: %v", project.Identifier, err)
				}

				query := gql.MutationQuery{
					Mutation: `
						mutation {
							addNamespace(input: [
								%s
							]) {
								numUids
							}
						}
					`,
				}

				for _, ns := range namespaces {
					if name, ok := defaults[ns.ID]; len(ns.Name) <= 0 && ok {
						ns.Name = name
					}

					query.Payload += fmt.Sprintf(
						`{ name: "%s", identifier: %d, inLanguage: { identifier: "%s" }, isPartOf: { identifier: "%s" }}`,
						ns.Name,
						ns.ID,
						project.InLanguage.Identifier,
						project.Identifier)
				}

				if err := client.Run(ctx, graphql.NewRequest(fmt.Sprintf(query.Mutation, query.Payload)), &map[string]interface{}{}); err != nil {
					log.Panicf("mutation for '%s' failed: %v", project.Identifier, err)
				}
			}
		}()
	}

	for _, project := range res.QueryProject {
		queue <- project
	}

	close(queue)
	wg.Wait()
}
