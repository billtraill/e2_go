package main

import "fmt"

// BoardLocation : describes the position of location on the board and what is currently at that location
type BoardLocation struct {
	//x, y         int   // the x y location of the position - static
	tile *Tile // pointer to current tile at
	// rotation     int   // current rotation of tile piece. Note corners and edge pieces have static rotations!
	positionType byte //
	edgePairMap  tileEdgePairMap
	//edgePair     edgePairID // the edgepair this location could take -
	edgePairList *tileEdgePairList
	traverseNext BoardPosition
}

// Board : holds the description of the board and any tiles that may currently be placed apon it
type Board struct {
	loc           [][]BoardLocation
	width, height int
}

// BoardPosition refers to a position on the board.
type BoardPosition struct {
	x, y int
}

// nextPosition returns the next position on the board to solve.
func (board Board) nextPosition(pos BoardPosition) BoardPosition {
	return board.loc[pos.y][pos.x].traverseNext

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
		s = s + "  "
		for x := range board.loc[y] {
			loc := &board.loc[y][x]

			if loc.tile != nil {
				s = s + fmt.Sprintf("%2v ", loc.tile.tileNumber)
			} else {
				s = s + fmt.Sprintf(".  ")
			}

		}
		s = s + "\n"
	}
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
		board.loc[yp][xp].traverseNext = BoardPosition{x, y}
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

func setRowByRowTraversal() {
	/*	pos.x++
		if pos.x >= board.width {
		  pos.x = 0
		  pos.y++
		}
		return pos */
}

func (board Board) setTraversal() {
	xp, yp := 0, 0
	x, y := 0, 1
	minY := 0
	maxY := board.height - 1
	minX := 0
	maxX := board.width - 1

	for {
		fmt.Println(x, y)
		board.loc[yp][xp].traverseNext = BoardPosition{x, y}
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

func (board *Board) createBoard(tileSet TileSet) error {
	board.width = tileSet.width
	board.height = tileSet.height
	board.loc = make([][]BoardLocation, tileSet.height)
	// loop over the rows allocating the slice for each row
	for y := range board.loc {
		board.loc[y] = make([]BoardLocation, tileSet.width)
	}
	// Set position types
	for y := range board.loc {
		for x := range board.loc[y] {
			loc := &board.loc[y][x]
			if x == 0 && y == 0 { // top right
				loc.positionType = 'C'
				//loc.rotation = 0
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap
				//loc.edgePair = calcEdgePairID(0, 0) //can be deleted
				//
				// This is required to place 1st tile in top corner
				//
				loc.edgePairList = loc.edgePairMap[calcEdgePairID(0, 0)]
				loc.edgePairList.needCount++
			} else if x == 0 && y == tileSet.height-1 { // bottom right
				loc.positionType = 'C'
				//loc.rotation = 3
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap
			} else if x == tileSet.width-1 && y == 0 { // top left
				loc.positionType = 'C'
				//loc.rotation = 1
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap
			} else if x == tileSet.width-1 && y == tileSet.height-1 { // bottom right
				loc.positionType = 'C'
				//loc.rotation = 2
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap
			} else if x == 0 { // left  edge of board
				loc.positionType = 'E'
				//loc.rotation = 0
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap
			} else if x == tileSet.width-1 { // right edge of board
				loc.positionType = 'E'
				//loc.rotation = 2
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap
			} else if y == 0 { // top edge of board
				loc.positionType = 'E'
				//loc.rotation = 1
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap
			} else if y == tileSet.height-1 { // top edge of board
				loc.positionType = 'E'
				//loc.rotation = 3
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap
			} else { // someplace in the middile of board
				loc.positionType = 'N'
				loc.edgePairMap = tileSet.normalTilesEdgePairsMap
				// rotations are dynamic in the middile of the board so no need to set here.
			}
		}
	}
	board.setTraversal()
	return nil
}

func (board *Board) placeTile(tile *Tile, pos BoardPosition) {
	loc := &board.loc[pos.y][pos.x]
	loc.tile = tile
	/*
		if loc.positionType == 'E' || loc.positionType == 'C' {
			// given rotation and board rotation should be the same, just do a quick check?
			if loc.rotation != rotation {
				log.Fatal(fmt.Sprintf("Tile %v, at position %v does not match rotation. Board:%v Tile %v", tile, pos, loc.rotation, rotation))
			}
		}
	*/
}

func (board *Board) removeTile(tile *Tile, pos BoardPosition) {
	loc := &board.loc[pos.y][pos.x]
	loc.tile = nil
	// loc.rotation = rotation
}

func (board *Board) getEdgePairIDForLocation(pos BoardPosition) edgePairID {
	var a, b side
	if pos.x == 0 {
		a = 0
	} else {
		a = board.loc[pos.y][pos.x-1].tile.sides[(board.loc[pos.y][pos.x-1].tile.rotation+2)%4] // OK
	}
	if pos.y == 0 {
		b = 0
	} else {
		b = board.loc[pos.y-1][pos.x].tile.sides[(board.loc[pos.y-1][pos.x].tile.rotation+3)%4]
	}
	return calcEdgePairID(a, b)
}
