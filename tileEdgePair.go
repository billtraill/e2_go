package main

import (
	"fmt"
	"log"
	"sort"
)

const maxEdgePairListSize = 20 // Probably only 17

// tileEdgePairList an list that refers to each tile with a given edgePairID
type tileEdgePairList struct { // equiv to EPListType
	tiles            [maxEdgePairListSize]*tileAndRotation // pointer to current tile at
	noTiles          int
	availableNoTiles int
	needCount        int
	edgePairID       edgePairID
	tileType         byte
}
type tileAndRotation struct {
	tile                            *Tile
	tilepositionInEdgePairListIndex int // we need to know this as when we move the tiles in the edgePosition lists this tells us the lists index in the tile! This does not change
	rotation                        int // tile rotaton note this is its (4-index)%4
	// Dynamic elements....
	previousPosition int // when we remove a tile, this holds its the position we removed it from. Need for restore.
}

type tileEdgePairMap map[edgePairID]*tileEdgePairList // Map of EdgePairID to tileEdgePairList

func sortedKeys(m tileEdgePairMap) []int {
	var keys []int
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	return keys
}

func (m tileEdgePairMap) String() string {
	s := ""
	keys := sortedKeys(m)
	for _, k := range keys {
		key := edgePairID(k)
		s = s + fmt.Sprintln("k:", edgePairDescription(key), "v:", m[key])
	}
	return s
}

func (edgePairList *tileEdgePairList) String() string {

	s := fmt.Sprintf("%v: %v %v (", edgePairDescription(edgePairList.edgePairID), tileTypeDescription(edgePairList.tileType), edgePairList.availableNoTiles)
	for i := 0; i < edgePairList.availableNoTiles; i++ {
		s = s + fmt.Sprintf("t%v r%v i%v,", edgePairList.tiles[i].tile.tileNumber, edgePairList.tiles[i].rotation, edgePairList.tiles[i].tilepositionInEdgePairListIndex)
	}
	s = s + ")  **  "
	//
	for i := edgePairList.availableNoTiles; i < edgePairList.noTiles; i++ {
		s = s + fmt.Sprintf("t%v r%v [pp%v],", edgePairList.tiles[i].tile.tileNumber, edgePairList.tiles[i].rotation, edgePairList.tiles[i].previousPosition)
	}
	return s
}

// addTile adds a tile to an edgePairList, not the tileRotaton and index in tile edgePairLists is effectively the same
// rotations 0..3 match the position in the array 0..3
func (edgePairList *tileEdgePairList) addTile(tile *Tile, tileRotation int) {
	var tileAndRotation tileAndRotation

	tileAndRotation.tilepositionInEdgePairListIndex = tileRotation

	tileAndRotation.tile = tile
	//tileAndRotation.previousPosition = -1   // only used for debugging when printing out info
	tileAndRotation.rotation = tileRotation // (4 - tileRotation) % 4

	edgePairList.tiles[edgePairList.noTiles] = &tileAndRotation
	// Each tile tracks which edgePairList it is in and also its position in that list
	// so that when a tile is placed it can remove itself from those available in the list
	tile.edgePairLists[tileRotation] = edgePairList
	tile.positionInEdgePairList[tileRotation] = edgePairList.noTiles

	edgePairList.noTiles++
	edgePairList.availableNoTiles = edgePairList.noTiles
}

func (edgePairList *tileEdgePairList) removeTile(positionInList int, rotation int) {
	tileToRemove := edgePairList.tiles[positionInList]
	//fmt.Println("removeTile:removing tile no:", tileToRemove.tile.tileNumber, " in position:", positionInList, "from list:", edgePairList)
	// Remember it position (used when we restore it)
	tileToRemove.previousPosition = positionInList
	// get position of last tile in the list - we are going to swap the one we are removing with this one!
	positionLastTileInList := edgePairList.availableNoTiles - 1
	// copy the tile
	swapTile := edgePairList.tiles[positionLastTileInList] // remember tile at end of the list ...
	//fmt.Println("removeTile: swapped tile before position in list amended:", swapTile.tile, "Rotation:", swapTile.rotation)
	// move the tile we are removing to this position.
	edgePairList.tiles[positionLastTileInList] = tileToRemove
	// move the tile that was last in list to the place we took out the one we were removing
	edgePairList.tiles[positionInList] = swapTile
	swapTile.tile.positionInEdgePairList[swapTile.tilepositionInEdgePairListIndex] = positionInList // note if we do this after next line it breaks if this is last element in list!
	//tileToRemove.tile.positionInEdgePairList[rotation] = -1                                         // this is just for debug purposes. We don;t really care about its position when its been 'removed'
	//fmt.Println("removeTile:*removed tile state:", tileToRemove.tile)
	// decrement the number of tiles available in the list
	edgePairList.availableNoTiles--

	// Now we need to tell the tile that we swapped its new position in the edgePairList

	//fmt.Println("removeTile: swapped tile after position in list amended :", swapTile.tile)
	//fmt.Println("removeTile: list after removal:", edgePairList)
	//fmt.Println("removeTile: removed tile state:", tileToRemove.tile)
}

// restoreTile resores the last removed tile from the list
// it is located one after the end of the list
//
func (edgePairList *tileEdgePairList) restoreTile() { // HERE
	// get the  tile at one behond the "end of the list" that is going to be restored
	tileToRestore := edgePairList.tiles[edgePairList.availableNoTiles]
	// get previous position of that tile
	positionToRestoreTo := edgePairList.tiles[edgePairList.availableNoTiles].previousPosition
	// copy what was at that location
	swapTile := edgePairList.tiles[positionToRestoreTo] // remember tile at end of the list ...
	edgePairList.tiles[positionToRestoreTo] = tileToRestore
	edgePairList.tiles[edgePairList.availableNoTiles] = swapTile

	// Now we need to tell the tiles that we swapped its new position in the edgePairList
	swapTile.tile.positionInEdgePairList[swapTile.tilepositionInEdgePairListIndex] = edgePairList.availableNoTiles
	tileToRestore.tile.positionInEdgePairList[tileToRestore.tilepositionInEdgePairListIndex] = positionToRestoreTo

	edgePairList.availableNoTiles++
}

func createEdgePairLists(tiles tileArray, tileType byte) tileEdgePairMap {
	m := make(tileEdgePairMap)
	//fmt.Println("Building List:")
	for _, tile := range tiles {
		for i, edgePair := range tile.edgePairs {
			if _, ok := m[edgePair]; !ok {
				var edgePairList tileEdgePairList
				edgePairList.edgePairID = edgePair // this is just to make debugging easier :-)
				edgePairList.tileType = tileType
				m[edgePair] = &edgePairList
			}
			var edgePairList = m[edgePair] // note the suttle difference in type from above, this is a ponter, the other was the struct
			if edgePairList.tileType != tile.tileType {
				log.Fatalln("Missmatching tile type and list type!")
			}
			edgePairList.addTile(tile, i) // note the index and rotation are one and the same!
		}
	}
	return m
}
