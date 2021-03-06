package connect

import (
	"errors"
	"fmt"
	"strings"
)

const (
	WHITE = 1 << iota
	BLACK
	WHITE_CONNECTED
	BLACK_CONNECTED
)

type colorFlags struct {
	color_flag     int8
	connected_flag int8
}

var black_flags = colorFlags{
	color_flag:     BLACK,
	connected_flag: BLACK_CONNECTED,
}

var white_flags = colorFlags{
	color_flag:     WHITE,
	connected_flag: WHITE_CONNECTED,
}

type coord struct {
	x int
	y int
}

type board struct {
	height int
	width  int
	fields [][]int8
}

func newBoard(lines []string) (board, error) {
	if len(lines) < 1 {
		return board{}, errors.New("No lines given")
	}
	height := int(len(lines))
	if len(lines[0]) < 1 {
		return board{}, errors.New("First line is empty string")
	}
	width := int(len(lines[0]))
	// This trick for 2D arrays comes from Effective Go
	fields := make([][]int8, height)
	fieldsBacker := make([]int8, height*width)
	for i := range fields {
		fields[i], fieldsBacker = fieldsBacker[:width], fieldsBacker[width:]
	}
	for y, line := range lines {
		for x, c := range line {
			switch c {
			case 'X':
				fields[y][x] = BLACK
			case 'O':
				fields[y][x] = WHITE
			}
			// No need for default, zero value already means no stone
		}
	}
	board := board{
		height: height,
		width:  width,
		fields: fields,
	}
	return board, nil
}

// Whether there is a stone of the given color at the given location.
//
// Returns both whether there is a stone of the correct color and
// whether the connected flag was set for it.
func (b board) at(c coord, cf colorFlags) (bool, bool) {
	f := b.fields[c.y][c.x]
	return f&cf.color_flag == cf.color_flag,
		f&cf.connected_flag == cf.connected_flag
}

func (b board) markConnected(c coord, cf colorFlags) {
	b.fields[c.y][c.x] |= cf.connected_flag
}

func (b board) validCoord(c coord) bool {
	return c.x >= 0 && c.x < b.width && c.y >= 0 && c.y < b.height
}

func (b board) neighbours(c coord) []coord {
	coords := make([]coord, 0, 6)
	dirs := []coord{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {-1, 1}, {1, -1}}
	for _, dir := range dirs {
		nc := coord{x: c.x + dir.x, y: c.y + dir.y}
		if b.validCoord(nc) {
			coords = append(coords, nc)
		}
	}
	return coords
}

func (b board) startCoords(cf colorFlags) []coord {
	if cf.color_flag == WHITE {
		coords := make([]coord, b.width)
		for i := 0; i < b.width; i++ {
			coords[i] = coord{x: i, y: 0}
		}
		return coords
	} else {
		coords := make([]coord, b.height)
		for i := 0; i < b.height; i++ {
			coords[i] = coord{x: 0, y: i}
		}
		return coords
	}
}

func (b board) isTargetCoord(c coord, cf colorFlags) bool {
	if cf.color_flag == WHITE {
		return c.y == b.height-1
	} else {
		return c.x == b.width-1
	}
}

func (b board) evaluate(c coord, cf colorFlags) bool {
	stone, connected := b.at(c, cf)
	if stone && !connected {
		b.markConnected(c, cf)
		if b.isTargetCoord(c, cf) {
			return true
		}
		for _, nc := range b.neighbours(c) {
			if b.evaluate(nc, cf) {
				return true
			}
		}
	}
	return false
}

// Helper for debugging.
func (b board) dump() {
	for y := 0; y < b.height; y++ {
		spaces := strings.Repeat(" ", y)
		chars := make([]string, b.width)
		for x := 0; x < b.width; x++ {
			if b.fields[y][x]&WHITE == WHITE {
				if b.fields[y][x]&WHITE_CONNECTED == WHITE_CONNECTED {
					chars[x] = "O"
				} else {
					chars[x] = "o"
				}
			} else if b.fields[y][x]&BLACK == BLACK {
				if b.fields[y][x]&BLACK_CONNECTED == BLACK_CONNECTED {
					chars[x] = "X"
				} else {
					chars[x] = "x"
				}
			} else {
				chars[x] = "."
			}
		}
		fmt.Printf("%s%s\n", spaces, strings.Join(chars, " "))
	}
}

// ResultOf evaluates the board and return the winner, "black" or
// "white". If there's no winnner ResultOf returns "".
func ResultOf(lines []string) (string, error) {
	board, err := newBoard(lines)
	if err != nil {
		return "", err
	}
	for _, c := range board.startCoords(black_flags) {
		if board.evaluate(c, black_flags) {
			return "black", nil
		}
	}
	for _, c := range board.startCoords(white_flags) {
		if board.evaluate(c, white_flags) {
			return "white", nil
		}
	}
	return "", nil
}
