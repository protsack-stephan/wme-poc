package main

import (
	"context"
	"log"

	"github.com/protsack-stephan/wme-poc/pkg/gql"

	"github.com/machinebox/graphql"
)

func main() {
	ctx := context.Background()
	client := gql.NewClient()
	query := gql.MutationQuery{
		Mutation: `
			mutation {
				addLicense(input: [
					{ identifier: "CC-BY-SA-1.0", name: "Creative Commons Attribution Share Alike 1.0 Generic" }
					{ identifier: "CC-BY-SA-2.0", name: "Creative Commons Attribution Share Alike 2.0 Generic" }
					{ identifier: "CC-BY-SA-2.0-UK", name: "Creative Commons Attribution Share Alike 2.0 England and Wales" }
					{ identifier: "CC-BY-SA-2.1-JP", name: "Creative Commons Attribution Share Alike 2.1 Japan" }
					{ identifier: "CC-BY-SA-2.5", name: "Creative Commons Attribution Share Alike 2.5 Generic" }
					{ identifier: "CC-BY-SA-3.0", name: "Creative Commons Attribution Share Alike 3.0 Unported" }
					{ identifier: "CC-BY-SA-3.0-AT", name: "Creative Commons Attribution-Share Alike 3.0 Austria" }
					{ identifier: "CC-BY-SA-4.0", name: "Creative Commons Attribution Share Alike 4.0 International" }
				]) {
					numUids
				}
			}
		`,
	}

	if err := client.Run(ctx, graphql.NewRequest(query.Mutation), &map[string]interface{}{}); err != nil {
		log.Panicf("mutation failed: %v", err)
	}
}
