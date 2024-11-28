package formerfast

import (
	"container/heap"
	"math"
	"sync"
)

type State struct {
	Board    *Board
	Moves    []uint8 // Sequence of moves to reach this state
	Estimate float32 // h: Heuristic value
	Priority float32 // f: Steps + Estimate
}

type PriorityQueue []*State

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*State))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	state := old[n-1]
	*pq = old[0 : n-1]
	return state
}

// Thread-safe priority queue
type SafePriorityQueue struct {
	pq    PriorityQueue
	mutex sync.Mutex
}

func (spq *SafePriorityQueue) Len() int {
	spq.mutex.Lock()
	defer spq.mutex.Unlock()
	return len(spq.pq)
}

func (spq *SafePriorityQueue) Push(state *State) {
	spq.mutex.Lock()
	heap.Push(&spq.pq, state)
	spq.mutex.Unlock()
}

func (spq *SafePriorityQueue) Pop() *State {
	spq.mutex.Lock()
	defer spq.mutex.Unlock()
	if len(spq.pq) > 0 {
		return heap.Pop(&spq.pq).(*State)
	}
	return nil
}

func SolveBoardUsingAStar(board *Board, maxThreads int, heuristicTuning float32) []uint8 {
	spq := &SafePriorityQueue{pq: make(PriorityQueue, 0)}
	heap.Init(&spq.pq)

	initialState := &State{
		Board:    board.Copy(),
		Moves:    []uint8{},
		Estimate: board.heuristic(heuristicTuning),
		Priority: board.heuristic(heuristicTuning),
	}

	var wg sync.WaitGroup
	var returnOnce sync.Once

	numThreadsSemephore := make(chan struct{}, maxThreads)

	result := make(chan []uint8, 1)
	done := make(chan struct{})

	spq.Push(initialState)

	for {
		select {
		case <-done:
			return <-result
		default:
			current := spq.Pop()

			if current == nil {
				isQueueEmpty := spq.Len() == 0
				if !isQueueEmpty {
					continue
				}
				wg.Wait() // wait for unfinnished threads, and check queue again
				isQueueStillEmpty := spq.Len() == 0
				// if queue is empty we did not find a solution
				if isQueueStillEmpty {
					return nil
				}
				// We are consuming items to fast from the queue,
				// and need to continue
				continue
			}

			numThreadsSemephore <- struct{}{}
			wg.Add(1)
			go func(state *State) {
				defer wg.Done()
				defer func() { <-numThreadsSemephore }()

				if state.Board.isBoardEmpty() {
					result <- state.Moves
					returnOnce.Do(func() { close(done) })
					return
				}

				possibleClicks := state.Board.GetPossibleClicks()

				for _, pos := range possibleClicks {
					nextBoard := state.Board.Copy()
					nextBoard.RemoveBricksIterative(pos)
					nextBoard.Gravity()

					estimatedDistance := nextBoard.heuristic(heuristicTuning)

					nextMoves := append([]uint8{}, state.Moves...)
					nextMoves = append(nextMoves, pos)

					nextState := &State{
						Board:    nextBoard,
						Moves:    nextMoves,
						Estimate: estimatedDistance,
						Priority: float32(len(nextMoves)) + estimatedDistance,
					}

					spq.Push(nextState)
				}
			}(current)
		}

	}
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
