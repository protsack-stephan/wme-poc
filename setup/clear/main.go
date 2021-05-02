package main

import (
	"context"
	"dgraph-tutorial/pkg/gql"
	"log"

	"github.com/machinebox/graphql"
)

func main() {
	ctx := context.Background()

	queries := []string{
		// `mutation {
		// 	deleteLanguage(filter: {}) {
		// 		numUids
		// 	}
		// }`,
		// `mutation {
		// 	deleteProject(filter: {}) {
		// 		numUids
		// 	}
		// }`,
		// `mutation {
		// 	deletePage(filter: {}) {
		// 		numUids
		// 	}
		// }`,
		`mutation {
			deleteNamespace(filter: {}) {
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
