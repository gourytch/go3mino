package main

import (
	"fmt"
	"os"

	"github.com/gourytch/go3mino"
)

func Step(board *go3mino.Board, piece go3mino.Trimino) {
	if !board.JustPutIt(piece) {
		fmt.Printf("piece %v cannot be placed\n", piece)
		os.Exit(0)
	}
	fmt.Printf("piece %v placed\n", piece)
	fmt.Println(board)
	fmt.Println()
}

func test() {
	board := go3mino.NewBoard(2)
	fmt.Println("empty board")
	fmt.Println(board)
	// board is 5x5
	/*
	      0      U: [0,1,5]
	     5 1     C: [1,3,5]->[5,1,3]
	   5 5 1 1   L: [3,4,5]->[5,3,4]
	  4 3 3 3 2  R: [1,2,3]
	*/
	C := go3mino.Trimino{1, 3, 5}
	U := go3mino.Trimino{0, 1, 5}
	L := go3mino.Trimino{3, 4, 5}
	R := go3mino.Trimino{1, 2, 3}
	board.PlaceFirst(C.RotateL())
	fmt.Println("first piece placed.")
	fmt.Println(board)
	Step(board, U)
	Step(board, L)
	Step(board, R)
}

func game() {
	const noreturn bool = true
	board := go3mino.NewBoard(2)
	pool := go3mino.AllTriminos()
	hand := go3mino.Triminos{}
	var pick go3mino.Trimino
	var ok bool
	steps := 0
	defer func() {
		fmt.Printf("game finished after %d steps.\n", steps)
		fmt.Printf("%d pieces on board\n", board.GetNumPieces())
		fmt.Printf("hand: %v\n", hand)
		fmt.Printf("pool: %v\n", pool)
		fmt.Printf("check: %d\n", board.GetNumPieces()+len(hand)+len(pool))
	}()
	for i := 0; i < 6; i++ {
		pool, pick, ok = pool.Pick()
		if !ok {
			fmt.Println("something went wrong")
			return
		}
		hand = append(hand, pick)
	}
	hand, pick, ok = hand.Pick()
	if !ok {
		fmt.Println("something went wrong")
		return
	}
	fmt.Printf("put first %v\n", pick)
	board.PlaceFirst(pick)
	steps++

	for {
		{
			n1 := board.GetNumPieces()
			n2 := len(pool)
			n3 := len(hand)
			fmt.Printf("check: board(%d)+pool(%d)+hand(%d) = %d\n", n1, n2, n3, n1+n2+n3)
			fmt.Printf("hand: %v\n", hand)
		}
		fmt.Println(board)
		if len(hand) == 0 {
			fmt.Println("hand is empty")
			if !noreturn {
				return
			}
		} else {
			ok = false
			seen := go3mino.Triminos{}
			for {
				if len(hand) == 0 {
					fmt.Println("no unseen pieces in the hand")
					break
				}
				hand, pick, _ = hand.Pick()
				ok = board.JustPutIt(pick)
				if ok {
					fmt.Printf("put %v\n", pick)
					steps++
					break
				} else {
					seen = append(seen, pick)
				}
			}
			hand = append(hand, seen...) // move all seen back to hand
			if ok {
				continue // next step
			}
		}
		fmt.Println("go fishing...")
		for {
			if len(pool) == 0 {
				fmt.Println("pool is empty and I have no valid step")
				return
			}
			// get piece from pool
			pool, pick, _ = pool.Pick()
			fmt.Printf("get from pool %v ... \n", pick)
			ok = board.JustPutIt(pick)
			if ok {
				fmt.Println("... put it on the board")
				steps++
				break
			} else {
				fmt.Println("... keep it in hand")
				hand = append(hand, pick)
			}
		}
	}
}

func main() {
	game()
}
