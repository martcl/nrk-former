package formerfast

import (
	"container/heap"
	"fmt"
	"math"
)

type State struct {
	Board    *Board
	Moves    []uint8 // Sequence of moves to reach this state
	Estimate float32 // h: Heuristic value
	Priority float32 // f: Steps + Estimate
}

type PriorityQueue []*State

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// Smaller Priority value means higher priority
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*State))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func SolveBoardUsingAStar(board *Board, heuristicTuning float32) []uint8 {
	pq := &PriorityQueue{}

	initialState := &State{
		Board:    board.Copy(),
		Moves:    []uint8{},
		Estimate: board.heuristic(heuristicTuning),
		Priority: board.heuristic(heuristicTuning),
	}

	heap.Init(pq)
	heap.Push(pq, initialState)

	i := 0
	for pq.Len() > 0 {
		current := heap.Pop(pq).(*State)

		// Goal check: is the board empty?
		if current.Board.isBoardEmpty() {
			fmt.Printf("Solved with %d iterations", i)
			return current.Moves
		}

		possibleClicks := current.Board.GetPossibleClicks()

		if i%10000 == 0 && i != 0 {
			fmt.Printf("Iteration %d, moves: %d, estimate: %f, remaining: %d. Heapsize: %d\n", i, len(current.Moves), current.Estimate, len(possibleClicks), len(*pq))
			current.Board.PrintBoard()
		}

		for _, pos := range possibleClicks {
			nextBoard := current.Board.Copy()
			nextBoard.RemoveBricksIterative(pos)
			nextBoard.Gravity()

			estimatedDistance := nextBoard.heuristic(heuristicTuning)

			nextMoves := append([]uint8{}, current.Moves...)
			nextMoves = append(nextMoves, pos)

			nextState := &State{
				Board:    nextBoard,
				Moves:    nextMoves,
				Estimate: estimatedDistance,
				Priority: float32(len(nextMoves)) + estimatedDistance,
			}

			heap.Push(pq, nextState)
		}

		i++
	}

	return nil
}

// Estimates amount of clicks needed to win the game from the current game state
// The assumptions is that when there a lot of options to choose from many of them
// are not real clicks you need to click since blocks merge over time.
// The fewer click options we have, the higher chache is it that the click is real.
// TODO: improve this to get better estimate to goal
func (board *Board) heuristic(heuristicTuning float32) float32 {
	possibleClicks := board.GetPossibleClicks()
	return float32(math.Log(float64(len(possibleClicks)))) * heuristicTuning
}
