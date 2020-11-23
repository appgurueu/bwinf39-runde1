# **Lösungen 38. Bundeswettbewerb Informatik**

# Lars Müller

## Team: appguru, ID: 945

Die folgenden vier der fünf Aufgaben wurden von mir bearbeitet und finden sich in den jeweiligen Ordnern:

* `a1-Woerter-aufrauemen`
* `a2-Dreieckspuzzle`
* `a3-Tobis-Turnier`
* `a5-Wichteln`

Die Struktur jedes Ordners ist wie folgt:

* Dokumentation als `Dokumentation.md` ("Quellcode" in Markdown, eher ungeeignet) und `Dokumentation.pdf` (zum Ansehen)
* Go-Quellcode als `main.go`, ausführbare Datei für Linux als `main`
  * Kompilieren: `go build` im Aufgabenordner
  * Ausführen: `./main <pfad>` oder `go run main.go <pfad>` im Aufgabenordner
    * Etwa: `cd a1-Woerter-aufrauemen && go run main.go beispieldaten/raetsel0.txt`
* Weitere Beispieldaten in `beispieldaten` fortnummeriert
* Lösungen als Textdateien in `loesungen`, gleiche Namen wie die `beispieldaten`

Die Aufgabenstellungen wurden in den Ordnern beibehalten.

Die versteckten Ordner `.docs` und `.git` sind irrelevant.

### Umgebung

#### Betriebssystem

* Ubuntu 20.04 LTS (Linux, 64-Bit)

#### Programmiersprache

* Go 1.15

In Go werden Arrays, deren Länge nicht zur Compile-time feststeht, als "Slices" ("Stücke") bezeichnet.

Als Notation für eine Abbildung über Arrays, Slices oder Maps wird "\[Schlüssel] = Wert" verwendet.

#### Editor

* Visual Studio Code mit installiertem Go-Plugin geeignet
