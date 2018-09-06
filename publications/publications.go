package publications

type publication int

// These constants represent various publications
const (
	WashingtonPost publication = iota + 1
	NewYorkTimes
)

// HostPublicationMap maps hostnames to their respective publications
var HostPublicationMap = map[string]publication{
	"www.nytimes.com":        NewYorkTimes,
	"www.washingtonpost.com": WashingtonPost,
}
