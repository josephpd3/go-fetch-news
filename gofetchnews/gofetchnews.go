package gofetchnews

import (
    "fmt"
    "github.com/josephpd3/gofetchnews/refs"
    "log"
    "net/http"
    "net/url"
    "time"

    "github.com/PuerkitoBio/goquery"
    "github.com/josephpd3/gofetchnews/publications"
)

// FetchReferences takes a variable number of URLs for news articles and
// attempts to retrieve linked references from within them
func FetchReferences(articleURLs ...string) {
    httpClient := http.Client{
        Timeout: 30 * time.Second,
    }

    // Collect all channels in a list to merge after parsing all URLs for processing
    var allReferences []<-chan refs.Reference

    for _, articleURL := range articleURLs {
        // Fetch a goquery.Document of the article HTML
        article := fetchArticle(httpClient, articleURL)

        // Parse the URL to retrieve the hostname and map to a known publication
        parsedURL, err := url.Parse(articleURL)
        if err != nil {
            log.Fatal(err)
        }

        pub, ok := publications.HostPublicationMap[parsedURL.Hostname()]
        if !ok {
            log.Fatal("Publication not recognized!")
        }

        // Depending on the publication, fetch references in that particular site's article format
        switch pub {
        case publications.WashingtonPost:
            references, err := publications.FetchWaPoReferences(article)
            if err != nil {
                log.Fatal(err)
            } else {
                allReferences = append(allReferences, references)
            }
        case publications.NewYorkTimes:
            references, err := publications.FetchNYTReferences(article)
            if err != nil {
                log.Fatal(err)
            } else {
                allReferences = append(allReferences, references)
            }
        default:
            log.Fatal("Parser not written for publication!")
        }

    }

    printReferences(refs.MergeOutboundReferenceChannels(allReferences...))
}

// Internal helper for fetching a goquery.Document given an http.Client and a URL
func fetchArticle(httpClient http.Client, url string) *goquery.Document {
    response, err := httpClient.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    document, err := goquery.NewDocumentFromResponse(response)
    if err != nil {
        log.Fatal(err)
    }
    return document
}

// Internal helper for printing references as a test
func printReferences(ch <-chan refs.Reference) {
    for ref := range ch {
        fmt.Println(ref)
    }
}
