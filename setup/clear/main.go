package main

import (
	"context"
	"log"

	"github.com/protsack-stephan/wme-poc/pkg/gql"

	"github.com/machinebox/graphql"
)

func main() {
	ctx := context.Background()

	queries := []string{
		`mutation {
			deleteLanguage(filter: {}) {
				numUids
			}
		}`,
		`mutation {
			deleteProject(filter: {}) {
				numUids
			}
		}`,
		`mutation {
			deleteNamespace(filter: {}) {
				numUids
			}
		}`,
		`mutation {
			deleteLicense(filter: {}) {
				numUids
			}
		}`,
		`mutation {
			deleteArticleBody(filter: {}) {
				numUids
			}
		}`,
		`mutation {
			deleteQID(filter: {}) {
				numUids
			}
		}`,
		`mutation {
			deletePage(filter: {}) {
				numUids
			}
		}`,
	}

	client := gql.NewClient()

	for _, query := range queries {
		req := graphql.NewRequest(query)

		if err := client.Run(ctx, req, &map[string]interface{}{}); err != nil {
			log.Fatal(err)
		}
	}
}
