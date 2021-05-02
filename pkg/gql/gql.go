package gql

import "github.com/machinebox/graphql"

// NewClient get new graphql client
func NewClient() *graphql.Client {
	return graphql.NewClient("http://localhost:8080/graphql")
}

// MutationQuery query to mutate the info
type MutationQuery struct {
	Mutation string
	Payload  string
}
