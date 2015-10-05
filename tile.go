package main

import (
	"fmt"
	"log"
	"time"
)

type edgePairID int
type side int
type tileArray []*Tile

// Tile : Holds all the  attributes of a tile
// gets initialised once and no further changes done to it.
type Tile struct {
	// Static attributes
	tileNumber    int                  // the order the tile was read from the file, starting at 1
	sides         [4]side              // values of each of the sides. Used when calculating Edge pairs
	tileType      byte                 // E is edge, C is corner, N is normal
	edgePairs     [4]edgePairID        // Four edge pairs, adjacent edges
	edgePairLists [4]*tileEdgePairList // note order of list implies the rotation of the tile from its normalised positon
	// dynamic values .... these get changed as we run ....
	positionInEdgePairList [4]int // this tracks where the tile is currently in the edgePairLists - this changes as we remove/add tiles to lists
}

func tileTypeDescription(t byte) string {
	var desc string
	switch t {
	case 'C':
		desc = "Corner"
	case 'E':
		desc = "Edge"
	case 'N':
		desc = "Normal"
	default:
		desc = "Unknown"
	}

	return fmt.Sprintf(desc)
}

const tileEdgePairShift = 5   // This will allow 5 bits of information for the edge colour ... if more than 2^5 edge types then increase this
const tileEdgePairMask = 0x1F // 5 bits worth

func edgePairDescription(egdePair edgePairID) string {
	a := egdePair >> tileEdgePairShift
	b := egdePair & tileEdgePairMask
	return fmt.Sprintf("(%v %v)", a, b)
}

func edgePairsDescription(edgePairs [4]edgePairID) string {
	var desc string
	for _, v := range edgePairs {
		desc = desc + edgePairDescription(v)
	}
	return desc
}

func (tiles tileArray) String() string {
	s := ""
	for _, v := range tiles {
		s = s + fmt.Sprintln(v)
	}
	return s
}

func (tile Tile) String() string {
	return fmt.Sprintf("Tile No:%v  %v Sides:%v EdgePairs:%v PositionInEPList:%v",
		tile.tileNumber, tileTypeDescription(tile.tileType), tile.sides, edgePairsDescription(tile.edgePairs), tile.positionInEdgePairList)
}

// int EPid = createEPId(nodes[nodeId].sides[s],nodes[nodeId].sides[(s+1)%4]);

// EdgePairID Given a pair of edges of a tile returns the ID of the pair of them. Used to match pairs of edges
func calcEdgePairID(e1 side, e2 side) edgePairID {
	return edgePairID((int(e1) << tileEdgePairShift) + int(e2)) // stricly speaking we shoud mask e2 but assuming that shift is big enough
}

func (tile *Tile) setEdgePairs() {
	for i := range tile.sides {
		tile.edgePairs[i] = calcEdgePairID(tile.sides[i], tile.sides[(i+1)%4])
	}
}

// normaliseEdges normalises the tile so the border edges are first in the array
func (tile *Tile) normaliseEdges() {
	idx := 0
	for i, v := range tile.sides {
		if v == 0 {
			idx = i
			break
		}
	}
	var tmpSides [4]side
	for i2 := range tile.sides {
		tmpSides[i2] = tile.sides[(i2+idx)%4]
	}
	tile.sides = tmpSides
}

func (tile *Tile) setTileType() {
	edgeCount := 0
	for _, v := range tile.sides {
		if v == 0 {
			edgeCount++
		}
	}
	switch edgeCount {
	case 0:
		tile.tileType = 'N'
	case 1:
		tile.tileType = 'E'
	case 2:
		tile.tileType = 'C'
	default:
		fmt.Println("Tile has more than 2 side edges !!!!")
	}
}

func (tile *Tile) setTileProperties() {
	// Determind type of the tile
	tile.setTileType()
	// Normalise the tile so the boarder edges are fist
	if tile.tileType == 'E' || tile.tileType == 'C' {
		tile.normaliseEdges()
	}
	//
	tile.setEdgePairs()
}

//
// This does a quick check to see if valid edgePairs are in the adjacent locations
// and if they are check if there are any tiles available in thoes lists and if they
// are do a preallocation by incrementing their need count.
// This ensures we detect that we have not created impossible edge constraint as early as possible
//
func reserveDownPosition(pos BoardPosition) bool {
	var edgePairList *tileEdgePairList

	var ok bool
	if pos.y+1 < board.height {
		edgePairID := board.getEdgePairIDForLocation(BoardPosition{pos.x, pos.y + 1})
		edgePairList, ok = board.loc[pos.y+1][pos.x].edgePairMap[edgePairID]
		if !ok {
			return false
		}

		if edgePairList.needCount >= edgePairList.availableNoTiles {
			return false
		}
		board.loc[pos.y+1][pos.x].edgePairList = edgePairList
		edgePairList.needCount++
	}
	return true
}
func clearReserveDownPosition(pos BoardPosition) {
	if pos.y+1 < board.height {
		board.loc[pos.y+1][pos.x].edgePairList.needCount--
		board.loc[pos.y+1][pos.x].edgePairList = nil // not nessary but should catch any bugs!
	}
	return
}

func reserveAcrossPosition(pos BoardPosition) bool {
	var edgePairList *tileEdgePairList
	var ok bool
	if pos.y == 0 && pos.x+1 < board.width {
		edgePairID := board.getEdgePairIDForLocation(BoardPosition{pos.x + 1, pos.y})
		edgePairList, ok = board.loc[pos.y][pos.x+1].edgePairMap[edgePairID]
		if !ok {
			return false
		}

		if edgePairList.needCount >= edgePairList.availableNoTiles {
			return false
		}
		board.loc[pos.y][pos.x+1].edgePairList = edgePairList
		edgePairList.needCount++
	}
	return true
}
func clearReserveAcrossPosition(pos BoardPosition) {
	if pos.y == 0 && pos.x+1 < board.width {
		board.loc[pos.y][pos.x+1].edgePairList.needCount--
		board.loc[pos.y][pos.x+1].edgePairList = nil // not nessary but should catch any bugs!
	}
	return
}

//
// placeTileOnBoard is the main solver function. It recursively places tiles down on the board
// checking ahead to see if valid solutions are still possible from the remaining tiles
// the main datastructure used is the edgePairLists, each tile has 4 associated lists, one for
// each of its rotations. Each edgepair has a list of all the tile that have this combination of
// edges.
func (tile *Tile) placeTileOnBoard(pos BoardPosition, rotation int, progress int) bool {

	//fmt.Println("Placing tile:", tile.tileNumber, "at position:", pos, "rotation:", rotation)
	// remove the tile from the lists
	tile.edgePairLists[0].removeTile(tile.positionInEdgePairList[0])
	tile.edgePairLists[1].removeTile(tile.positionInEdgePairList[1])
	tile.edgePairLists[2].removeTile(tile.positionInEdgePairList[2])
	tile.edgePairLists[3].removeTile(tile.positionInEdgePairList[3])

	// place tile on the board
	board.placeTile(tile, rotation, pos)
	if reserveDownPosition(pos) {
		if reserveAcrossPosition(pos) {

			// get next location to move to
			//nextPos := board.nextPosition(pos)
			nextPos := board.loc[pos.y][pos.x].traverseNext
			//fmt.Println("Next Position:", nextPos)

			//fmt.Println("Next edgePairID:", edgePairDescription(edgePairID))
			if progress >= highest_progress {
				fmt.Println(board)
				highest_progress = progress
				fmt.Println("Placed:", progress, time.Now().Format(time.RFC850))
				if progress == (board.width * board.height) {
					fmt.Println(board)
					log.Fatalln("finished solution ") // TODO Print out proper solution
					return true
				}
			}

			edgePairList := board.loc[nextPos.y][nextPos.x].edgePairList

			// Iterates over all the tiles in the list...
			for i := 0; i < edgePairList.availableNoTiles; i++ {
				nexTtile := edgePairList.tiles[i].tile
				nexTtilerotation := edgePairList.tiles[i].rotation
				// Travers to next position on board
				finished := nexTtile.placeTileOnBoard(nextPos, nexTtilerotation, progress+1)
				if finished {
					return true
				}
			}
			clearReserveAcrossPosition(pos)
		}
		clearReserveDownPosition(pos)
	}

	// remove from board
	//fmt.Println("removeTile :", tile.tileNumber, "Pos:", pos)
	board.removeTile(tile, rotation, pos)
	// restore tile to its edge pair lists, has to be done in the reverse they were added
	// to deal with the fact that some times have the same edge pair list more than once !
	tile.edgePairLists[3].restoreTile()
	tile.edgePairLists[2].restoreTile()
	tile.edgePairLists[1].restoreTile()
	tile.edgePairLists[0].restoreTile()
	return false
}
