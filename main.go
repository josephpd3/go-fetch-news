package main

import (
	"github.com/josephpd3/gofetchnews/gofetchnews"
)

var testPages = []string{
	"https://www.washingtonpost.com/world/national-security/trump-administration-to-circumvent-court-limits-on-detention-of-child-migrants/2018/09/06/181d376c-b1bd-11e8-a810-4d6b627c3d5d_story.html",
	"https://www.nytimes.com/2018/09/05/business/media/new-york-times-trump-anonymous.html",
}

func main() {
	gofetchnews.FetchReferences(testPages...)
}
