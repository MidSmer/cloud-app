package common

import (
	"math/big"
	"math/rand"
	"time"
)

// second 17bit      1  0  0  1  1  1  0  0  1  1  0  1  0  0  1  0  1
// day              1  0  1  0  1  0  0  0  1  0  1  1  0  1  1  0  0
// rand            1  1  0  1  0  0  1  0  1  0  1  0  0  0  1  1  0 

const (
	initSecond int64 = 1630000000 // start date：2021-08-27
)

var (
	charList = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '-', '_'}
)

func NewShortId() (string, error) {
	totalSec := time.Now().Unix() - initSecond
	day := totalSec / (24 * 60 * 60)
	sec := totalSec % (24 * 60 * 60)

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Uint64()

	var shortId []byte
	for i := 0; i < 9; i++ {
		var bits big.Int
		for o := 0; o < 6; o++ {
			var bit uint
			index := i*2 + (o / 3)

			if o%3 == 0 {
				bit = uint((sec >> index) & 1)
			} else if o%3 == 1 {
				bit = uint((day >> index) & 1)
			} else if o%3 == 2 {
				bit = uint((randNum >> index) & 1)
			}

			bits.SetBit(&bits, o, bit)
		}

		shortId = append([]byte{charList[bits.Int64()]}, shortId...)
	}

	return string(shortId), nil
}
