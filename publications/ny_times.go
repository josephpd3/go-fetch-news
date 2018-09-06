/*
Package publications parses articles from various publications

This file is specific to The New York Times
*/
package publications

import (
    "errors"

    "github.com/PuerkitoBio/goquery"
    "github.com/josephpd3/gofetchnews/refs"
)

// FetchNYTReferences fetches all references in a New York Times article
// If there are issues parsing the article's body, an error will be returned
func FetchNYTReferences(article *goquery.Document) (<-chan refs.Reference, error) {
    ch := make(chan refs.Reference)

    articleBody := article.Find("article#story").First()
    // Check to see if we have a match on our article body selector
    if len(articleBody.Nodes) != 0 {
        go parseNYTArticleBody(articleBody, ch)
        return ch, nil
    }

    return nil, errors.New("could not parse article")
}

// Internal helper to parse paragraphs in article and delegate to parseNYTParagraph
func parseNYTArticleBody(articleBody *goquery.Selection, ch chan<- refs.Reference) {
    articleBody.Find("p").Each(refs.CallWithRefChannel(refs.GQChannelAndCallback{
        Channel:  ch,
        Callback: parseNYTParagraph,
    }))
    close(ch)
}

// Internal helper to parse links out of paragraphs to construct References to send back via channel
func parseNYTParagraph(index int, element *goquery.Selection, ch chan<- refs.Reference) {
    paragraphText := element.Text()
    element.Find("a").Each(func(idx int, e *goquery.Selection) {
        href, exists := e.Attr("href")
        if exists {
            ch <- refs.Reference{
                Href:    href,
                Text:    e.Text(),
                Context: paragraphText,
            }
        }
    })
}
