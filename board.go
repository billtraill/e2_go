package main

import "fmt"

// BoardLocation : describes the position of location on the board and what is currently at that location
type BoardLocation struct {
	//x, y         int   // the x y location of the position - static
	tile         *Tile // pointer to current tile at
	rotation     int   // current rotation of tile piece. Note corners and edge pieces have static rotations!
	positionType byte  //
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

func (board Board) nextPosition(pos BoardPosition) BoardPosition {
	return board.loc[pos.y][pos.x].traverseNext
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

/*
func (board Board) nextPositionDiagonal(pos BoardPosition) BoardPosition {
	pos.y--
	if y < 0 {
		pos.y = pos.x + 1
		if pos.y == board.height {
			pos.y = board.height - 1

		} else {
			pos.x = 0
		}

	} else {
		pos.x = pos.x + 1
	}

	return pos
}
*/
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

func tileLine(tile *Tile, rotation int, line int) string {
	s := ""
	if tile == nil {
		s = "      "
	} else if line == 0 {
		s = fmt.Sprintf("   %2v    ", tile.sides[(rotation+1)%4])
	} else if line == 1 {
		s = fmt.Sprintf("%2v %2v %2v ", tile.sides[(rotation)%4], tile.tileNumber, tile.sides[(rotation+2)%4])
	} else if line == 2 {
		s = fmt.Sprintf("   %2v    ", tile.sides[(rotation+3)%4])
	}
	return s
}

func boardwithSides(board Board) string {
	s := ""
	for y := range board.loc {
		for line := 0; line < 4; line++ {
			for x := range board.loc[y] {
				loc := &board.loc[y][x]
				s = s + tileLine(loc.tile, loc.rotation, line)
			}
			s = s + "\n"
		}
	}
	return s
}

func (board Board) String() string {
	s := ""
	for y := range board.loc {
		for x := range board.loc[y] {
			loc := &board.loc[y][x]
			s = s + fmt.Sprintf("%v ", boardLocationTypeDescription(loc.positionType))
		}
		s = s + "  "
		for x := range board.loc[y] {
			loc := &board.loc[y][x]
			s = s + fmt.Sprintf("%v ", loc.rotation)
		}
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
	s = s + boardwithSides(board)
	return s
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
				loc.rotation = 0
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap
				//loc.edgePair = calcEdgePairID(0, 0) //can be deleted
				//
				// This is required to place 1st tile in top corner
				//
				loc.edgePairList = loc.edgePairMap[calcEdgePairID(0, 0)]
				loc.edgePairList.needCount++
			} else if x == 0 && y == tileSet.height-1 { // bottom right
				loc.positionType = 'C'
				loc.rotation = 3
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap
			} else if x == tileSet.width-1 && y == 0 { // top left
				loc.positionType = 'C'
				loc.rotation = 1
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap
			} else if x == tileSet.width-1 && y == tileSet.height-1 { // bottom right
				loc.positionType = 'C'
				loc.rotation = 2
				loc.edgePairMap = tileSet.cornerTilesEdgePairsMap
			} else if x == 0 { // left  edge of board
				loc.positionType = 'E'
				loc.rotation = 0
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap
			} else if x == tileSet.width-1 { // right edge of board
				loc.positionType = 'E'
				loc.rotation = 2
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap
			} else if y == 0 { // top edge of board
				loc.positionType = 'E'
				loc.rotation = 1
				loc.edgePairMap = tileSet.edgeTilesEdgePairsMap
			} else if y == tileSet.height-1 { // top edge of board
				loc.positionType = 'E'
				loc.rotation = 3
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

func (board *Board) placeTile(tile *Tile, rotation int, pos BoardPosition) {
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
	loc.rotation = rotation
}

func (board *Board) removeTile(tile *Tile, rotation int, pos BoardPosition) {
	loc := &board.loc[pos.y][pos.x]
	loc.tile = nil
	// loc.rotation = rotation
}

func (board *Board) getEdgePairIDForLocation(pos BoardPosition) edgePairID {
	var a, b side
	if pos.x == 0 {
		a = 0
	} else {
		a = board.loc[pos.y][pos.x-1].tile.sides[(board.loc[pos.y][pos.x-1].rotation+2)%4] // OK
	}
	if pos.y == 0 {
		b = 0
	} else {
		b = board.loc[pos.y-1][pos.x].tile.sides[(board.loc[pos.y-1][pos.x].rotation+3)%4]
	}
	return calcEdgePairID(a, b)
}
