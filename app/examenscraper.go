package app

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strings"
)

/**
<div class="event" data-id="2787">
	<div class="head">
		<a href="#">
			<span class="arrow">
				<i class="fa fa-caret-right hide-on-open"></i>
				<i class="fa fa-caret-down show-on-open"></i>
			</span>
			<time datetime="2021-06-05">05-06</time>
			<span class="name">2.13 HSC De Volharding</span>
			<span class="type hide-on-mobile">VZH/IGP/SPH</span>
		</a>
	</div>
	<div class="results hidden">
		<div class="meta">
			<div class="column">
				<strong><i class="fa fa-map-marker"></i> Amersfoort, 2.13</strong>
			</div>
			<div class="column">
				<strong><i class="fa fa-trophy"></i> VZH/IGP/SPH</strong>
			</div>
			<div class="column">
				<strong><i class="fa fa-flag"></i> KM(s) : C Janssen</strong>
			</div>
		</div>
		<div class="holder">
			<!-- Wordt door JS gevuld -->
		</div>
	</div>
</div>
*/

var examenHoofdpaginaUrlFormat = "https://nbg-hondensport.nl/examens/uitslagen/%d/%d"
var wedstrijdHoofdpaginaUrlFormat = "https://nbg-hondensport.nl/wedstrijden/uitslagen/%d/%d"

var locatiePattern = regexp.MustCompile(`(.*),\s+.*`)


func ScrapeExamenInfo(jaar, maand int) []*Gebeurtenis {
	url := fmt.Sprintf(examenHoofdpaginaUrlFormat, jaar, maand)
	gebeurtenissen := scrapeGebeurtenisInfo(url)
	for _, gebeurtenis := range gebeurtenissen {
		gebeurtenis.Soort = "examen"
	}
	log.Println(len(gebeurtenissen), "examens gevonden in", maand, jaar)
	return gebeurtenissen
}

func ScrapeWedstrijdInfo(jaar, maand int) []*Gebeurtenis {
	url := fmt.Sprintf(wedstrijdHoofdpaginaUrlFormat, jaar, maand)
	gebeurtenissen := scrapeGebeurtenisInfo(url)
	for _, gebeurtenis := range gebeurtenissen {
		gebeurtenis.Soort = "wedstrijd"
	}
	log.Println(len(gebeurtenissen), "wedstrijden gevonden in", maand, jaar)
	return gebeurtenissen
}

func scrapeGebeurtenisInfo(url string) []*Gebeurtenis {
	// Request the HTML page.

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Parse the document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Scrape the data
	gebeurtenissen := make([]*Gebeurtenis, 0, 0)
	doc.Find("div.event").Each(func(i int, s *goquery.Selection) {
		gebeurtenis := new(Gebeurtenis)

		gId, idExists := s.Attr("data-id")
		if idExists {
			gebeurtenis.Id = strings.TrimSpace(gId)
		}


		d, dateExists := s.Find("time").Attr("datetime")
		if dateExists {
			gebeurtenis.Datum = strings.TrimSpace(d)
		}

		clubName := s.Find("span.name").Text()
		if strings.Contains(strings.ToLower(clubName), "opgeheven") {
			gebeurtenis.VerenigingId = "0.00"
		} else {
			idSubmatch := clubNamePattern.FindStringSubmatch(strings.TrimSpace(clubName))
			if len(idSubmatch) == 3 {
				gebeurtenis.VerenigingId = idSubmatch[1]
			}
		}

		column1 := strings.TrimSpace(s.Find("div.meta div.column:nth-child(1)").Text())
		locSubmatch := locatiePattern.FindStringSubmatch(column1)
		if len(locSubmatch) == 2 {
			gebeurtenis.Locatie = locSubmatch[1]
		}

		column2 := strings.TrimSpace(s.Find("div.meta div.column:nth-child(2)").Text())
		if len(column2) > 0 {
			gebeurtenis.Certs = column2
		}

		column3 := strings.TrimSpace(s.Find("div.meta div.column:nth-child(3)").Text())
		if len(column3) > 8 { // prefix should be: 'KM(s) : '
			gebeurtenis.Keurmeesters = make([]string, 0, 0)
			for _, keurmeester := range strings.Split(column3[8:], ", ") {
				gebeurtenis.Keurmeesters = append(gebeurtenis.Keurmeesters, keurmeester)
			}
		}

		gebeurtenissen = append(gebeurtenissen, gebeurtenis)
	})

	return gebeurtenissen
}
