package former

import (
	"encoding/json"
)

type Scale struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Origin struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Sprite struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	X          float64  `json:"x"`
	Y          float64  `json:"y"`
	Depth      int      `json:"depth"`
	Scale      Scale    `json:"scale"`
	Origin     Origin   `json:"origin"`
	FlipX      bool     `json:"flipX"`
	FlipY      bool     `json:"flipY"`
	Rotation   int      `json:"rotation"`
	Alpha      int      `json:"alpha"`
	Visible    bool     `json:"visible"`
	BlendMode  int      `json:"blendMode"`
	TextureKey string   `json:"textureKey"`
	FrameKey   string   `json:"frameKey"`
	Data       struct{} `json:"data"`
}

type GemData struct {
	GemColor   string `json:"gemColor"`
	Sprite     Sprite `json:"sprite"`
	IsEmpty    bool   `json:"isEmpty"`
	TintColor  int    `json:"tintColor"`
	BlastColor int    `json:"blastColor"`
}

func LoadBoard(jsonData string) (*Board, error) {
	var data [][]GemData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}

	height := len(data)
	width := len(data[0])
	board := &Board{
		Width:  width,
		Height: height,
		Bricks: make([]*Brick, width*height),
	}

	gemColorToType := map[string]int{
		"sirkel":   pink,
		"firkant":  blue,
		"pil":      green,
		"diagonal": orange,
	}

	for y, row := range data {
		for x, gem := range row {
			index := y*board.Width + x
			if gem.IsEmpty {
				board.Bricks[index] = nil
			} else {
				brickType := gemColorToType[gem.GemColor]
				board.Bricks[index] = &Brick{
					Type: brickType,
				}
			}
		}
	}

	return board, nil
}
