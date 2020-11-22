# Lösung Aufgabe 2 "Dreieckspuzzle"

## Lösungsidee

Wir nutzen eine Brute-Force, um mögliche Lösungen des Puzzles auszuprobieren.

Hierfür unterteilen wir ein gelöstes Puzzle zunächst in "Eckteile" und "Kernteile":

![Illustration](Illustration.png "Eckteile (türkis) und Kernteile (blau)")

Den Kern wiederum unterteilen wir in "Randfiguren" (äußere Figuren des Kerns) und "Kernfiguren" (innere Figuren des Kerns).

Zunächst probieren wir alle möglichen Kerne aus:

1. Wir beginnen mit einer der drei Seiten eines Teils
2. Wir suchen ein *anderes* Teil, das einen zum dort abgebildeten passenden Figurenteil besitzt
3. Wir bestimmen die gegenüberliegende Seite des passenden Teils, und suchen nun dafür wieder passende Teile (zurück zu Schritt 1)

So lange, bis die 6 Kernteile zusammen kommen. Dann muss noch die letzte Seite des letzten Teils zur ersten Seite des ersten Teils passen, und der Kern ist geschlossen!

Ein geschlossener Kern reicht allerdings noch nicht. Die übrigen Teile müssen noch als Eckteile an die Randfiguren des Kerns passen - ist dies gegeben, haben wir eine Lösung für das Puzzle!

### Korrektheit

Abgesehen von Implementierungsfehlern könnte die Brute-Force nur falsche Resultate liefern, wenn nicht ausreichend viele Möglichkeiten betrachtet würden. Dies ist hier aber nicht der Fall.

### Komplexität

Wir überlegen uns als obere Grenze: Für das erste Kernteil bestehen 9 Wahlmöglichkeiten, für das zweite nur noch 8, usw.

Es gibt also maximal $\frac{9!}{3!} = 60.480$ Möglichkeiten, die Kernteile und ihre Reihenfolge auszuwählen.

Jedes Kernteil kann nun maximal zwei passende Seiten besitzen, also gibt es für jedes Kernteil nochmal zwei Drehmöglichkeiten. 

Somit kommt man auf maximal $\frac{9!}{3!} \cdot 2^6 = 60.480 \cdot 64 = 3.870.720$ mögliche Kernanordnungen.

Für jede dieser Kernanordnungen müssen noch Ecken probiert werden. Hierbei gibt es zwei Möglichkeiten, den Kern zu drehen, und für jede $3! = 6$ infragekommende Eckanordnungen.

Tatsächlich brauchen wir als "Startteil" des Kerns nur 4 Teile auszuprobieren (unter 4 Teilen muss ein Kernteil sein), was diese Anzahl wiederum halbiert. Andererseits müssen wir für diese jeweils alle 3 Seiten probieren, was wiederum zu einer Verdreifachung führt.

Schließlich erhält man $\frac{9!}{3!} \cdot 2^6 \cdot 2 \cdot 3! \div 2 \cdot 3 = 9! \cdot 2^6 \cdot 3 = 69.672.960$ maximal auszuprobierende Lösungen.

Selbst wenn dieser theoretische worst-case - insofern er überhaupt möglich ist - einträte, würde das Programm schätzungsweise noch in sinnvoller Zeit terminieren (siehe verstrichene Zeiten).

## Umsetzung

Implementierung in der modernen und performanten Programmiersprache [Go](https://golang.org).

### Kompilieren

`go build` (erzeugt `a5-Wichteln`) oder `go build main.go` (erzeugt `main`)

### Verwendung

`go run main.go <pfad>` oder `./main <pfad>`

Beispiel: `./main beispieldaten/puzzle0.txt`

### Ausgabe

Figuren werden in den jeweiligen Buchstaben im Alphabet umgewandelt (1: A, 2: B usw.). Die unteren Figurenteile (negative Zahlen) werden dabei mit Kleinbuchstaben repräsentiert (-1: a, -2: b usw.).

```
Teile:
<Teile als ASCII-Art>

Lösung:
<Gelöstes Puzzle als ASCII-Art>

Zeit verstrichen: <Verstrichene Zeit in Millisekunden> ms
```

oder

```
Puzzle unlösbar
Zeit verstrichen: <Verstrichene Zeit in Millisekunden> ms
```

### Bibliotheken

* `fmt`: Ausgabe & Formattierung
* `io/ioutil`: Einlesen der Datei
* `os`: Programmargumente
* `strconv`: String/Integer-Konversion
* `strings`: Auftrennen von Text
* `time`: Zeitmessung

### Typen

#### `Figur`

Alias für `int8` für bessere Lesbarkeit.

#### `Teil`

Puzzleteil: 3-Figuren-Array.

#### `Kernteil`

Kernteil: Seite, and die das nächste Kernteil "anzudocken" hat, Teil-ID, und Zeiger auf vorangehendes Kernteil (Knoten einer Single-Linked-List)

#### `Eckteil`

Eckteil: Seite, mit der es an zur Randfigur passt, und Teil-ID.

#### `VerwendeteTeile`

Bitflag für verwendete Teile: Alias für `uint16`

##### `Verwendet(teil uint8) bool`

Prüft, ob das Teil mit der jeweiligen Nummer schon verwendet ist.

##### `Verwende(teil uint8) VerwendeteTeile`

Gibt neuen Flag zurück mit verwendetem Teil

### Ablauf

Wir implementieren die Brute-Force mit zwei rekursiven Funktionen: `probiereKerne` und `probiereEcken`.

Erstere nimmt als Argumente:

* Die zur "Startfigur" (Figur gewählte Seite erstes Teil) passende [`Figur`]
* Einen Zeiger auf das vorangehende Kernteil
* [`VerwendeteTeile`]
* Die Anzahl bisher aneindergehangener Kernteile

In jedem Schritt werden dann alle neun Teile durchgegangen. Für die noch nicht Verwendeten wird geprüft, ob sie an das letzte Kernteil passen; wenn ja, wird mit einer Variante des Kerns mit angehangenem Teil weiterprobiert. Hierbei wird die Anzahl aneindergehangener Teile natürlich um eins erhöht, das Teil als verwendet markiert, und ein neues Kernteil erzeugt, welches als Vorgänger den Zeiger auf das aktuelle Kernteil erhält.

Gestartet wird die Rekursion für jedes der ersten 4 Teile und für jede der drei Seiten. Sie terminiert, sobald ein kompletter Kern mit 6 Teilen (dafür die Anzahl) und passendem letzten Teil (dafür `passendZurStartfigur`) gefunden ist. Dann wird die zweite Rekursion gestartet.

Diese nimmt:

* Einen Versatz (Startwerte: 0 (erster Aufruf), 1 (zweiter Aufruf))
* [`VerwendeteTeile`] (Startwert: von Kernteil-Rekursion verwendete Teile)
* Eckteil-Slice (Startwert: leere Slice)

Die Randfiguren, zu denen passende Ecken gefunden werden müssen werden durch Go's "lexical scope" zugänglich.

Diese geht ebenfalls noch nicht verwendete Teile durch - die drei Übrigen - und versucht, diese als Eckteile einzusetzen. Passt eines, wird es an eine Kopie der Slice angehangen und es wird rekursiv weiterprobiert - bis schließlich alle drei Ecken bestimmt sind.
Dann erfolgt nach einem einfachen "sanity-check" die Ausgabe. Diese ist algorithmisch nicht besonders interessant; das Ausgabeformat für das gelöst Puzzle ist als String in ASCII-Art gegeben, mit Platzhaltern für die "Figurengruppen" erste Ecke, zweite Ecke, dritte Ecke, Randfiguren und Kernfiguren. Um die Ausgabe einer Figur auf ein Zeichen zu beschränken, werden Buchstaben anstatt von Zahlen verwendet (siehe [Ausgabe]).

Da Figurengruppen im String nicht unbedingt die richtige Reihenfolge haben, müssen hierfür Slices die Umordnung der Elemente angeben (siehe [Quellcode] Z. 170-184).

Schließlich werden Teile, Lösung und die verstrichene Zeit ausgegeben, woraufhin das Programm mit `os.Exit` beendet werden kann und muss, so dass die Brute-Force, nach erfolgreicher Lösungsfindung, stoppt.

Wird keine Lösung gefunden, kommt es nicht zur Programmbeendigung und der Code am Programmende wird ausgeführt. Hier wird zunächst ausgegeben, dass keine Lösung möglich ist, zusammen mit der verstrichenen Zeit.

## Quellcode

### `main.go`

```file:go
main.go
```

## Beispiele

In `loesungen` als Textdateien mit gleichem Namen wie die Aufgabe.

### `puzzle0.txt`

```file:
loesungen/puzzle0.txt
```

### `puzzle1.txt`

```file:
loesungen/puzzle1.txt
```

### `puzzle2.txt`

```file:
loesungen/puzzle2.txt
```

### `puzzle3.txt`

```file:
loesungen/puzzle3.txt
```

### `puzzle4.txt`

Zusätzliches Beispiel:

```file:
beispieldaten/puzzle4.txt
```

```file:
loesungen/puzzle4.txt
```
