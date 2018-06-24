package engine

import (
	"errors"
	// "fmt"
	// "time"

	dragon "github.com/Bubblyworld/dragontoothmg"
)

type SearchStatsT struct {
	Nodes uint64           // #nodes visited
	Mates uint64           // #true terminal nodes
	NonLeafs uint64        // #non-leaf nodes
	Killers uint64         // #nodes with killer move available
	KillerCuts uint64      // #nodes with killer move cut
	DeepKillers uint64     // #nodes with deep killer move available
	DeepKillerCuts uint64  // #nodes with deep killer move cut
	QNodes uint64          // #nodes visited in qsearch
	QMates uint64          // #true terminal nodes in qsearch
	QNonLeafs uint64       // #non-leaf qnodes 
	QKillers uint64        // #qnodes with killer move available
	QKillerCuts uint64     // #qnodes with killer move cut
	QDeepKillers uint64    // #nodes with deep killer move available
	QDeepKillerCuts uint64 // #nodes with deep killer move cut
	QPats uint64           // #qnodes with stand pat best
	QPatCuts uint64        // #qnodes with stand pat cut
	QShallows uint64       // #qnodes where we reached full depth - i.e. likely failed to quiesce
	QPrunes uint64         // #qnodes where we reached full depth - i.e. likely failed to quiesce
}
	

// Configuration options
type SearchAlgorithmT int
const (
	MiniMax SearchAlgorithmT = iota
	AlphaBeta
)
var SearchAlgorithm = AlphaBeta
var SearchDepth = 6
var UseQSearch = true
var QSearchDepth = 6
var UseKillerMoves = true
var UseDeepKillerMoves = true
var UseDeltaEval = true

func SearchAlgorithmString() string {
	switch SearchAlgorithm {
	case MiniMax:
		return "MiniMax"
	case AlphaBeta:
		return "AlphaBeta"
	default:
		SearchAlgorithm = AlphaBeta
		return "AlphaBeta"
	}
}

const MaxDepth = 1024
const NoMove dragon.Move = 0

func Search(board *dragon.Board) (dragon.Move, EvalCp, SearchStatsT, error) {
	var deepKillers [MaxDepth]dragon.Move
	var searchStats SearchStatsT
	var bestMove = NoMove
	var eval EvalCp = 0

	switch SearchAlgorithm {
	case MiniMax:
		bestMove, eval = minimax(board, /*depthToGo*/SearchDepth, /*depthFromRoot*/0, StaticEval(board), &searchStats)

	case AlphaBeta:
		bestMove, eval = alphabeta(board, /*depthToGo*/SearchDepth, /*depthFromRoot*/0, BlackCheckMateEval, WhiteCheckMateEval, StaticEval(board), NoMove, deepKillers[:], &searchStats)

	default:
		return NoMove, 0, searchStats, errors.New("bot: unrecognised search algorithm")
	}

	if bestMove == NoMove {
		return NoMove, 0, searchStats, errors.New("bot: no legal move found in search")
	}

	return bestMove, eval, searchStats, nil
}

// Return the eval for stalemate or checkmate from white perspective.
// Only valid if there are no legal moves.
func mateEval(board *dragon.Board, depthFromRoot int) EvalCp {
	if board.OurKingInCheck() {
		// checkmate - closer to root is better
		if board.Wtomove {
			return BlackCheckMateEval + EvalCp(depthFromRoot)
		}
		
		return WhiteCheckMateEval - EvalCp(depthFromRoot)
	}
	// stalemate
	return DrawEval
}

// Return the new static eval - either by fast delta, or by full evaluation, depending on configuration
func getStaticEval(board* dragon.Board, oldStaticEval EvalCp, move dragon.Move, moveInfo *dragon.MoveApplication, isWhiteMove bool) EvalCp {
	if UseDeltaEval {
		// Much faster
		return oldStaticEval + EvalDelta(move, moveInfo, isWhiteMove)
	} else {
		// Much more dependable :P
		return StaticEval(board)
	}
}

// Return the best eval attainable through minmax from the given
//   position, along with the move leading to the principal variation.
func minimax(board *dragon.Board, depthToGo int, depthFromRoot int, staticEval EvalCp, stats *SearchStatsT) (dragon.Move, EvalCp) {

	stats.Nodes++
	stats.NonLeafs++

	// Generate all legal moves - thanks dragontoothmg!
	legalMoves := board.GenerateLegalMoves()

	// Check for checkmate or stalemate
	if len(legalMoves) == 0 {
		stats.Mates++
		return NoMove, mateEval(board, depthFromRoot)
	}
	
	isWhiteMove := board.Wtomove
	var bestMove = NoMove
	var bestEval EvalCp = WhiteCheckMateEval
	if isWhiteMove {
		bestEval = BlackCheckMateEval
	}

	for _, move := range legalMoves {
		// Make the move
		moveInfo := board.Apply2(move)

		newStaticEval := getStaticEval(board, staticEval, move, moveInfo, isWhiteMove)

		// Get the (deep) eval
		var eval EvalCp
		if depthToGo <= 1 {
			stats.Nodes++
			// Ignore mate check to avoid generating moves at all leaf nodes
			eval = newStaticEval
		} else {
			_, eval = minimax(board, depthToGo-1, depthFromRoot+1, newStaticEval, stats)
		}

		// Take back the move
		moveInfo.Unapply()

		// If we're white, we try to maximise our eval. If we're black, we try to
		// minimise our eval.
		if isWhiteMove {
			// Strictly > to match alphabeta
			if eval > bestEval {
				bestEval, bestMove = eval, move
			}
		} else {
			// Strictly < to match alphabeta
			if eval < bestEval {
				bestEval, bestMove = eval, move
			}
		}
	}

	return bestMove, bestEval
}


// Move the killer move to the front of the legal moves list, if it's in the legal moves list
func prioritiseKillerMove(legalMoves []dragon.Move, killer dragon.Move) {
	if killer != NoMove {
		for i := 0; i < len(legalMoves); i++ {
			if legalMoves[i] == killer {
				legalMoves[0], legalMoves[i] = killer, legalMoves[0]
				break
			}
		}
	}
}

// Return the best eval attainable through alpha-beta from the given position (with killer-move hint), along with the move leading to the principal variation.
func alphabeta(board *dragon.Board, depthToGo int, depthFromRoot int, alpha EvalCp, beta EvalCp, staticEval EvalCp, killer dragon.Move, deepKillers []dragon.Move, stats *SearchStatsT) (dragon.Move, EvalCp) {

	stats.Nodes++
	stats.NonLeafs++

	// Generate all legal moves - thanks dragontoothmg!
	legalMoves := board.GenerateLegalMoves()

	// Check for checkmate or stalemate
	if len(legalMoves) == 0 {
		stats.Mates++
		return NoMove, mateEval(board, depthFromRoot)
	}

	usingDeepKiller := false
	if UseKillerMoves {
		if killer == NoMove && UseDeepKillerMoves {
			usingDeepKiller = true
			killer = deepKillers[depthFromRoot]
		}
		// Place killer-move first if it's there
		prioritiseKillerMove(legalMoves, killer)
		if legalMoves[0] == killer {
			if usingDeepKiller {
				stats.DeepKillers++
			} else {
				stats.Killers++
			}
		}
	}

	// Would be smaller with negalpha-beta but this is simple
	if board.Wtomove {
		// White to move - maximise eval with beta cut-off
		var bestMove = NoMove
		var bestEval EvalCp = BlackCheckMateEval
		childKiller := NoMove
		
		for _, move := range legalMoves {
			// Make the move
			moveInfo := board.Apply2(move)
			
			newStaticEval := getStaticEval(board, staticEval, move, moveInfo, /*isWhiteMove*/true)
			
			// Get the (deep) eval
			var eval EvalCp
			if depthToGo <= 1 {
				stats.Nodes ++
				if UseQSearch {
					// Quiesce
					childKiller, eval = qsearchAlphabeta(board, QSearchDepth, depthFromRoot+1, alpha, beta, newStaticEval, childKiller, deepKillers, stats)
				} else {
					eval = newStaticEval
				}
			} else {
				childKiller, eval = alphabeta(board, depthToGo-1, depthFromRoot+1, alpha, beta, newStaticEval, childKiller, deepKillers, stats)
			}

			// Take back the move
			moveInfo.Unapply()
			
			// We're white - maximise our eval.
			// Note - this MUST be strictly > because we fail-soft AT the current best evel - beware!
			if eval > bestEval {
				bestEval, bestMove = eval, move
			}

			if alpha < bestEval {
				alpha = bestEval
			}

			// Note that this is aggressive, and we fail-soft AT the parent's best eval - be very ware!
			if beta <= alpha {
				if bestMove == killer {
					if usingDeepKiller {
						stats.DeepKillerCuts++
					} else {
						stats.KillerCuts++
					}
				}
				// beta cut-off
 				deepKillers[depthFromRoot] = bestMove
				return bestMove, bestEval
			}
		}
		
		deepKillers[depthFromRoot] = bestMove
		return bestMove, bestEval
	} else {
		// Black to move - minimise eval with alpha cut-off
		var bestMove = NoMove
		var bestEval EvalCp = WhiteCheckMateEval
		childKiller := NoMove
		
		for _, move := range legalMoves {
			// Make the move
			moveInfo := board.Apply2(move)
			
			newStaticEval := getStaticEval(board, staticEval, move, moveInfo, /*isWhiteMove*/false)
			
			// Get the (deep) eval
			var eval EvalCp
			if depthToGo <= 1 {
				stats.Nodes++
				if UseQSearch {
					// Quiesce
					childKiller, eval = qsearchAlphabeta(board, QSearchDepth, depthFromRoot+1, alpha, beta, newStaticEval, childKiller, deepKillers, stats)
				} else {
					eval = newStaticEval
				}
			} else {
				childKiller, eval = alphabeta(board, depthToGo-1, depthFromRoot+1, alpha, beta, newStaticEval, childKiller, deepKillers, stats)
			}
			
			// Take back the move
			moveInfo.Unapply()
			
			// We're black - minimise our eval.
			// Note - this MUST be strictly < because we fail-soft AT the current best evel - beware!
			if eval < bestEval {
				bestEval, bestMove = eval, move
			}

			if bestEval < beta {
				beta = bestEval
			}

			// Note that this is aggressive, and we fail-soft AT the parent's best eval - be very ware!
			if beta <= alpha {
				if bestMove == killer {
					if usingDeepKiller {
						stats.DeepKillerCuts++
					} else {
						stats.KillerCuts++
					}
				}
				// alpha cut-off
				deepKillers[depthFromRoot] = bestMove
				return bestMove, bestEval
			}
		}

		deepKillers[depthFromRoot] = bestMove
		return bestMove, bestEval
	}
}

// Quiescence search - differs from full search as follows:
//   - we only look at captures and promotions - we could/should also possibly look at check evasion, but check detection is currently expensive
//   - we consider 'standing pat' - i.e. do alpha/beta cutoff according to the node's static eval (TODO)
// TODO - implement more efficient generation of 'noisy' moves in dragontoothmg
// TODO - better static eval if we bottom out without quescing, e.g. static exchange evaluation (SEE)
func qsearchAlphabeta(board *dragon.Board, qDepthToGo int, depthFromRoot int, alpha EvalCp, beta EvalCp, staticEval EvalCp, killer dragon.Move, deepKillers []dragon.Move, stats *SearchStatsT) (dragon.Move, EvalCp) {

	stats.QNodes++
	stats.QNonLeafs++

	// Stand pat - equivalent to considering the null move as a valid move.
	// Essentially the player to move doesn't _have_ to make a capture - (we assume that there is a non-capture move available.)
	if board.Wtomove {
		if alpha < staticEval {
			alpha = staticEval
		}

		// Note that this is aggressive, and we fail-soft AT the parent's best eval - be very ware!
		if beta <= alpha {
			stats.QPats++
			stats.QPatCuts++
			return NoMove, staticEval
		}

	} else {
		if staticEval < beta {
			beta = staticEval
		}
		
		// Note that this is aggressive, and we fail-soft AT the parent's best eval - be very ware!
		if beta <= alpha {
			stats.QPats++
			stats.QPatCuts++
			return NoMove, staticEval
		}
	}

	// Generate all noisy legal moves - thanks dragontoothmg!
	legalMoves, isInCheck := board.GenerateLegalMoves2(/*onlyCapturesPromosCheckEvasion*/true)

	// No noisy mvoes
	if len(legalMoves) == 0 {
		// Check for checkmate or stalemate
		if isInCheck {
			stats.QMates++
			return NoMove, mateEval(board, depthFromRoot) // TODO checks for mate again expensively
		} else {
			// Already quiesced - just return static eval
			stats.QShallows++
			return NoMove, staticEval
		}
	}

	usingDeepKiller := false
	if UseKillerMoves {
		if killer == NoMove && UseDeepKillerMoves {
			usingDeepKiller = true
			killer = deepKillers[depthFromRoot]
		}
		// Place killer-move first if it's there
		prioritiseKillerMove(legalMoves, killer)
		if legalMoves[0] == killer {
			if usingDeepKiller {
				stats.QDeepKillers++
			} else {
				stats.QKillers++
			}
		}
	}

	// Would be smaller with negalpha-beta but this is simple
	if board.Wtomove {
		// White to move - maximise eval with beta cut-off
		var bestMove = NoMove
		var bestEval EvalCp = BlackCheckMateEval
		childKiller := NoMove
		
		for _, move := range legalMoves {
			// Make the move
			moveInfo := board.Apply2(move)
			
			newStaticEval := getStaticEval(board, staticEval, move, moveInfo, /*isWhiteMove*/true)
			
			// Get the (deep) eval
			var eval EvalCp
			if qDepthToGo <= 1 {
				stats.QNodes++
				stats.QPrunes++
				// Ignore mate check to avoid generating moves at all leaf nodes
				eval = newStaticEval
			} else {
				childKiller, eval = qsearchAlphabeta(board, qDepthToGo-1, depthFromRoot+1, alpha, beta, newStaticEval, childKiller, deepKillers, stats)
			}

			// Take back the move
			moveInfo.Unapply()
			
			// We're white - maximise our eval.
			// Note - this MUST be strictly > because we fail-soft AT the current best evel - beware!
			if eval > bestEval {
				bestEval, bestMove = eval, move
			}

			if alpha < bestEval {
				alpha = bestEval
			}

			// Note that this is aggressive, and we fail-soft AT the parent's best eval - be very ware!
			if beta <= alpha {
				if bestMove == killer {
					if usingDeepKiller {
						stats.QDeepKillerCuts++
					} else {
						stats.QKillerCuts++
					}
				}
				// beta cut-off
 				deepKillers[depthFromRoot] = bestMove
				return bestMove, bestEval
			}
		}

		stats.QPats++
		deepKillers[depthFromRoot] = bestMove
		return bestMove, bestEval
	} else {
		// Black to move - minimise eval with alpha cut-off
		var bestMove = NoMove
		var bestEval EvalCp = WhiteCheckMateEval
		childKiller := NoMove
		
		for _, move := range legalMoves {
			// Make the move
			moveInfo := board.Apply2(move)
			
			newStaticEval := getStaticEval(board, staticEval, move, moveInfo, /*isWhiteMove*/false)
			
			// Get the (deep) eval
			var eval EvalCp
			if qDepthToGo <= 1 {
				stats.QNodes++
				stats.QPrunes++
				// Ignore mate check to avoid generating moves at all leaf nodes
				eval = newStaticEval
			} else {
				childKiller, eval = qsearchAlphabeta(board, qDepthToGo-1, depthFromRoot+1, alpha, beta, newStaticEval, childKiller, deepKillers, stats)
			}
			
			// Take back the move
			moveInfo.Unapply()
			
			// We're black - minimise our eval.
			// Note - this MUST be strictly < because we fail-soft AT the current best evel - beware!
			if eval < bestEval {
				bestEval, bestMove = eval, move
			}

			if bestEval < beta {
				beta = bestEval
			}

			// Note that this is aggressive, and we fail-soft AT the parent's best eval - be very ware!
			if beta <= alpha {
				if bestMove == killer {
					if usingDeepKiller {
						stats.QDeepKillerCuts++
					} else {
						stats.QKillerCuts++
					}
				}
				// alpha cut-off
 				deepKillers[depthFromRoot] = bestMove
				return bestMove, bestEval
			}
		}

		stats.QPats++
		deepKillers[depthFromRoot] = bestMove
		return bestMove, bestEval
	}
}
