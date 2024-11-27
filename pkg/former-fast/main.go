package formerfast

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math/rand"
)

const (
	orange = iota
	green
	pink
	blue

	empty // does not have a board state, but represents 0 in any board state pos
)

type BrickType = uint8

// We store the board state for each color in a uint64 since
// we know that the board is always 7*9 = 63 size. If a bit
// in the board is 0, there is not that color there and if the
// bit is 1 there is that color there. If nether one of the colors
// have a 1 bit in that position, then there is nothing there.
type Board struct {
	State [4]uint64
}

// the pos is a single uint64 type and we use that to find out
// witch color is there
func (b *Board) GetBrick(pos uint8) (BrickType, error) {
	if pos >= 63 {
		return empty, fmt.Errorf("position out of bounds")
	}
	mask := uint64(1) << pos
	if b.State[green]&mask != 0 {
		return green, nil
	} else if b.State[blue]&mask != 0 {
		return blue, nil
	} else if b.State[orange]&mask != 0 {
		return orange, nil
	} else if b.State[pink]&mask != 0 {
		return pink, nil
	} else {
		return empty, fmt.Errorf("no brick at position")
	}
}

func (b *Board) RemoveBricksIterative(pos uint8) {
	stack := []uint8{pos}
	visited := make(map[uint8]bool)

	brickType, _ := b.GetBrick(pos)

	for len(stack) > 0 {
		pos = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[pos] {
			continue
		}

		brick, err := b.GetBrick(pos)
		if err != nil || brick == empty || brick != brickType {
			continue
		}

		visited[pos] = true

		mask := uint64(1) << pos
		switch brick {
		case blue:
			b.State[blue] &= ^mask
		case green:
			b.State[green] &= ^mask
		case orange:
			b.State[orange] &= ^mask
		case pink:
			b.State[pink] &= ^mask
		}

		directions := []int64{-7, -1, 1, 7} // UP, LEFT, RIGHT, DOWN
		for _, d := range directions {
			neighbor := int64(pos) + d
			if neighbor < 0 || neighbor >= 63 {
				continue
			}
			if (pos%7 == 0 && d == -1) || (pos%7 == 6 && d == 1) {
				continue
			}
			stack = append(stack, uint8(neighbor))
		}
	}
}

func (b *Board) GetPossibleClicks() []uint8 {
	clicks := []uint8{}
	visited := make(map[uint8]bool)

	for pos := uint8(7*9) - 1; pos != 0; pos-- {
		if visited[pos] {
			continue
		}

		brick, _ := b.GetBrick(pos)

		if brick == empty {
			visited[pos] = true
			continue
		}

		b.MarkConnectedBricks(pos, brick, visited)
		clicks = append(clicks, pos)
	}
	return clicks
}

func (b *Board) MarkConnectedBricks(pos uint8, brickType BrickType, visited map[uint8]bool) {
	stack := []uint8{pos}

	for len(stack) > 0 {
		pos = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[pos] {
			continue
		}

		brick, err := b.GetBrick(pos)
		if err != nil || brick == empty || brick != brickType {
			continue
		}

		visited[pos] = true

		directions := []int64{-7, -1, 1, 7} // UP, LEFT, RIGHT, DOWN
		for _, d := range directions {
			neighbor := int64(pos) + d
			if neighbor < 0 || neighbor >= 63 {
				continue
			}
			if (pos%7 == 0 && d == -1) || (pos%7 == 6 && d == 1) {
				continue
			}
			stack = append(stack, uint8(neighbor))
		}
	}
}

// TODO optemize this
func (b *Board) Gravity() {
	for x := 0; x < 7; x++ {
		// Collect bricks in the column
		stack := []BrickType{}
		for y := 0; y < 9; y++ {
			pos := uint8(y*7 + x)
			brick, err := b.GetBrick(pos)
			if err == nil && brick != empty {
				stack = append(stack, brick)
			}
		}

		// Clear the column
		for y := 0; y < 9; y++ {
			pos := uint8(y*7 + x)
			mask := uint64(1) << pos
			for c := orange; c <= blue; c++ {
				b.State[c] &= ^mask
			}
		}

		// Place bricks from the bottom up
		yIndex := 8
		for i := len(stack) - 1; i >= 0; i-- {
			pos := uint8(yIndex*7 + x)
			brick := stack[i]
			mask := uint64(1) << pos
			b.State[brick] |= mask
			yIndex--
		}
	}
}

// Hash generates a unique hash for the current board state.
// It includes all 64 bits of each uint64 in the State array.
func (b *Board) Hash() uint32 {
	h := fnv.New32a()

	buf := make([]byte, 0, 32)

	for _, state := range b.State {
		temp := make([]byte, 8)
		binary.LittleEndian.PutUint64(temp, state)
		buf = append(buf, temp...)
	}

	h.Write(buf)

	return h.Sum32()
}

func (b *Board) Copy() *Board {
	newBoard := &Board{
		State: b.State,
	}
	return newBoard
}

func (b *Board) PrintBoard() {
	fmt.Println("--- Board ---")
	brickSymbols := map[BrickType]string{
		orange: "O",
		green:  "G",
		pink:   "P",
		blue:   "B",
	}
	for y := 0; y < 9; y++ {
		for x := 0; x < 7; x++ {
			pos := uint8(y*7 + x)
			brick, err := b.GetBrick(pos)
			if err == nil && brick != empty {
				fmt.Printf("%s ", brickSymbols[brick])
			} else {
				fmt.Print("# ")
			}
		}
		fmt.Println()
	}
	fmt.Println("--------------")
}

func CreateRadomBoard(height int, width int) (*Board, error) {
	bricks := [4]uint64{}

	for pos := 0; pos < height*width; pos++ {
		color := rand.Intn(4)
		bricks[color] |= uint64(1) << pos
	}

	return &Board{
		State: bricks,
	}, nil
}

func (board *Board) isBoardEmpty() bool {
	return board.State[orange]|board.State[green]|board.State[pink]|board.State[blue] == 0
}
