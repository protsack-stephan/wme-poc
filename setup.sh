echo "Initializing schema..."
curl -X POST localhost:8080/admin/schema --data-binary '@schema.graphql'
echo "
Fetching Languages and Projects..."
go run setup/projects/main.go
echo "Fething Licenses..."
go run setup/licencses/main.go
echo "Fetching Namespaces..."
go run setup/namespaces/main.go
echo "Fetching Pages..."
go run setup/pages/main.go