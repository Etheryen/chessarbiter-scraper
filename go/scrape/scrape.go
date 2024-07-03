package scrape

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/etheryen/chessarbiter-scraper/utils"
	"golang.org/x/net/html"
)

const baseURL string = "https://www.chessarbiter.com/index.php"

type SimpleTournament struct {
	Name        string
	Date        string
	Location    string
	TimeControl string
	Status      string
	Href        string
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

	return tournaments, nil
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
	status := getStatusFromTR(tr)
	name := tr.FirstChild.NextSibling.FirstChild.FirstChild.Data
	href := getHrefFromTR(tr)
	location := getLocationFromTR(tr)
	timeControl := tr.LastChild.LastChild.FirstChild.Data

	return SimpleTournament{
		Date:        date,
		Name:        name,
		Status:      status,
		Href:        href,
		Location:    location,
		TimeControl: timeControl,
	}
}

func getStatusFromTR(tr *html.Node) string {
	element := tr.FirstChild.LastChild.FirstChild.FirstChild

	if element != nil {
		return element.Data
	}
	return "zakończony"
}

func getLocationFromTR(tr *html.Node) string {
	location := tr.FirstChild.NextSibling.LastChild.FirstChild.Data
	location = strings.Split(location, " [")[0]
	location = strings.TrimSpace(location)

	if !strings.Contains(location, " ") {
		location = utils.ToCapitalizedUtf8(location)
	}

	return location
}

func getHrefFromTR(tr *html.Node) string {
	for _, attr := range tr.FirstChild.NextSibling.FirstChild.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}
