package refs

import (
	"fmt"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Reference refers to a link within an article and its context (containing paragraph or sentence)
type Reference struct {
	Source  string
	Href    string
	Text    string
	Context string
}

func (r Reference) String() string {
	return fmt.Sprintf("Reference{ Source: '%v', Href: '%v', Text: '%v', Context: '%v'}", r.Source, r.Href, r.Text, r.Context)
}

// GQChannelAndCallback wraps an output channel and a callback for use in GoQuery selections
type GQChannelAndCallback struct {
	Channel  chan<- Reference
	Callback func(int, *goquery.Selection, chan<- Reference)
}

// CallWithRefChannel creates a function that injects an input channel implicitly via a closure
func CallWithRefChannel(channelAndCallback GQChannelAndCallback) func(int, *goquery.Selection) {
	return func(index int, element *goquery.Selection) {
		channelAndCallback.Callback(index, element, channelAndCallback.Channel)
	}
}

// MergeOutboundReferenceChannels merges output channels of References
// to create a single output channel of References
// - - -
// Courtesy of Francesc Campoy:
// https://medium.com/justforfunc/two-ways-of-merging-n-channels-in-go-43c0b57cd1de
// - - -
func MergeOutboundReferenceChannels(cs ...<-chan Reference) <-chan Reference {
	out := make(chan Reference)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan Reference) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
