package main

import (
	"fmt"
)

// Tile : Holds all the  static attributes of a tile
// gets initialised once and no further changes done to it.
type Tile struct {
	tileNumber int // the order the tile was read from the file, starting at 1
	sides      [4]int
	// rotations  [4]int
	tileType  byte   // E is edge, C is corner, N is normal
	edgePairs [4]int // Four edge pairs, adjacent edges
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

func edgePairDescription(egdePair int) string {
	a := egdePair >> tileEdgePairShift
	b := egdePair & tileEdgePairMask
	return fmt.Sprintf("(%v %v)", a, b)
}

func edgePairsDescription(edgePairs [4]int) string {
	var desc string
	for _, v := range edgePairs {
		desc = desc + edgePairDescription(v)
	}
	return desc
}

func (tile Tile) String() string {
	return fmt.Sprintf("Tile No:%v  %v Sides:%v EdgePairs:%v \n",
		tile.tileNumber, tileTypeDescription(tile.tileType), tile.sides, edgePairsDescription(tile.edgePairs))
}

// int EPid = createEPId(nodes[nodeId].sides[s],nodes[nodeId].sides[(s+1)%4]);

// EdgePairID Given a pair of edges of a tile returns the ID of the pair of them. Used to match pairs of edges
func EdgePairID(e1 int, e2 int) int {
	return ((e1 << tileEdgePairShift) + e2) // stricly speaking we shoud mask e2 but assuming that shift is big enough
}

func (tile *Tile) setEdgePairs() {
	for i := range tile.sides {
		tile.edgePairs[i] = EdgePairID(tile.sides[i], tile.sides[(i+1)%4])
	}
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
	tile.setEdgePairs()
	// Determind type of the tile
	tile.setTileType()
}
