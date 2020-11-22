# Lösung Aufgabe 3 "Tobis Turnier"

## Lösungsidee

Die drei Turniervarianten werden nach gegebener Spezifikation implementiert und zahlreiche Durchläufe simuliert.

Entsprechend ist die eigentliche Arbeit hauptsächlich **Umsetzung** der Spezifikation.

### Genauigkeit

Abhängig von der gleichmäßigen Streuung der Zufallszahlen des verwendeten Zufallsgenerators und der Anzahl der simulierten Durchläufe.

### Komplexität

Sei $n$ die Anzahl der Spieler. Dann benötigen die Turniervarianten jeweils folgend viele Spiele:

* Liga: $\sum_{i=1}^{n-1} i = \frac{n(n-1)}{2}$
* K.O.: $\sum_{i=1}^{log_2(n)} \frac{n}{2^i} = n$
* K.O.x5: 5-mal so viele Spiele wie K.O.: $5 \cdot n$

Offensichtlich ist Liga am aufwendigsten mit quadratisch vielen Spielen ($S(n) = O(n^2)$). K.O. und K.O.x5 dahingegen benötigen nur linear viele ($S(n) = O(n)$).

## Umsetzung

Implementierung in der modernen und performanten Programmiersprache [Go](https://golang.org).

### Kompilieren

`go build` (erzeugt `a5-Wichteln`) oder `go build main.go` (erzeugt `main`)

### Verwendung

`go run main.go <pfad>` oder `./main <pfad>`

Beispiel: `./main beispieldaten/spielstaerken1.txt`

### Ausgabe

```
Liga: <Siege Bester Spieler bei 1.000.000 Liga-Turnieren in Prozent> %
K.O.: <Siege Bester Spieler bei 1.000.000 K.O.-Turnieren in Prozent> %
K.O.x5: <Siege Bester Spieler bei 1.000.000 K.O.x5-Turnieren in Prozent> %
Zeit verstrichen: <Verstrichene Zeit in Sekunden> s

```

Das Programm lässt sich nach dem EVA-Prinzip in Eingabe, Verarbeitung und Ausgabe gliedern:

### Bibliotheken

* Eingabe:
  * `os`: Programmargumente zum Erhalten des Pfades
  * `io/ioutil`: Hilfsbibliothek ("utility") zum Einlesen der Datei mit den Spielstärken.
  * `strings`: Unterteilen der Datei in Zeilen ("split")
  * `strconv`: Umwandlung von Strings in Zahlen
* Verarbeitung:
  * `time`: "Seeden" des ansonsten determinierten Zufallsgenerators von Go, Zeitmessung
  * `math/rand`: Zufallsgenerator
* Ausgabe: `fmt`: Formattierung. Nötig für Ausgabe von Zahlen (implementiert Funktionalität u.a. wie C's `printf`)

### Eingabe

Das erste Argument wird als Dateipfad verstanden. Die Datei wird eingelesen und in Zeilen unterteilt. Die Spielstärken werden mit einer Schleife zeilenweise zu Zahlen konvertiert und in einer natürliche-Zahlen-"Slice" fester Länge gespeichert.

### Verarbeitung

Zuerst wird ein "sanity-check" mit den Beispieldaten durchgeführt: Es darf nur einen besten Spieler geben, ansonsten bricht das Programm aufgrund fehlerhafter Eingaben ab.

Zentral ist zunächst eine Funktion, die den Gewinner eines einzigen Spiels bestimmt. Diese ist nach Spezifikation:

* Erster Spieler gewinnt, wenn seine Kugel gezogen wird
  * Die gezogene Kugel ist eine Zufallszahl von 1 bis zu den addierten Spielstärken
  * Diese ist eine Kugel des ersten Spielers, wenn sie <= der Spielstärke des ersten Spielers ist
    * Anschaulich: Von den nummerierten Kugeln gehören die ersten Spielstärke-viele Kugeln dem ersten Spieler, alle "danach" dem Zweiten
* Sonst gewinnt der zweite Spieler

#### Liga

In der Implementation von Liga muss nach Spezifikation jeder Spieler gegen jeden anderen einmal spielen.
Entsprechend geht man alle Spieler durch. Für jeden Spieler iteriert man dann über alle Anderen mit einer höheren Spielernummer als Kontrahenten (zwei geschachtelte Schleifen). An einem einfachen Beispiel mit drei Spielern 1, 2, 3 wird sofort klar, wieso dies funktioniert:

1. Betrachte Spieler 1
  1. Spiele gegen Spieler 2
  2. Spiele gegen Spieler 3
2. Betrachte Spieler 2
  3. Gegen Spieler 1 muss nicht mehr gespielt werden
  4. Spiele Gegen Spieler 3
3. Betrachte Spieler 3
  5. Gegen keinen Spieler muss noch gespielt werden

Für die Siege der Spieler wird eine Slice angelegt mit \[Spielernummer] = Siege. Nach Simulieren aller Spiele wird der Spieler mit den meisten Siegen bestimmt (einfache Maximumsuche), wobei nach Spezifikation bei Gleichstand der Spieler mit der kleineren Spielernummer gewinnt.

#### K.O.

Wir implementieren K.O. rekursiv:

##### Rekursionsanfang

**Sieger des Turniers** ist derjenige, der siegt, wenn wir den kompletten Turnierplan als Ausschnitt wählen.

##### Rekursiver Aufruf

Wir betrachten einen Ausschnitt / Teil des Turnierplans.
Sieger des Ausschnittes ist derjenige Spieler, der im Spiel zwischen dem Sieger der linken Hälfte des Ausschnitts und dem Sieger der rechten Hälfte des Ausschnitt siegt.

##### Rekursionsende

Umfasst der betrachtete Ausschnitt nur zwei Spieler, ist der Sieger derjenige, der im Spiel der beiden siegt.

#### K.O.x5

Wir nutzen die K.O.-Implementation, ersetzen aber die Funktion, die entscheidet, wer ein Spiel gewinnt:

Anstatt eines einzigen Spiels werden bis zu fünf Spiele simuliert. Hat der erste Spieler drei gewonnen, gewinnt er und es wird abgebrochen. Kommt dies nicht vor, gewinnt sein Kontrahent.

#### Simulation

Schließlich werden alle Turniere mit einer einfachen Schleife eine Million Male gespielt. Dabei wird eine Zählvariable erhöht, wenn der beste Spieler gewinnt.

### Ausgabe

Mit `fmt` wird nach jeder Simulation eine Zeile im Format `<Turniervariante>: <Siege erster Spieler in Prozent> %` ausgegeben.

Am Ende der Programmausführung steht die Ausgabe der verstrichenen Zeit.

## Quellcode

### **`main.go`**

```file:go
main.go
```

## Beispiele

### `spielstaerken1.txt`

```file:
loesungen/spielstaerken1.txt
```

### `spielstaerken2.txt`

```file:
loesungen/spielstaerken2.txt
```

### `spielstaerken3.txt`

```file:
loesungen/spielstaerken3.txt
```

### `spielstaerken4.txt`

```file:
loesungen/spielstaerken4.txt
```

### `spielstaerken5.txt`

Zusätzliches Beispiel:

```file:
beispieldaten/spielstaerken5.txt
```

```file:
loesungen/spielstaerken5.txt
```

### Fazit

* K.O.x5 weist immer bessere Siegesquoten des besten Spielers auf als einfaches K.O., da die Wahrscheinlichkeit, dass der schlechtere Spieler in einem Aufeinandertreffen gewinnt, weiter gesenkt wird, indem dieser den Großteil Spiele gewinnen müsste. Insgesamt ist K.O.x5 genauer als K.O.
* Durchschnittlich erweist sich K.O.x5 auch im Vergleich zur Liga als zuverlässiger (35 % vs. 64 %, 21 % vs. 37 %, 31.5 % vs. 31 %, 11.5 % vs. 8 %, 99.8 % vs. 99.7 %).
* Im Falle vieler ähnlich starker Kontrahenten (`spielstaerken4.txt`) hat der beste Spieler es schwieriger, sich im K.O.x5 zu behaupten, da der Verlust eines einzigen Aufeinandertreffens reicht, damit er verliert. Hier ist Liga etwas zuverlässiger (11.5 % vs. 8 %).
* Insgesamt empfehle ich Tobi entsprechend **K.O.x5**