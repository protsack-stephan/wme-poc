package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/protsack-stephan/wme-poc/models"
	"github.com/protsack-stephan/wme-poc/pkg/gql"

	"github.com/machinebox/graphql"
	"github.com/protsack-stephan/mediawiki-api-client"
	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
)

type projectQuery struct {
	GetProject *models.Project `json:"getProject"`
}

func main() {
	var batch int
	var project string
	var workers int
	var license string
	var url string
	flag.StringVar(&project, "p", "simplewiki", "Database name for the project")
	flag.IntVar(&workers, "w", 10, "Number of workers")
	flag.StringVar(&license, "l", "CC-BY-SA-4.0", "License identifier")
	flag.IntVar(&batch, "b", 50, "Number of pages in batch")
	flag.StringVar(&url, "u", "https://simple.wikipedia.org", "Project URL")
	flag.Parse()

	ctx := context.Background()
	dg := gql.NewClient()
	nsQuery := `
		query {
			getProject(identifier: "%s") {
				namespaces(filter: {identifier: {eq: 0}}) {
					id
					identifier
					name
				}
			}
		}
	`
	ns := &models.Namespace{}

	if err := dg.Run(ctx, graphql.NewRequest(fmt.Sprintf(nsQuery, project)), &projectQuery{&models.Project{Namespaces: []*models.Namespace{ns}}}); err != nil {
		log.Panicf("ns not found: %v", err)
	}

	mwiki := mediawiki.NewClient(url)
	dcl := dumps.NewClient()
	wg := new(sync.WaitGroup)
	titles := make(chan []string, workers)
	queue := []string{}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for ptitles := range titles {
				data, err := mwiki.PagesData(ctx, ptitles...)

				if err != nil {
					log.Printf("error fetching the data: %v\n", err)
				}

				query := gql.MutationQuery{
					Mutation: `
						mutation {
							addPage(input: [
								%s
							]) {
								numUids
							}
						}
					`,
				}

				for title, meta := range data {
					page := &models.Page{
						Name:         meta.Title,
						Identifier:   meta.PageID,
						Version:      meta.LastRevID,
						DateModified: meta.Touched.Format(time.RFC3339),
						URL:          fmt.Sprintf("%s/wiki/%s", url, title),
						Namespace: &models.Namespace{
							ID: ns.ID,
						},
						InLanguage: &models.Language{
							Identifier: meta.PageLanguage,
						},
						MainEntity: &models.QID{
							Identifier: meta.Pageprops.WikibaseItem,
						},
						IsPartOf: &models.Project{
							Identifier: project,
						},
					}

					htmlBody, err := mwiki.PageHTML(ctx, title, page.Version)

					if err != nil {
						return
					}

					body := struct {
						ArticleBody string `json:"articleBody"`
					}{
						string(htmlBody),
					}

					htmlj, err := json.Marshal(body)

					if err != nil {
						log.Println(err)
					}

					mainEntity := ""

					if len(page.MainEntity.Identifier) > 0 {
						mainEntity = fmt.Sprintf(`mainEntity: { identifier: "%s" }, `, page.MainEntity.Identifier)
					}

					pp := fmt.Sprintf(
						`{ name: "%s", identifier: %d, version: %d, dateModified: "%s", url: "%s", namespace: { id: "%s" }, inLanguage: { identifier: "%s" },%s %s, encodingFormat: "%s", isPartOf: { identifier: "%s"} }`,
						strings.ReplaceAll(page.Name, `"`, `\"`),
						page.Identifier,
						page.Version,
						page.DateModified,
						strings.ReplaceAll(page.URL, `"`, `\"`),
						page.Namespace.ID,
						page.InLanguage.Identifier,
						mainEntity,
						strings.TrimSuffix(strings.Replace(string(htmlj), `{"articleBody":`, `articleBody:`, 1), `}`),
						"text/html",
						page.IsPartOf.Identifier)

					query.Payload += pp
				}

				if err := dg.Run(ctx, graphql.NewRequest(fmt.Sprintf(query.Mutation, query.Payload)), &map[string]interface{}{}); err != nil {
					log.Println(err)
				}
			}
		}()
	}

	dcl.PageTitles(ctx, project, time.Now().UTC(), func(p *dumps.Page) {
		if len(queue) >= batch {
			tmp := make([]string, len(queue))
			copy(tmp, queue)
			queue = []string{}
			titles <- tmp
		}

		queue = append(queue, p.Title)
	})

	titles <- queue

	close(titles)
	wg.Wait()

}

func minify(data string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' {
			return -1
		}

		return r
	}, data)
}
