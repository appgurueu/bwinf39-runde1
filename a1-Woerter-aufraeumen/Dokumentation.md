# Lösung Aufgabe 1 "Wörter aufräumen"

## Lösungsidee

Zunächst muss man überlegen, wann ein Wort "passt". Dies kann nur gegeben sein, wenn an den richtigen Stellen die richtigen Buchstaben stehen, und das Wort die passende Länge besitzt. Dadurch erhalten wir pro Lücke mehrere passende Wörter. Kommt für eine Lücke nur ein Wort infrage, muss dieses eingesetzt werden. Das Wort kann dann gestrichen werden: Es ist kann nicht mehr als passendes Wort infrage kommen. Dadurch wiederholt sich dann für weitere Wörter diese Situation, so lange, bis alle Wörter eingesetzt sind.

### Korrektheit

Damit für eine Lücke eindeutig feststeht, welches Wort einzusetzen ist - damit also die Lösung des Rätsels eindeutig ist -, darf nur ein Wort "passen". Dieses muss dann wieder zu eindeutig passenden Wörtern führen, so lange, bis alle Lücken gefüllt sind.

### Komplexität

Das Finden der passenden Wörter hat im Worst-Case quadratische Komplexität: Für jede der $n$ Lücken müssen alle $n$ Wörter geprüft werden, ob sie passen.

Das Einsetzen hat ebenfalls quadratische Komplexität: Es müssen $n$ Lücken gefüllt werden. Bei bis zu $n - 1$ weiteren Lücken muss das Wort als passendes Wort gestrichen werden.

Insgesamt ergibt sich also eine **Komplexität von $O(n^2)$**

## Umsetzung

Implementierung in der modernen und performanten Programmiersprache [Go](https://golang.org).

### Kompilieren

`go build` (erzeugt `a5-Wichteln`) oder `go build main.go` (erzeugt `main`)

### Verwendung

`go run main.go <pfad>` oder `./main <pfad>`

Beispiel: `./main beispieldaten/raetsel0.txt`

### Ausgabe

```
Lösung: <Satz mit eingesetzten Wörtern>
Zeit verstrichen: <Verstrichene Zeit in Millisekunden> ms

```

### Bibliotheken

* `fmt`: Ausgabe & Formattierung
* `io/ioutil`: Einlesen der Rätseldatei
* `os`: Programmargumente
* `strings`: Auftrennen ("splitten") von Text
* `time`: Zeitmessung
* `unicode`: `IsLetter`-Funktion

### Typen

#### `Text`

Rune-Slice: repräsentiert einen `string` in UTF-32

Zwar speicher-ineffizient im Vergleich zu Go UTF-8 `string`s, aber praktischer zu manipulieren und Buchstaben zu extrahieren:

```go
utf8 := "Äpfel"
utf32 := Text(utf8)
// In UTF-8 sind Buchstaben wie "ä" zwei Zeichen
println(len(utf8)) // 6
// In UTF-32 ist ein Buchstabe genau ein Zeichen
println(len(utf32)) // 5
// Entsprechend fällt auch die Extraktion leichter aus
println(utf32[0] == 'Ä') // true
// Ebenso wie die Manipulation
raetsel := Text("Ä____ sind wohlschmeckend")
copy(raetsel[0:6], utf32)
println(string(raetsel)) // Äpfel sind wohlschmeckend
```

#### `Wort`

`Text` und die Anzahl der Vorkommen in der Wortliste.

#### `WortMitLuecken`

Auszufüllendes Lückenwort:

* Position im Lückentext
* Länge
* Gegebene Buchstaben als \[position] = UTF-32-Zeichen
* Infragekommende ("passende") Wörter als \[id] = true

### Ablauf

Die Programmausführung beginnt mit dem Einlesen der durch die Programmargumente angegeben Datei und dem Zerlegen in zwei Zeilen.
Die erste Zeile - der Lückentext - wird in einen UTF-32-Text konvertiert, die Wörter werden aus der zweiten Zeile extrahiert in dem man diese nach Leerzeichen "splittet". Mit einer Abbildung von Wort nach ID wird sichergestellt, das Wörter nur einmal in den "Wortindex", der Wörter nach Länge und ID kategorisiert, aufgenommen werden - bei mehrmaligem Vorkommen wird einfach die Anzahl erhöht (siehe [`Wort`]). Diese ID dient als Verweis auf das Wort (andere Programmiersprachen wie etwa [Lua](https://lua.org) übernehmen dies sogar für einen; dies erlaubt dann Stringvergleiche in konstanter Zeit und beschleunigt Maps mit String-Keys - und spart zusätzlich noch Speicher).

Das Bestimmen passender Wörter ist relativ einfach: Mit zwei geschachtelten Schleifen prüft man für alle Wörter der passenden Länge, ob sie an jeder gegebenen Stelle den richtigen Buchstaben besitzen.

Weiter wird ein Index für Lücken nach passendem Wort angelegt.
Die zum Einsetzen verwendete Funktion nimmt nun eine Lücke ([`WortMitLuecken`]) und die ID eines einzusetzenden Wortes. Das Wort wird dann eingesetzt, indem der entsprechende Bereich der Slice überschrieben wird. Die Anzahl der noch einzusetzenden Vorkommen wird dann verringert, sinkt sie auf 0 wird das Wort gestrichen: Überall, wo das Wort gepasst hätte - hier verwenden wir oben genannten Index für Lücken nach Wort - streichen wir das Wort. Passt dann nur noch ein Wort für eine Lücke, rufen wir rekursiv die Einsetzen-Funktion für diese Lücke und das passende Wort auf. Diese rekursiven Aufrufe halten so lange an, bis alle Wörter eingesetzt sind.

"Wörter mit Lücken" werden aus der ersten Zeile extrahiert, indem diese zeichenweise durchgegangen ist. Buchstaben (nach `unicode.IsLetter`, also auch Unterstützung für andere Sprachen) und Leerzeichen sind Teil von Lückenworten, alles andere nicht. Wird ein Lückenwort durch ein Zeichen terminiert, kann die Suche nach passenden Wörtern gestartet werden, die eventuell schon das direkte Einsetzen des Wortes zur Folge haben kann.

Für die Ausgabe wird schließlich der UTF-32 Lösungstext in einen UTF-8 string umgewandelt. Zusätzlich wird noch die verstrichene Zeit ausgegeben.

## Quellcode

### `main.go`

```file:go
main.go
```

## Beispiele

### `raetsel0.txt`

```file:
loesungen/raetsel0.txt
```

### `raetsel1.txt`

```file:
loesungen/raetsel1.txt
```

### `raetsel2.txt`

```file:
loesungen/raetsel2.txt
```

### `raetsel3.txt`

```file:
loesungen/raetsel3.txt
```

### `raetsel4.txt`

```file:
loesungen/raetsel4.txt
```

### `raetsel5.txt`

Zusätzliches Beispiel:

```file:
beispieldaten/raetsel5.txt
```

```file:
loesungen/raetsel5.txt
```
