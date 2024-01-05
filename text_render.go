package go3mino

import (
	"strconv"
	"strings"
)

func (board *Board) Render() []string {
	var half = map[bool][5]string{
		false: {
			`,,,,,`,
			`,,,,`,
			`,,,`,
			`,,`,
			`,`,
		},
		true: {
			`.`,
			`..`,
			`...`,
			`....`,
			`.....`,
		},
	}

	var empty = map[bool][5]string{
		false: {
			`.`,
			`...`,
			`.....`,
			`.......`,
			`.........`,
		},
		true: {
			`,,,,,,,,,`,
			`,,,,,,,`,
			`,,,,,`,
			`,,,`,
			`,`,
		},
	}
	var cell = map[bool][5]string{
		false: {
			`^`,
			`/A\`,
			`/   \`,
			`/C   B\`,
			`/_______\`,
		},
		true: {
			`\~~~~~~~/`,
			`\A   B/`,
			`\   /`,
			`\C/`,
			`v`,
		},
	}

	var result []string = nil
	for row := 0; row < board.Height; row++ {
		for line := 0; line < 5; line++ {
			var bld strings.Builder
			filler := half[row%2 == 0][line]
			bld.WriteString(filler)
			for col := 0; col < board.Width; col++ {
				p := board.Map[row*board.Width+col]
				v_type := row%2 == col%2
				if p.IsFree() {
					bld.WriteString(empty[v_type][line])
				} else {
					a := strconv.Itoa(p.GetNode(0))
					b := strconv.Itoa(p.GetNode(1))
					c := strconv.Itoa(p.GetNode(2))
					bld.WriteString(
						strings.ReplaceAll(
							strings.ReplaceAll(
								strings.ReplaceAll(
									cell[v_type][line],
									`A`, a),
								`B`, b),
							`C`, c))
				}
			}
			bld.WriteString(filler)
			result = append(result, bld.String())
		}
	}
	return result
}

func (board *Board) String() string {
	return strings.Join(board.Render(), "\n")
}
