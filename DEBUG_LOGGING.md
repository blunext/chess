# Debug Logging

Silnik tworzy dwa pliki logów:

## game.log
Główny log z ruchami i wynikami search:
```
12:21:10 | M/S: Qa2b2 | Sc: -31080 cp | Ns: 10822027 | T: 16.93s
```

## debug.log
Szczegółowy log iterative deepening (każda głębokość):
```
18:17:58 | M/S: START      | Sc: moves=27 first=a2b2 TT=48MB | T: 5s
18:17:58 | M/S: D1:a2a1    | Sc: -340     | Ns: 190      | T: 20ms
18:17:58 | M/S: D2:a2a1    | Sc: -345     | Ns: 1469     | T: 180ms
18:17:58 | M/S: ACCEPT     | Sc: d2=a2a1  | Accepted depth 2 result
18:17:58 | M/S: D3:a2a1    | Sc: -345     | Ns: 3972     | T: 410ms
18:17:58 | M/S: ACCEPT     | Sc: d3=a2a1  | Accepted depth 3 result
18:18:02 | M/S: D4:a2a1    | Sc: -335     | Ns: 41419    | T: 4.66s
18:18:02 | M/S: ACCEPT     | Sc: d4=a2a1  | Accepted depth 4 result
18:18:02 | M/S: TIMECUT    | Sc: 3.7x     | Time cutoff (would take too long)
18:18:02 | M/S: FINAL:a2a1 | Sc: -335     | Final chosen move
```

### Informacje w START:
- **moves=N** - ile ruchów legalnych do oceny
- **first=X** - pierwszy ruch po sortowaniu (oceniany jako pierwszy)
- **TT=NMB** - rozmiar transposition table

### Informacje o zmianie ruchu:
Gdy najlepszy ruch zmienia się między głębokościami:
```
D3:g7g6 CHANGED(a2a1->g7g6) | Sc: +50
```

### Informacje w REJECT:
Gdy search jest odrzucony (timeout w trakcie):
```
REJECT | Sc: d5_stopped prev=a2a1 | (poprzedni zaakceptowany ruch)
```

### Typy wpisów w debug.log:

- **START** - Początek search z time limitem
- **BOOK:move** - Ruch z książki otwarć
- **DN:move** - Wynik z głębokości N (np. D1:a2a1)
  - Status: **OK** (normalne zakończenie) lub **STOPPED** (timeout)
- **ACCEPT** - Zaakceptowano wynik z danej głębokości
- **REJECT** - Odrzucono wynik (search był przerwany)
- **MATE** - Znaleziono mata, przerywamy search
- **TIMECUT** - Przerwano bo następna iteracja zajęłaby za długo
- **FINAL:move** - Ostateczny wybrany ruch

### Jak używać do debugowania:

Jeśli silnik zagra dziwny ruch (np. a2b2 z score -31080), debug.log pokaże:
1. Jakie ruchy były najlepsze na każdej głębokości
2. Czy search był przerwany (STOPPED)
3. Czy któryś wynik został odrzucony (REJECT)
4. Jaki był finalny wynik (FINAL)

To pozwoli zidentyfikować czy:
- Silnik faktycznie wybrał zły ruch (FINAL != oczekiwany)
- Search był przerwany podczas iteracji
- Wyniki zmieniały się między głębokościami

### Włączanie/wyłączanie:

Debug logging jest włączony domyślnie w `uci/uci.go:Start()`.
Aby wyłączyć, zakomentuj linię:
```go
// uci.session.SetDebugLogger(debugLogger)
```
