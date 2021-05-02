package main

import (
	"context"
	"dgraph-tutorial/models"
	"dgraph-tutorial/pkg/gql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/machinebox/graphql"
	"github.com/protsack-stephan/mediawiki-api-client"
	dumps "github.com/protsack-stephan/mediawiki-dumps-client"
)

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
	mwiki := mediawiki.NewClient(url)
	wg := new(sync.WaitGroup)
	dcl := dumps.NewClient()
	titles := make(chan []string, workers)
	dg := gql.NewClient()
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
							Identifier: 0,
						},
						InLanguage: &models.Language{
							Identifier: meta.PageLanguage,
						},
						MainEntity: &models.QID{
							Identifier: meta.Pageprops.WikibaseItem,
						},
						ArticleBody: &models.ArticleBody{},
						IsPartOf: &models.Project{
							Identifier: project,
						},
					}

					errs := make(chan error, 2)

					go func() {
						data, err := mwiki.PageHTML(ctx, title, page.Version)

						if err != nil {
							errs <- err
							return
						}

						page.ArticleBody.HTML = minify(string(data))
						errs <- nil
					}()

					go func() {
						data, err := mwiki.PageWikitext(ctx, title, page.Version)

						if err != nil {
							errs <- err
							return
						}

						page.ArticleBody.Wiktext = minify(string(data))
						errs <- nil
					}()

					for i := 0; i < 2; i++ {
						err := <-errs

						if err != nil {
							log.Printf("page with title '%s' failed: %v\n", title, err)
						}
					}

					mainEntity := ""

					if len(page.MainEntity.Identifier) > 0 {
						mainEntity = fmt.Sprintf(`mainEntity: { identifier: "%s" }, `, page.MainEntity.Identifier)
					}

					body, err := json.Marshal(page.ArticleBody)

					if err != nil {
						log.Printf("error during matshal: %v\n", err)
					}

					articleBody := strings.Replace(strings.Replace(string(body), `"html":`, ` html:`, 1), `"wikitext":`, ` wikitext:`, 1)

					pp := fmt.Sprintf(
						`{ name: "%s", identifier: %d, version: %d, dateModified: "%s", url: "%s", namespace: { identifier: %d }, inLanguage: { identifier: "%s" },%s articleBody: %s, isPartOf: { identifier: "%s"} }`,
						strings.ReplaceAll(page.Name, `"`, `\"`),
						page.Identifier,
						page.Version,
						page.DateModified,
						strings.ReplaceAll(page.URL, `"`, `\"`),
						page.Namespace.Identifier,
						page.InLanguage.Identifier,
						mainEntity,
						articleBody,
						page.IsPartOf.Identifier)

					query.Payload += pp
				}

				if err := dg.Run(ctx, graphql.NewRequest(fmt.Sprintf(query.Mutation, query.Payload)), &map[string]interface{}{}); err != nil {
					log.Println(err)
					log.Println(query.Payload)
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
