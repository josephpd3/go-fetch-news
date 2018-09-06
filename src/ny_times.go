/*
Package nytimes parses References from NYT Articles
*/
package nytimes

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const testPage = "https://www.nytimes.com/2018/09/05/business/media/new-york-times-trump-anonymous.html"

type Reference struct {
	Href    string
	Text    string
	Context string
}

func (r Reference) String() string {
	return fmt.Sprintf("Reference{ Href: '%v', Text: '%v', Context: '%v'}", r.Href, r.Text, r.Context)
}

func parseStory(index int, element *goquery.Selection) {
	element.Find("p").Each(parseParagraph)
}

func parseParagraph(index int, element *goquery.Selection) {
	paragraphText := element.Text()
	element.Find("a").Each(func(idx int, e *goquery.Selection) {
		href, exists := e.Attr("href")
		if exists {
			fmt.Println(Reference{
				Href:    href,
				Text:    e.Text(),
				Context: paragraphText,
			})
		}
	})
}

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Get(testPage)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	document.Find("article#story").First().Each(parseStory)
}
