package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// Figur - Ganzzahl von -128 bis 127
type Figur int8

// Teil - Drei Figuren
type Teil [3]Figur

// Kernteil - Seite, Teil, und Vorheriges
type Kernteil struct {
	seite      uint8
	teil       uint8
	vorheriges *Kernteil
}

// Eckteil - Seite und Teil
type Eckteil struct {
	seite uint8
	teil  uint8
}

// VerwendeteTeile - Flag
type VerwendeteTeile uint16

// Verwendet - Gibt zurück, ob ein Teil verwendet ist
func (teile VerwendeteTeile) Verwendet(teil uint8) bool {
	return teile&(1<<teil) > 0
}

// Verwende - Verwendet ein Teil und gibt Flag zurück
func (teile VerwendeteTeile) Verwende(teil uint8) VerwendeteTeile {
	return teile | (1 << teil)
}

func main() {
	// Zeitmessung
	nanos := time.Now().UnixNano()
	verstricheneZeit := func() {
		fmt.Println("Zeit verstrichen:", float64(time.Now().UnixNano()-nanos)/1e6, "ms")
	}
	// Eingabe
	if len(os.Args) != 2 {
		println("Verwendung: <pfad>")
		return
	}
	text, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(text), "\n")
	// Teile einlesen
	teile := make([]Teil, 9)
	for l := 2; l < 11; l++ {
		var teil Teil
		figuren := strings.Split(lines[l], " ")
		for f := 0; f < 3; f++ {
			teilF, _ := strconv.Atoi(figuren[f])
			teil[f] = Figur(teilF)
		}
		teile[l-2] = teil
	}
	// Rekursive Funktion, die "Kerne" erzeugt und probiert
	var probiereKerne func(Figur, *Kernteil, VerwendeteTeile, uint8)
	probiereKerne = func(passendZurStartfigur Figur, vorheriges *Kernteil, verwendeteTeile VerwendeteTeile, anzahlKernteile uint8) {
		if anzahlKernteile == 6 {
			// Abbruchbedingung: Sechs Kernteile
			// Probieren: Erster Test: Passt die Seite des Kernteils zur ersten Seite (ist der Kern geschlossen?)
			if teile[vorheriges.teil][vorheriges.seite] != passendZurStartfigur {
				return
			}
			// Randfiguren des Kerns
			randfiguren := [6]Figur{}
			cursor := vorheriges
			for i := range randfiguren {
				// Umgekehrte Reihenfolge
				randfiguren[5-i] = teile[cursor.teil][(cursor.seite+2)%3]
				if i < 5 {
					cursor = cursor.vorheriges
				}
			}
			// Rekursive Funktion, die alle möglichen Eckanordnungen probiert
			var probiereEcken func(uint8, VerwendeteTeile, []Eckteil)
			probiereEcken = func(versatz uint8, verwendeteTeile VerwendeteTeile, ecken []Eckteil) {
				if len(ecken) == 3 {
					if verwendeteTeile != 0b111111111 {
						// "Sanity-check": Am Ende müssen alle Teile verwendet worden sein
						panic("Teile mehrfach verwendet")
					}
					// 3 Passende Ecken wurden gefunden, das Puzzle ist gelöst!
					// Ausgabe:
					// Obere Figurenteile sind Großbuchstaben, untere Kleinbuchstaben
					figurenOben, figurenUnten := [27]rune{}, [27]rune{}
					for f := 0; f < 27-8; f++ {
						figurenOben[f], figurenUnten[f] = rune('A'+f), rune('a'+f)
					}
					// Gibt für eine Figur den Buchstaben zurück
					figur := func(num Figur) rune {
						if num < 0 {
							return figurenUnten[-num]
						}
						return figurenOben[num]
					}
					// Teile ausgeben, platzsparend nebeneinander
					fmt.Println("Teile:")
					teileFmt := []string{"", "", "", ""}
					for _, teil := range teile {
						for i, zeile := range []string{
							`   /-\    `,
							fmt.Sprintf(`  /%c %c\   `, figur(teil[0]), figur(teil[1])),
							fmt.Sprintf(` /  %c  \  `, figur(teil[2])),
							`/-------\ `,
						} {
							teileFmt[i] += zeile
						}
					}
					fmt.Println(strings.Join(teileFmt, "\n"))
					fmt.Println()
					// Lösung ausgeben
					fmt.Println("Lösung:")
					format := []rune(`           /-\
          /0 0\
         /  0  \
        /-------\
       / \  3  / \
      /3 4\4 4/4 3\
     /  4  \ /  4  \
    /---------------\
   / \  4  / \  4  / \
  /2 2\3 4/4 4\4 3/1 1\
 /  2  \ /  3  \ /  1  \
/-----------------------\
`)
					// Unterteilung in 5 Figurengruppen: Eckfiguren (0-2), Randfiguren des Kerns (3), Kernfiguren (4)
					figuren := [5][]Figur{}
					// Ecken einsetzen
					for i, ecke := range ecken {
						teil := teile[ecke.teil]
						figuren[i] = []Figur{teil[ecke.seite], teil[(ecke.seite+1)%3], teil[(ecke.seite+2)%3]}
					}
					if versatz == 0 {
						// Kein Versatz: Randfiguren und Kernteile sind in richtiger Reihenfolge
						figuren[3] = randfiguren[:]
						cursor = vorheriges
					} else {
						// Versatz von 1: Randfiguren und Kernteile müssen um eins verschoben werden
						// Dabei wird das erste Element zum neuen Letzten
						figuren[3] = append(randfiguren[1:], randfiguren[0])
						// "cursor" ist ein einfach verlinkter Listenknoten, der auf das Erste Element zeigt
						// Dessen Nachfolger wird jetzt der Zweite Knoten
						cursor.vorheriges = vorheriges
					}
					// Kernfiguren einsetzen
					figuren[4] = make([]Figur, 12)
					for i := 5; i >= 0; i-- {
						teil := teile[cursor.teil]
						// Pro Teil immer Zwei Figuren
						figuren[4][i*2] = teil[(cursor.seite+1)%3]
						figuren[4][i*2+1] = teil[cursor.seite]
						cursor = cursor.vorheriges
					}
					// Im String sind die Figurengruppen nicht in der richtigen Reihenfolge:
					// - Die Ecken müssen noch gedreht werden
					// - Die Randfiguren sind kreisförmig angeordnet, und nicht von oben nach unten - links nach rechts
					// - Die Kernfiguren ebenfalls
					reihenfolge := [5][]int{
						// Drehung der Ecken
						{1, 2, 0},
						{0, 1, 2},
						{2, 0, 1},
						// Kreisform Randfiguren
						{0, 5, 1, 4, 2, 3},
						// Kreisform Kernfiguren
						{11, 0, 1, 2, 10, 3, 9, 4, 8, 7, 6, 5},
					}
					// N-tes Element jeder Figurengruppe
					n := [5]int{}
					for i, c := range format {
						// Platzhalter für Figurengruppen sind die jeweiligen Zahlen
						if c >= '0' && c <= '4' {
							// Zahl 48 - 52 (ASCII) in Zahl 0 - 4 konvertieren, "c" gibt Figurengruppe an
							c -= '0'
							// N-te Figur aus der passenden Figurengruppe in der richtigen Reihenfolge
							format[i] = figur(figuren[c][reihenfolge[c][n[c]]])
							// Nächstes Element der Figurengruppe beim nächsten Platzhalter
							n[c]++
						}
					}
					// Ausgeben
					fmt.Println(string(format))
					verstricheneZeit()
					// Programm beenden
					os.Exit(0)
				}
				// Passende Figur ist anderer Teil der Randfigur an entsprechender Stelle
				passendeFigur := -randfiguren[uint8(2*len(ecken))+versatz]
				for teil := uint8(0); teil < 9; teil++ {
					if verwendeteTeile.Verwendet(teil) {
						// Teil schon verwendet
						continue
					}
					for seite, figur := range teile[teil] {
						// Seiten probieren
						if figur == passendeFigur {
							// Erstelle Kopie der Ecken & füge neue Ecke hinzu
							eckenKopie := make([]Eckteil, len(ecken)+1)
							for i, ecke := range ecken {
								eckenKopie[i] = ecke
							}
							eckenKopie[len(ecken)] = Eckteil{uint8(seite), teil}
							// Probiere weitere Ecken
							probiereEcken(versatz, verwendeteTeile.Verwende(teil), eckenKopie)
							break
						}
					}
				}
			}
			// Ecken probieren, mit Versatz 0 und 1
			probiereEcken(0, verwendeteTeile, []Eckteil{})
			probiereEcken(1, verwendeteTeile, []Eckteil{})
			return
		}
		for teil := uint8(0); teil < 9; teil++ {
			if verwendeteTeile.Verwendet(teil) {
				// Teil schon verwendet
				continue
			}
			// Passende Seiten ermitteln
			passendeFigur := -teile[vorheriges.teil][vorheriges.seite]
			passendeSeiten := []uint8{}
			for seite, figur := range teile[teil] {
				if figur == passendeFigur {
					passendeSeiten = append(passendeSeiten, uint8(seite))
				}
			}
			// Alle Seiten passen: Seiten sind identisch, Drehung ist egal
			if len(passendeSeiten) == 3 {
				passendeSeiten = []uint8{0}
			}
			for _, seite := range passendeSeiten {
				// Kernteil erzeugen. Seite ist hierbei die, an die das nächste Kernteil anbinden muss.
				kernteil := &Kernteil{(seite + 2) % 3, teil, vorheriges}
				// Neues Kernteil übergeben, Teil verwenden, Anzahl Kernteile erhöhen
				probiereKerne(passendZurStartfigur, kernteil, verwendeteTeile.Verwende(teil), anzahlKernteile+1)
			}
		}
	}
	for teil := uint8(0); teil < 5; teil++ {
		for seite := uint8(0); seite < 3; seite++ {
			// Unter 4 Teilen muss eines dabei sein, dass keine Ecke, sondern ein Kernteil ist
			// Starte Brute-Force mit Kernteil
			probiereKerne(-teile[teil][seite], &Kernteil{(seite + 2) % 3, teil, nil}, VerwendeteTeile(0).Verwende(teil), 1)
			// TODO break
		}
	}
	// Kein Abbruch ist erfolgt: Das Puzzle konnte nicht gelöst werden
	fmt.Println("Puzzle unlösbar")
	verstricheneZeit()
}
