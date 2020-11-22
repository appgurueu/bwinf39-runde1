# Lösung Aufgabe 5 "Wichteln"

## Lösungsidee

Wir verwenden einen optimierenden Algorithmus. Die Schüler werden in einer Reihenfolge durchgegangen. In dieser werden zuerst erste Wünsche erfüllt, dann Zweite, und schließlich Dritte.

Diese Reihenfolge verbessern wir solange möglich, indem wir immer zwei Schüler vertauschen, die neue Reihenfolge ausprobieren, und beibehalten falls diese besser ist.

### Korrektheit

Wenn eine Reihenfolge nicht mehr verbesserbar ist, muss sie die eine optimale Reihenfolge sein.

Es bleibt zu zeigen, dass jede optimale Geschenkeverteilung auch als Reihenfolge darstellbar ist.

Bei einer optimalen Geschenkeverteilung wird jedem Schüler:

* Der 1. Wunsch
* Oder der 2. Wunsch
* Oder der 3. Wunsch
* Oder kein Wunsch

erfüllt.

Damit nun einem Schüler sein erster Wunsch erfüllt wird, brauchen wir ihn nur an eine Stelle zu stellen, wo das Geschenk noch nicht vergeben sein wird.

Analog können wir auch zweite und dritte Wunsch erfüllen, indem wir den betreffenden Schüler an einer Stelle in der Reihenfolge einsetzen, wo sein erster Wunsch bzw. sein erster und sein zweiter Wunsch schon vergeben sein werden (da Schüler, die weiter vorne in der Reihenfolge stehen, diese erhalten).

### Komplexität

* Das Ausprobieren einer Reihenfolge ist in linearer Laufzeit $O(n)$ möglich.
* Alle möglichen Täusche sind $\sum_{i=1}^{n-1} i = \frac{n(n-1)}{2} = O(n^2)$
* Somit ergibt sich eine **kubische Komplexität von $O(n^3)$**

## Umsetzung

Implementierung in der modernen und performanten Programmiersprache [Go](https://golang.org).

### Kompilieren

`go build` (erzeugt `a5-Wichteln`) oder `go build main.go` (erzeugt `main`)

### Verwendung

`go run main.go <pfad>` oder `./main <pfad>`

Beispiel: `./main beispieldaten/wichteln1.txt`

### Ausgabe

```
Lösung:
<In Textrichtung Ausgabe der erhaltenen Geschenke getrennt mit Komma und Leerzeichen, aufsteigend nach Nummer des erhaltenden Schülers>

Erfüllte Wünsche: <Anzahl erfüllte 1. Wünsche>, <Anzahl erfüllte 2. Wünsche>, <Anzahl erfüllte 3. Wünsche>

Zeit verstrichen: <Verstrichene Zeit in Sekunden> s

```

### Bibliotheken

* `fmt`: Ausgabe, Formattierung
* `io/ioutil`: Einlesen der Datei
* `os`: Programmargumente
* `regexp`: Regulärer Ausdruck zum Extrahieren der Zahlen
* `strconv`: String/Integer-Konversion
* `strings`: Auftrennen des Dateiinhalts nach Zeilen
* `time`: Zeitmessung

### Typen

#### `Verteilung`

* Wert als vergleichbare Zahl: Anzahl erfüllter Wünsche so kodiert, dass zwei Verteilungen über ihren Wert vergleichbar sind
* Erhaltene Geschenke als "Slice" mit \[Schülernummer - 1] = Geschenknummer
* Vergebene Geschenke als "Slice" mit \[Geschenknummer - 1] = `true` wenn vergeben, sonst `false`

### Eingabe

Das erste Programmargument ist der Pfad zur Aufgabendatei. Diese wird gelesen und an Zeilenumbrüchen aufgetrennt.
Die Wünsche jedes Schülers werden mithilfe eines einfachen regulären Ausdrucks als Strings extrahiert und zu Zahlen konvertiert.
Schließlich erhält man eine Slice mit \[Schülernummer-1] = 3-er-Array{Geschenknummer 1. Wunsch - 1, Geschenknummer 2. Wunsch - 1, Geschenknummer 3. Wunsch - 1}

### Verarbeitung

Eine Funktion probiert die aktuelle Reihenfolge aus, indem zuerst erste, dann zweite, und schließlich dritte Wünsche in der Reihenfolge erfüllt werden. Hierfür werden vergebene und erhaltene Geschenke mit einer [`Verteilung`] nachgehalten. Der Wert wird als Zahl (`uint64`) zur Basis $n + 1$ mit $n$ = Anzahl Schüler dargestellt. Hierbei stehen erste Wünsche an erster Stelle, 2. an 2. und 3. an 3.: $wert = (n+1)^2 \cdot erfuellteErsteWuensche + (n+1) \cdot erfuellteZweiteWuensche + erfuellteDritteWuensche$. Der Vergleich zweier Verteilungen wird somit zu einem einfachen Zahlenvergleich.

Wir beginnen mit der Einlesereihenfolge der Schüler als Startreihenfolge.

Dann probieren wir solange Swaps (Täusche) aus, bis keiner der möglichen Swaps mehr zu einer Verbesserung des Wertes der aktuellen Reihenfolge führt.
Einen Swap machen wir rückgängig, wenn sich herausstellt, dass dieser zu keiner Verbesserung geführt hat. Insgesamt ergeben sich 3 geschachtelte Schleifen.

### Ausgabe

Zunächst wird die Verteilung komplettiert: Schüler, die keinen Wunsch erfüllt bekommen haben, erhalten die erstbesten freien Geschenke (in Reihenfolge der Geschenkenummern). Hierfür gehen wir die erhaltenen Geschenke nach Schüler durch, finden leer ausgehende Schüler, und gehen dann die Geschenke durch, anfangend nach den schon durchgegangen Geschenken, bis wir ein noch nicht vergebenes finden. Dieses erhält der Schüler.

Schließlich geben wir "einfach" die erhaltenen Geschenke aus. Diese trennen wir in einer Zeile mit Leerzeichen und Komma, sonst durch Zeilenumbrüche. Hierbei sorgen wir für Zeilen, die nicht länger als 80 Zeichen werden.

Die Anzahl der erfüllten Wünsche extrahieren wir aus dem Wert der Verteilung mittels Division mit Rest und geben diese zusammen mit der verstrichenen Zeit ebenfalls aus.

## Quellcode

### **`main.go`**

```file:go
main.go
```

## Beispiele

### `wichteln1.txt`

```file:
loesungen/wichteln1.txt
```

### `wichteln1.txt`

```file:
loesungen/wichteln1.txt
```

### `wichteln2.txt`

```file:
loesungen/wichteln2.txt
```

### `wichteln3.txt`

```file:
loesungen/wichteln3.txt
```

### `wichteln4.txt`

```file:
loesungen/wichteln4.txt
```

### `wichteln5.txt`

```file:
loesungen/wichteln5.txt
```

### `wichteln6.txt`

```file:
loesungen/wichteln6.txt
```

### `wichteln7.txt`

```file:
loesungen/wichteln7.txt
```

### `wichteln8.txt`

Zusätzliches Beispiel:

```file:
beispieldaten/wichteln8.txt
```

```file:
loesungen/wichteln8.txt
```
