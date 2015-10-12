package main

import (
	"fmt"

	"github.com/davidminor/uint128"
)

// BoardLocation : describes the position of location on the board and what is currently at that location
type BoardLocation struct {
	x, y         int            // the x y location of the position - static
	positionType byte           //
	traverseNext *BoardLocation // traversal order - setup when board is created
	traversePrev *BoardLocation // points to previous location
	left         *BoardLocation // points to location on left
	down         *BoardLocation // points to location on down
	up           *BoardLocation //
	right        *BoardLocation
	// Dynamtic parts of a location on the board
	tile           *Tile           // pointer to current tile at
	edgePairMap    tileEdgePairMap // this is map of all the edgepair lists valid for this location (corner, edge or normal ) this has to be like this as a corner and edge edgepair can look the same!
	edgePairList   *tileEdgePairList
	index          int    // current index in edgepair list
	listSize       int    // just used for debug, used to record the size of the edge pair list on the current location - used to show structure/progress in debug
	noTimesVisited uint64 // number of times this location has been visited.

	// used for conposite board traversal coroutine of get next edge pair
	edgePairChan chan int
}

// Board : holds the description of the board and any tiles that may currently be placed apon it
type Board struct {
	loc           [][]BoardLocation
	width, height int
	cTilePlaced   uint128.Uint128 // Tracks which tiles have been placed on compisite board
}

// boardLocationTypeDescription for a given position type returns a string "describing" that position.
func boardLocationTypeDescription(t byte) string {
	var desc string
	switch t {
	case 'C':
		desc = "C"
	case 'E':
		desc = "E"
	case 'N':
		desc = "N"
	default:
		desc = "U"
	}
	return desc
}

// boardshowTilesPlaced gives a string representation of all the tiles currently placed on the board
func boardshowTilesPlaced(board Board) string {
	s := ""
	for y := range board.loc {
		for line := 0; line < 4; line++ {
			for x := range board.loc[y] {
				loc := &board.loc[y][x]
				if loc.tile != nil {
					s = s + loc.tile.tileLine(line)
				} else {
					s = s + "      "
				}

			}
			s = s + "\n"
		}
	}
	return s
}

// String  returns a string describing the board
// Useful for seeing progress of how solution is progressing
func (board Board) String() string {
	var combinations uint64
	combinations = 1
	s := ""

	for y := range board.loc {

		for x := range board.loc[y] {
			loc := &board.loc[y][x]
			s = s + fmt.Sprintf("%v ", boardLocationTypeDescription(loc.positionType))
		}

		/*
					s = s + "  "
					for x := range board.loc[y] {
						loc := &board.loc[y][x]
			      if loc.
						s = s + fmt.Sprintf("%v ", tile.rotation)
					} */
		/*
			s = s + "  "
			for x := range board.loc[y] {
				loc := &board.loc[y][x]

				if loc.tile != nil {
					s = s + fmt.Sprintf("%2v ", loc.tile.tileNumber)
				} else {
					s = s + fmt.Sprintf(".  ")
				}
			}
			s = s + "  "
		*/
		for x := range board.loc[y] {
			loc := &board.loc[y][x]

			if loc.tile != nil {
				s = s + fmt.Sprintf("%3v/%3v ", loc.listSize, loc.index)
				combinations = combinations * uint64(loc.listSize)
			} else {
				s = s + fmt.Sprintf(".       ")
			}
		}
		s = s + "  "
		for x := range board.loc[y] {
			loc := &board.loc[y][x]
			s = s + fmt.Sprintf("%16v ", loc.noTimesVisited)
		}
		s = s + "\n"
	}
	s = s + fmt.Sprintf("Current no of combinations: %v\n", combinations)
	s = s + boardshowTilesPlaced(board)
	return s
}

// setDiagonalTraversal sets a diagonal traversal of the board. It assumes the starting
// point is (0,0) (top left of board)
// If the traversal order is changed then the way to set lookforward constraints also needs to change
// this is currently setup in placeTileOnBoard function in the reserveDownPosition and reserveAcrossPosition
// function calls  (and their associated reversal ones clearReserveAcrossPosition,clearReserveDownPosition)
func (board Board) setDiagonalTraversal() {
	xp, yp := 0, 0
	x, y := 0, 1
	minY := 0
	maxY := board.height - 1
	minX := 0
	maxX := board.width - 1

	for {
		fmt.Println(x, y)
		board.loc[yp][xp].traverseNext = &board.loc[y][x]
		xp, yp = x, y
		if x == board.width-1 && y == board.height-1 {
			// board.loc[yp][xp].traverseNext = BoardPosition{-1, -1}
			break
		}
		y--
		if y < minY {
			y = x + 1
			if y > maxY {
				y = maxY
				minX++
			}
			x = minX
		} else {
			x++
			if x > maxX {
				y = x + 1
				if y > maxY {
					y = maxY
					minX++
				}
				x = minX
			}
		}

	}

	return

}

func (board Board) setRowByRowTraversal() {
	xp, yp := 0, 0
	x, y := 1, 0
	for {
		//fmt.Println(x, y)
		board.loc[yp][xp].traverseNext = &board.loc[y][x]
		board.loc[y][x].traversePrev = &board.loc[yp][xp]
		xp, yp = x, y
		if x == board.width-1 && y == board.height-1 {
			// board.loc[yp][xp].traverseNext = BoardPosition{-1, -1}
			break
		}
		x++
		if x >= board.width {
			x = 0
			y++
		}

	}
}

func (board Board) setTraversal() {
	//board.setDiagonalTraversal()
	board.setRowByRowTraversal()
}

func (board *Board) createBoard(tileSet TileSet, width int, height int) error {
	fmt.Println("createBoard: Board size", width, height)
	board.width = width
	board.height = height
	board.loc = make([][]BoardLocation, board.height)
	// loop over the rows allocating the slice for each row
	for y := range board.loc {
		board.loc[y] = make([]BoardLocation, board.width)
	}
	// Set position types
	for y := range board.loc {
		for x := range board.loc[y] {
			//fmt.Println("createBoard: setting up point", x, y)
			loc := &board.loc[y][x]
			loc.x = x // not really used for much
			loc.y = y
			// create channels for each location, used by compisite board traversal to get get valid tile/rotation for a location.
			board.loc[y][x].edgePairChan = make(chan int, 100) // just a guess at size
			if x == 0 && y == 0 {                              // top left
				loc.positionType = 'C'
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap
				loc.left = nil
				loc.up = nil
				loc.right = &board.loc[y][x+1]
				loc.down = &board.loc[y+1][x]

			} else if x == 0 && y == board.height-1 { // bottom left
				loc.positionType = 'C'
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap

				loc.left = nil
				loc.up = &board.loc[y-1][x]
				loc.right = &board.loc[y][x+1]
				loc.down = nil
			} else if x == board.width-1 && y == 0 { // top right
				loc.positionType = 'C'
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap

				loc.left = &board.loc[y][x-1]
				loc.up = nil
				loc.right = nil
				loc.down = &board.loc[y+1][x]
			} else if x == board.width-1 && y == board.height-1 { // bottom right
				loc.positionType = 'C'
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap

				loc.left = &board.loc[y][x-1]
				loc.up = &board.loc[y-1][x]
				loc.right = nil
				loc.down = nil

			} else if x == 0 { // left  edge of board
				loc.positionType = 'E'
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap

				loc.left = nil
				loc.up = &board.loc[y-1][x]
				loc.right = &board.loc[y][x+1]
				loc.down = &board.loc[y+1][x]

			} else if x == board.width-1 { // right edge of board
				loc.positionType = 'E'
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap

				loc.left = &board.loc[y][x-1]
				loc.up = &board.loc[y-1][x]
				loc.right = nil
				loc.down = &board.loc[y+1][x]

			} else if y == 0 { // top edge of board
				loc.positionType = 'E'
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap

				loc.left = &board.loc[y][x-1]
				loc.up = nil
				loc.right = &board.loc[y][x+1]
				loc.down = &board.loc[y+1][x]
			} else if y == board.height-1 { // bottom  edge of board
				loc.positionType = 'E'
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap

				loc.left = &board.loc[y][x-1]
				loc.up = &board.loc[y-1][x]
				loc.right = &board.loc[y][x+1]
				loc.down = nil
			} else { // someplace in the middile of board
				loc.positionType = 'N'
				loc.edgePairMap = tileSet.normalTilesEdgePairsMap

				loc.left = &board.loc[y][x-1]
				loc.up = &board.loc[y-1][x]
				loc.right = &board.loc[y][x+1]
				loc.down = &board.loc[y+1][x]
			}
		}
	}
	board.setTraversal()
	return nil
}

func (loc *BoardLocation) getEdgePairIDForLocation() edgePairID {
	var a, b side
	if loc.left == nil {
		a = 0
	} else {
		a = loc.left.tile.sides[(loc.left.tile.rotation+2)%4] // OK
	}
	if loc.up == nil {
		b = 0
	} else {
		b = loc.up.tile.sides[(loc.up.tile.rotation+3)%4]
	}
	return calcEdgePairID(a, b)
}

func (loc *BoardLocation) getCompositeEdgePairIDForLocation() edgePairID {
	var a, b side
	if loc.left == nil {
		a = 0
	} else {
		a = loc.left.tile.sides[(loc.left.tile.rotation+2)%4] // OK
		a = cTileSideSwap(a)
	}
	if loc.up == nil {
		b = 0
	} else {
		b = loc.up.tile.sides[(loc.up.tile.rotation+3)%4]
		b = cTileSideSwap(b)
	}
	return calcEdgePairID(a, b)
}

func (loc *BoardLocation) getCompositeEdgePairIDForLocationAssumingGivenTileIsOnLeft(tile *Tile, rotation int) edgePairID {
	var a, b side
	if loc.left == nil {
		a = 0
	} else {
		a = tile.sides[(rotation+2)%4] // OK
		a = cTileSideSwap(a)
	}
	if loc.up == nil {
		b = 0
	} else {
		b = loc.up.tile.sides[(loc.up.tile.rotation+3)%4]
		b = cTileSideSwap(b)
	}
	return calcEdgePairID(a, b)
}
