package main

import (
	"fmt"
	"time"

	formerfast "github.com/martcl/nrk-former/pkg/former-fast"
)

func main() {
	board := formerfast.CreateBoardFromDate(time.Now())

	// better to start high, then make it smaller. high ~ 6, low ~ 3
	heuristicTuning := 3.0
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
