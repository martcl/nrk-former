package main

import (
	"fmt"
	"log"
	"os"

	"github.com/martcl/nrk-former/pkg/former"
)

func main() {

	boardBytes, err := os.ReadFile("tests/27-11-2024.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	board, err := former.LoadBoard(string(boardBytes))
	if err != nil {
		log.Fatalf("Error loading board: %v", err)
	}

	board.PrintBoard()

	heuristicTuning := 2.3
	solution := former.SolveBoardUsingAStar(board, heuristicTuning)
	fmt.Printf("\nSolution sequence, len %d:\n", len(solution))

	// Apply the solution to verify
	for i, click := range solution {
		// board.RemoveBricksIterative(click.X, click.Y, board.Bricks[click.Y*board.Width+click.X].Type)
		// board.Gravity()
		fmt.Printf("%d. (x:%d, y:%d)\n", i, click.X, click.Y)
		// board.PrintBoard()

	}

}
