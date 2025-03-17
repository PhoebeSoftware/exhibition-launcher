package torrent

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Game1337x struct {
	Title     string
	Uploader  string
	Downloads int
	Date      string
	Magnet    string
}

func Scrape_1337x(item string) []string {
	billCollector := colly.NewCollector()
	results := []string{}

	billCollector.OnHTML("tbody tr td:first-of-type a:nth-of-type(2)", func(e *colly.HTMLElement) {
		results = append(results, e.Attr("href"))
	})

	billCollector.Visit("https://1337x.to/sort-category-search/" + strings.Replace(item, " ", "%20", -1) + "/Games/time/desc/1/")

	billCollector.Wait()
	return results
}

func Get_1337x_data(torrent_page string) Game1337x {
	billCollector := colly.NewCollector()

	results := Game1337x{
		Title:     "Unknown",
		Uploader:  "Unknown",
		Downloads: 0,
		Date:      "Unknown",
	}

	billCollector.OnHTML(".box-info-heading.clearfix", func(e *colly.HTMLElement) {
		results.Title = strings.Trim(e.Text, " ")
	})

	billCollector.OnHTML("li span a", func(e *colly.HTMLElement) {
		results.Uploader = e.Text
	})

	billCollector.OnHTML(".l308ffcf452c7ea53cc6f51251333f2e075003256.clearfix ul:nth-of-type(3) li:first-of-type span", func(e *colly.HTMLElement) {
		downloads, err := strconv.Atoi(e.Text)
		if err != nil {
			return
		}

		results.Downloads = downloads
	})

	billCollector.OnHTML(".l308ffcf452c7ea53cc6f51251333f2e075003256.clearfix ul:nth-of-type(3) li:nth-of-type(3) span", func(e *colly.HTMLElement) {
		results.Date = e.Text
	})

	billCollector.OnHTML(".l2629102922252783c610a2483b8b889ceff42c45 l8117a79035b80d6ca4aff12a3eb37266989069cb ld24b061224b0748215ecd65621ae0e7e3254b325", func(e *colly.HTMLElement) {
		results.Magnet = e.Attr("href")
	})

	billCollector.Visit(fmt.Sprintf("https://1337x.to%s", torrent_page))
	billCollector.Wait()
	return results
}
