package go3mino

import (
	"fmt"
	"math/rand"
)

/*************************************************************
 * Triominos:
 * [0,0,0], [0,0,1], [0,0,2], [0,0,3], [0,0,4], [0,0,5],
 * [0,1,1], [0,1,2], [0,1,3], [0,1,4], [0,1,5],
 * [0,2,2], [0,2,3], [0,2,4], [0,2,5],
 * [0,3,3], [0,3,4], [0,3,5],
 * [0,4,4], [0,4,5],
 * [0,5,5],
 *
 * [1,1,1], [1,1,2], [1,1,3], [1,1,4], [1,1,5],
 * [1,2,2], [1,2,3], [1,2,4], [1,2,5],
 * [1,3,3], [1,3,4], [1,3,5],
 * [1,4,4], [1,4,5],
 * [1,5,5],
 *
 * [2,2,2], [2,2,3], [2,2,4], [2,2,5],
 * [2,3,3], [2,3,4], [2,3,5],
 * [2,4,4], [2,4,5],
 * [2,5,5],
 *
 * [3,3,3], [3,3,4], [3,3,5],
 * [3,4,4], [3,4,5],
 * [3,5,5],
 *
 * [4,4,4], [4,4,5],
 * [4,5,5],
 *
 * [5,5,5]
 *
 */

const FREE_CELL = -1

type Trimino [3]int
type Triminos []Trimino

type Side [2]int

// type Triomino struct {
// 	Nodes   [3]int // [0,0,0] .. [5,5,5]
// 	Variant int    // [0,1,2]^, [0,1,2]v, [2,0,1]^, [2,0,1]v, [1,2,0]^, [1,2,0]v
// 	// OffsetRow int
// 	// OffsetCol int
// 	// Links     [3]*Triomino
// }

func AllTriminos() Triminos {
	result := Triminos{}
	for i := 0; i < 6; i++ {
		for j := i; j < 6; j++ {
			for k := j; k < 6; k++ {
				result = append(result, Trimino{i, j, k})
			}
		}
	}
	return result
}

func (t Triminos) Pick() (Triminos, Trimino, bool) {
	N := len(t)
	if N == 0 {
		return t, FreeTrimino(), false
	}
	choice := rand.Intn(N)
	pick := t[choice]
	return append(t[:choice], t[choice+1:]...), pick, true
}

func FreeTrimino() Trimino {
	return Trimino{FREE_CELL, FREE_CELL, FREE_CELL}
}

func (t Trimino) IsFree() bool {
	return t[0] == FREE_CELL
}

func (t Trimino) IsNormalized() bool {
	return t[0] <= t[1] && t[1] <= t[2]
}

func (t Trimino) Rotate() Trimino {
	return Trimino{t[1], t[2], t[0]}
}

func (t Trimino) RotateL() Trimino {
	return Trimino{t[2], t[0], t[1]}
}

func (t Trimino) Normalize() Trimino {
	q := t
	for !q.IsNormalized() {
		q = q.Rotate()
	}
	return q
}

func (t Trimino) GetNode(ix int) int {
	switch ix {
	case 0, 1, 2:
		return t[ix]
	default:
		panic(fmt.Sprintf("GetSide(%d) bad side number", ix))
	}
}

func (t Trimino) GetSide(ix int) Side {
	switch ix {
	case 0:
		return Side{t[0], t[1]}
	case 1:
		return Side{t[1], t[2]}
	case 2:
		return Side{t[2], t[0]}
	default:
		panic(fmt.Sprintf("GetSide(%d) bad side number", ix))
	}
}

func (s Side) IsFree() bool {
	return s[0] == FREE_CELL
}

func (s Side) Reverse() Side {
	return Side{s[1], s[0]}
}

func (s Side) IsAligned(v Side) bool {
	if s.IsFree() || v.IsFree() {
		return true
	}
	return s[0] == v[0] && s[1] == v[1]
}
