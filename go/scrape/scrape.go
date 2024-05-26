package scrape

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

const baseURL string = "https://www.chessarbiter.com/index.php"

type SimpleTournament struct {
	Name     string
	Date     string
	Location string
	Type     string
	Status   string
	Href     string
}

func GetByYearMonth(year int, month int) ([]SimpleTournament, error) {
	urlStruct, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	queries := urlStruct.Query()

	queries.Add("rok", fmt.Sprint(year))
	queries.Add("miesiac", fmt.Sprint(month))

	urlStruct.RawQuery = queries.Encode()

	root, err := getRoot(urlStruct.String())
	if err != nil {
		return nil, err
	}

	tournaments := processHtml(root)

	fmt.Println(tournaments)

	return nil, nil
}

func getRoot(urlString string) (*html.Node, error) {
	resp, err := http.Get(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return html.Parse(resp.Body)
}

func processHtml(n *html.Node) []SimpleTournament {
	var tournaments []SimpleTournament

	if isRightElement(n) {
		t := getSimpleFromTR(n)

		tournaments = append(tournaments, t)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tournaments = append(tournaments, processHtml(c)...)
	}

	return tournaments
}

func isRightElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "tr" &&
		n.FirstChild.Data != "th" && isInRightTable(n)
}

func isInRightTable(n *html.Node) bool {
	if n != nil && n.Parent != nil && n.Parent.Parent != nil {
		for _, attr := range n.Parent.Parent.Attr {
			if attr.Key == "style" && attr.Val == "tbl" {
				return true
			}
		}
	}
	return false
}

func getSimpleFromTR(tr *html.Node) SimpleTournament {
	date := tr.FirstChild.FirstChild.Data
	status := tr.FirstChild.LastChild.FirstChild.FirstChild.Data
	name := tr.FirstChild.NextSibling.FirstChild.FirstChild.Data
	href := getHrefFromTR(tr)

	// TODO: get location as first word from another element

	return SimpleTournament{
		Date:   date,
		Name:   name,
		Status: status,
		Href:   href,
	}
}

func getHrefFromTR(tr *html.Node) string {
	for _, attr := range tr.FirstChild.NextSibling.FirstChild.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}
