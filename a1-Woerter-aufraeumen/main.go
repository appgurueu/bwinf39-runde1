package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unicode"
)

// Text - Unicode-Zeichen-Array, um keine UTF-8-Probleme zu bekommen
type Text []rune

// Wort - Text und Anzahl Vorkommen
type Wort struct {
	text   Text
	anzahl uint
}

// WortMitLuecken - Anfang und Ende im Lückentext, gegebene Buchstaben, und infragekommende Wörter
type WortMitLuecken struct {
	start              uint
	length             uint
	gegebeneBuchstaben map[uint]rune
	passendeWoerter    map[uint]bool
}

func main() {
	nanos := time.Now().UnixNano()
	// Einlesen
	if len(os.Args) != 2 {
		println("Verwendung: <pfad>")
		return
	}
	text, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(text), "\n")
	aufgabe := Text(lines[0])
	_woerter := strings.Split(lines[1], " ")
	// Verarbeiten
	woerterNachLaenge := map[uint]map[uint]*Wort{}
	bekannteWoerter := map[string]uint{}
	for id, _wort := range _woerter {
		wort := Text(_wort)
		length := uint(len(wort))
		bekannt := bekannteWoerter[_wort]
		if bekannt != 0 {
			// Falls bekannt, einfach Anzahl erhöhen
			woerterNachLaenge[length][bekannt-1].anzahl++
		} else {
			// Ansonsten Wort einfügen in Index nach Länge
			if woerterNachLaenge[length] == nil {
				woerterNachLaenge[length] = map[uint]*Wort{}
			}
			woerterNachLaenge[length][uint(id)] = &Wort{wort, 1}
			// Wort als bekannt setzen. ID mit einem Offset von 1, da bei einem miss 0 zurückgegeben wird, was von der ID 0 unterschieden werden muss
			bekannteWoerter[_wort] = uint(id) + 1
		}
	}
	// Lücken, bei denen ein Wort passen würde, nach Wort-ID. Da die Wort-IDs von 0 - n sind, kann eine slice (array) verwendet werden
	lueckenNachWort := make([]map[*WortMitLuecken]bool, len(_woerter))
	for index := range lueckenNachWort {
		lueckenNachWort[index] = map[*WortMitLuecken]bool{}
	}
	// Aktuelles Wort mit Lücken
	var wortMitLuecken *WortMitLuecken
	// Setzt ein Wort an einer Stelle an
	var setzeWortEin func(wortMitLuecken *WortMitLuecken, id uint)
	setzeWortEin = func(wortMitLuecken *WortMitLuecken, id uint) {
		// Ersetzt mit dem richtigen Wort die entsprechende Stelle im Text
		copy(aufgabe[wortMitLuecken.start:wortMitLuecken.start+wortMitLuecken.length], Text(_woerter[id]))
		// Anzahl verringern
		woerterNachLaenge[wortMitLuecken.length][id].anzahl--
		if woerterNachLaenge[wortMitLuecken.length][id].anzahl == 0 {
			// Anzahl 0, Wort steht nicht mehr zur Verfügung
			delete(woerterNachLaenge[wortMitLuecken.length], id)
			delete(lueckenNachWort[id], wortMitLuecken)
			// Überall, wo das Wort infragekommt...
			for lueckenwort := range lueckenNachWort[id] {
				// ... Wort streichen
				delete(lueckenwort.passendeWoerter, id)
				if len(lueckenwort.passendeWoerter) == 1 {
					// Nur noch ein infragekommendes Wort
					for id := range lueckenwort.passendeWoerter {
						// Einsetzen!
						setzeWortEin(lueckenwort, id)
					}
				}
			}
			// Fertig mit dem Wort: Keine Lücken dürfen es noch als Kandidaten haben!
			lueckenNachWort[id] = nil
		}
	}
	findePassendeWoerter := func() {
		var id uint
		// Sucht ein passendes Wort
	wortSuche:
		for _id, wort := range woerterNachLaenge[wortMitLuecken.length] {
			for index, gegebenerBuchstabe := range wortMitLuecken.gegebeneBuchstaben {
				if wort.text[index] != gegebenerBuchstabe {
					// Buchstabe an einer Stelle passt nicht: Wort kommt nicht infrage
					continue wortSuche
				}
			}
			id = _id
			// Passendes Wort merken
			wortMitLuecken.passendeWoerter[id] = true
			lueckenNachWort[id][wortMitLuecken] = true
		}

		if len(wortMitLuecken.passendeWoerter) == 1 {
			// Nur ein passendes Wort: Einsetzen!
			setzeWortEin(wortMitLuecken, id)
		}
	}
	for pos, char := range aufgabe {
		istBuchstabe := unicode.IsLetter(char)
		istGesucht := char == '_'
		if istBuchstabe || istGesucht {
			// Teil eines Wortes mit Lücken
			if wortMitLuecken == nil {
				// Initialisierung falls erster Buchstabe / erste Lücke des Wortes
				wortMitLuecken = &WortMitLuecken{uint(pos), 0, map[uint]rune{}, map[uint]bool{}}
			}
			if istBuchstabe {
				// Gegebenen Buchstaben eintragen
				wortMitLuecken.gegebeneBuchstaben[uint(pos)-wortMitLuecken.start] = char
			}
		} else if wortMitLuecken != nil {
			wortMitLuecken.length = uint(pos) - wortMitLuecken.start
			findePassendeWoerter()
			wortMitLuecken = nil
		}
	}
	if wortMitLuecken != nil {
		wortMitLuecken.length = uint(len(aufgabe)) - wortMitLuecken.start
		findePassendeWoerter()
	}
	// lueckenNachWort nach Wörtern durchgehen, für die nur eine einzige Lücke infrage kommt, ist nicht nötig
	// Ausgabe der Lösung
	fmt.Println("Lösung:", string(aufgabe))
	fmt.Println("Zeit verstrichen:", float64(time.Now().UnixNano()-nanos)/1e6, "ms")
}
