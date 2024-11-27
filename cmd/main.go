package main

import (
	"fmt"
	"log"
	"os"

	formerfast "github.com/martcl/nrk-former/pkg/former-fast"
)

func main() {
	boardBytes, err := os.ReadFile("tests/27-11-2024.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	board, err := formerfast.LoadBoard(string(boardBytes))
	if err != nil {
		log.Fatalf("Error loading board: %v", err)
	}

	board.PrintBoard()

	heuristicTuning := 5.0
	solution := formerfast.SolveBoardUsingAStar(board, float32(heuristicTuning))
	fmt.Printf("\nSolution sequence, len %d:\n", len(solution))

	// Apply the solution to verify
	for i, pos := range solution {
		board.RemoveBricksIterative(pos)
		board.Gravity()
		fmt.Printf("%d. (x: %d, y:%d)\n", i, pos%7, (pos / 7))
		board.PrintBoard()
	}
}
