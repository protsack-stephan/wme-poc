package main

import (
	"context"
	"dgraph-tutorial/models"
	"fmt"
	"log"

	"dgraph-tutorial/pkg/gql"

	"github.com/machinebox/graphql"
	mediawiki "github.com/protsack-stephan/mediawiki-api-client"
)

const siteURL = "https://en.wikipedia.org"

func main() {
	ctx := context.Background()
	mwiki := mediawiki.NewClient(siteURL)
	smx, err := mwiki.Sitematrix(ctx)

	if err != nil {
		log.Panicf("sitematrix fetch error: %v", err)
	}

	queries := []*gql.MutationQuery{
		{
			Mutation: `
				mutation {
					addLanguage(input: [
						%s
					]) {
						numUids
					}
				}
			`,
		},
		{
			Mutation: `
				mutation {
					addProject(input: [
						%s
					]) {
						numUids
					}
				}
			`,
		},
	}

	for _, proj := range smx.Projects {
		lang := &models.Language{
			Name:          proj.Name,
			Identifier:    proj.Code,
			AlternateName: proj.Localname,
		}

		queries[0].Payload += fmt.Sprintf(`{name:"%s",identifier:"%s",alternateName:"%s"} `, lang.Name, lang.Identifier, lang.AlternateName)

		for _, site := range proj.Site {
			if !site.Closed {
				proj := &models.Project{
					Name:       site.Sitename,
					Identifier: site.DBName,
					URL:        site.URL,
					InLanguage: &models.Language{
						Identifier: lang.Identifier,
					},
				}

				queries[1].Payload += fmt.Sprintf(
					`{ name: "%s", identifier: "%s", url: "%s", inLanguage: { identifier: "%s" } }`,
					proj.Name,
					proj.Identifier,
					proj.URL,
					proj.InLanguage.Identifier)
			}
		}
	}

	client := gql.NewClient()

	for _, query := range queries {
		req := graphql.NewRequest(fmt.Sprintf(query.Mutation, query.Payload))

		res := map[string]interface{}{}

		if err := client.Run(ctx, req, &res); err != nil {
			log.Fatal(err)
		}
	}
}
