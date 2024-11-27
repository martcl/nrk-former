package former

import (
	"fmt"
	"hash/fnv"

	"math/rand"
)

const (
	green = iota
	blue
	orange
	pink
)

type BrickType = int

type Brick struct {
	Type BrickType
}

type Board struct {
	Height int
	Width  int
	Bricks []*Brick
}

type Click struct {
	X int
	Y int
}

type ClickGroup struct {
	Click Click
	Size  int
	Type  BrickType
}

func (b Board) GetBrick(x int, y int) (*Brick, error) {
	if x > b.Width || y > b.Height {
		return nil, fmt.Errorf("coordinates out of bounds")
	}
	if x < 0 || y < 0 {
		return nil, fmt.Errorf("coordinates out of bounds")
	}

	index := y*(b.Width) + x

	if index < 0 || index >= len(b.Bricks) {
		return nil, fmt.Errorf("coordinates out of bounds")
	}
	if brick := b.Bricks[index]; brick != nil {
		return brick, nil
	}

	return nil, fmt.Errorf("no brick at pos (%d, %d)", x, y)
}

func (b Board) RemoveBrick(x int, y int) error {
	if x > b.Width || y > b.Height {
		return fmt.Errorf("coordinates out of bounds")
	}
	if x < 0 || y < 0 {
		return fmt.Errorf("coordinates out of bounds")
	}

	index := y*(b.Width) + x

	if index < 0 || index >= len(b.Bricks) {
		return fmt.Errorf("coordinates out of bounds")
	}

	b.Bricks[index] = nil

	return nil
}

func (b Board) Gravity() error {
	for x := 0; x < b.Width; x++ {
		stack := []*Brick{}
		for y := 0; y < b.Height; y++ {
			index := y*b.Width + x
			if b.Bricks[index] != nil {
				stack = append(stack, b.Bricks[index])
			}
		}

		for y := b.Height - 1; y >= 0; y-- {
			index := y*b.Width + x
			if len(stack) > 0 {
				b.Bricks[index] = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			} else {
				b.Bricks[index] = nil
			}
		}
	}
	return nil
}

func CreateBoard(height int, width int) (Board, error) {
	bricks := make([]*Brick, height*width)

	for i := 0; i < height*width; i++ {
		bricks[i] = &Brick{Type: BrickType(rand.Intn(4))}
	}

	return Board{
		Height: height,
		Width:  width,
		Bricks: bricks,
	}, nil
}

func (b Board) PrintBoard() {
	fmt.Println("--- Board ---")
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			index := y*b.Width + x
			if b.Bricks[index] != nil {
				fmt.Printf("%d ", b.Bricks[index].Type)
			} else {
				fmt.Print("# ")
			}
		}
		fmt.Println()
	}
	fmt.Println("--------------")
}

func GetPossibleSectorClicks(board *Board) []ClickGroup {
	var clicks []ClickGroup
	visited := make(map[int]bool)
	for x := 0; x < board.Width; x++ {
		for y := board.Height - 1; y >= 0; y-- {
			brick, err := board.GetBrick(x, y)
			if brick == nil || err != nil {
				continue
			}
			index := y*board.Width + x
			if visited[index] {
				continue
			}

			groupSize := board.GetConnectedBricks(x, y, brick.Type, visited)
			if groupSize > 0 {
				clicks = append(clicks, ClickGroup{Click{x, y}, groupSize, brick.Type})
			}
		}
	}
	return clicks
}

func (b *Board) Hash() uint32 {
	h := fnv.New32a()
	for _, brick := range b.Bricks {
		var val byte = 0
		if brick != nil {
			val = byte(brick.Type)
		}
		h.Write([]byte{val})
	}
	return h.Sum32()
}

func (b *Board) RemoveBricksIterative(x, y, brickType int) {
	stack := []int{x, y}
	visited := make(map[int]bool)

	for len(stack) > 0 {
		x, y = stack[len(stack)-2], stack[len(stack)-1]
		stack = stack[:len(stack)-2]

		index := y*b.Width + x
		if visited[index] {
			continue
		}

		brick, err := b.GetBrick(x, y)
		if err != nil || brick == nil || brick.Type != brickType {
			continue
		}

		visited[index] = true

		b.Bricks[index] = nil

		directions := [][2]int{
			{0, -1}, // Up
			{0, 1},  // Down
			{-1, 0}, // Left
			{1, 0},  // Right
		}
		for _, d := range directions {
			nx, ny := x+d[0], y+d[1]
			if nx >= 0 && nx < b.Width && ny >= 0 && ny < b.Height {
				stack = append(stack, nx, ny)
			}
		}
	}
}

func (b *Board) GetConnectedBricks(x, y, brickType int, visited map[int]bool) int {
	stack := []int{x, y}
	count := 0
	for len(stack) > 0 {
		x, y = stack[len(stack)-2], stack[len(stack)-1]
		stack = stack[:len(stack)-2]

		index := y*b.Width + x

		if visited[index] {
			continue
		}

		brick, err := b.GetBrick(x, y)
		if err != nil || brick == nil || brick.Type != brickType {
			continue
		}

		visited[index] = true
		count++

		directions := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		for _, d := range directions {
			nx, ny := x+d[0], y+d[1]
			if nx >= 0 && nx < b.Width && ny >= 0 && ny < b.Height {
				stack = append(stack, nx, ny)
			}
		}
	}
	return count
}

func (b *Board) Copy() *Board {
	newBoard := &Board{
		Height: b.Height,
		Width:  b.Width,
	}
	newBoard.Bricks = make([]*Brick, len(b.Bricks))
	copy(newBoard.Bricks, b.Bricks)
	return newBoard
}
