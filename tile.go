package main

import (
	"fmt"
	"log"
)

type edgePairID int
type side int
type tileArray []*Tile

// Tile : Holds all the  attributes of a tile
// gets initialised once and no further changes done to it.
type Tile struct {
	// Static attributes
	tileNumber int // the order the tile was read from the file, starting at 1
	sides      [4]side
	// rotations  [4]int
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

func (tile *Tile) removeTileFromEdgePairLists() {
	//fmt.Println("removeTileFromEdgePairLists before", tile)

	for r := 0; r < 4; r++ {
		edgePairList := tile.edgePairLists[r]
		//fmt.Println("Removing from list associated with index/rotation  :", r, edgePairList)
		p := tile.positionInEdgePairList[r]
		edgePairList.removeTile(p, r)
		//fmt.Println("**Removing from list associated with index/rotation:", r, edgePairList)
	}
	//fmt.Println("removeTileFromEdgePairLists before", tile)
}

func (tile *Tile) restoreTileToEdgePairLists() {
	//fmt.Println("restoreTileToEdgePairLists before:", tile)

	for r := 3; r >= 0; r-- {
		edgePairList := tile.edgePairLists[r]
		// p := tile.positionInEdgePairList[r] // TODO  check if this is the right position
		//fmt.Println("Restoring to list associated with index/rotation   :", r, edgePairList)
		edgePairList.restoreTile()
		//fmt.Println("**Restoring to list associated with index/rotation :", r, edgePairList)
	}
	//fmt.Println("restoreTileToEdgePairLists after:", tile)

}

func (tile *Tile) placeTileOnBoard(pos BoardPosition, rotation int, progress int) bool {

	//fmt.Println("Placing tile:", tile.tileNumber, "at position:", pos, "rotation:", rotation)
	// remove the tile from the lists
	tile.removeTileFromEdgePairLists()
	//fmt.Println(tileSet.cornerTilesEdgePairsMap)
	//fmt.Println(tileSet.edgeTilesEdgePairsMap)
	//fmt.Println(tileSet.normalTilesEdgePairsMap)

	// place tile on the board
	board.placeTile(tile, rotation, pos)

	if progress == (board.width * board.height) {
		fmt.Println(board)
		log.Fatalln("finished solution ") // TODO Print out proper solution
		return true
	}
	// get next location to move to

	nextPos := board.nextPosition(pos)
	//fmt.Println("Next Position:", nextPos)
	// set the edgePairIDs of the adjacent tiles
	edgePairID := board.getEdgePairIDForLocation(nextPos)
	//fmt.Println("Next edgePairID:", edgePairDescription(edgePairID))
	if progress > highest_progress {
		fmt.Println(board)
		highest_progress = progress
	}

	//os.Stdout.Sync()
	//positionType := board.loc[nextPos.y][nextPos.x].positionType
	var ok bool
	var edgePairList *tileEdgePairList
	edgePairList, ok = board.loc[nextPos.y][nextPos.x].edgePairMap[edgePairID]
	/*
		if positionType == 'C' {
			edgePairList, ok = tileSet.cornerTilesEdgePairsMap[edgePairID]
		} else if positionType == 'E' {
			edgePairList, ok = tileSet.edgeTilesEdgePairsMap[edgePairID]
		} else {
			edgePairList, ok = tileSet.normalTilesEdgePairsMap[edgePairID]
		}
	*/
	if ok {
		//fmt.Println("Iterating over the following edgePairlist")
		//fmt.Println(edgePairList)
		for i := 0; i < edgePairList.availableNoTiles; i++ {
			nexTtile := edgePairList.tiles[i].tile
			nexTtilerotation := edgePairList.tiles[i].rotation
			// Travers to next position on board
			finished := nexTtile.placeTileOnBoard(nextPos, nexTtilerotation, progress+1)
			if finished {
				return true
			}
		}
	} else {
		//fmt.Println("Unable to find list for edgePairID", edgePairDescription(edgePairID), "Position:", nextPos)
	}
	// remove from board
	//fmt.Println("removeTile :", tile.tileNumber, "Pos:", pos)
	board.removeTile(tile, rotation, pos)
	// restore
	//fmt.Println("restoreTileToAvailableLists :", tile.tileNumber)
	//fmt.Println(tileSet.cornerTilesEdgePairsMap)
	//fmt.Println(tileSet.edgeTilesEdgePairsMap)
	//fmt.Println(tileSet.normalTilesEdgePairsMap)
	tile.restoreTileToEdgePairLists()
	//fmt.Println("Backtracking")
	return false
}
