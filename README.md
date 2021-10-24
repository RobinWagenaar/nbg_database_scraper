# Scraper voor de website van de Nederlandse Bond voor Gebruikshonden (NBG)

De data op de website van de NBG (https://nbg-hondensport.nl/nbg-hondensport) mag 
dan misschien publiek zijn; erg handig doorzoekbaar is het niet. Als je wil zoeken 
naar alle honden waarmee een geleider ooit een examen/wedstrijd heeft gelopen, of 
plaintext zoeken op naam van de hond, dan kom je bedrogen uit. 

Mooie reden voor mij om een beetje te prutsen met Golang en MongoDB. Dit scriptje 
downloaded binnen enkele minuten alle wedstrijden/examenresultaten. Deze documenten 
worden vervolgens opgeslagen in een lokale MongoDB, zodat deze makkelijk offline 
doorzoekbaar zijn. 

Let wel, dit zijn mijn eerste stapjes in de wondere wereld van Golang. Dus garantie
tot aan de deur :-)