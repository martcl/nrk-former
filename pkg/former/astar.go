package former

import (
	"container/heap"
	"fmt"
	"math"
)

type State struct {
	Board     *Board
	Moves     []Click // Sequence of moves to reach this state
	Steps     int     // g: Number of steps taken
	Estimate  float64 // h: Heuristic value
	Priority  float64 // f: Steps + Estimate
	StateHash uint32  // Unique identifier for the board state
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

var possibleClickCache = map[uint32][]ClickGroup{}

// heuristic_tuning: good 6 - 3
func SolveBoardUsingAStar(board *Board, heuristicTuning float64) []Click {
	initialHash := board.Hash()
	pq := &PriorityQueue{}

	initialState := &State{
		Board:     board.Copy(),
		Moves:     []Click{},
		Steps:     0,
		Estimate:  heuristic(board, heuristicTuning),
		Priority:  heuristic(board, heuristicTuning),
		StateHash: initialHash,
	}

	heap.Init(pq)
	heap.Push(pq, initialState)

	i := 0
	for pq.Len() > 0 {

		current := heap.Pop(pq).(*State)

		// Goal check: is the board empty?
		if isBoardEmpty(current.Board) {
			return current.Moves
		}

		var clickGroups []ClickGroup
		if newGroup, exists := possibleClickCache[current.Board.Hash()]; exists {
			clickGroups = newGroup
		} else {
			clickGroups = GetPossibleSectorClicks(current.Board)
			possibleClickCache[board.Hash()] = clickGroups
		}

		if i%10000 == 0 {
			fmt.Printf("Iteration %d, moves: %d, esitmate: %f, remainding: %d\n", i, current.Steps, current.Estimate, len(clickGroups))
		}

		for _, clickGroup := range clickGroups {
			nextBoard := current.Board.Copy()
			nextBoard.RemoveBricksIterative(clickGroup.Click.X, clickGroup.Click.Y, clickGroup.Type)
			nextBoard.Gravity()

			nextHash := nextBoard.Hash()

			heuristic := heuristic(nextBoard, heuristicTuning)

			// Add cut off huristic to save memmory
			// if current.Steps+heuristic > current.Priority {
			// 	continue
			// }

			nextMoves := append([]Click{}, current.Moves...)
			nextMoves = append(nextMoves, clickGroup.Click)

			nextState := &State{
				Board:     nextBoard,
				Moves:     nextMoves,
				Steps:     current.Steps + 1,
				Estimate:  heuristic,
				Priority:  float64(current.Steps) + 1 + heuristic,
				StateHash: nextHash,
			}

			heap.Push(pq, nextState)
		}
		i++
	}

	return nil
}

// Estimate the number of moves to clear the board
// Since we don't know exactly what it takes to clear the board
// we need to create a good function to guess the distance
// to the goal. Current function tries to simulate that there will
// be less clicks needed if there are more click options and more
// percentage of clicks are real if there are few options.
// Improving this will make the program better!
func heuristic(board *Board, heuristicTuning float64) float64 {
	if clicks, exists := possibleClickCache[board.Hash()]; exists {
		return math.Log(float64((len(clicks) - 1))) * heuristicTuning
	}
	clickGroups := GetPossibleSectorClicks(board)
	possibleClickCache[board.Hash()] = clickGroups
	return math.Log(float64((len(clickGroups) - 1))) * heuristicTuning
}

func isBoardEmpty(board *Board) bool {
	for _, brick := range board.Bricks {
		if brick != nil {
			return false
		}
	}
	return true
}
