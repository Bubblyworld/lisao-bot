// Stolen shamelessly from dragontooth

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	dragon "github.com/dylhunn/dragontoothmg"

	"clanpj/lisao/engine"
)

var VersionString = "0.0h Pikachu 1" + "CPU 5ply " + runtime.GOOS + "-" + runtime.GOARCH

func main() {
	uciLoop()
}

func uciLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	board := dragon.ParseFen(dragon.Startpos) // the game board
	// used for communicating with search routine
	haltchannel := make(chan bool)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		if len(tokens) == 0 { // ignore blank lines
			continue
		}
		switch strings.ToLower(tokens[0]) {
		case "uci":
			fmt.Println("id name Lisao", VersionString)
			fmt.Println("id author Clan PJ")
			// fmt.Println("option name Hash type spin default", transtable.DefaultTtableSize, "min 8 max 65536")
			// fmt.Println("option name SearchThreads type spin default", search.DefaultSearchThreads, "min 1 max 128")
			// fmt.Println("option name DrawVal_Contempt_Centipawns type spin default",
			// 	eval.DefaultDrawScore, "min", search.NegInf, "max", search.PosInf)
			fmt.Println("uciok")
		case "isready":
			fmt.Println("readyok")
		case "ucinewgame":
			//transtable.Initialize(transtable.DefaultTtableSize)
			// reset the board, in case the GUI skips 'position' after 'newgame'
			board = dragon.ParseFen(dragon.Startpos)
			// reset the history map
		case "quit":
			return
		case "setoption":
			if len(tokens) != 5 || tokens[1] != "name" || tokens[3] != "value" {
				fmt.Println("info string Malformed setoption command")
				continue
			}
			switch strings.ToLower(tokens[2]) {
			// case "hash":
			// 	res, err := strconv.Atoi(tokens[4])
			// 	if err != nil {
			// 		fmt.Println("info string Hash value is not an int (", err, ")")
			// 		continue
			// 	}
			// 	fmt.Println("info string Changed table size. Clearing and reloading table...")
			// 	transtable.DefaultTtableSize = res // reset the size and reload the table
			// 	transtable.Initialize(transtable.DefaultTtableSize)
			// case "searchthreads":
			// 	res, err := strconv.Atoi(tokens[4])
			// 	if err != nil {
			// 		fmt.Println("info string Number of threads is not an int (", err, ")")
			// 		continue
			// 	}
			// 	search.DefaultSearchThreads = res
			// case "DrawVal_Contempt_Centipawns":
			// 	res, err := strconv.Atoi(tokens[4])
			// 	if err != nil {
			// 		fmt.Println("info string DrawVal_Contempt_Centipawns is not an int (", err, ")")
			// 		continue
			// 	}
			// 	fmt.Println("info string Changed contempt factor to", res, "centipawns.")
			// 	eval.DefaultDrawScore = int16(res)
			default:
				fmt.Println("info string Unknown UCI option", tokens[2])
			}
		case "go":
			goScanner := bufio.NewScanner(strings.NewReader(line))
			goScanner.Split(bufio.ScanWords)
			goScanner.Scan() // skip the first token
			var wtime, btime, winc, binc int
			var infinite bool
			var err error
			for goScanner.Scan() {
				nextToken := strings.ToLower(goScanner.Text())
				switch nextToken {
				case "infinite":
					infinite = true
					continue
				case "wtime":
					if !goScanner.Scan() {
						fmt.Println("info string Malformed go command option wtime")
						continue
					}
					wtime, err = strconv.Atoi(goScanner.Text())
					if err != nil {
						fmt.Println("info string Malformed go command option; could not convert wtime")
						continue
					}
				case "btime":
					if !goScanner.Scan() {
						fmt.Println("info string Malformed go command option btime")
						continue
					}
					btime, err = strconv.Atoi(goScanner.Text())
					if err != nil {
						fmt.Println("info string Malformed go command option; could not convert btime")
						continue
					}
				case "winc":
					if !goScanner.Scan() {
						fmt.Println("info string Malformed go command option winc")
						continue
					}
					winc, err = strconv.Atoi(goScanner.Text())
					if err != nil {
						fmt.Println("info string Malformed go command option; could not convert winc")
						continue
					}
				case "binc":
					if !goScanner.Scan() {
						fmt.Println("info string Malformed go command option binc")
						continue
					}
					binc, err = strconv.Atoi(goScanner.Text())
					if err != nil {
						fmt.Println("info string Malformed go command option; could not convert binc")
						continue
					}
				default:
					fmt.Println("info string Unknown go subcommand", nextToken)
					continue
				}
			}
			stop := false
			go uciSearch(&board, haltchannel, &stop)
			// TODO (rpj) - work out how this unblocks in the case of infinite time???
			if wtime != 0 && btime != 0 && !infinite { // If times are specified
				var ourtime, opptime, ourinc, oppinc int
				if board.Wtomove {
					ourtime, opptime, ourinc, oppinc = wtime, btime, winc, binc
				} else {
					ourtime, opptime, ourinc, oppinc = btime, wtime, binc, winc
				}
				allowedTime := uciCalculateAllowedTimeMs(&board, ourtime, opptime, ourinc, oppinc)
				go uciSearchTimeout(haltchannel, allowedTime, &stop)
			}
		// case "secretparam": // secret parameters used for optimizing the evaluation function
		// 	res, _ := strconv.Atoi(tokens[2])
		// 	switch tokens[1] {
		// 	case "BishopPairBonus":
		// 		eval.BishopPairBonus = res
		// 	case "DiagonalMobilityBonus":
		// 		eval.DiagonalMobilityBonus = res
		// 	case "OrthogonalMobilityBonus":
		// 		eval.OrthogonalMobilityBonus = res
		// 	case "DoubledPawnPenalty":
		// 		eval.DoubledPawnPenalty = res
		// 	case "PassedPawnBonus":
		// 		eval.PassedPawnBonus = res
		// 	case "IsolatedPawnPenalty":
		// 		eval.IsolatedPawnPenalty = res

		// 	default:
		// 		if tokens[1][0:14] == "PawnTableStart" {
		// 			idx := tokens[1][14:len(tokens[1])]
		// 			square, _ := strconv.Atoi(idx)
		// 			val, _ := strconv.Atoi(tokens[2])
		// 			eval.PawnTableStart[square] = val
		// 		} else if tokens[1][0:14] == "KingTableStart" {
		// 			idx := tokens[1][14:len(tokens[1])]
		// 			square, _ := strconv.Atoi(idx)
		// 			val, _ := strconv.Atoi(tokens[2])
		// 			eval.KingTableStart[square] = val
		// 		} else if tokens[1][0:15] == "CentralizeTable" {
		// 			idx := tokens[1][15:len(tokens[1])]
		// 			square, _ := strconv.Atoi(idx)
		// 			val, _ := strconv.Atoi(tokens[2])
		// 			eval.CentralizeTable[square] = val
		// 		} else if tokens[1][0:16] == "KnightTableStart" {
		// 			idx := tokens[1][16:len(tokens[1])]
		// 			square, _ := strconv.Atoi(idx)
		// 			val, _ := strconv.Atoi(tokens[2])
		// 			eval.KnightTableStart[square] = val
		// 		} else {
		// 			fmt.Println("Unknown secret param")
		// 		}
		// 	}
		case "stop":
			haltchannel <- true // TODO(dylhunn): stop deadlock on double stop
		case "position":
			posScanner := bufio.NewScanner(strings.NewReader(line))
			posScanner.Split(bufio.ScanWords)
			posScanner.Scan() // skip the first token
			if !posScanner.Scan() {
				fmt.Println("info string Malformed position command")
				continue
			}
			//search.HistoryMap = make(map[uint64]int) // reset the history map
			if strings.ToLower(posScanner.Text()) == "startpos" {
				board = dragon.ParseFen(dragon.Startpos)
				//search.HistoryMap[board.Hash()]++ // record that this state has occurred
				posScanner.Scan() // advance the scanner to leave it in a consistent state
			} else if strings.ToLower(posScanner.Text()) == "fen" {
				fenstr := ""
				for posScanner.Scan() && strings.ToLower(posScanner.Text()) != "moves" {
					fenstr += posScanner.Text() + " "
				}
				if fenstr == "" {
					fmt.Println("info string Invalid fen position")
					continue
				}
				board = dragon.ParseFen(fenstr)
				//search.HistoryMap[board.Hash()]++ // record that this state has occurred
			} else {
				fmt.Println("info string Invalid position subcommand")
				continue
			}
			if strings.ToLower(posScanner.Text()) != "moves" {
				continue
			}
			for posScanner.Scan() { // for each move
				moveStr := strings.ToLower(posScanner.Text())
				legalMoves := board.GenerateLegalMoves()
				var nextMove dragon.Move
				found := false
				for _, mv := range legalMoves {
					if mv.String() == moveStr {
						nextMove = mv
						found = true
						break
					}
				}
				if !found { // we didn't find the move, but we will try to apply it anyway
					fmt.Println("info string Move", moveStr, "not found for position", board.ToFen())
					var err error
					nextMove, err = dragon.ParseMove(moveStr)
					if err != nil {
						fmt.Println("info string Contingency move parsing failed")
						continue
					}
				}
				board.Apply(nextMove)
				//search.HistoryMap[board.Hash()]++
			}
		default:
			fmt.Println("info string Unknown command:", line)
		}
	}
}

// Lightweight wrapper around Lisao Search.
// Prints the results (bestmove). TODO PV, stats
// TODO - plumb timing and halt stuff properly
func uciSearch(board *dragon.Board, halt <-chan bool, stop *bool) {
	fmt.Println("info searching...")

	// Ignore timing and just call the fixed depth search
	bestMove, _ := engine.Search(board)

	fmt.Println("info got best move", &bestMove)

	// Wait for the stop signal and print the result
	//*stop = <-halt
	fmt.Println("bestmove", &bestMove)
}

// After a certain period of time, sends a signal to halt the search, unless it has already been sent.
// If the sleep time is 0, does nothing.
// The bool pointer alreadyStopped should be the same as the one given to Search().
func uciSearchTimeout(halt chan<- bool, ms int, alreadyStopped *bool) {
	if ms == 0 {
		return
	}
	fmt.Printf("info sleeping for %d ms\n", ms)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	fmt.Printf("info woke after %d ms\n", ms)
	if !(*alreadyStopped) { // don't send the halt signal if the search has already been stopped
		halt <- true
	}
}

// Simple strategy - use 1/16th of the remaining time (which we currently ignore anyway :D )
func uciCalculateAllowedTimeMs(b *dragon.Board, ourtimeMs int, opptimeMs int, ourincMs int, oppincMs int) int {
	result := ourtimeMs / 16
	if result <= 0 {
		return ourincMs
	}
	return result
}