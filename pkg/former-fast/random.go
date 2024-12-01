package formerfast

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"time"
)

func CreateHashFunction() func(string) float64 {
	var seed float64 = 4022871197

	return func(input string) float64 {
		for _, chr := range input {
			seed += float64(chr)
			product := 0.02519603282416938 * seed

			seed = float64(uint32(product))
			product -= float64(seed)

			product *= float64(seed)

			seed = float64(uint32(product))

			product -= float64(seed)
			seed += float64(uint32(product * 4294967296))

		}
		return float64(uint32(seed)) * 2.3283064365386963e-10
	}
}

type RandomState struct {
	state1      float64
	state2      float64
	fraction    float64
	integerPart int64
}

func (r *RandomState) Next() float64 {
	const multiplier = 2091639.0
	const norm = 2.3283064365386963e-10 // 1 / 2^32

	product := multiplier*r.state1 + float64(r.integerPart)*norm
	r.state1 = r.state2
	r.state2 = r.fraction
	r.fraction = product - float64(int64(product)) // Fractional part
	r.integerPart = int64(product)                 // Integer part as the new carry
	return r.fraction
}

func InitializeRandomState(seedString string) RandomState {
	hashFunction := CreateHashFunction()

	state1 := hashFunction(" ")
	state2 := hashFunction(" ")
	fraction := hashFunction(" ")
	integerPart := int64(1)

	state1 -= hashFunction(seedString)
	if state1 < 0 {
		state1 += 1
	}
	state2 -= hashFunction(seedString)
	if state2 < 0 {
		state2 += 1
	}
	fraction -= hashFunction(seedString)
	if fraction < 0 {
		fraction += 1
	}

	return RandomState{
		state1,
		state2,
		fraction,
		integerPart,
	}
}

func md5Sum(date string) string {
	hash := md5.Sum([]byte(date))
	return hex.EncodeToString(hash[:])
}

func CreateBoardFromDate(date time.Time) *Board {
	const offset = -31

	// for some reason they have the dates mixed up
	// and have some sort of offset. The API said it
	// was -60, but that is not true. I think it is more
	// like -31 days from the current date.
	// TODO: wait and see how the seed changes over time to verify

	newDate := date.AddDate(0, 0, offset)
	dateString := newDate.Format("02012006")

	randomState := InitializeRandomState(md5Sum(dateString))

	board, _ := CreateBoardWithPseudoRandom(7, 9, randomState)

	return board
}

func CreateBoardWithPseudoRandom(height int, width int, randomState RandomState) (*Board, error) {
	bricks := [4]uint64{}

	colorMap := map[int]int{
		0: orange,
		1: pink,
		2: green,
		3: blue,
	}

	for pos := 0; pos < height*width; pos++ {
		color := colorMap[int(math.Floor(randomState.Next()*4))]
		bricks[color] |= uint64(1) << pos
	}

	return &Board{
		State: bricks,
	}, nil
}
