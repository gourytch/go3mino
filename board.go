package go3mino

import (
	"fmt"
	"math/rand"
)

type Board struct {
	Gap    int // size of empty border around the placed pieces
	Width  int
	Height int
	Map    []Trimino // Width * Height
}

func NewBoard(gap int) *Board {
	sz := 2*gap + 1
	board := &Board{
		Gap:    gap,
		Width:  sz,
		Height: sz,
		Map:    make([]Trimino, sz*sz),
	}
	for i := 0; i < board.Height*board.Width; i++ {
		board.Map[i] = FreeTrimino()
	}
	return board
}

func (board *Board) GetPieces() Triminos {
	result := Triminos{}
	for _, piece := range board.Map {
		if !piece.IsFree() {
			result = append(result, piece.Normalize())
		}
	}
	return result
}

func (board *Board) GetNumPieces() int {
	result := 0
	for _, piece := range board.Map {
		if !piece.IsFree() {
			result++
		}
	}
	return result
}

func (board *Board) IsEmpty() bool {
	for _, piece := range board.Map {
		if !piece.IsFree() {
			return false
		}
	}
	return true
}

func (board *Board) Adjust() {
	var min_x int
	var max_x int
	var min_y int
	var max_y int
	found := false
	for y := 0; y < board.Height; y++ {
		for x := 0; x < board.Width; x++ {
			p := board.Map[y*board.Width+x]
			if !p.IsFree() {
				if !found {
					min_x = x
					max_x = x
					min_y = y
					max_y = y
					found = true
				} else {
					if x < min_x {
						min_x = x
					} else if x > max_x {
						max_x = x
					}
					if y < min_y {
						min_y = y
					} else if y > max_y {
						max_y = y
					}
				}
			}
		}
	}
	if !found {
		fmt.Println("Adjust: empty board")
		return
	}
	fmt.Printf("Adjust: rows=[%d..%d], cols=[%d..%d]\n", min_y, max_y, min_x, max_x)
	new_height := board.Height
	new_width := board.Width
	dx := 0
	dy := 0
	// check left side
	d := board.Gap - min_x
	if d > 0 {
		if d%2 != 0 { // must be even
			d++
		}
		dx = d
		new_width += d
		fmt.Printf("Adjust: move right to %d, new width=%d\n", dx, new_width)
	}
	// check right side
	d = board.Gap - (board.Width - 1 - max_x)
	if d > 0 {
		new_width += d
		fmt.Printf("expand right side to %d, to %d\n", d, new_width)
		if new_width%2 == 0 { // must be odd
			new_width++
			fmt.Printf("... and adjust again, to %d\n", new_width)
		}
	}
	// check top side
	d = board.Gap - min_y
	if d > 0 {
		if d%2 != 0 { // must be even
			d++
		}
		dy = d
		new_height += d
		fmt.Printf("Adjust: move bottom to %d, new height=%d\n", dy, new_height)
	}
	// check bottom side
	d = board.Gap - (board.Height - 1 - max_y)
	if d > 0 {
		new_height += d
	}
	if new_height == board.Height && new_width == board.Width {
		fmt.Println("Adjust is not needed")
		return
	}
	fmt.Printf("Adjust: [%d,%d] -> [%d,%d]\n", board.Height, board.Width, new_height, new_width)
	new_map := make([]Trimino, new_height*new_width)
	// FIXME rewrite next and eliminate a double assignment
	for i := 0; i < new_height*new_width; i++ {
		new_map[i] = FreeTrimino()
	}
	for y := 0; y < board.Height; y++ {
		for x := 0; x < board.Width; x++ {
			new_map[(y+dy)*new_width+(x+dx)] = board.Map[y*board.Width+x]
		}
	}
	board.Height = new_height
	board.Width = new_width
	board.Map = new_map
}

func (board *Board) GetPiece(row, col int) Trimino {
	if row < 1 || row >= board.Height-2 {
		return FreeTrimino()
	}
	if col < 1 || col >= board.Width-2 {
		return FreeTrimino()
	}
	pos := row*board.Width + col
	return board.Map[pos]
}

func (board *Board) CanPlace(row, col int, piece Trimino) (ret bool) {
	// defer func() {
	// 	fmt.Printf("CanPlace(row=%d, col=%d, piece=%v) returned %v\n", row, col, piece, ret)
	// }()
	if row < 0 || row >= board.Height {
		return false
	}
	if col < 0 || col >= board.Width {
		return false
	}
	p := board.GetPiece(row, col)
	if !p.IsFree() {
		return false
	}
	if row%2 == col%2 { // v-shape
		/*
		   V-shape: even(row) == even(col):
		              ^
		             / \
		            /   \
		           /-1;+0\
		           ~~~~~~~
		         ^ _______ ^
		        / \\+0;+0// \
		       /   \\   //   \
		      /+0;-1\\ //+0;+1\
		      ~~~~~~~ v ~~~~~~~
		      E0 vs Cell[-1;+0]E1'
		      E1 vs Cell[+0;+1]E2'
		      E2 vs Cell[+0;-1]E0'

		           0      U: [0,1,5]
		          5 1     C: [1,3,5]->[5,1,3]
		        5 5 1 1   L: [3,4,5]->[5,3,4]
		       4 3 3 3 2  R: [1,2,3]
		   by node:
		       C[0] == U[2] == L[0]
		       C[1] == U[1] == R[0]
		       C[2] == L[1] == R[2]
		*/
		l := board.GetPiece(row, col-1)
		r := board.GetPiece(row, col+1)
		u := board.GetPiece(row-1, col)
		if l.IsFree() && r.IsFree() && u.IsFree() {
			return false
		}
		// fmt.Printf("check piece %v at V-pos %d;%d vs L=%v R=%v U=%v\n", piece, row, col, l, r, u)

		if !u.IsFree() {
			if piece[0] != u[2] || piece[1] != u[1] {
				return false
			}
		}
		if !l.IsFree() {
			if piece[0] != l[0] || piece[2] != l[1] {
				return false
			}
		}
		if !r.IsFree() {
			if piece[1] != r[0] || piece[2] != r[2] {
				return false
			}
		}
		// fmt.Printf("... match!\n")
		return true
	} else {
		/*
			        A-shape: even(row) != even(col):
			               ______ ^ _______
			              \+0;-1// \\+0;+1/
			               \   //   \\   /
			                \ //+0;+0\\ /
			                 v ~~~~~~~ v
			                   _______
			                   \+0;+1/
			                    \   /
			                     \ /
			                      v
			          E0 vs Cell[+0;+1]E2'
			          E1 vs Cell[+1;+0]E0'
			          E2 vs Cell[+0;-1]E1'
				      4 5 5 5 0   L: [3,4,5] -> [4,5,3]
				       3 3 1 1    C: [1,3,5] -> [5,1,3]
				         3 1      R: [0,1,5] -> [5,0,1]
				          2       D: [1,2,3] -> [3,1,2]
				      by node:
				       C[0] == L[1] == R[0]
				       C[1] == R[2] == D[1]
				       C[2] == L[2] == D[0]
		*/
		l := board.GetPiece(row, col-1)
		r := board.GetPiece(row, col+1)
		d := board.GetPiece(row+1, col)
		if l.IsFree() && r.IsFree() && d.IsFree() {
			return false
		}
		// fmt.Printf("check piece %v at A-pos %d;%d vs L=%v, R=%v, D=%v\n", piece, row, col, l, r, d)
		if !l.IsFree() {
			if piece[0] != l[1] || piece[2] != l[2] {
				return false
			}
		}
		if !r.IsFree() {
			if piece[0] != r[0] || piece[1] != r[2] {
				return false
			}
		}
		if !d.IsFree() {
			if piece[1] != d[1] || piece[2] != d[0] {
				return false
			}
		}
		// fmt.Printf("... match!\n")
		return true
	}
}

func (board *Board) PlaceFirst(piece Trimino) {
	pos := (board.Height/2)*board.Width + (board.Width / 2)
	board.Map[pos] = piece
}

func (board *Board) Place(row, col int, piece Trimino) {
	if !board.CanPlace(row, col, piece) {
		panic(fmt.Sprintf("piece %v does not fit at pos (%d,%d)", piece, row, col))
	}
	if row < 0 || row >= board.Height {
		panic(fmt.Sprintf("row=%d is out of bounds [0..%d)", row, board.Height))
	}
	if col < 0 || col >= board.Width {
		panic(fmt.Sprintf("col=%d is out of bounds [0..%d)", col, board.Width))
	}
	pos := row*board.Width + col
	if !board.Map[pos].IsFree() {
		panic(fmt.Sprintf("position[row=%d,col=%d is not free", row, col))
	}
	board.Map[pos] = piece
	board.Adjust()
}

type Position struct {
	Row int
	Col int
}

type Positions []Position

func (board *Board) ScanForPlaces(piece Trimino) Positions {
	result := Positions{}
	for row := 0; row < board.Height; row++ {
		for col := 0; col < board.Width; col++ {
			if board.CanPlace(row, col, piece) {
				result = append(result, Position{row, col})
			}
		}
	}
	return result
}

type CheckedStep struct {
	Piece     Trimino
	Positions Positions
}

type CheckedSteps [3]CheckedStep

func (board *Board) CheckAll(piece Trimino) CheckedSteps {
	c := piece.Normalize()
	var result CheckedSteps
	for i := 0; i < 3; i++ {
		result[i].Piece = c
		result[i].Positions = board.ScanForPlaces(c)
		c = c.Rotate()
	}
	return result
}

func (board *Board) JustPutIt(piece Trimino) bool {
	steps := board.CheckAll(piece)
	num := 0
	for i := 0; i < len(steps); i++ {
		num += len(steps[i].Positions)
	}
	if num == 0 {
		fmt.Printf("no available positions for piece %v\n", piece)
		return false
	}
	choice := rand.Intn(num)
	fmt.Printf("%d available positions for piece %v, choose #%d\n", num, piece, choice)
	i := 0
	for choice >= len(steps[i].Positions) {
		choice -= len(steps[i].Positions)
		i++
	}
	fmt.Printf("... place %v at %v\n", steps[i].Piece, steps[i].Positions[choice])
	board.Place(steps[i].Positions[choice].Row, steps[i].Positions[choice].Col, steps[i].Piece)
	return true
}
