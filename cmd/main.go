package main

import (
	"fmt"

	formerfast "github.com/martcl/nrk-former/pkg/former-fast"
)

func main() {

	// Can use the current date to generate the board
	// but are to unstable since I can't figure out how
	// to predict the excat date. Better to just use the
	// seed we are given for now to generate board.
	// board := formerfast.CreateBoardFromDate(time.Now())

	// using seed to generate the board
	randomState := formerfast.InitializeRandomState("cff00d616484462eb325f50a5c0cd6a3")
	board, _ := formerfast.CreateBoardWithPseudoRandom(7, 9, randomState)

	// better to start high, then make it smaller. high ~ 6, low ~ 3
	heuristicTuning := 3.4
	numThreads := 12

	fmt.Printf("[info] Distance tuning variable: %f\n", heuristicTuning)
	fmt.Printf("[info] Number of threads: %d\n", numThreads)

	board.PrintBoard()

	solution := formerfast.SolveBoardUsingAStar(board, numThreads, float32(heuristicTuning))

	fmt.Printf("\nFound solution with length: %d\n", len(solution))

	// Apply the solution to verify
	for i, pos := range solution {
		// board.RemoveBricksIterative(pos)
		// board.Gravity()
		fmt.Printf("click %d. (x: %d, y:%d)\n", i, pos%7, (pos / 7))
		// board.PrintBoard()
	}
}
