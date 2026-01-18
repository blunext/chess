# Chess Engine Roadmap

> **ZASADA: Pracujemy iteracyjnie - jedna mała rzecz na raz, nie wiele naraz.**

---

# Move Generator Roadmap

## ✅ Iteracja 1: Move Struct + Sliding Pieces
- [x] Struktura `Move` w `board/move.go`
- [x] `GenerateSlidingMoves()` (Goniec, Wieża, Hetman)
- [x] Optymalizacja: bit-scanning zamiast `ToSlice()`
- [x] Optymalizacja: usunięcie `filterColor()`

## ✅ Magic Bitboards
- [x] Generator magic numbers (`magic/generate.go`)
- [x] Testy weryfikujące poprawność (`magic/magic_test.go`)
- [x] Integracja z `GenerateMoves()` - O(1) lookup dla sliding pieces
- [x] Funkcje `rookAttacks()`, `bishopAttacks()`

## ✅ Iteracja 2: Skoczek + Król
- [x] Dodanie skoczka do `GenerateMoves`
- [x] Dodanie króla do `GenerateMoves` (bez roszad)
- [x] Testy

## ✅ Iteracja 3: Bicia
- [x] Rozszerzenie logiki o ruchy na pola przeciwnika
- [x] Ustawienie `Move.Captured` dla Bishop/Rook/Queen
- [x] Funkcja `pieceAt()` do wykrywania typu zbitej figury
- [x] Bicia dla Knight/King

## ✅ Optymalizacje
- [x] Pre-alokacja slice: `make([]Move, 0, 64)`
- [x] Cache `ourPieces`, `enemyPieces`, `allPieces`
- [x] Jedna alokacja zamiast 5 (append pattern)
- [x] Benchmarki w `bench/moves_test.go`
- [ ] `PieceMoves` jako array zamiast map
- [ ] Kompaktowa struktura Move (uint32)
- [ ] Unikanie switch w hot loop

## ✅ Iteracja 4: Piony (podstawowe)
- [x] Ruch 1 pole do przodu
- [x] Ruch 2 pola z pozycji startowej
- [x] Bicia ukośne
- [x] Obsługa białych i czarnych pionów
- [x] Zapobieganie wrap-around na krawędziach (fileA/fileH masks)

## ✅ Iteracja 5: Rozszerzenie struktury Move
- [x] Pole `Promotion Piece` (Q/R/B/N)
- [x] Flaga `Flags` z `FlagEnPassant` i `FlagCastling`
- [x] UCI notation: `ToUCI()` (e2e4, e7e8q)
- [x] Zaktualizowany `String()` z obsługą promocji i flag

## ✅ Iteracja 6: Piony (specjalne)
- [x] Promocja (generowanie 4 ruchów: Q/R/B/N)
- [x] En passant (bicie w przelocie)

## ✅ Iteracja 7: Roszady
- [x] Kingside (O-O)
- [x] Queenside (O-O-O)
- [x] Sprawdzenie praw (CastleSide flags)
- [x] Sprawdzenie blokad (pola między K-R puste)
- [x] Sprawdzenie ataków (król nie przechodzi przez szach)

## ✅ Iteracja 8: Legalność
- [x] `MakeMove()` / `UnmakeMove()` - wykonanie i cofnięcie ruchu
- [x] `GenerateLegalMoves()` - filtr ruchów pozostawiających króla w szachu
- [x] Obsługa wszystkich specjalnych ruchów (en passant, roszady, promocje)

## ✅ Iteracja 9: Perft (weryfikacja)
- [x] Zliczanie ruchów na głębokość N
- [x] Porównanie z known perft results (Initial, Kiwipete, Position3)
- [x] Debug: divide (znaleziono i naprawiono en passant wrap-around bug) (perft per move)

---

# Engine Roadmap

## ✅ Iteracja 10: Sprawdzanie szacha
- [x] `isSquareAttacked(sq, byColor)` - czy pole jest atakowane
- [x] `isInCheck()` - czy król jest w szachu
- [x] Wykorzystanie magic BB do szybkiego sprawdzania ataków sliding pieces
- [x] Prekomputowane tablice ataków dla skoczka i króla

## Iteracja 11: Ocena pozycji
- [x] Materiał (wartości figur: P=100, N=320, B=330, R=500, Q=900)
- [x] Piece-Square Tables (PST)
- [x] Struktura pionów (zdwojone, izolowane, przechodzące)
- [ ] Aktywność figur (mobilność)
- [ ] Kontrola przestrzeni (Space bonus) - szczegóły poniżej
- [ ] Tuning PST (Piece-Square Tables) - szczegóły poniżej
- [ ] Bezpieczeństwo króla (King Safety) - szczegóły poniżej

### Bezpieczeństwo króla (King Safety)

> **Cel:** Karanie pozycji z wystawionym królem, brakującą osłoną pionkową lub pod atakiem.

#### Komponenty (od najprostszego):

**1. Pawn Shield (osłona pionkowa)** ✅
- [x] Bonus za pionki przed oroszowanym królem (+10 za 2 linię, +5 za 3 linię)
- [x] Kara za brakujące pionki w osłonie (-25)
- [x] Rozpoznawanie pozycji oroszowanej (król na g1/h1 lub a1/b1/c1)

**2. Open Files Near King (otwarte linie)** ✅
- [x] -25 za pół-otwartą linię obok króla
- [x] -40 za pełną otwartą linię

**3. Game Phase Scaling (skalowanie fazy gry)** ✅
- [x] Redukcja king safety gdy brak hetmanów (dzielenie przez 4)

**4. Uncastled King Penalty (nieoroszowany król)** ✅
- [x] -50 kara za króla na kolumnach d/e w middlegame

**5. King Tropism (opcjonalne - na później)**
- [ ] Bonus za bliskość figur atakujących do króla przeciwnika
- [ ] Wagi: Hetman (2x), Wieża (1x), Skoczek (1.5x), Goniec (1x)

**6. Pawn Storm (opcjonalne - na później)**
- [ ] Kara za pionki przeciwnika zbliżające się do naszego króla

### Kontrola przestrzeni (Space Bonus)

> **Cel:** Nagradzanie za kontrolę przestrzeni przez zaawansowane pionki. Silnik powinien preferować ruchy pionami które dają kontrolę nad centrum i terytorium przeciwnika.

**1. Central Pawn Bonus**
- [ ] Bonus za pionki na d4/e4/d5/e5 (+20-30 cp)
- [ ] Mniejszy bonus za c4/c5/f4/f5 (+10-15 cp)

**2. Space Calculation**
- [ ] Zliczanie pól kontrolowanych za linią pionów
- [ ] Bonus skalowany z liczbą figur (więcej figur = przestrzeń ważniejsza)
- [ ] Typowe wartości: +0.5 cp za każde kontrolowane pole

**3. Pawn Advancement Bonus**
- [ ] Bonus za zaawansowane pionki (rank 4-6) poza passed pawn bonus
- [ ] Skalowanie: +5 za rank 4, +10 za rank 5, +15 za rank 6

### Tuning PST (Piece-Square Tables)

> **Problem:** Obecne pawnPST bazuje na "Simplified Evaluation Function" która karze centralne pionki na początkowych pozycjach (d2/e2 = -20!), co powoduje że silnik preferuje ruchy figurami zamiast pionami.

**1. Pawn PST Fix** ✅
- [x] Usunąć negatywne wartości dla d2/e2 (było -20 → teraz +5)
- [x] Usunąć negatywne wartości dla c3/f3 (było -10 → teraz 0/+10)
- [x] Zwiększyć bonusy dla zaawansowanych pionów (rank 4-6)

**2. Middlegame vs Endgame PST** ✅
- [x] Osobne tablice dla middlegame i endgame (PeSTO tables)
- [x] Interpolacja między fazami gry (tapered eval, gamePhase 0-24)
- [x] W endgame król powinien być aktywny (egKingTable)

## Iteracja 12: Search
- [x] Minimax
- [x] Alpha-beta pruning
- [x] Move ordering (captures first → MVV-LVA)
- [x] Zobrist hashing
- [x] Opening book (Polyglot format)
- [x] Quiescence search (kontynuacja przeszukiwania dla bić)
- [x] Iterative deepening
- [x] Transposition table
- [ ] ~~Null Move Pruning~~ (wyłączone - patrz sekcja "Search pruning")

## ✅ Iteracja 13: Time Management (podstawowy)
- [x] Iterative Deepening (pogłębianie przeszukiwania: 1, 2, 3...)
- [x] Parsowanie wtime/btime/winc/binc w UCI
- [x] Przerwanie search gdy czas się kończy (timeout check co N węzłów)
- [x] Alokacja czasu (prosta heurystyka: czas/30)

## Iteracja 13b: Time Management (zaawansowany)

> **Priorytety:** 🔴 Krytyczne → 🟡 Ważne → 🟢 Nice-to-have

### ✅ Emergency Buffer
- [x] Odejmij 200ms od dostępnego czasu jako rezerwę na lag sieciowy
- Problem: Komunikacja z serwerem ma opóźnienie, silnik może przekroczyć czas
- Rozwiązanie: `allocated = max(calculated - 200ms, 50ms)`

### 🟡 Move Overhead (UCI Option)
- [ ] Opcja `Move Overhead` (margines czasowy konfigurowalny przez użytkownika)
- Lichess/Arena pozwalają ustawić (zwykle 100-300ms)
- Format: `option name Move Overhead type spin default 100 min 0 max 1000`

### 🟡 Soft/Hard Time Limit
- [ ] Soft limit: "spróbuj skończyć do X ms" (można kontynuować jeśli jest czas)
- [ ] Hard limit: "bezwzględnie przerwij przed Y ms"
- Przykład: soft=2000ms, hard=2800ms → jeśli skończę depth 6 w 1800ms, mogę spróbować depth 7

### 🟢 Adaptacyjna alokacja (Smart Time)
- [ ] Jedyny legalny ruch → zagraj natychmiast (0ms)
- [ ] Score stability: krótszy czas gdy wynik stabilny przez 3 głębokości
- [ ] Position complexity: więcej czasu na skomplikowane pozycje (dużo bić/szachów)
- [ ] Otwarcie (pierwsze 15-20 ruchów): mniej czasu (mamy książkę otwarć)

### 🟢 Pondering
- [ ] Myślenie w czasie przeciwnika (`go ponder`)
- [ ] Obsługa `ponderhit` (przeciwnik zagrał przewidziany ruch)
- [ ] Wymaga: predykcji najbardziej prawdopodobnej odpowiedzi

## ✅ Iteracja 14: UCI Options
- [x] Obsługa `setoption name X value Y`
- [x] Opcja `Hash` (rozmiar transposition table w MB)
- [ ] Opcja `Threads` (liczba wątków - przygotowanie pod multi-threading)
- [ ] Opcja `UCI_ShowWDL` (pokazywanie Win/Draw/Loss)
- [ ] Ponder (`go ponder`, `ponderhit`, `stop`)

---

# Dalsze optymalizacje (po podstawowej wersji)

## Search pruning

### Null Move Pruning (wyłączone - wymaga poprawek)
> **Problem:** Podstawowa implementacja NMP powoduje błędne odcinanie linii na głębszych poziomach (depth 7+), co prowadzi do złych wyborów ruchów.

- [x] Podstawowa implementacja (R=2, depth >= 3)
- [ ] **Verification search** - po null move cutoff, weryfikuj wynik pełnym przeszukiwaniem
- [ ] **Static eval check** - NMP tylko gdy staticEval >= beta
- [ ] **Dynamiczne R** - redukcja zależna od głębokości: R = 2 + depth/6
- [ ] **Threat detection** - nie rób NMP gdy są oczywiste groźby
- [ ] Re-enable po implementacji verification search

### Inne techniki pruning
- [ ] Late Move Reductions (LMR)
- [ ] Aspiration windows
- [ ] Principal Variation Search (PVS)
- [ ] Futility pruning

## Search Extensions

> **Cel:** Przedłużanie przeszukiwania w krytycznych sytuacjach, aby nie przegapić taktyki.

- [x] **Check Extensions** - +1 ply gdy pozycja jest w szachu (najważniejsze!)
- [ ] **Single Reply Extensions** - +1 ply gdy jest tylko jeden legalny ruch
- [ ] **Recapture Extensions** - +1 ply przy odbiciu na tym samym polu
- [ ] **Passed Pawn Extensions** - +1 ply dla promocji pionów przechodzących

## Regression Testing (ochrona przed błędami)

> **Cel:** Wykrywanie regresji po zmianach - czy silnik nadal gra poprawnie?

### ✅ Perft (move generator)
- [x] Perft dla 6 standardowych pozycji (Initial, Kiwipete, Position 3-6)
- [x] Głębokości 1-4 (szybkie), 5-6 (slow tests)
- [x] Weryfikacja en passant, roszad, promocji

### ✅ Tactical Test Suite (search + eval)
- [x] Mate in 1-3 (10+ pozycji) - silnik MUSI znaleźć mata
- [x] Win material (10+ pozycji) - widelce, związania, odkryte ataki
- [x] WAC subset (35 pozycji) - klasyczne pozycje z "Win At Chess"
- [x] Defensive positions (5+ pozycji) - musi bronić, nie stracić materiału
- [x] Test runner: sprawdza czy silnik znajduje bestMove w limicie głębokości/czasu

### WAC Failures to Investigate
> Te pozycje failują - zbadać czy to bug w silniku czy problem z konwersją SAN→UCI
> Pozycje zakomentowane w `engine/tactical_test.go`

- [ ] **WAC.002**: Engine finds `c4c3`, expected `b3b2` (Rxb2) - endgame pawn capture

> ✅ Naprawione po wyłączeniu hasMateInOne (22x speedup):
> - WAC.003, WAC.007, WAC.022, WAC.040, WAC.083
>
> ✅ Naprawione po poprawkach TT:
> - WAC.009 - problem był w kolejności TT probe vs check extension
>   - Bug: TT probe używał depth PRZED check extension, Store używał depth PO extension
>   - Fix: przenieść check extension PRZED TT probe
>   - Dodatkowo: poprawiona logika TT bounds (nie modyfikować alpha/beta, tylko cutoff)
>
> 📊 Status testów (2026-01-18):
> - `TestTacticalSuite` (depth-based): 35/35 (100%)
> - `TestTacticalSuiteWithTime` (1s limit): >70% threshold - PASS

### Tactical Positions to Verify
> Pozycje które wymagają ręcznej weryfikacji - czy FEN i oczekiwany ruch są poprawne?
> Zakomentowane w `engine/tactical_test.go`

- [ ] **Knight fork: King and Queen** - Nd5 NIE atakuje K na e8 ani Q na d8!
  - FEN: `r1bqk2r/pppp1ppp/2n2n2/4p3/1bB1P3/2N2N2/PPPP1PPP/R1BQK2R w KQkq - 0 1`
  - Nd5 atakuje: f6 (skoczek), b4 (goniec) - to nie jest royal fork
  - Potrzeba: znaleźć prawdziwy royal fork (skoczek atakuje K i Q jednocześnie)

#### Jak weryfikować pozycje taktyczne:
1. Wczytaj FEN w lichess.org/editor lub chess.com/analysis
2. Sprawdź czy oczekiwany ruch jest legalny
3. Sprawdź czy ruch faktycznie realizuje opisaną taktykę (fork, pin, etc.)
4. Zweryfikuj z silnikiem (Stockfish) czy to najlepszy ruch

### Search Determinism
- [ ] Fixed-depth tests: ten sam depth = ten sam ruch i score
- [ ] Benchmark positions z zapisanymi expected values
- [ ] Wykrywanie czy "optymalizacja" przypadkiem nie zmienia wyników

### Self-play Tournament (opcjonalne)
- [ ] Nowa wersja vs stara wersja (100+ partii)
- [ ] Statystyczna weryfikacja że siła gry nie spadła
- [ ] Narzędzie: cutechess-cli lub własny skrypt

## Move ordering zaawansowane

> **Cel:** Lepsze sortowanie ruchów = szybsze cutoffs = głębsze przeszukiwanie

- [ ] **Killer moves** - 2 sloty na głębokość dla ruchów które spowodowały cutoff
- [ ] **History heuristic** - tablica [from][to] z punktami za dobre ruchy
- [ ] Countermove heuristic
- [ ] SEE (Static Exchange Evaluation) dla sortowania bić

## Quiescence Search Improvements

> **Cel:** Lepsze wykrywanie taktyki w quiescence search

### ✅ Obecna implementacja
- [x] Przeszukiwanie tylko bić do "spokojnej" pozycji
- [x] Stand-pat evaluation

### Ulepszenia (priorytetyzowane)
- [ ] ~~**Mate threat detection**~~ (wyłączone - 22x overhead, patrz Search Extensions)
  - Implementacja w quiescence była zbyt kosztowna (hasMateInOne w każdym węźle)
  - Alternatywa: Mate Threat Extensions w main search (patrz niżej)
- [x] **Check evasion** - kontynuuj gdy w szachu (nie kończ quiescence)
- [ ] **Delta pruning** - obcinaj bicia które nie mogą poprawić alpha

## Search Extensions (rozszerzenia)

> **Cel:** Przedłużanie przeszukiwania w krytycznych sytuacjach

- [x] **Check Extensions** - +1 ply gdy pozycja jest w szachu
- [ ] **Single Reply Extensions** - +1 ply gdy jest tylko jeden legalny ruch
- [ ] **Recapture Extensions** - +1 ply przy odbiciu na tym samym polu
- [ ] **Passed Pawn Extensions** - +1 ply dla promocji pionów przechodzących
- [ ] **Mate Threat Extensions** - +1 ply gdy przeciwnik grozi matem
  - **Rekomendowane** zamiast mate detection w quiescence (22x mniejszy overhead)
  - Sprawdź raz na węzeł w main search, nie w każdym węźle quiescence
  - Użyj prostej heurystyki: czy ostatni ruch dał szach lub zaatakował króla?

# Multi-Session Support (Iteracja 14b)

> **Cel:** Możliwość grania wielu partii równolegle w osobnych goroutynach

## ✅ Implementacja
- [x] Struktura `Session` z własnym TT i RNG
- [x] Przeniesienie globalnego `TT` do `Session`
- [x] `Search()` jako metoda na `Session`
- [x] UCI tworzy `Session` per gra

## Współdzielone (read-only, bezpieczne):
- PST tables (pawnPST, knightPST, ...)
- fileMasks, adjacentFileMasks
- pieceValues
- OpeningBook
- magic bitboards

## Per-session (izolowane):
- `TT *TranspositionTable`
- `bookRng *rand.Rand`

---

# Parallelizacja (Iteracja 15)

> **Cel:** Wykorzystanie wielu rdzeni CPU dla większej mocy obliczeniowej

## Wymagane wcześniej (blokery)
- [x] Multi-Session Support (Iteracja 14b)
- [ ] Iterative Deepening (dla Lazy SMP)

## Etapy implementacji

### Etap 1: Root-level parallelism (🟢 Łatwy)
- [ ] Każdy ruch z root position w osobnej goroutynie
- [ ] Kopiowanie `Position` dla każdej goroutyny
- [ ] Zbieranie wyników przez channel
- [ ] ~10-20% speedup

### Etap 2: Shared Transposition Table (🟡 Średni)
- [ ] `sync.RWMutex` dla TT lub lock-free z atomic
- [ ] Wątki współdzielą wyniki przeszukiwania
- [ ] Unikanie duplikacji pracy

### Etap 3: Lazy SMP (🟡 Średni)
- [ ] N wątków przeszukuje to samo drzewo równolegle
- [ ] Różne parametry (depth +/- 1) dla diversity
- [ ] Współdzielona TT synchronizuje wyniki
- [ ] ~50-70% speedup przy 4 wątkach

### Etap 4: YBWC / Young Brothers Wait Concept (🔴 Trudny)
- [ ] Pierwszy ruch sekwencyjnie, reszta równolegle
- [ ] Lepsza efektywność pruning w parallel
- [ ] Wymaga bardziej złożonej synchronizacji

---

# Przyszłość (poza obecnym scopem)

- [ ] Syzygy tablebases (końcówki)
- [ ] NNUE (ewaluacja siecią neuronową)
