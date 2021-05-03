# Structure Content With DGraph POC

This is a simple setup to try and organize `simplewiki` data inside graph database that can be accessed through `graphql`. 

To run this project you need `docker` and `docker-compose` installed on your machine.

Getting started:
```bash
docker-compose up #wait till server fully starts might take upt to 30s

./setup.sh #last step of this process (fetching pages) will take couple of hours
```

After that's done you can access [http://localhost:8080/graphql](http://localhost:8080/graphql) with your favorite `graphql` client ([Altair](https://altair.sirmuel.design/) for example).

You can find the schema in `schema.graphql` file. Than will give you the idea of properties you can fetch.