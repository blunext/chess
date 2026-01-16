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
- [ ] Struktura pionów (zdwojone, izolowane, przechodzące)
- [ ] Aktywność figur (mobilność)
- [ ] Bezpieczeństwo króla

## Iteracja 12: Search
- [x] Minimax
- [ ] Alpha-beta pruning
- [ ] Iterative deepening
- [ ] Move ordering (bicia, killer moves)
- [ ] Quiescence search
- [ ] Zobrist hashing (wymagane dla TT)
- [ ] Transposition table

## Iteracja 13: Time Management
## Iteracja 13: Time Management
- [ ] Podstawowy time control w UCI (parsowanie wtime/btime)
- [ ] Iterative Deepening (pogłębianie przeszukiwania: 1, 2, 3...)
- [ ] Przerwanie search gdy czas się kończy (timeout check)
- [ ] Alokacja czasu (prosta heurystyka: czas/40 lub czas/20)

## Iteracja 14: UCI Options
- [ ] Obsługa `setoption name X value Y`
- [ ] Opcja `Hash` (rozmiar transposition table w MB)
- [ ] Opcja `Threads` (liczba wątków - przygotowanie pod multi-threading)
- [ ] Opcja `Move Overhead` (margines czasowy)
- [ ] Opcja `UCI_ShowWDL` (pokazywanie Win/Draw/Loss)
- [ ] Ponder (`go ponder`, `ponderhit`, `stop`)

---

# Dalsze optymalizacje (po podstawowej wersji)

## Search pruning
- [ ] Null move pruning
- [ ] Late Move Reductions (LMR)
- [ ] Aspiration windows
- [ ] Principal Variation Search (PVS)
- [ ] Futility pruning

## Move ordering zaawansowane
- [ ] History heuristic
- [ ] Countermove heuristic
- [ ] SEE (Static Exchange Evaluation) dla sortowania bić

---

# Przyszłość (poza obecnym scopem)

- [ ] Opening book
- [ ] Syzygy tablebases (końcówki)
- [ ] NNUE (ewaluacja siecią neuronową)
- [ ] Multi-threaded search (Lazy SMP)
- [ ] Pondering (myślenie w czasie przeciwnika)
