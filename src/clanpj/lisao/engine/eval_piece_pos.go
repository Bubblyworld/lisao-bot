// Evaluation using piece values and fixed piece position values

package engine

import (
	// "fmt"
	"math/bits"

	dragon "github.com/Bubblyworld/dragontoothmg"
)

// Piece values
const nothingVal = EvalCp(0)
var pawnVal = EvalCp(100)
var knightVal = EvalCp(300)
var bishopVal = EvalCp(300)
var rookVal = EvalCp(500)
var queenVal = EvalCp(900)
const kingVal = EvalCp(0)

// Register the piece values as eval config params
func init() {
	//RegisterConfigParamEvalCpDefault("pawn-piece-val", &pawnVal) since this is supposed to be a centi-pawn metric, we leave pawn eval pinned at 100
	RegisterConfigParamEvalCpDefault("knight-piece-val", &knightVal)
	RegisterConfigParamEvalCpDefault("bishop-piece-val", &bishopVal)
	RegisterConfigParamEvalCpDefault("rook-piece-val", &rookVal)
	RegisterConfigParamEvalCpDefault("queen-piece-val", &queenVal)
}

var pieceVals = [7]EvalCp{
	nothingVal,
	pawnVal,
	knightVal,
	bishopVal,
	rookVal,
	queenVal,
	kingVal}

var nothingPosVals = [64]int8{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0}

// Stolen from SunFish (tables inverted to reflect dragon pos ordering)
var whitePawnPosVals = [64]int8{
	0, 0, 0, 0, 0, 0, 0, 0,
	-31, 8, -7, -37, -36, -14, 3, -31,
	-22, 9, 5, -11, -10, -2, 3, -19,
	-26, 3, 10, 9, 6, 1, 0, -23,
	-17, 16, -2, 15, 14, 0, 15, -13,
	7, 29, 21, 44, 40, 31, 44, 7,
	78, 83, 86, 73, 102, 82, 85, 90,
	0, 0, 0, 0, 0, 0, 0, 0}

func init() { RegisterConfigParamInt8ArrayDefault("white-pawn-pos-val", whitePawnPosVals[:]) }

var whiteKnightPosVals = [64]int8{
	-74, -23, -26, -24, -19, -35, -22, -69,
	-23, -15, 2, 0, 2, 0, -23, -20,
	-18, 10, 13, 22, 18, 15, 11, -14,
	-1, 5, 31, 21, 22, 35, 2, 0,
	24, 24, 45, 37, 33, 41, 25, 17,
	10, 67, 1, 74, 73, 27, 62, -2,
	-3, -6, 100, -36, 4, 62, -4, -14,
	-66, -53, -75, -75, -10, -55, -58, -70}

func init() { RegisterConfigParamInt8ArrayDefault("white-knight-pos-val", whiteKnightPosVals[:]) }

var whiteBishopPosVals = [64]int8{
	-7, 2, -15, -12, -14, -15, -10, -10,
	19, 20, 11, 6, 7, 6, 20, 16,
	14, 25, 24, 15, 8, 25, 20, 15,
	13, 10, 17, 23, 17, 16, 0, 7,
	25, 17, 20, 34, 26, 25, 15, 10,
	-9, 39, -32, 41, 52, -10, 28, -14,
	-11, 20, 35, -42, -39, 31, 2, -22,
	-59, -78, -82, -76, -23, -107, -37, -50}

func init() { RegisterConfigParamInt8ArrayDefault("white-bishop-pos-val", whiteBishopPosVals[:]) }

var whiteRookPosVals = [64]int8{
	-30, -24, -18, 5, -2, -18, -31, -32,
	-53, -38, -31, -26, -29, -43, -44, -53,
	-42, -28, -42, -25, -25, -35, -26, -46,
	-28, -35, -16, -21, -13, -29, -46, -30,
	0, 5, 16, 13, 18, -4, -9, -6,
	19, 35, 28, 33, 45, 27, 25, 15,
	55, 29, 56, 67, 55, 62, 34, 60,
	35, 29, 33, 4, 37, 33, 56, 50}

func init() { RegisterConfigParamInt8ArrayDefault("white-rook-pos-val", whiteRookPosVals[:]) }

var whiteQueenPosVals = [64]int8{
	-39, -30, -31, -13, -31, -36, -34, -42,
	-36, -18, 0, -19, -15, -15, -21, -38,
	-30, -6, -13, -11, -16, -11, -16, -27,
	-14, -15, -2, -5, -1, -10, -20, -22,
	1, -16, 22, 17, 25, 20, -13, -6,
	-2, 43, 32, 60, 72, 63, 43, 2,
	14, 32, 60, -10, 20, 76, 57, 24,
	6, 1, -8, -104, 69, 24, 88, 26}

func init() { RegisterConfigParamInt8ArrayDefault("white-queen-pos-val", whiteQueenPosVals[:]) }

var whiteKingPosVals = [64]int8{
	17, 30, -3, -14, 6, -1, 40, 18,
	-4, 3, -14, -50, -57, -18, 13, 4,
	-47, -42, -43, -79, -64, -32, -29, -32,
	-55, -43, -52, -28, -51, -47, -8, -50,
	-55, 50, 11, -4, -19, 13, 0, -49,
	-62, 12, -57, 44, -67, 28, 37, -31,
	-32, 10, 55, 56, 56, 55, 10, 3,
	4, 54, 47, -99, -99, 60, 83, -62}

func init() { RegisterConfigParamInt8ArrayDefault("white-king-pos-val", whiteKingPosVals[:]) }

// From - https://chessprogramming.wikispaces.com/Simplified+evaluation+function - (tables inverted to reflect dragon pos ordering)
// Added some deliberate jitter to add artificial discrimination into stuck end-games (like Fine 70)
var whiteKingEndgamePosVals = [64]int8{
	-50, -33, -30, -28, -28, -30, -33, -50,
	-32, -30, 4, 0, 0, 4, -30, -32,
	-30, -12, 18, 29, 29, 18, -12, -30,
	-29, -10, 31, 41, 41, 31, -10, -29,
	-29, -10, 30, 40, 40, 30, -10, -29,
	-30, -12, 17, 29, 29, 17, -12, -30,
	-32, -19, -9, 0, 0, -9, -19, -32,
	-50, -41, -30, -20, -20, -30, -41, -50}

func init() { RegisterConfigParamInt8ArrayDefault("white-king-eg-val", whiteKingEndgamePosVals[:]) }

var whitePiecePosVals = [7]*[64]int8{
	&nothingPosVals,
	&whitePawnPosVals,
	&whiteKnightPosVals,
	&whiteBishopPosVals,
	&whiteRookPosVals,
	&whiteQueenPosVals,
	&whiteKingPosVals}

// Stolen from SunFish
var blackPawnPosVals = [64]int8{
	0, 0, 0, 0, 0, 0, 0, 0,
	78, 83, 86, 73, 102, 82, 85, 90,
	7, 29, 21, 44, 40, 31, 44, 7,
	-17, 16, -2, 15, 14, 0, 15, -13,
	-26, 3, 10, 9, 6, 1, 0, -23,
	-22, 9, 5, -11, -10, -2, 3, -19,
	-31, 8, -7, -37, -36, -14, 3, -31,
	0, 0, 0, 0, 0, 0, 0, 0}

func init() { RegisterConfigParamInt8ArrayDefault("black-pawn-pos-val", blackPawnPosVals[:]) }

var blackKnightPosVals = [64]int8{
	-66, -53, -75, -75, -10, -55, -58, -70,
	-3, -6, 100, -36, 4, 62, -4, -14,
	10, 67, 1, 74, 73, 27, 62, -2,
	24, 24, 45, 37, 33, 41, 25, 17,
	-1, 5, 31, 21, 22, 35, 2, 0,
	-18, 10, 13, 22, 18, 15, 11, -14,
	-23, -15, 2, 0, 2, 0, -23, -20,
	-74, -23, -26, -24, -19, -35, -22, -69}

func init() { RegisterConfigParamInt8ArrayDefault("black-knight-pos-val", blackKnightPosVals[:]) }

var blackBishopPosVals = [64]int8{
	-59, -78, -82, -76, -23, -107, -37, -50,
	-11, 20, 35, -42, -39, 31, 2, -22,
	-9, 39, -32, 41, 52, -10, 28, -14,
	25, 17, 20, 34, 26, 25, 15, 10,
	13, 10, 17, 23, 17, 16, 0, 7,
	14, 25, 24, 15, 8, 25, 20, 15,
	19, 20, 11, 6, 7, 6, 20, 16,
	-7, 2, -15, -12, -14, -15, -10, -10}

func init() { RegisterConfigParamInt8ArrayDefault("black-bishop-pos-val", blackBishopPosVals[:]) }

var blackRookPosVals = [64]int8{
	35, 29, 33, 4, 37, 33, 56, 50,
	55, 29, 56, 67, 55, 62, 34, 60,
	19, 35, 28, 33, 45, 27, 25, 15,
	0, 5, 16, 13, 18, -4, -9, -6,
	-28, -35, -16, -21, -13, -29, -46, -30,
	-42, -28, -42, -25, -25, -35, -26, -46,
	-53, -38, -31, -26, -29, -43, -44, -53,
	-30, -24, -18, 5, -2, -18, -31, -32}

func init() { RegisterConfigParamInt8ArrayDefault("black-rook-pos-val", blackRookPosVals[:]) }

var blackQueenPosVals = [64]int8{
	6, 1, -8, -104, 69, 24, 88, 26,
	14, 32, 60, -10, 20, 76, 57, 24,
	-2, 43, 32, 60, 72, 63, 43, 2,
	1, -16, 22, 17, 25, 20, -13, -6,
	-14, -15, -2, -5, -1, -10, -20, -22,
	-30, -6, -13, -11, -16, -11, -16, -27,
	-36, -18, 0, -19, -15, -15, -21, -38,
	-39, -30, -31, -13, -31, -36, -34, -42}

func init() { RegisterConfigParamInt8ArrayDefault("black-queen-pos-val", blackQueenPosVals[:]) }

var blackKingPosVals = [64]int8{
	4, 54, 47, -99, -99, 60, 83, -62,
	-32, 10, 55, 56, 56, 55, 10, 3,
	-62, 12, -57, 44, -67, 28, 37, -31,
	-55, 50, 11, -4, -19, 13, 0, -49,
	-55, -43, -52, -28, -51, -47, -8, -50,
	-47, -42, -43, -79, -64, -32, -29, -32,
	-4, 3, -14, -50, -57, -18, 13, 4,
	17, 30, -3, -14, 6, -1, 40, 18}

func init() { RegisterConfigParamInt8ArrayDefault("black-king-pos-val", blackKingPosVals[:]) }

// From - https://chessprogramming.wikispaces.com/Simplified+evaluation+function
// Added some deliberate jitter to add artificial discrimination into stuck end-games (like Fine 70)
var blackKingEndgamePosVals = [64]int8{
	-50, -41, -30, -20, -20, -30, -41, -50,
	-32, -19, -9, 0, 0, -9, -19, -32,
	-30, -12, 17, 29, 29, 17, -12, -30,
	-29, -10, 30, 40, 40, 30, -10, -29,
	-29, -10, 31, 41, 41, 31, -10, -29,
	-30, -12, 18, 29, 29, 18, -12, -30,
	-32, -30, 4, 0, 0, 4, -30, -32,
	-50, -33, -30, -28, -28, -30, -33, -50}

func init() { RegisterConfigParamInt8ArrayDefault("black-king-eg-pos-val", blackKingEndgamePosVals[:]) }

var blackPiecePosVals = [7]*[64]int8{
	&nothingPosVals,
	&blackPawnPosVals,
	&blackKnightPosVals,
	&blackBishopPosVals,
	&blackRookPosVals,
	&blackQueenPosVals,
	&blackKingPosVals}

var colorPiecePosVals = [2]*[7]*[64]int8{
	&whitePiecePosVals,
	&blackPiecePosVals}



// Cheap part of static eval by opportunistic delta eval.
// Doing the easy case first and falling back to full eval until someone's more keen
func NegaStaticPiecePosEvalOrder0Fast(board *dragon.Board, prevEval0 EvalCp, moveInfo *dragon.BoardSaveT) EvalCp {
	// If the moving piece is a king, or we have captured a non-pawn, then just do full eval0.
	// King move includes castling which is extra tricky.
	// (This should be a small minority of moves)
	if moveInfo.FromPiece == dragon.King || endGameRatioChangesWithCapturePiece(moveInfo.CapturePiece) {
		// Full re-eval
		return negaStaticPiecePosEvalOrder0(board)
	} else {
		// Delta eval
		return prevEval0 + negaDeltaPiecePosEvalOrder0(board, moveInfo)
	}
}

// Cheap part  - O(0) by delta eval - of static eval from the perspective of the player to move
func negaStaticPiecePosEvalOrder0(board *dragon.Board) EvalCp {
	staticEval0 := StaticPiecePosEvalOrder0(board)

	if board.Colortomove == dragon.White {
		return staticEval0
	}
	return -staticEval0
}

// O(0) eval delta from the given move - from the perspective of the player to move
// Doesn't handle all cases = see NegaStaticEvalOrder0Fast(...)
func negaDeltaPiecePosEvalOrder0(board *dragon.Board, moveInfo *dragon.BoardSaveT) EvalCp {
	// This is after the move, so inverted from what you might expect
	captureColor := board.Colortomove
	fromToColor := dragon.Black ^ captureColor // TODO export this from dragontoothmg
	
	fromDelta := pieceMoveDelta(moveInfo.FromPiece, moveInfo.FromLoc, fromToColor)
	toDelta := pieceMoveDelta(moveInfo.ToPiece, moveInfo.ToLoc, fromToColor)
	captureDelta := pieceMoveDelta(moveInfo.CapturePiece, moveInfo.CaptureLoc, captureColor)
		
	// The player to move is the opposite of the player who last moved, so this is inverted from what you might expect
	return fromDelta - toDelta - captureDelta
}

func pieceMoveDelta(piece dragon.Piece, loc uint8, color dragon.ColorT) EvalCp {
	return pieceVals[piece] + EvalCp(colorPiecePosVals[color][piece][loc])
	
}

// Cheap part  - O(0) by delta eval - of static eval from white's perspective.
// This is full evaluation - we prefer to do much cheaper delta evaluation.
func StaticPiecePosEvalOrder0(board *dragon.Board) EvalCp {
	whitePiecesEval := piecesEval(&board.Bbs[dragon.White])
	blackPiecesEval := piecesEval(&board.Bbs[dragon.Black])

	piecesEval := whitePiecesEval - blackPiecesEval

	endGameRatio := endGameRatioByPiecesCount(board)

	whitePiecesPosEval := piecesPosVal(&board.Bbs[dragon.White], &whitePiecePosVals, &whiteKingEndgamePosVals, endGameRatio)
	blackPiecesPosEval := piecesPosVal(&board.Bbs[dragon.Black], &blackPiecePosVals, &blackKingEndgamePosVals, endGameRatio)

	piecesPosEval := whitePiecesPosEval - blackPiecesPosEval

	return piecesEval + piecesPosEval
}

const EVAL0_ONLY = false

// Expensive part - O(n) even with delta eval - of static eval from white's perspective.
func StaticPiecePosEvalOrderN(board *dragon.Board) EvalCp {

	endGameRatio := endGameRatioByPiecesCount(board)

	pawnExtrasEval := pawnExtrasVal(board)
	kingProtectionEval := kingProtectionVal(board, endGameRatio)
	bishopPairEval := bishopPairVal(board)
	endgameEval := endgameVal(board)

	orderNEval := pawnExtrasEval + kingProtectionEval + bishopPairEval + endgameEval

	if EVAL0_ONLY {
		orderNEval = EvalCp(0)
	}

	orderNEval += StaticPositionEvalOrderN(board)

	// Clamp it to the absolute bounds
	if orderNEval > MaxAbsStaticEvalOrderN {
		orderNEval = MaxAbsStaticEvalOrderN
	}
	if orderNEval < -MaxAbsStaticEvalOrderN {
		orderNEval = -MaxAbsStaticEvalOrderN
	}

	return orderNEval
}

// Expensive part - O(n) even with delta eval - of static eval from white's perspective.
func StaticPositionEvalOrderN(board *dragon.Board) EvalCp {
	var posEval PosEvalT
	InitPosEval(board, &posEval)
	return posEval.calcInfluenceEval() + posEval.calcSpaceEval()
}

// Sum of individual piece evals
func piecesEval(bitboards *dragon.Bitboards) EvalCp {
	eval := int(pawnVal) * bits.OnesCount64(bitboards[dragon.Pawn])
	eval += int(bishopVal) * bits.OnesCount64(bitboards[dragon.Bishop])
	eval += int(knightVal) * bits.OnesCount64(bitboards[dragon.Knight])
	eval += int(rookVal) * bits.OnesCount64(bitboards[dragon.Rook])
	eval += int(queenVal) * bits.OnesCount64(bitboards[dragon.Queen])

	return EvalCp(eval)
}

// Return true iff the given capture piece can possibly affect the end-game ratio
func endGameRatioChangesWithCapturePiece(capturePiece dragon.Piece) bool {
	return capturePiece != dragon.Nothing && capturePiece != dragon.Pawn
}

func nonPawnsCount(board *dragon.Board) int {
	allBW := board.Bbs[dragon.White][dragon.All] | board.Bbs[dragon.Black][dragon.All]
	pawnsBW := board.Bbs[dragon.White][dragon.Pawn] | board.Bbs[dragon.Black][dragon.Pawn]
	// Includes kings
	nonPawnsBW := allBW & ^pawnsBW

	return bits.OnesCount64(nonPawnsBW)
}

// Transition smoothly from King starting pos table to king end-game table between these total piece counts.
// Note these are counts of black plus white pieces excluding pawns and including kings.
const EndGamePiecesCountHi = 8
const EndGamePiecesCountLo = 4

// TODO - add the above two to config?

func endGameRatioForCount(count int) float64 {
	if count > EndGamePiecesCountHi {
		return 0.0
	}

	if count < EndGamePiecesCountLo {
		return 1.0
	}

	return float64(EndGamePiecesCountHi-count) / float64(EndGamePiecesCountHi-EndGamePiecesCountLo)
}

// To what extent are we in end game; from 0.0 (not at all) to 1.0 (definitely)
func endGameRatioByPiecesCount(board *dragon.Board) float64 {
	count := nonPawnsCount(board)
	return endGameRatioForCount(count)
}

// Return the endgame ratio's before and after the last capture move, presuming the last capture
//   was a non-pawn, i.e. non-pawn pieces count has decreased by 1
func endGameRatioByPiecesCountBeforeAndAfterCapture(board *dragon.Board) (float64, float64) {
	countAfter := nonPawnsCount(board)
	countBefore := countAfter - 1
	return endGameRatioForCount(countBefore), endGameRatioForCount(countAfter)
}


// Sum of piece position values
//   endGameRatio is a number between 0.0 and 1.0 where 1.0 means we're in end-game
func piecesPosVal(bitboards *dragon.Bitboards, piecePosVals *[7]*[64]int8, kingEndgamePosVals *[64]int8, endGameRatio float64) EvalCp {
	eval := pieceTypePiecesPosVal(bitboards[dragon.Pawn], piecePosVals[dragon.Pawn])
	eval += pieceTypePiecesPosVal(bitboards[dragon.Bishop], piecePosVals[dragon.Bishop])
	eval += pieceTypePiecesPosVal(bitboards[dragon.Knight], piecePosVals[dragon.Knight])
	eval += pieceTypePiecesPosVal(bitboards[dragon.Rook], piecePosVals[dragon.Rook])
	eval += pieceTypePiecesPosVal(bitboards[dragon.Queen], piecePosVals[dragon.Queen])

	kingStartEval := pieceTypePiecesPosVal(bitboards[dragon.King], piecePosVals[dragon.King])
	kingEndgameEval := pieceTypePiecesPosVal(bitboards[dragon.King], kingEndgamePosVals)

	kingEval := (1.0-endGameRatio)*float64(kingStartEval) + endGameRatio*float64(kingEndgameEval)

	return eval + EvalCp(kingEval)
}

// Sum of piece position values for a particular type of piece
func pieceTypePiecesPosVal(bitmask uint64, piecePosVals *[64]int8) EvalCp {
	var eval EvalCp = 0

	for bitmask != 0 {
		pos := bits.TrailingZeros64(bitmask)
		// (Could also use firstBit-1 trick to clear the bit)
		firstBit := uint64(1) << uint(pos)
		bitmask = bitmask ^ firstBit

		eval += EvalCp(piecePosVals[pos])
	}

	return eval
}

// Passed pawn bonuses
// TODO rationalise these with pawn pos vals
var pp2 int8 = 7
var pp3 int8 = 13
var pp4 int8 = 20
var pp5 int8 = 28
var pp6 int8 = 37

func init() {
	RegisterConfigParamInt8Default("passed-pawn-rank-2", &pp2)
	RegisterConfigParamInt8Default("passed-pawn-rank-3", &pp3)
	RegisterConfigParamInt8Default("passed-pawn-rank-4", &pp4)
	RegisterConfigParamInt8Default("passed-pawn-rank-5", &pp5)
	RegisterConfigParamInt8Default("passed-pawn-rank-6", &pp6)
}

var whitePassedPawnPosVals = [64]int8{
	0, 0, 0, 0, 0, 0, 0, 0,
	pp2, pp2, pp2, pp2, pp2, pp2, pp2, pp2,
	pp3, pp3, pp3, pp3, pp3, pp3, pp3, pp3,
	pp4, pp4, pp4, pp4, pp4, pp4, pp4, pp4,
	pp5, pp5, pp5, pp5, pp5, pp5, pp5, pp5,
	pp6, pp6, pp6, pp6, pp6, pp6, pp6, pp6,
	0, 0, 0, 0, 0, 0, 0, 0, // a 7th rank pawn is always passed, so covered by the pawn-pos-val
	0, 0, 0, 0, 0, 0, 0, 0}

var blackPassedPawnPosVals = [64]int8{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, // a 7th rank pawn is always passed, so covered by the pawn-pos-val
	pp6, pp6, pp6, pp6, pp6, pp6, pp6, pp6,
	pp5, pp5, pp5, pp5, pp5, pp5, pp5, pp5,
	pp4, pp4, pp4, pp4, pp4, pp4, pp4, pp4,
	pp3, pp3, pp3, pp3, pp3, pp3, pp3, pp3,
	pp2, pp2, pp2, pp2, pp2, pp2, pp2, pp2,
	0, 0, 0, 0, 0, 0, 0, 0}

// Bonus for pawns protecting pawns
// Disabled since it actually seems worse for root position search
var pProtPawnVal int8 = 0

// Bonus for pawns protecting pieces
var pProtPieceVal int8 = 10

// Penalty per doubled pawn
var doubledPawnPenalty int8 = -5

func init() {
	RegisterConfigParamInt8Default("pawn-prot-pawn", &pProtPawnVal)
	RegisterConfigParamInt8Default("pawn-prot-piece", &pProtPieceVal)
	RegisterConfigParamInt8Default("doubled-pawn", &doubledPawnPenalty)
}

// Pawn extras
func pawnExtrasVal(board *dragon.Board) EvalCp {
	wPawns := board.Bbs[dragon.White][dragon.Pawn]
	bPawns := board.Bbs[dragon.Black][dragon.Pawn]

	// Passed pawns
	wPawnScope := WPawnScope(wPawns)
	bPawnScope := BPawnScope(bPawns)

	wPassedPawns := wPawns & ^bPawnScope
	bPassedPawns := bPawns & ^wPawnScope

	wPPVal := pieceTypePiecesPosVal(wPassedPawns, &whitePassedPawnPosVals)
	bPPVal := pieceTypePiecesPosVal(bPassedPawns, &blackPassedPawnPosVals)

	// Pawns protected by pawns
	wPawnAtt := WPawnAttacks(wPawns)
	wPawnsProtectedByPawns := wPawnAtt & wPawns
	wPProtPawnsVal := bits.OnesCount64(wPawnsProtectedByPawns) * int(pProtPawnVal)

	bPawnAtt := BPawnAttacks(bPawns)
	bPawnsProtectedByPawns := bPawnAtt & bPawns
	bPProtPawnsVal := bits.OnesCount64(bPawnsProtectedByPawns) * int(pProtPawnVal)

	// Pieces protected by pawns
	wPieces := board.Bbs[dragon.White][dragon.All] & ^wPawns
	wPiecesProtectedByPawns := wPawnAtt & wPieces
	wPProtPiecesVal := bits.OnesCount64(wPiecesProtectedByPawns) * int(pProtPieceVal)

	bPieces := board.Bbs[dragon.Black][dragon.All] & ^bPawns
	bPiecesProtectedByPawns := bPawnAtt & bPieces
	bPProtPiecesVal := bits.OnesCount64(bPiecesProtectedByPawns) * int(pProtPieceVal)

	// Doubled pawns
	wPawnTelestop := NFill(N(wPawns))
	wDoubledPawns := wPawnTelestop & wPawns
	wDoubledPawnVal := bits.OnesCount64(wDoubledPawns) * int(doubledPawnPenalty)

	bPawnTelestop := SFill(S(bPawns))
	bDoubledPawns := bPawnTelestop & bPawns
	bDoubledPawnVal := bits.OnesCount64(bDoubledPawns) * int(doubledPawnPenalty)

	return (wPPVal - bPPVal) +
		EvalCp(wPProtPawnsVal-bPProtPawnsVal) +
		EvalCp(wPProtPiecesVal-bPProtPiecesVal) +
		EvalCp(wDoubledPawnVal-bDoubledPawnVal)
}

type KingProtectionT uint8

const (
	NoProtection KingProtectionT = iota
	QSideProtection
	KSideProtection
)

// Which white king positions qualify for protection eval - index 0 is square A1, index 63 is square H8
var wKingProtectionTypes = [64]KingProtectionT{
	QSideProtection, QSideProtection, QSideProtection, NoProtection, NoProtection, NoProtection, KSideProtection, KSideProtection,
	QSideProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, KSideProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection}

// Bitboard locations of white king protecting pieces indexes by protection type
var wKingProtectionBbs = [3]uint64{
	0x0,                // NoProtection
	0x0007070000000000, // QSideProtection
	0x00e0e00000000000} // KSideProtection

// Which black king positions qualify for protection eval
var bKingProtectionTypes = [64]KingProtectionT{
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection,
	QSideProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, NoProtection, KSideProtection,
	QSideProtection, QSideProtection, QSideProtection, NoProtection, NoProtection, NoProtection, KSideProtection, KSideProtection}

// Bitboard locations of black king protecting pieces indexes by protection type
var bKingProtectionBbs = [3]uint64{
	0x0,                // NoProtection
	0x0000000000070700, // QSideProtection
	0x0000000000e0e000} // KSideProtection

// Bonus for pieces that are protecting the king
var kingProtectorVal int8 = 8

// Additional bonus for pawns that are protecting the king
var kingPawnProtectorVal int8 = 11

func init() {
	RegisterConfigParamInt8Default("king-prot", &kingProtectorVal)
	RegisterConfigParamInt8Default("king-pawn-prot", &kingPawnProtectorVal)
}

// Naive king protection - count pieces around the king if the king is in the corner
// From White's perspective
func kingProtectionVal(board *dragon.Board, endGameRatio float64) EvalCp {
	if endGameRatio == 1.0 {
		return 0
	}

	wBbs := &board.Bbs[dragon.White]
	wKingPos := bits.TrailingZeros64(wBbs[dragon.King])
	wKingProtectionType := wKingProtectionTypes[wKingPos]
	wKingProtectionBb := wKingProtectionBbs[wKingProtectionType]

	wNonKingPieces := wBbs[dragon.All] & ^wBbs[dragon.King]
	wKingProtectors := wNonKingPieces & wKingProtectionBb

	wKingPawnProtectors := wBbs[dragon.Pawn] & wKingProtectionBb

	wKingProtectionVal := bits.OnesCount64(wKingProtectors)*int(kingProtectorVal) + bits.OnesCount64(wKingPawnProtectors)*int(kingPawnProtectorVal)

	bBbs := &board.Bbs[dragon.Black]
	bKingPos := bits.TrailingZeros64(bBbs[dragon.King])
	bKingProtectionType := bKingProtectionTypes[bKingPos]
	bKingProtectionBb := bKingProtectionBbs[bKingProtectionType]

	bNonKingPieces := bBbs[dragon.All] & ^bBbs[dragon.King]
	bKingProtectors := bNonKingPieces & bKingProtectionBb

	bKingPawnProtectors := bBbs[dragon.Pawn] & bKingProtectionBb

	bKingProtectionVal := bits.OnesCount64(bKingProtectors)*int(kingProtectorVal) + bits.OnesCount64(bKingPawnProtectors)*int(kingPawnProtectorVal)

	// King protection in end-game is irrelevant
	return EvalCp(float64(wKingProtectionVal-bKingProtectionVal) * (1.0 - endGameRatio))
}

// Bishop pair bonuses
// From White's perspective
func bishopPairVal(board *dragon.Board) EvalCp {
	return bishopPairColorVal(&board.Bbs[dragon.White]) - bishopPairColorVal(&board.Bbs[dragon.Black])
}

const blackSquares = uint64(0x5555555555555555)
const whiteSquares = uint64(0xaaaaaaaaaaaaaaaa)
var bishopPairBonus = EvalCp(80)
var bishopPairProxBonus = EvalCp(3)

func init() {
	RegisterConfigParamEvalCpDefault("bishop-pair", &bishopPairBonus)
	RegisterConfigParamEvalCpDefault("bishop-pair-proximity", &bishopPairProxBonus)
}

func bishopPairColorVal(bitboards *dragon.Bitboards) EvalCp {
	bishopsBonus := EvalCp(0)
	bishops := bitboards[dragon.Bishop]
	blackBishops := bishops & blackSquares
	whiteBishops := bishops & whiteSquares
	if blackBishops != 0 && whiteBishops != 0 {
		bishopsBonus += bishopPairBonus

		// Prefer Bishops to hunt close together (wildly speculative but it looks pretty).
		// We assume there's only one bishop of each colour - but really if you are going to under-promote to a bishop you deserve a bogus bishopsBonus.
		bBishopPos := bits.TrailingZeros64(blackBishops)
		wBishopPos := bits.TrailingZeros64(whiteBishops)

		bishopsBonus += bishopPairProxBonus*EvalCp(2-kingWalkDistance(bBishopPos, wBishopPos))
	}

	return bishopsBonus
}

// Return the rank and file of a position
func rankFile(pos int) (int, int) {
	return (pos >> 3), (pos & 7)
}

// King-walk distance - I can't remember what the formal term is - basically max(delta(x), delta(y))
func kingWalkDistance(pos1 int, pos2 int) int {
	pos1Rank, pos1File := rankFile(pos1)
	pos2Rank, pos2File := rankFile(pos2)

	rankDiff8 := absDiff8[pos1Rank][pos2Rank]
	fileDiff8 := absDiff8[pos1File][pos2File]

	return int(max8[rankDiff8][fileDiff8])
}

// Absolute difference between two numbers in range [0,7]
var absDiff8 = [8][8]uint8{
	{0, 1, 2, 3, 4, 5, 6, 7}, // [0][n]
	{1, 0, 1, 2, 3, 4, 5, 6}, // [1][n]
	{2, 1, 0, 1, 2, 3, 4, 5}, // [2][n]
	{3, 2, 1, 0, 1, 2, 3, 4}, // [3][n]
	{4, 3, 2, 1, 0, 1, 2, 3}, // [4][n]
	{5, 4, 3, 2, 1, 0, 1, 2}, // [5][n]
	{6, 5, 4, 3, 2, 1, 0, 1}, // [6][n]
	{7, 6, 5, 4, 3, 2, 1, 0}} // [7][n]

// Max of two numbers in range [0,7]
var max8 = [8][8]uint8{
	{0, 1, 2, 3, 4, 5, 6, 7}, // [0][n]
	{1, 1, 2, 3, 4, 5, 6, 7}, // [1][n]
	{2, 2, 2, 3, 4, 5, 6, 7}, // [2][n]
	{3, 3, 3, 3, 4, 5, 6, 7}, // [3][n]
	{4, 4, 4, 4, 4, 5, 6, 7}, // [4][n]
	{5, 5, 5, 5, 5, 5, 6, 7}, // [5][n]
	{6, 6, 6, 6, 6, 6, 6, 7}, // [6][n]
	{7, 7, 7, 7, 7, 7, 7, 7}} // [7][n]

// End-game bonuses/penalties from White's perspective
func endgameVal(board *dragon.Board) EvalCp {
	return endgameColorVal(board, dragon.White) - endgameColorVal(board, dragon.Black)
}

func oppColor(color dragon.ColorT) dragon.ColorT {
	return dragon.Black ^ color
}

var endgameKingPawnProxBonus = EvalCp(13)
var loneMinorPiecePenalty = -EvalCp(190)
var endgameKingMajorProxBonus = EvalCp(17)

func init() {
	RegisterConfigParamEvalCpDefault("endgame-king-pawn-proximity", &endgameKingPawnProxBonus)
	RegisterConfigParamEvalCpDefault("lone-minor-piece", &loneMinorPiecePenalty)
	RegisterConfigParamEvalCpDefault("endgame-king-major-proximity", &endgameKingMajorProxBonus)
}

// TODO (rpj) I'm a bit concerned that these various eval categories will lead to sudden changes in the eval.
func endgameColorVal(board *dragon.Board, color dragon.ColorT) EvalCp {
	myBbs := board.Bbs[color]
	oppBbs := board.Bbs[oppColor(color)]

	myPawns := myBbs[dragon.Pawn]
	myKings := myBbs[dragon.King]
	myAll := myBbs[dragon.All]

	myPieces := myAll & ^(myPawns | myKings)

	// Case 1 - only pawns: we want to get the king next to a pawn to optimise promo potential
	if myPieces == 0 && myPawns != 0 {
		var fwdPawnPos int // the most advanced pawn
		if color == dragon.White {
			fwdPawnPos = 63 - bits.LeadingZeros64(myPawns)
		} else {
			fwdPawnPos = bits.TrailingZeros64(myPawns)
		}
		kingPos := bits.TrailingZeros64(myKings)

		pawnEgBonus := endgameKingPawnProxBonus * EvalCp(3-kingWalkDistance(fwdPawnPos, kingPos))

		// fmt.Println("       color", color, "fwdPawnPos", fwdPawnPos, "kingPos", kingPos, "dist", kingWalkDistance(fwdPawnPos, kingPos), "bonus", pawnEgBonus)
		return pawnEgBonus
	}

	// Case 2 - no pawns: we want to get the king near a rook or queen and shrink the opponent's king box
	if myPieces != 0 && myPawns == 0 {
		myMajors := myBbs[dragon.Rook] | myBbs[dragon.Queen]
		if myMajors == 0 {
			// Bummer dude, only minor pieces.
			// If there's only one, then penalise the eval - even a pawn gives more hope.
			myMinors := myBbs[dragon.Bishop] | myBbs[dragon.Knight]
			nMinors := bits.OnesCount64(myMinors)
			if nMinors == 1 {
				return loneMinorPiecePenalty
			}
			// TODO(rpj) How do we encourage minor piece checkmate?
		} else {
			// Encourage rook-style box shrinking mate by rewarding a small opponent king box,
			//   and king proximity to the major piece.
			majorPos := bits.TrailingZeros64(myMajors)
			kingPos := bits.TrailingZeros64(myKings)
			oppKingPos := bits.TrailingZeros64(oppBbs[dragon.King])
			majorRank, majorFile := rankFile(majorPos)
			kingRank, kingFile := rankFile(majorPos)
			oppKingRank, oppKingFile := rankFile(oppKingPos)

			// No bonus if opposition king is on a check line
			majorCheckMateBonus := EvalCp(0)
			if oppKingRank != majorRank && oppKingFile != majorFile {
				// No proximity bonus if our king is outside the opposition king's box
				kingSouthOfMajor := kingRank < majorRank
				kingWestOfMajor := kingFile < majorFile
				oppKingSouthOfMajor := oppKingRank < majorRank
				oppKingWestOfMajor := oppKingFile < majorFile
				
				kingMajorProxBonus := EvalCp(0)
				// TODO(rpj) - this is wrong for some cases of our king and our major on the same rank or file
				if kingSouthOfMajor != oppKingSouthOfMajor || kingWestOfMajor != oppKingWestOfMajor {
					kingMajorProxBonus = endgameKingMajorProxBonus * EvalCp(2-kingWalkDistance(majorPos, kingPos))
				}
				
				// Caculate opposition king square size
				kingBoxNRanks := majorRank
				if !oppKingSouthOfMajor {
					kingBoxNRanks = 7-majorRank
				}
				kingBoxNFiles := majorFile
				if !oppKingWestOfMajor {
					kingBoxNFiles = 7-majorFile
				}

				majorCheckMateBonus = kingMajorProxBonus + EvalCp(7*7 - kingBoxNRanks*kingBoxNFiles) // smaller box is better
			}
			return majorCheckMateBonus
		}
	}
	return EvalCp(0)
}