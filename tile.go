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
	tileNumber         int                  // the order the tile was read from the file, starting at 1
	sides              [4]side              // values of each of the sides. Used when calculating Edge pairs
	tileType           byte                 // E is edge, C is corner, N is normal
	duplicateEdgePairs bool                 // flag to indicate if there are duplicate Edge Pairs it means for these tiles we cannot manage their edge pair lists concurrently
	edgePairs          [4]edgePairID        // Four edge pairs, adjacent edges
	edgePairLists      [4]*tileEdgePairList // note order of list implies the rotation of the tile from its normalised positon
	// dynamic values .... these get changed as we run ....
	positionInEdgePairList [4]int // this tracks where the tile is currently in the edgePairLists - this changes as we remove/add tiles to lists
	rotation               int
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

// tileLine is used by boardshowTilesPlaced, it is designed to show a printed representation of a tile over 3 lines.
// It returns a string for each line for a tile requested, showing its current rotation on the board
func (tile *Tile) tileLine(line int) string {
	s := ""
	if tile == nil {
		s = "      "
	} else if line == 0 {
		s = fmt.Sprintf("   %2v    ", tile.sides[(tile.rotation+1)%4])
	} else if line == 1 {
		s = fmt.Sprintf("%2v %2v %2v ", tile.sides[(tile.rotation)%4], tile.tileNumber, tile.sides[(tile.rotation+2)%4])
	} else if line == 2 {
		s = fmt.Sprintf("   %2v    ", tile.sides[(tile.rotation+3)%4])
	}
	return s
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
	// Just a bit of information
	tile.duplicateEdgePairs = false
	for i := range tile.edgePairs {
		for j := range tile.edgePairs {
			if i != j {
				if tile.edgePairs[i] == tile.edgePairs[j] {
					tile.duplicateEdgePairs = true
					fmt.Println("Tile:", tile.tileNumber, "Has duplicate edgepairs!")
				}
			}
		}
	}
	//tile.duplicateEdgePairs = true
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
func reserveDownPosition(loc *BoardLocation) bool {
	var edgePairList *tileEdgePairList

	var ok bool
	if loc.down != nil {
		edgePairID := loc.down.getEdgePairIDForLocation()
		edgePairList, ok = loc.down.edgePairMap[edgePairID]
		if !ok {
			return false
		}

		if edgePairList.needCount >= edgePairList.availableNoTiles {
			return false
		}
		loc.down.edgePairList = edgePairList
		edgePairList.needCount++
	}
	return true
}
func clearReserveDownPosition(loc *BoardLocation) {
	if loc.down != nil {
		loc.down.edgePairList.needCount--
		loc.down.edgePairList = nil // not nessary but should catch any bugs!
	}
	return
}

func reserveAcrossPosition(loc *BoardLocation) bool {
	var edgePairList *tileEdgePairList
	var ok bool
	if loc.up == nil && loc.right != nil {
		edgePairID := loc.right.getEdgePairIDForLocation()
		edgePairList, ok = loc.right.edgePairMap[edgePairID]
		if !ok {
			return false
		}

		if edgePairList.needCount >= edgePairList.availableNoTiles {
			return false
		}
		loc.right.edgePairList = edgePairList
		edgePairList.needCount++
	}
	return true
}
func clearReserveAcrossPosition(loc *BoardLocation) {
	if loc.up == nil && loc.right != nil {
		loc.right.edgePairList.needCount--
		loc.right.edgePairList = nil // not nessary but should catch any bugs!
	}
	return
}

//
// placeTileOnBoard is the main solver function. It recursively places tiles down on the board
// checking ahead to see if valid solutions are still possible from the remaining tiles
// the main datastructure used is the edgePairLists, each tile has 4 associated lists, one for
// each of its rotations. Each edgepair has a list of all the tile that have this combination of
// edges.
var complete = make(chan int, 5)

func (loc *BoardLocation) placeTileOnBoard(progress int) bool {

	//fmt.Println("Placing tile:", tile.tileNumber, "at position:", loc.x, loc.y, "rotation:", tile.rotation)
	// remove the tile from the lists

	loc.tile.edgePairLists[0].removeTile(loc.tile.positionInEdgePairList[0])
	loc.tile.edgePairLists[1].removeTile(loc.tile.positionInEdgePairList[1])
	loc.tile.edgePairLists[2].removeTile(loc.tile.positionInEdgePairList[2])
	loc.tile.edgePairLists[3].removeTile(loc.tile.positionInEdgePairList[3])

	// place tile on the board
	//loc.tile = tile

	//if reserveDownPosition(loc) {
	//if reserveAcrossPosition(loc) {

	//fmt.Println("Next edgePairID:", edgePairDescription(edgePairID))
	if progress >= highestProgress {
		fmt.Println(board)
		highestProgress = progress
		fmt.Println("Placed:", progress, time.Now().Format(time.RFC850))
		if progress == (board.width * board.height) {
			fmt.Println(board)
			log.Fatalln("finished solution ") // TODO Print out proper solution
			return true
		}
	}
	// get next location to move to
	nextPos := loc.traverseNext
	//fmt.Println("Next Position:", nextPos.x, nextPos.y)

	edgePairID := nextPos.getEdgePairIDForLocation()
	//edgePairList := nextPos.edgePairList

	edgePairList, ok := nextPos.edgePairMap[edgePairID]
	if ok { // there is an edge pair mapping for this location ...
		// Iterates over all the tiles in the list...
		for i := 0; i < edgePairList.availableNoTiles; i++ {
			nexTtile := edgePairList.tiles[i].tile
			// set the tiles rotation
			nexTtile.rotation = edgePairList.tiles[i].rotationForEdgePair
			// place tile on the board in the next position and traverse to it
			nextPos.tile = nexTtile
			// Travers to next position on board
			finished := nextPos.placeTileOnBoard(progress + 1)
			nextPos.tile = nil
			if finished {
				return true
			}
		}
	}

	//	clearReserveAcrossPosition(loc)
	//}
	//clearReserveDownPosition(loc)
	//}

	// restore tile to its edge pair lists, has to be done in the reverse they were added
	// to deal with the fact that some times have the same edge pair list more than once !

	loc.tile.edgePairLists[3].restoreTile()
	loc.tile.edgePairLists[2].restoreTile()
	loc.tile.edgePairLists[1].restoreTile()
	loc.tile.edgePairLists[0].restoreTile()

	return false
}

func traverseBoard() {
	var nextPos *BoardLocation
	var loc *BoardLocation
	var edgePairID edgePairID
	var progress int
	var edgePairList *tileEdgePairList
	var ok bool

	var highestProgress int

	// need to have current location set to 1st tile to start TODO
	loc = &board.loc[0][0]
	loc.edgePairList = loc.edgePairMap[calcEdgePairID(0, 0)]
	loc.index = 0

	progress = 0
	highestProgress = 0

	for {
		// are there still tiles left to try in the edge pair list on the current location
		if loc.index < loc.edgePairList.availableNoTiles {
			loc.tile = loc.edgePairList.tiles[loc.index].tile
			loc.tile.rotation = loc.edgePairList.tiles[loc.index].rotationForEdgePair

			if progress >= highestProgress {
				fmt.Println(board)
				highestProgress = progress
				fmt.Println("Placed:", progress, time.Now().Format(time.RFC850))
				if progress == (board.width*board.height)-1 {
					fmt.Println(board)
					fmt.Println("finished solution ") // TODO Print out proper solution
					return
				}
			}

			// Now see if there is a valid EP in next location
			nextPos = loc.traverseNext
			edgePairID = nextPos.getEdgePairIDForLocation()

			edgePairList, ok = nextPos.edgePairMap[edgePairID]
			if ok && edgePairList.availableNoTiles > 0 { // valid edgepair list and it has tiles on it. Note the edgepair might be the same as the current one TODO work out if this has any impacts.
				progress++
				// remove the tile from the EP lists
				loc.tile.edgePairLists[0].removeTile(loc.tile.positionInEdgePairList[0])
				loc.tile.edgePairLists[1].removeTile(loc.tile.positionInEdgePairList[1])
				loc.tile.edgePairLists[2].removeTile(loc.tile.positionInEdgePairList[2])
				loc.tile.edgePairLists[3].removeTile(loc.tile.positionInEdgePairList[3])
				// move to next position on board
				nextPos.edgePairList = edgePairList
				nextPos.index = 0
				loc = nextPos
				continue
			} else { // move onto next tile in list
				loc.index++
				continue
			}
		} else {
			// backtrack
			//fmt.Println("Backtracking")
			progress--
			// traverse to previous tile
			loc = loc.traversePrev
			// remove the tile in the previous location ...
			loc.tile.edgePairLists[3].restoreTile()
			loc.tile.edgePairLists[2].restoreTile()
			loc.tile.edgePairLists[1].restoreTile()
			loc.tile.edgePairLists[0].restoreTile()
			//loc.tile = nil // probably not required but handy, used when we print out progress ...
			loc.index++

		}
	}

}
