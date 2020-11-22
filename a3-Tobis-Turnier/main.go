package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

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
	anzahl, err := strconv.Atoi(lines[0])
	if err != nil {
		panic(err)
	}
	// Spielstärken, bester Spieler
	spielstaerken := make([]int, anzahl)
	besterSpieler := 0
	for index := range spielstaerken {
		// Atoi = Text zu Zahl
		spielstaerke, err := strconv.Atoi(lines[index+1])
		if err != nil {
			panic(err)
		}
		spielstaerken[index] = spielstaerke
		if spielstaerke > spielstaerken[besterSpieler] {
			besterSpieler = index
		}
	}
	for index, spielstaerke := range spielstaerken {
		if index != besterSpieler && spielstaerke == spielstaerken[besterSpieler] {
			// "Sanity-check": Es darf nur einen besten Spieler geben
			panic("Mehrere beste Spieler!")
		}
	}
	// Random "seeden" - ansonsten ist Go-Random determiniert
	rand.Seed(time.Now().UnixNano())

	// Gibt zurück, ob Spieler 1 gewonnen hat
	spieler1Gewinnt := func(spieler1, spieler2 int) bool {
		if rand.Intn(spielstaerken[spieler1]+spielstaerken[spieler2]) < spielstaerken[spieler1] {
			return true
		}
		return false
	}

	// Gewinner eines Spiels, Gibt Spielernummer zurück
	gewinner := func(spieler1, spieler2 int) int {
		if spieler1Gewinnt(spieler1, spieler2) {
			return spieler1
		}
		return spieler2
	}

	// Eine Runde Liga: Gibt 1 zurück, wenn der beste Spieler gewonnen hat, sonst 0
	liga := func() int {
		// Slice der Siege
		siege := make([]int, anzahl)
		// Jeder gegen jeden
		for spieler1 := range spielstaerken {
			for spieler2 := spieler1 + 1; spieler2 < len(spielstaerken); spieler2++ {
				// Gewinner erhält den Sieg
				siege[gewinner(spieler1, spieler2)]++
			}
		}
		// Sieger ermitteln: Spieler von kleiner zu großer Spielernummer durchgehen
		meisteSiege := 0
		for spieler, anzahlSiege := range siege {
			// Nur bei mehr Siegen neuer Sieger: Bei gleich vielen bleibt es der mit der kleineren Spielernummer
			if siege[meisteSiege] < anzahlSiege {
				meisteSiege = spieler
			}
		}
		if meisteSiege == besterSpieler {
			// Bester Spieler hat gesiegt
			return 1
		}
		return 0
	}

	// Gibt eine Funktion zurück, die eine Runde K.O. simuliert
	koVariante := func(gewinner func(int, int) int) func() int {
		return func() int {
			// Turnierplan erstellen
			turnierplan := make([]int, anzahl)
			for index := range turnierplan {
				turnierplan[index] = index
			}
			// Mischen (verwendet Fisher-Yates)
			rand.Shuffle(len(turnierplan), func(i, j int) {
				turnierplan[i], turnierplan[j] = turnierplan[j], turnierplan[i]
			})
			// Rekursiv Sieger eines "Bereiches" des Turnierplans ermitteln
			var sieger func(int, int) int
			sieger = func(start, ende int) int {
				diff := ende - start
				if diff == 1 {
					// Linke & rechte Hälfte umfassen nur einen Spieler: Gegeneinander antreten lassen!
					return gewinner(turnierplan[start], turnierplan[ende])
				}
				mitte := start + diff/2
				// Es spielt der Gewinner der linken Hälfte gegen den der rechten Hälfte
				return gewinner(sieger(start, mitte), sieger(mitte, ende))
			}
			if sieger(0, anzahl-1) == besterSpieler {
				// 1 zurückgeben, wenn der beste Spieler gewonnen hat
				return 1
			}
			// Sonst 0
			return 0
		}
	}

	// Einfache K.O.-Variante: Ein Spiel entscheidet
	ko := koVariante(gewinner)

	// K.O. x5: "Best of 5"
	ko5 := koVariante(func(spieler1, spieler2 int) int {
		// Siege des ersten Spielers
		siegeSpieler1 := 0
		for i := 0; i < 5; i++ {
			if spieler1Gewinnt(spieler1, spieler2) {
				siegeSpieler1++
				if siegeSpieler1 == 3 {
					// 3. Sieg, Spieler 1 hat gewonnen!
					return spieler1
				}
			}
		}
		return spieler2
	})

	// Tester für Turniervariante: Lässt viele Simulationen laufen
	anzahlLaeufe := int(1e6)
	testeTurniervariante := func(name string, runde func() int) {
		siegeBesterSpieler := 0
		for laeufe := 0; laeufe < anzahlLaeufe; laeufe++ {
			// Spiele eine Runde!
			siegeBesterSpieler += runde()
		}
		// Ausgeben
		fmt.Println(name+":", (float32(siegeBesterSpieler)/float32(anzahlLaeufe))*100.0, "%")
	}

	// Varianten testen
	testeTurniervariante("Liga", liga)
	testeTurniervariante("K.O.", ko)
	testeTurniervariante("K.O.x5", ko5)

	fmt.Println("Zeit verstrichen:", float32(time.Now().UnixNano()-nanos)/1e9, "s")
}
