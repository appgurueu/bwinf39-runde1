package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Verteilung - Wert als vergleichbare Zahl, von Schülern erhaltene Geschenke, vergebene Geschenke
type Verteilung struct {
	wert               uint64
	erhalteneGeschenke []uint
	vergebeneGeschenke []bool
}

func main() {
	// Zeitmessung
	nanos := time.Now().UnixNano()
	// Eingabe
	if len(os.Args) != 2 {
		println("Verwendung: <pfad>")
		return
	}
	// Einlesen
	text, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(text), "\n")
	_anzahl, err := strconv.Atoi(lines[0])
	if err != nil {
		panic(err)
	}
	// Anzahl Schüler
	anzahl := uint(_anzahl)
	// Basis für schnellen Vergleich
	basis := uint64(anzahl + 1)
	// Schüler: [Nummer - 1] = {Nummer 1. Wunsch - 1, Nummer 2. Wunsch - 1, Nummer 3. Wunsch - 1}
	schueler := make([][3]uint, anzahl)
	for s := uint(0); s < anzahl; s++ {
		var wunsch [3]uint
		// Regulärer Ausdruck: Alle Zahlen (= Wünsche) aus der Zeile extrahieren
		wunschStr := regexp.MustCompile("[0-9]+").FindAllString(lines[s+1], 3)
		for w := 0; w < 3; w++ {
			// Wünsche einlesen
			wunschW, err := strconv.Atoi(wunschStr[w])
			if err != nil {
				panic(err)
			}
			wunschW--
			wunsch[w] = uint(wunschW)
		}
		schueler[s] = wunsch
	}
	// Verfahren
	// Reihenfolge, in der Schüler wünsche erhalten
	reihenfolge := make([]uint, anzahl)
	// Probiert eine Reihenfolge, gibt "Wert" der resultierenden Verteilung als Zahl zurück
	probiereReihenfolge := func() (verteilung Verteilung) {
		// Wert
		verteilung.wert = 0
		// Vergebene Geschenke: [Geschenk] = Vergeben?
		verteilung.vergebeneGeschenke = make([]bool, anzahl)
		// Erhaltene Geschenke: [Schüler] = Geschenk
		verteilung.erhalteneGeschenke = make([]uint, anzahl)
		stelle := basis * basis
		for i := 0; i < 3; i++ {
			erfuellteWuensche := uint64(0)
			for _, s := range reihenfolge {
				// Schüler in Reihenfolge durchgehen
				wunsch := schueler[s][i]
				if verteilung.erhalteneGeschenke[s] != 0 || verteilung.vergebeneGeschenke[wunsch] {
					// Schüler schon "abgespeist" oder Geschenk schon vergeben
					continue
				}
				// Zähler erhöhen, Geschenk zuordnen & als vergeben markieren
				erfuellteWuensche++
				verteilung.erhalteneGeschenke[s] = wunsch + 1
				verteilung.vergebeneGeschenke[wunsch] = true
			}
			verteilung.wert += stelle * erfuellteWuensche
			stelle /= basis
		}
		return
	}
	// Startreihenfolge: 0 bis anzahl aufsteigend
	for i := range reihenfolge {
		reihenfolge[i] = uint(i)
	}
	// Startvergleichswerte für die Reihenfolge
	besteVerteilung := probiereReihenfolge()
	for {
		verbesserung := false
		for i := range reihenfolge {
			for j := range reihenfolge[i+1:] {
				// Probiere "swaps"
				reihenfolge[i], reihenfolge[j] = reihenfolge[j], reihenfolge[i]
				// Werte nach Swap
				andereVerteilung := probiereReihenfolge()
				if andereVerteilung.wert > besteVerteilung.wert {
					// Bessere Verteilung! Aktualisieren
					besteVerteilung = andereVerteilung
					verbesserung = true
				} else {
					// Ansonsten: Swap wieder rückgängig machen
					reihenfolge[i], reihenfolge[j] = reihenfolge[j], reihenfolge[i]
				}
			}
		}
		// ...solange die Reihenfolge verbessert werden kann
		if !verbesserung {
			break
		}
	}

	// Beste Verteilung komplettieren: "Leer ausgehenden" (kein Wunsch erfüllt) Schülern erstbeste Geschenke geben
	var i uint
	for s, g := range besteVerteilung.erhalteneGeschenke {
		if g == 0 {
			// Kein Wunsch erfüllt, kein Geschenk zugeordnet
			for {
				if !besteVerteilung.vergebeneGeschenke[i] {
					// Erstbestes übriges Geschenk erhalten
					besteVerteilung.erhalteneGeschenke[s] = i + 1
					i++
					break
				}
				i++
			}
		}
	}

	// Ausgabe
	fmt.Println("Lösung:")
	// Erste Zeile der Lösung beginnt mit dem Geschenk, das der erste Schüler erhält
	zeile := strconv.Itoa(int(besteVerteilung.erhalteneGeschenke[0]))
	for _, g := range besteVerteilung.erhalteneGeschenke[1:] {
		geschenk := strconv.Itoa(int(g))
		if len(zeile)+len(geschenk)+2 > 80 {
			// Passt nicht mehr in Zeile: Maximale Zeilenlänge von 80
			fmt.Println(zeile)
			// Neue Zeile beginnen
			zeile = ""
		} else {
			// In gleicher Zeile anhängen mit trennendem Komma
			zeile += ", "
		}
		// Zugeordnetes Geschenk zur Ausgabe hinzufügen
		zeile += geschenk
	}
	if zeile != "" {
		// Letzte Zeile ausgeben
		fmt.Println(zeile)
	}
	// Wert der Verteilung (in erfüllten Wünschen) & verstrichene Zeit ausgeben
	// Dafür Rückrechnung:
	// - wert/basis² = erfüllte 1. Wünsche
	// - (wert/basis) % basis = erfüllte 2. Wünsche
	// - wert % basis = erfüllte 3. Wünsche
	fmt.Printf(`
Erfüllte Wünsche: %d, %d, %d

Zeit verstrichen: %v s
`, besteVerteilung.wert/(basis*basis), (besteVerteilung.wert/basis)%basis, besteVerteilung.wert%basis, float32(time.Now().UnixNano()-nanos)/1e9)
}
