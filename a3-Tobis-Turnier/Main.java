import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.function.BiFunction;
import java.util.concurrent.ThreadLocalRandom;
import java.util.function.BiConsumer;
import java.util.function.Function;
import java.util.function.Supplier;

/**
 *
 * @author appguru
 */
public class Main {
    static class RecursiveBiFunction<T, U, V> implements BiFunction<T, U, V> {
        public BiFunction<T, U, V> function;

        @Override
        public V apply(T t, U u) {
            return function.apply(t, u);
        }
        
    }
    public static void main(String[] args) throws IOException, Exception {
        long millis = System.currentTimeMillis();
        if (args.length != 1) {
            System.out.println("Verwendung: <pfad>");
            return;
        }
        String text = new String(Files.readAllBytes(Path.of(args[0])));
        String[] zeilen = text.split("\n");
        int anzahl = Integer.parseInt(zeilen[0]);
        int[] spielstaerken = new int[anzahl];
        int _besterSpieler = 0;
        for (int i = 0; i < anzahl; i++) {
            int spielstaerke = Integer.parseInt(zeilen[i + 1]);
            spielstaerken[i] = spielstaerke;
            if (spielstaerke > spielstaerken[_besterSpieler]) {
                _besterSpieler = i;
            }
        }
        final int besterSpieler = _besterSpieler;
        for (int i = 0; i < anzahl; i++) {
            if (i != besterSpieler && spielstaerken[i] == spielstaerken[besterSpieler]) {
                throw new Exception("Mehrere beste Spieler");
            }
        }
        BiFunction<Integer, Integer, Boolean> spieler1Gewinnt = (Integer spieler1, Integer spieler2) -> ThreadLocalRandom.current().nextInt(0, spielstaerken[spieler1] + spielstaerken[spieler2]) < spielstaerken[spieler1];
        BiFunction<Integer, Integer, Integer> gewinner = (Integer spieler1, Integer spieler2) -> spieler1Gewinnt.apply(spieler1, spieler2) ? spieler1 : spieler2;
        Supplier<Integer> liga = () -> {
            int[] siege = new int[anzahl];
            for (int spieler1 = 0; spieler1 < anzahl; spieler1++) {
                for (int spieler2 = spieler1 + 1; spieler2 < anzahl; spieler2++) {
                    siege[gewinner.apply(spieler1, spieler2)]++;
                }
            }
            int meisteSiege = 0;
            for (int spieler = 0; spieler < anzahl; spieler++) {
                if (siege[meisteSiege] < siege[spieler]) {
                    meisteSiege = spieler;
                }
            }
            if (meisteSiege == besterSpieler) {
                return 1;
            }
            return 0;
        };
        Function<BiFunction<Integer, Integer, Integer>, Supplier<Integer>> koVariante = vGewinner -> {
            return () -> {
                int[] turnierplan = new int[anzahl];
                for (int i = 0; i < anzahl; i++) {
                    turnierplan[i] = i;
                }
                for (int i = 0; i < anzahl; i++) {
                    int j = ThreadLocalRandom.current().nextInt(i, anzahl);
                    int back = turnierplan[j];
                    turnierplan[j] = turnierplan[i];
                    turnierplan[i] = back;
                }
                RecursiveBiFunction<Integer, Integer, Integer> sieger = new RecursiveBiFunction();
                sieger.function = (Integer start, Integer ende) -> {
                    int diff = ende - start;
                    if (diff == 1) {
                        return vGewinner.apply(turnierplan[start], turnierplan[ende]);
                    }
                    int mitte = start + diff / 2;
                    return vGewinner.apply(sieger.apply(start, mitte), sieger.apply(mitte, ende));
                };
                if (sieger.apply(0, anzahl - 1) == besterSpieler) {
                    return 1;
                }
                return 0;
            };
        };
        Supplier<Integer> ko = koVariante.apply(gewinner);
        Supplier<Integer> ko5 = koVariante.apply((spieler1, spieler2) -> {
            int siegeSpieler1 = 0;
            for (int i = 0; i < 5; i++) {
                if (spieler1Gewinnt.apply(spieler1, spieler2)) {
                    siegeSpieler1++;
                    if (siegeSpieler1 == 3) {
                        return spieler1;
                    }
                }
            }
            return spieler2;
        });
        int anzahlLaeufe = 1000_000;
        BiConsumer<String, Supplier<Integer>> testeTurniervariante = (String name, Supplier<Integer> runde) -> {
            int siegeBesterSpieler = 0;
            for (int laeufe = 0; laeufe < anzahlLaeufe; laeufe++) {
                siegeBesterSpieler += runde.get();
            }
            System.out.println(name + ": " + (siegeBesterSpieler * 100 / (float) anzahlLaeufe) + " %");
        };
        testeTurniervariante.accept("Liga", liga);
        testeTurniervariante.accept("K.O.", ko);
        testeTurniervariante.accept("K.O.x5", ko5);
        System.out.println("Zeit verstrichen: " + ((System.currentTimeMillis() - millis) / (float) 1000) + " s");
    }

}
