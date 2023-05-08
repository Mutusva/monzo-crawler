package crawler

import (
	"context"
	"fmt"
	monzo_interview "github.com/Mutusva/monzo-webcrawler"
	"github.com/Mutusva/monzo-webcrawler/worker"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"time"
)

type htmlCrawler struct {
	Seed []string
}

func (h *htmlCrawler) Start(processExternal bool) error {
	visited := make(map[string]bool)
	filters := urlFilters(h.Seed, processExternal)

	queue := h.Seed
	urlWorker := worker.NewWorker(25, queue, processUrl)
	results := urlWorker.GetResultChan()
	go urlWorker.Run(filters, visited)
	displayResults(results)
	return nil
}

func displayResults(results <-chan map[string][]string) {
	for result := range results {
		for k := range result {
			fmt.Printf("-----------------links on %s url-------------------\n", k)
			fmt.Printf("%v \n", result[k])
		}
	}
}

// processUrl makes 'an' Http request to the url
func processUrl(curUrl string, filters []string) ([]string, error) {
	pageLinks := make([]string, 0)
	parsedUrl, err := ParseUrl(curUrl)
	if err != nil {
		fmt.Println("Could not parse Url: ", err)
		return pageLinks, err
	}

	req, err := http.NewRequest(http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		fmt.Printf("%v", err)
		return pageLinks, err
	}

	// Cancel the request if it does not respond in time
	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%v", err)
		return pageLinks, err
	}

	// if response status not 200 return
	if resp.StatusCode != http.StatusOK {
		return pageLinks, err
	}

	// Parse the response body as HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return pageLinks, err
	}

	err = resp.Body.Close()
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return pageLinks, err
	}

	pageLinks = findLinks(doc, filters, parsedUrl.Scheme, parsedUrl.Host)

	return pageLinks, nil
}

func ParseUrl(curUrl string) (*url.URL, error) {
	parsedUrl, err := url.Parse(curUrl)
	if err != nil {
		fmt.Println("Could not parse Url: ", err)
		return &url.URL{}, err
	}

	return parsedUrl, nil
}

// urlFilters get the hosts for seed url to only crawl a seed's domain
// Without going to external links or subdomains
func urlFilters(seed []string, processExternal bool) []string {
	filters := make([]string, 0)
	if processExternal {
		return filters
	}

	for _, u := range seed {
		url, err := ParseUrl(u)
		if err != nil {
			continue
		}
		filters = append(filters, url.Host)
	}

	return filters
}

// Find all links on a given HTML node
func findLinks(n *html.Node, filters []string, scheme, curHost string) []string {
	var links []string

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				url := attr.Val
				// Get all relative urls
				if len(url) > 0 && string(url[0]) == "/" {
					url = scheme + "://" + curHost + url
				}
				parseUrl, err := ParseUrl(url)
				if err != nil {
					continue
				}

				if len(filters) > 0 {
					if Contains(filters, parseUrl.Host) {
						links = append(links, parseUrl.String())
					}
				} else {
					links = append(links, parseUrl.String())
				}

			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, findLinks(c, filters, scheme, curHost)...)
	}
	return links
}

func New(seed []string) monzo_interview.Crawler {
	return &htmlCrawler{
		Seed: seed,
	}
}
