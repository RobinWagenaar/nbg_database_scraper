package app

type Vereniging struct {
	Id             string
	Naam           string
	Stad           string
	TelefoonNummer string
	Website        string
	EmailAdres     string
	NbgPagina      string
	Certificaten   []string
}

type Geleider struct {
	Id string
	Naam string
	VerenigingId string
}

type Hond struct {
	Id string
	Naam string
	Ras string
}

type Gebeurtenis struct {
	Id string
	Soort string
	VerenigingId string
	Datum string
	Locatie string
	Keurmeesters []string
	Certs string
	Uitslagen []*Uitslag
}

type Uitslag struct {
	Geleider *Geleider
	Hond *Hond
	Certificaat string
	ScoreA, ScoreB, ScoreC string
	Totaal string
	Kwalificatie string
}