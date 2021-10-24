package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

/**
<tr>
	<td class="nr">1</td>
	<td>
		<a href="lid/94865-v-van-gilse" class="tooltip" data-tooltip-title="Meer over deze geleider">
			V. van Gilse
		</a>
	</td>
	<td>
		2.18
	</td>
	<td>
		<a href="hond/2246-homer-simpson-vom-hexenwald" class="tooltip" data-tooltip-title="Meer over deze hond">Homer Simpson vom Hexenwald</a>
	</td>
	<td>
		OV
	</td>
	<td>
		IGPIII
	</td>
	<td class="score">78</td>
	<td class="score">82</td>
	<td class="score">80</td>
	<td class="score">240</td>
	<td class="score">Goed</td>
</tr>
*/

var examenDetailEndpointUrl = "https://nbg-hondensport.nl/examens/uitslagen/getResults"
var wedstrijdDetailEndpointUrl = "https://nbg-hondensport.nl/wedstrijden/uitslagen/getResults"

func ScrapeExamenUitslagen(examenId string) []*Uitslag{
	return scrapeUitslagen(examenDetailEndpointUrl, examenId)
}

func ScrapeWedstrijdUitslagen(wedstrijdId string) []*Uitslag{
	return scrapeUitslagen(wedstrijdDetailEndpointUrl, wedstrijdId)
}

func scrapeUitslagen(url, id string) []*Uitslag {
	body := []byte(fmt.Sprintf("eventId=%s", id))

	client := &http.Client{}
	req, _ := http.NewRequest("POST", examenDetailEndpointUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseObject := new(JsonExamenDetailResponse)
	json.Unmarshal(responseData, &responseObject)
	html := fmt.Sprintf("<html><body>%s</body></html>", responseObject.Html)
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(html)))


	uitslagen := make([]*Uitslag, 0, 0)
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		// uitslag
		uitslag := new(Uitslag)
		uitslagen = append(uitslagen, uitslag)

		// geleider
		geleider := new(Geleider)
		uitslag.Geleider = geleider
		geleiderAnchor := s.Find("td:nth-child(2) a")

		if geleiderAnchor.Length() != 0 {
			geleiderNaam := strings.TrimSpace(geleiderAnchor.Text())
			if len(geleiderNaam) > 0 {
				geleider.Naam = geleiderNaam
			}

			gid, gidExists := geleiderAnchor.Attr("href")
			if gidExists {
				geleider.Id = strings.TrimSpace(gid)[4:]
			}
		}else{
			geleider.Naam = strings.TrimSpace(s.Find("td:nth-child(2)").Text())
		}

		// geleider vereniging lidmaatschap
		geleiderVerenigingId := strings.TrimSpace(s.Find("td:nth-child(3)").Text())
		if len(geleiderVerenigingId) > 0 {
			geleider.VerenigingId = geleiderVerenigingId
		} else {
			geleider.VerenigingId = "0.00" //opgeheven
		}

		// hond
		hond := new(Hond)
		uitslag.Hond = hond
		hondAnchor := s.Find("td:nth-child(4) a")

		if hondAnchor.Length() != 0 {
			hondNaam := strings.TrimSpace(hondAnchor.Text())
			if len(hondNaam) > 0 {
				hond.Naam = hondNaam
			}

			hid, hidExists := hondAnchor.Attr("href")
			if hidExists {
				hond.Id = strings.TrimSpace(hid)[5:]
			}
		} else {
			hond.Naam = strings.TrimSpace(s.Find("td:nth-child(4)").Text())
		}

		hondRas := strings.TrimSpace(s.Find("td:nth-child(5)").Text())
		if len(hondRas) > 0 {
			hond.Ras = hondRas
		}

		// examensoort
		certificaat := strings.TrimSpace(s.Find("td:nth-child(6)").Text())
		if len(certificaat) > 0 {
			uitslag.Certificaat = certificaat
		}

		// scores
		uitslag.ScoreA = s.Find("td:nth-child(7)").Text()
		uitslag.ScoreB = s.Find("td:nth-child(8)").Text()
		uitslag.ScoreC = s.Find("td:nth-child(9)").Text()
		uitslag.Totaal = s.Find("td:nth-child(10)").Text()
		uitslag.Kwalificatie = s.Find("td:nth-child(11)").Text()
	})

	log.Println(len(uitslagen), "uitslagen gevonden voor event", id)
	return uitslagen
}


type JsonExamenDetailResponse struct {
	Status string `json:"status"`
	Html string `json:"html"`
}