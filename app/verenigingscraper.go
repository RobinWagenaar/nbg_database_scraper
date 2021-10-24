package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var jsonEndpointFormat = "https://nbg-hondensport.nl/regios/clubs?regionId=%d"
var clubNamePattern = regexp.MustCompile(`(\d.\d\d\d?|\w{2,8})\s+(.*)`)


func ScrapeVerenigingInfo(regio int) []*Vereniging{
	verenigingen := make([]*Vereniging, 0, 0)


	regionData := fetchRegionJson(regio)
	responseObject := new(JsonRegionResponse)
	json.Unmarshal(regionData, &responseObject)
	for _, club := range responseObject.Clubs {
		verenigingen = append(verenigingen, mapToVereniging(club))
	}

	log.Println(len(verenigingen), "actieve verenigingen gevonden in regio", regio)
	return verenigingen
}

func mapToVereniging(club JsonClub) *Vereniging {
	vereniging := new(Vereniging)
	parts := clubNamePattern.FindStringSubmatch(club.Name)
	vereniging.Id = parts[1]
	vereniging.Naam = parts[2]
	vereniging.NbgPagina = club.Url
	vereniging.Website = club.Website
	vereniging.EmailAdres = club.Email
	vereniging.Stad = club.City
	vereniging.TelefoonNummer = club.PhoneNumber
	return vereniging
}

func fetchRegionJson(regionId int) []byte {
	url := fmt.Sprintf(jsonEndpointFormat, regionId)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseData
}

type JsonClub struct {
	Name        string `json:"name"`
	City        string `json:"city"`
	PhoneNumber string `json:"phoneNumber"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Url         string `json:"url"`
}

type JsonRegion struct {
	Id     string `json:"id"`
	Number string `json:"number"`
	Name   string `json:"name"`
}

type JsonRegionResponse struct {
	JsonRegion `json:"region"`
	Clubs []JsonClub `json:"clubs"`
}