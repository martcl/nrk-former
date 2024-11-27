package formerfast

import (
	"encoding/json"

	"github.com/martcl/nrk-former/pkg/former"
)

func LoadBoard(jsonData string) (*Board, error) {
	var data [][]former.GemData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}

	board := &Board{}

	gemColorToType := map[string]BrickType{
		"sirkel":   pink,
		"firkant":  blue,
		"pil":      green,
		"diagonal": orange,
	}

	for y, row := range data {
		for x, gem := range row {
			if gem.IsEmpty {
				continue
			}
			brickType := gemColorToType[gem.GemColor]
			pos := uint8(y*7 + x)

			mask := uint64(1) << pos
			board.State[brickType] |= mask
		}
	}

	return board, nil
}
