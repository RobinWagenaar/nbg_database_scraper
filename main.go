package main

import (
	"fmt"
	"log"
	"nbgscraper/app"
	"time"
)

func main() {


	repo := new(app.MongoRepository)

	// initiële lijst van verenigingen
	for i := 1; i < 10; i++ {
	  log.Println("Bezig met scrapen van verenigingen in regio", i)
		verenigingen := app.ScrapeVerenigingInfo(i)
		for _, vereniging := range verenigingen {
			repo.InsertOrReplaceVereniging(vereniging)
		}
	}

	// lijst van examens en resultaten
	for jaar := time.Now().Year(); jaar >= 2000; jaar-- {
		log.Println("Bezig met scrapen van examens in jaar", jaar)
		for maand := 12; maand > 0; maand-- {
			examens := app.ScrapeExamenInfo(jaar, maand)
			for _, examen := range examens {
				examen.Uitslagen = app.ScrapeExamenUitslagen(examen.Id)
				repo.InsertOrReplaceGebeurtenis(examen)
			}
		}
	}

	// lijst van wedstrijden en resultaten
	for jaar := time.Now().Year(); jaar >= 2000; jaar-- {
		log.Println("Bezig met scrapen van wedstrijden in jaar", jaar)
		for maand := 12; maand > 0; maand-- {
			examens := app.ScrapeWedstrijdInfo(jaar, maand)
			for _, examen := range examens {
				examen.Uitslagen = app.ScrapeWedstrijdUitslagen(examen.Id)
				repo.InsertOrReplaceGebeurtenis(examen)
			}
		}
	}

	data, _ := json.MarshalIndent(gebeurtenissen,  "", "  ")
	fmt.Println(string(data))
	fmt.Println("Afgerond.")
}

