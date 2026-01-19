package engine

import (
	"testing"

	"chess/board"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate_InitialPosition(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	eval := Evaluate(pos)
	assert.Equal(t, 0, eval, "Initial position should be equal material")
}

func TestEvaluate_WhiteMissingPawn(t *testing.T) {
	// Initial position but white is missing e2 pawn
	// Material diff = -100, PST may vary slightly due to missing central pawn
	pos := board.CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.InDelta(t, -PawnValue, eval, 50, "White missing pawn should be around -100")
}

func TestEvaluate_BlackMissingPawn(t *testing.T) {
	// Initial position but black is missing e7 pawn
	pos := board.CreatePositionFormFEN("rnbqkbnr/pppp1ppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.InDelta(t, PawnValue, eval, 50, "Black missing pawn should be around +100")
}

func TestEvaluate_WhiteUpKnight(t *testing.T) {
	// White has extra knight (black missing g8 knight)
	pos := board.CreatePositionFormFEN("rnbqkb1r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.InDelta(t, KnightValue, eval, 50, "White up a knight should be around +320")
}

func TestEvaluate_WhiteUpRook(t *testing.T) {
	// Black missing a8 rook
	pos := board.CreatePositionFormFEN("1nbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQk - 0 1")
	eval := Evaluate(pos)
	assert.InDelta(t, RookValue, eval, 50, "White up a rook should be around +500")
}

func TestEvaluate_WhiteUpQueen(t *testing.T) {
	// Black missing queen
	pos := board.CreatePositionFormFEN("rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	// PeSTO queen value is ~1025 (middlegame), plus mobility and position adjustments
	assert.InDelta(t, 1025, eval, 150, "White up a queen should be around PeSTO queen value (~1025)")
}

func TestEvaluate_ComplexPosition(t *testing.T) {
	// FEN: 4kb2/pppppp2/8/8/8/8/PPPPP3/RNBQK3
	// White: K, Q, R, B, N, 5P = 900 + 500 + 330 + 320 + 500 = 2550
	// Black: K, B, 6P = 330 + 600 = 930
	// Material diff = 2550 - 930 = 1620, PST will adjust this
	pos := board.CreatePositionFormFEN("4kb2/pppppp2/8/8/8/8/PPPPP3/RNBQK3 w Q - 0 1")
	eval := Evaluate(pos)

	whiteExpected := QueenValue + RookValue + BishopValue + KnightValue + 5*PawnValue
	blackExpected := BishopValue + 6*PawnValue
	expected := whiteExpected - blackExpected

	assert.InDelta(t, expected, eval, 100, "Complex position should be around material diff")
}

func TestEvaluate_OnlyKings(t *testing.T) {
	pos := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/4K3 w - - 0 1")
	eval := Evaluate(pos)
	assert.Equal(t, 0, eval, "Just kings should be 0")
}

func TestPopCount(t *testing.T) {
	assert.Equal(t, 0, popCount(0))
	assert.Equal(t, 1, popCount(1))
	assert.Equal(t, 1, popCount(0x8000000000000000))
	assert.Equal(t, 8, popCount(0xFF))
	assert.Equal(t, 64, popCount(0xFFFFFFFFFFFFFFFF))
}

// PST-specific tests

func TestPST_KnightCenterBonus(t *testing.T) {
	// Knight on e4 (central) vs knight on a1 (corner)
	// e4 = index 28, a1 = index 0
	centerKnight := board.CreatePositionFormFEN("4k3/8/8/8/4N3/8/8/4K3 w - - 0 1")
	cornerKnight := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/N3K3 w - - 0 1")

	centerEval := Evaluate(centerKnight)
	cornerEval := Evaluate(cornerKnight)

	// Central knight should be worth more than corner knight
	assert.Greater(t, centerEval, cornerEval, "Knight on e4 should score higher than knight on a1")
	// Difference should be significant (central bonus vs edge penalty)
	assert.Greater(t, centerEval-cornerEval, 50, "PST difference should be significant")
}

func TestPST_PawnAdvancement(t *testing.T) {
	// Pawn on e6 (advanced) vs pawn on e2 (starting)
	advancedPawn := board.CreatePositionFormFEN("4k3/8/4P3/8/8/8/8/4K3 w - - 0 1")
	startingPawn := board.CreatePositionFormFEN("4k3/8/8/8/8/8/4P3/4K3 w - - 0 1")

	advancedEval := Evaluate(advancedPawn)
	startingEval := Evaluate(startingPawn)

	// Advanced pawn should be worth more
	assert.Greater(t, advancedEval, startingEval, "Advanced pawn should score higher")
}

func TestPST_RookOnSeventhRank(t *testing.T) {
	// Rook on 7th rank vs rook on 1st rank
	seventhRank := board.CreatePositionFormFEN("4k3/4R3/8/8/8/8/8/4K3 w - - 0 1")
	firstRank := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/R3K3 w - - 0 1")

	seventhEval := Evaluate(seventhRank)
	firstEval := Evaluate(firstRank)

	// Rook on 7th should be worth more
	assert.Greater(t, seventhEval, firstEval, "Rook on 7th rank should score higher")
}

func TestPST_KingEndgame(t *testing.T) {
	// In endgame (K+R only), central king is better (PeSTO endgame tables)
	castledKing := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/5RK1 w - - 0 1")
	centerKing := board.CreatePositionFormFEN("8/4k3/8/8/4K3/8/8/5R2 w - - 0 1")

	castledEval := Evaluate(castledKing)
	centerEval := Evaluate(centerKing)

	// In endgame with just K+R, central king should be preferred (PeSTO egKingTable)
	// Note: This is correct behavior - in endgames you want an active king!
	assert.Greater(t, centerEval, castledEval, "In endgame, central king should score higher")
}

func TestPST_SymmetricPosition(t *testing.T) {
	// Initial position should be exactly 0 due to symmetry
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	eval := Evaluate(pos)
	assert.Equal(t, 0, eval, "Symmetric initial position should evaluate to 0")
}

func TestPST_BlackMirroring(t *testing.T) {
	// White knight on e4 should have same PST value as black knight on e5
	// (mirrored position - both are centrally placed for their side)
	whiteKnightE4 := board.CreatePositionFormFEN("4k3/8/8/8/4N3/8/8/4K3 w - - 0 1")
	blackKnightE5 := board.CreatePositionFormFEN("4k3/8/8/4n3/8/8/8/4K3 w - - 0 1")

	whiteEval := Evaluate(whiteKnightE4)
	blackEval := Evaluate(blackKnightE5)

	// Should be roughly opposite (white positive, black negative, similar magnitude)
	// e4 for white = e5 for black (both are rank 4 from their perspective)
	assert.InDelta(t, whiteEval, -blackEval, 20, "Mirrored knight positions should have opposite scores")
}

// === King Safety Tests ===

func TestKingSafety_CastledKingWithPawnShield(t *testing.T) {
	// White king castled kingside with full pawn shield (f2,g2,h2)
	pos := board.CreatePositionFormFEN("r1bq1rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 w - - 0 1")

	whiteKS := kingSafety(pos, true)
	blackKS := kingSafety(pos, false)

	// Both have castled with intact pawn shields - should have similar scores
	assert.InDelta(t, whiteKS, blackKS, 30, "Both castled kings with pawn shields should have similar safety")
}

func TestKingSafety_BrokenPawnShield(t *testing.T) {
	// White king castled but g2 pawn is missing (g-file open)
	brokenShield := board.CreatePositionFormFEN("r1bq1rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/5NP1/PPPP1P1P/RNBQ1RK1 w - - 0 1")

	// Compare with intact shield
	intactShield := board.CreatePositionFormFEN("r1bq1rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 w - - 0 1")

	brokenKS := kingSafety(brokenShield, true)
	intactKS := kingSafety(intactShield, true)

	// Broken shield should be worse (pawn on g3 instead of g2)
	// Note: g3 is advanced, f2/h2 still intact
	assert.LessOrEqual(t, brokenKS, intactKS, "Broken pawn shield should give lower score")
}

func TestKingSafety_UncastledKingPenalty(t *testing.T) {
	// White king still in center (e1)
	uncastled := board.CreatePositionFormFEN("r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1")

	// Castled position
	castled := board.CreatePositionFormFEN("r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 w kq - 0 1")

	uncastledKS := kingSafety(uncastled, true)
	castledKS := kingSafety(castled, true)

	// Uncastled king should have worse safety
	assert.Less(t, uncastledKS, castledKS, "Uncastled king should have worse safety than castled")
}

func TestKingSafety_OpenFileNearKing(t *testing.T) {
	// King castled but h-file is completely open (no pawns)
	openFile := board.CreatePositionFormFEN("r1bq1rk1/pppp1pp1/2n2n2/2b1p2p/2B1P3/5N2/PPPP1PP1/RNBQ1RK1 w - - 0 1")

	ks := kingSafety(openFile, true)

	// Should have penalty for open file
	// Exact value depends on other factors, but should be negative overall
	assert.Less(t, ks, 0, "Open file near king should result in negative safety score")
}

func TestKingSafety_EndgameScaling(t *testing.T) {
	// Position without queens - king safety should be reduced
	noQueens := board.CreatePositionFormFEN("r1b2rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNB2RK1 w - - 0 1")

	// Position with queens
	withQueens := board.CreatePositionFormFEN("r1bq1rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 w - - 0 1")

	noQueensKS := kingSafety(noQueens, true)
	withQueensKS := kingSafety(withQueens, true)

	// Without queens, king safety should matter less
	// The absolute value should be smaller
	absNoQueens := noQueensKS
	if absNoQueens < 0 {
		absNoQueens = -absNoQueens
	}
	absWithQueens := withQueensKS
	if absWithQueens < 0 {
		absWithQueens = -absWithQueens
	}

	// Without enemy queen, king safety penalty/bonus is divided by 4
	assert.LessOrEqual(t, absNoQueens, absWithQueens, "King safety should matter less without queens")
}

func TestKingSafety_PawnShieldFunction(t *testing.T) {
	// Test pawnShield directly
	// King on g1 with perfect shield (f2,g2,h2)
	perfectShield := board.CreatePositionFormFEN("8/8/8/8/8/8/5PPP/6K1 w - - 0 1")
	kingSq := 6 // g1

	score := pawnShield(perfectShield, kingSq, true)

	// Should have bonus for all three pawns
	assert.Greater(t, score, 0, "Perfect pawn shield should give positive score")
}

func TestKingSafety_Queensidecastle(t *testing.T) {
	// White king castled queenside (c1) with pawn shield (a2,b2,c2)
	queenside := board.CreatePositionFormFEN("r3kbnr/ppp2ppp/2nqb3/3pp3/3PP3/2NQB3/PPP2PPP/2KR1BNR w kq - 0 1")

	ks := kingSafety(queenside, true)

	// Should still evaluate pawn shield for queenside castling
	// With intact pawns on a2,b2,c2
	assert.GreaterOrEqual(t, ks, -50, "Queenside castled king with shield should not be heavily penalized")
}

// === Pawn Structure Tests ===

func TestPawnStructure_DoubledPawns(t *testing.T) {
	// White has doubled pawns on e-file (e2 and e4)
	doubled := board.CreatePositionFormFEN("4k3/8/8/8/4P3/8/4P3/4K3 w - - 0 1")

	penalty := doubledPawns(doubled.Pawns & doubled.White)

	// Should have penalty for one doubled pawn
	assert.Equal(t, -DoubledPawnPenalty, penalty, "Doubled pawns should have penalty")
}

func TestPawnStructure_TriplePawns(t *testing.T) {
	// White has tripled pawns on e-file (e2, e3, e4)
	tripled := board.CreatePositionFormFEN("4k3/8/8/8/4P3/4P3/4P3/4K3 w - - 0 1")

	penalty := doubledPawns(tripled.Pawns & tripled.White)

	// Should have penalty for two extra pawns
	assert.Equal(t, -2*DoubledPawnPenalty, penalty, "Tripled pawns should have 2x penalty")
}

func TestPawnStructure_NoDoubledPawns(t *testing.T) {
	// White has pawns on different files
	noDoubled := board.CreatePositionFormFEN("4k3/8/8/8/3PP3/8/8/4K3 w - - 0 1")

	penalty := doubledPawns(noDoubled.Pawns & noDoubled.White)

	assert.Equal(t, 0, penalty, "No doubled pawns should have no penalty")
}

func TestPawnStructure_IsolatedPawn(t *testing.T) {
	// White has isolated pawn on a-file (no pawns on b-file)
	isolated := board.CreatePositionFormFEN("4k3/8/8/8/P7/8/4P3/4K3 w - - 0 1")

	penalty := isolatedPawns(isolated.Pawns & isolated.White)

	// a-pawn is isolated (no pawn on b-file), e-pawn is also isolated (no pawns on d or f files)
	assert.Equal(t, -2*IsolatedPawnPenalty, penalty, "Isolated pawns should have penalty")
}

func TestPawnStructure_ConnectedPawns(t *testing.T) {
	// White has connected pawns on d and e files
	connected := board.CreatePositionFormFEN("4k3/8/8/8/3PP3/8/8/4K3 w - - 0 1")

	penalty := isolatedPawns(connected.Pawns & connected.White)

	assert.Equal(t, 0, penalty, "Connected pawns should have no isolated penalty")
}

func TestPawnStructure_PassedPawn(t *testing.T) {
	// White has passed pawn on e5 (no black pawns on d,e,f files ahead)
	passed := board.CreatePositionFormFEN("4k3/8/8/4P3/8/8/8/4K3 w - - 0 1")

	bonus := passedPawns(passed.Pawns&passed.White, passed.Pawns&passed.Black, true)

	// Should have bonus for passed pawn (base + rank bonus)
	assert.Greater(t, bonus, 0, "Passed pawn should have bonus")
	// e5 is rank 4 (0-indexed), so bonus = 20 + 4*10 = 60
	expectedBonus := PassedPawnBonus + 4*PassedPawnRankBonus
	assert.Equal(t, expectedBonus, bonus, "Passed pawn bonus should match expected")
}

func TestPawnStructure_BlockedPawn(t *testing.T) {
	// White pawn on e4 is blocked by black pawn on e5
	blocked := board.CreatePositionFormFEN("4k3/8/8/4p3/4P3/8/8/4K3 w - - 0 1")

	bonus := passedPawns(blocked.Pawns&blocked.White, blocked.Pawns&blocked.Black, true)

	assert.Equal(t, 0, bonus, "Blocked pawn should not be passed")
}

func TestPawnStructure_PassedPawnAdvanced(t *testing.T) {
	// White has passed pawn on e7 (very advanced!)
	advanced := board.CreatePositionFormFEN("4k3/4P3/8/8/8/8/8/4K3 w - - 0 1")

	bonus := passedPawns(advanced.Pawns&advanced.White, advanced.Pawns&advanced.Black, true)

	// e7 is rank 6, so bonus = 20 + 6*10 = 80
	expectedBonus := PassedPawnBonus + 6*PassedPawnRankBonus
	assert.Equal(t, expectedBonus, bonus, "Advanced passed pawn should have large bonus")
}

func TestPawnStructure_BlackPassedPawn(t *testing.T) {
	// Black has passed pawn on e2 (very advanced for black!)
	blackPassed := board.CreatePositionFormFEN("4k3/8/8/8/8/8/4p3/4K3 w - - 0 1")

	bonus := passedPawns(blackPassed.Pawns&blackPassed.Black, blackPassed.Pawns&blackPassed.White, false)

	// e2 is rank 1, for black bonus = 20 + (7-1)*10 = 80
	expectedBonus := PassedPawnBonus + 6*PassedPawnRankBonus
	assert.Equal(t, expectedBonus, bonus, "Black advanced passed pawn should have large bonus")
}

func TestPawnStructure_PawnStructureFunction(t *testing.T) {
	// Test the main pawnStructure function combining all factors
	// White has: doubled pawns on e-file, but one is passed
	pos := board.CreatePositionFormFEN("4k3/8/8/4P3/8/8/4P3/4K3 w - - 0 1")

	whitePawnScore := pawnStructure(pos, true)
	blackPawnScore := pawnStructure(pos, false)

	// White should have net positive (passed pawn bonus > doubled penalty)
	// But both pawns are isolated (no pawns on d or f files)
	// So: doubled (-20) + isolated (-15*2 = -30) + passed (20 + 40 for e5) = +10
	assert.Less(t, blackPawnScore, whitePawnScore, "White with passed pawn should score better")
}

// === Space Bonus Tests ===

func TestSpaceBonus_CentralPawns(t *testing.T) {
	// White has pawn on e4 (central square)
	centralPawn := board.CreatePositionFormFEN("4k3/8/8/8/4P3/8/8/4K3 w - - 0 1")

	// White has pawn on a4 (edge square)
	edgePawn := board.CreatePositionFormFEN("4k3/8/8/8/P7/8/8/4K3 w - - 0 1")

	centralScore := spaceBonus(centralPawn, true)
	edgeScore := spaceBonus(edgePawn, true)

	// Central pawn should score higher (central bonus + rank 4 advancement)
	assert.Greater(t, centralScore, edgeScore, "Central pawn should score higher than edge pawn")
	// e4 gets CentralPawnBonus (25) + AdvancedPawnRank4 (5) = 30
	// a4 gets only AdvancedPawnRank4 (5)
	assert.Equal(t, CentralPawnBonus+AdvancedPawnRank4, centralScore, "e4 pawn should get central + advancement bonus")
	assert.Equal(t, AdvancedPawnRank4, edgeScore, "a4 pawn should get only advancement bonus")
}

func TestSpaceBonus_ExtendedCenter(t *testing.T) {
	// White has pawn on c4 (extended center)
	c4Pawn := board.CreatePositionFormFEN("4k3/8/8/8/2P5/8/8/4K3 w - - 0 1")

	score := spaceBonus(c4Pawn, true)

	// c4 gets ExtendedCenterBonus (15) + AdvancedPawnRank4 (5) = 20
	assert.Equal(t, ExtendedCenterBonus+AdvancedPawnRank4, score, "c4 pawn should get extended center + advancement bonus")
}

func TestSpaceBonus_D4E4Center(t *testing.T) {
	// White has ideal pawn center d4 + e4
	idealCenter := board.CreatePositionFormFEN("4k3/8/8/8/3PP3/8/8/4K3 w - - 0 1")

	score := spaceBonus(idealCenter, true)

	// d4 + e4: 2 * (CentralPawnBonus + AdvancedPawnRank4) = 2 * (25 + 5) = 60
	expectedScore := 2 * (CentralPawnBonus + AdvancedPawnRank4)
	assert.Equal(t, expectedScore, score, "d4+e4 pawns should get double central bonus")
}

func TestSpaceBonus_AdvancedPawnRanks(t *testing.T) {
	// Pawn on rank 4 (e4)
	rank4 := board.CreatePositionFormFEN("4k3/8/8/8/4P3/8/8/4K3 w - - 0 1")
	// Pawn on rank 5 (e5)
	rank5 := board.CreatePositionFormFEN("4k3/8/8/4P3/8/8/8/4K3 w - - 0 1")
	// Pawn on rank 6 (e6) - NOT in central squares anymore
	rank6 := board.CreatePositionFormFEN("4k3/8/4P3/8/8/8/8/4K3 w - - 0 1")

	score4 := spaceBonus(rank4, true)
	score5 := spaceBonus(rank5, true)
	score6 := spaceBonus(rank6, true)

	// Verify exact values:
	// e4: central (25) + rank4 (5) = 30
	// e5: central (25) + rank5 (10) = 35
	// e6: no central bonus + rank6 (15) = 15 (loses central bonus!)
	assert.Equal(t, CentralPawnBonus+AdvancedPawnRank4, score4, "e4 should get central + rank4 bonus")
	assert.Equal(t, CentralPawnBonus+AdvancedPawnRank5, score5, "e5 should get central + rank5 bonus")
	assert.Equal(t, AdvancedPawnRank6, score6, "e6 is not in central squares, only rank6 bonus")

	// e5 > e4 (both central, e5 more advanced)
	assert.Greater(t, score5, score4, "Rank 5 central pawn should score higher than rank 4")
	// e5 > e6 because e5 keeps central bonus, e6 loses it
	assert.Greater(t, score5, score6, "Central e5 pawn with bonus should beat non-central e6")
}

func TestSpaceBonus_BlackPawns(t *testing.T) {
	// Black pawn on d5 (central from black's perspective)
	blackCenter := board.CreatePositionFormFEN("4k3/8/8/3p4/8/8/8/4K3 w - - 0 1")

	score := spaceBonus(blackCenter, false)

	// d5 for black: central (25) + rank 4 from black's view (5) = 30
	assert.Equal(t, CentralPawnBonus+AdvancedPawnRank4, score, "Black d5 pawn should get central + advancement bonus")
}

func TestSpaceBonus_SymmetricPosition(t *testing.T) {
	// Symmetric position with central pawns
	pos := board.CreatePositionFormFEN("4k3/3pp3/8/8/8/8/3PP3/4K3 w - - 0 1")

	whiteScore := spaceBonus(pos, true)
	blackScore := spaceBonus(pos, false)

	// Symmetric position should have equal space scores
	assert.Equal(t, whiteScore, blackScore, "Symmetric position should have equal space bonus")
}

func TestSpaceBonus_ComplexPosition(t *testing.T) {
	// Position after 1.e4 e5 2.d4
	pos := board.CreatePositionFormFEN("rnbqkbnr/pppp1ppp/8/4p3/3PP3/8/PPP2PPP/RNBQKBNR b KQkq - 0 1")

	whiteScore := spaceBonus(pos, true)
	blackScore := spaceBonus(pos, false)

	// White has d4 + e4 (both central), black has e5 (central)
	// White: 2 pawns on central + rank4 = 2 * (25 + 5) = 60
	// Black: e5 (central) + rank 4 from black's view = 25 + 5 = 30
	assert.Greater(t, whiteScore, blackScore, "White with d4+e4 should have more space than black with e5")
}

// === Mobility Tests ===

func TestMobility_BlockedBishop(t *testing.T) {
	// Bishop blocked by own pawns vs bishop with open diagonals
	blockedBishop := board.CreatePositionFormFEN("4k3/8/8/8/8/8/PPP5/B3K3 w - - 0 1")
	openBishop := board.CreatePositionFormFEN("4k3/8/8/8/8/8/5PPP/4KB2 w - - 0 1")

	blockedMob := mobility(blockedBishop, true)
	openMob := mobility(openBishop, true)

	// Open bishop should have better mobility score
	assert.Greater(t, openMob, blockedMob, "Open bishop should have higher mobility than blocked bishop")
}

func TestMobility_KnightInCorner(t *testing.T) {
	// Knight in corner (2 moves) vs knight in center (8 moves)
	cornerKnight := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/N3K3 w - - 0 1")
	centerKnight := board.CreatePositionFormFEN("4k3/8/8/8/4N3/8/8/4K3 w - - 0 1")

	cornerMob := mobility(cornerKnight, true)
	centerMob := mobility(centerKnight, true)

	// Center knight should have better mobility
	assert.Greater(t, centerMob, cornerMob, "Central knight should have higher mobility than corner knight")
}

func TestMobility_InitialPosition(t *testing.T) {
	// At start, white has limited piece mobility (only knights can move)
	pos := board.CreatePositionFormFEN(board.InitialPosition)

	whiteMob := mobility(pos, true)
	blackMob := mobility(pos, false)

	// Both sides should have similar (negative) mobility due to blocked pieces
	assert.InDelta(t, whiteMob, blackMob, 10, "Initial position should have roughly equal mobility")
	// Both should be negative since pieces have fewer moves than base
	assert.Less(t, whiteMob, 0, "Initial position should have negative mobility (blocked pieces)")
}

func TestMobility_RookOpenFile(t *testing.T) {
	// Rook on open file (many moves) vs rook blocked by pawns
	openRook := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/R3K3 w - - 0 1")
	blockedRook := board.CreatePositionFormFEN("4k3/8/8/8/8/8/P7/R3K3 w - - 0 1")

	openMob := mobility(openRook, true)
	blockedMob := mobility(blockedRook, true)

	// Open rook should have more mobility
	assert.Greater(t, openMob, blockedMob, "Rook on open file should have higher mobility")
}

func TestMobility_QueenDeveloped(t *testing.T) {
	// Developed queen (center of board) vs queen on starting square
	developedQueen := board.CreatePositionFormFEN("4k3/8/8/8/3Q4/8/8/4K3 w - - 0 1")
	startingQueen := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/3QK3 w - - 0 1")

	developedMob := mobility(developedQueen, true)
	startingMob := mobility(startingQueen, true)

	// Developed queen should have more mobility
	assert.Greater(t, developedMob, startingMob, "Developed queen should have higher mobility")
}

// Benchmark for mobility calculation overhead
func BenchmarkEvaluate_WithMobility(b *testing.B) {
	pos := board.CreatePositionFormFEN("r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Evaluate(pos)
	}
}

func BenchmarkEvaluatePeSTO_WithoutMobility(b *testing.B) {
	pos := board.CreatePositionFormFEN("r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvaluatePeSTO(pos)
	}
}

// ============================================================================
// Development Bonus Tests
// ============================================================================

// TestDevelopment_InitialPosition tests that initial position has penalties for undeveloped pieces
func TestDevelopment_InitialPosition(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)

	whiteBonus := developmentBonus(pos, true)
	blackBonus := developmentBonus(pos, false)

	// Both sides have all minors undeveloped: 2 knights (-30) + 2 bishops (-20) = -50
	// But they also have castling rights: +10 (kingside) + 5 (queenside) = +15
	// So net should be around -35 per side
	assert.Less(t, whiteBonus, 0, "White should have development penalty in initial position")
	assert.Less(t, blackBonus, 0, "Black should have development penalty in initial position")
	assert.Equal(t, whiteBonus, blackBonus, "Both sides should have equal development penalty")

	t.Logf("Development bonus in initial position: White=%d, Black=%d", whiteBonus, blackBonus)
}

// TestDevelopment_DevelopedPieces tests that developed pieces get no penalty
func TestDevelopment_DevelopedPieces(t *testing.T) {
	// Position after 1.e4 e5 2.Nf3 Nc6 3.Bc4 Bc5 - Italian Game opening
	// White: Nf3 developed, Bc4 developed, Nb1 and Bc1 still home
	// Black: Nc6 developed, Bc5 developed, Ng8 and Bc8 still home
	pos := board.CreatePositionFormFEN("r1bqk1nr/pppp1ppp/2n5/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1")

	whiteBonus := developmentBonus(pos, true)
	blackBonus := developmentBonus(pos, false)

	// White: 1 knight undeveloped (-15), 1 bishop undeveloped (-10) = -25 + castling rights
	// Black: 1 knight undeveloped (-15), 1 bishop undeveloped (-10) = -25 + castling rights
	t.Logf("Italian Game development: White=%d, Black=%d", whiteBonus, blackBonus)

	// Position after both sides fully developed (no minors on starting squares)
	// 1.e4 e5 2.Nf3 Nc6 3.Bc4 Bc5 4.Nc3 Nf6
	pos2 := board.CreatePositionFormFEN("r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/2N2N2/PPPP1PPP/R1BQK2R w KQkq - 0 1")

	whiteBonus2 := developmentBonus(pos2, true)
	blackBonus2 := developmentBonus(pos2, false)

	// Now only Bc1 is undeveloped for white, Bc8 for black
	assert.Greater(t, whiteBonus2, whiteBonus, "More development should give better score")
	t.Logf("Four Knights development: White=%d, Black=%d", whiteBonus2, blackBonus2)
}

// TestDevelopment_EarlyQueen tests penalty for early queen development
func TestDevelopment_EarlyQueen(t *testing.T) {
	// Position where white plays Qa4 early (like in our problematic game)
	// Knights still on b1 and g1
	pos := board.CreatePositionFormFEN("r1bqkbnr/pppp1ppp/2n5/4p3/Q3P3/8/PPPP1PPP/RNB1KBNR w KQkq - 0 1")

	whiteBonus := developmentBonus(pos, true)
	blackBonus := developmentBonus(pos, false)

	// White: Queen moved with both knights undeveloped -> -20 extra penalty
	// White: 2 knights undeveloped (-30), 2 bishops undeveloped (-20), early queen (-20)
	// Black: 1 knight developed (Nc6), 1 knight undeveloped (-15), 2 bishops undeveloped (-20)
	assert.Less(t, whiteBonus, blackBonus, "White should have worse development (early queen)")

	t.Logf("Early queen development: White=%d, Black=%d", whiteBonus, blackBonus)
}

// TestDevelopment_BlockedCenterPawn tests penalty for blocked center pawns
func TestDevelopment_BlockedCenterPawn(t *testing.T) {
	// Position where d2 pawn is blocked by knight on d3
	pos := board.CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/3N4/PPPPPPPP/R1BQKBNR w KQkq - 0 1")

	whiteBonus := developmentBonus(pos, true)
	blackBonus := developmentBonus(pos, false)

	// White has knight on d3 blocking d2 pawn -> -15 penalty
	t.Logf("Blocked d2 pawn: White=%d, Black=%d", whiteBonus, blackBonus)
}

// TestDevelopment_CastlingRights tests bonus for maintaining castling rights
func TestDevelopment_CastlingRights(t *testing.T) {
	// Position with full castling rights
	pos1 := board.CreatePositionFormFEN("r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 0 1")
	whiteBonus1 := developmentBonus(pos1, true)

	// Same position but white lost kingside castling rights
	pos2 := board.CreatePositionFormFEN("r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w Qkq - 0 1")
	whiteBonus2 := developmentBonus(pos2, true)

	// Same position but white lost all castling rights
	pos3 := board.CreatePositionFormFEN("r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w kq - 0 1")
	whiteBonus3 := developmentBonus(pos3, true)

	assert.Greater(t, whiteBonus1, whiteBonus2, "Full castling rights should be better than partial")
	assert.Greater(t, whiteBonus2, whiteBonus3, "Partial castling rights should be better than none")

	t.Logf("Castling rights: Full=%d, KingsideOnly=%d, None=%d", whiteBonus1, whiteBonus2, whiteBonus3)
}

// TestDevelopment_Endgame tests that development bonus is disabled in endgame
func TestDevelopment_Endgame(t *testing.T) {
	// Endgame position with only kings and pawns
	pos := board.CreatePositionFormFEN("4k3/pppppppp/8/8/8/8/PPPPPPPP/4K3 w - - 0 1")

	whiteBonus := developmentBonus(pos, true)
	blackBonus := developmentBonus(pos, false)

	// Game phase is 0 (no pieces), so development bonus should be 0
	assert.Equal(t, 0, whiteBonus, "Development bonus should be 0 in endgame")
	assert.Equal(t, 0, blackBonus, "Development bonus should be 0 in endgame")
}

// TestDevelopment_GamePhaseScaling tests that development bonus scales with game phase
func TestDevelopment_GamePhaseScaling(t *testing.T) {
	// Full middlegame (all pieces)
	pos1 := board.CreatePositionFormFEN(board.InitialPosition)
	phase1 := calculateGamePhase(pos1)

	// Late middlegame (queens traded)
	pos2 := board.CreatePositionFormFEN("rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNB1KBNR w KQkq - 0 1")
	phase2 := calculateGamePhase(pos2)

	// Early endgame (only rooks and minor pieces)
	pos3 := board.CreatePositionFormFEN("r1b1kb1r/pppppppp/8/8/8/8/PPPPPPPP/R1B1KB1R w KQkq - 0 1")
	phase3 := calculateGamePhase(pos3)

	t.Logf("Game phases: Full=%d, QueensTraded=%d, MinorPiecesTraded=%d", phase1, phase2, phase3)

	assert.Equal(t, 24, phase1, "Initial position should have phase 24")
	assert.Equal(t, 16, phase2, "Queens traded should have phase 16")
}
