package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/davidminor/uint128"
)

const maxEdgePairListSize = 600 // Probably only 17

// tileEdgePairList an list that refers to each tile with a given edgePairID
type tileEdgePairList struct { // equiv to EPListType
	tiles              [maxEdgePairListSize]*tileAndRotation // pointer to current tile at
	totalNoTilesInList int                                   // for information only - used in construction of lists, not really used operationaly!
	availableNoTiles   int                                   // the number of tiles in the list available to be placed.
	needCount          int                                   // as we reserve elements on lookahead this is incremented, it cannot be bigger than availableNoTiles
	edgePairID         edgePairID
	tileType           byte
	//removeChan         chan int
	//restoreChan        chan int
	//responseChan       chan int
	cTilePresent uint128.Uint128 // bit mask for each tile present in this list
}
type tileAndRotation struct {
	tile                            *Tile
	tilepositionInEdgePairListIndex int // we need to know this as when we move the tiles in the edgePosition lists this tells us the lists index in the tile! This does not change
	rotationForEdgePair             int // tile rotaton note this is its (4-index)%4
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

	s := fmt.Sprintf("%v: %v Size:%v (", edgePairDescription(edgePairList.edgePairID), tileTypeDescription(edgePairList.tileType), edgePairList.availableNoTiles)
	for i := 0; i < edgePairList.availableNoTiles; i++ {
		s = s + fmt.Sprintf("t%v r%v i%v,", edgePairList.tiles[i].tile.tileNumber, edgePairList.tiles[i].rotationForEdgePair, edgePairList.tiles[i].tilepositionInEdgePairListIndex)
	}
	s = s + ")  **  "
	//
	for i := edgePairList.availableNoTiles; i < edgePairList.totalNoTilesInList; i++ {
		s = s + fmt.Sprintf("t%v r%v [pp%v],", edgePairList.tiles[i].tile.tileNumber, edgePairList.tiles[i].rotationForEdgePair, edgePairList.tiles[i].previousPosition)
	}
	return s
}

// addTile adds a tile to an edgePairList, not the tileRotaton and index in tile edgePairLists is effectively the same
// rotations 0..3 match the position in the array 0..3
func (edgePairList *tileEdgePairList) addTile(tile *Tile, tileRotation int) {
	var tileAndRotation tileAndRotation
	if edgePairList.totalNoTilesInList == maxEdgePairListSize {
		log.Fatalln("Exceeded static edgePair list size", maxEdgePairListSize, tile)
	}
	tileAndRotation.tilepositionInEdgePairListIndex = tileRotation

	tileAndRotation.tile = tile
	//tileAndRotation.previousPosition = -1   // only used for debugging when printing out info
	tileAndRotation.rotationForEdgePair = tileRotation // (4 - tileRotation) % 4

	edgePairList.tiles[edgePairList.totalNoTilesInList] = &tileAndRotation
	// Each tile tracks which edgePairList it is in and also its position in that list
	// so that when a tile is placed it can remove itself from those available in the list
	tile.edgePairLists[tileRotation] = edgePairList
	tile.positionInEdgePairList[tileRotation] = edgePairList.totalNoTilesInList

	edgePairList.totalNoTilesInList++
	edgePairList.availableNoTiles = edgePairList.totalNoTilesInList

	edgePairList.cTilePresent = edgePairList.cTilePresent.Or(tile.cTileUsed) // track the primitive tiles present in list
}

//var swapTile *tileAndRotation

func (edgePairList *tileEdgePairList) removeTile(positionInList int) {

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

	// decrement the number of tiles available in the list
	edgePairList.availableNoTiles--

}

var complete = make(chan int, 5) // used by the concurrent versions of remove/restore - not used in current solutions as it worked out slower than sequential
func (edgePairList *tileEdgePairList) goRemoveTile(positionInList int) {

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

	// decrement the number of tiles available in the list
	edgePairList.availableNoTiles--

	complete <- 1
}

// restoreTile resores the last removed tile from the list
// it is located one after the end of the list
//
func (edgePairList *tileEdgePairList) restoreTile() {

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
func (edgePairList *tileEdgePairList) goRestoreTile() {

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

	complete <- 1
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
				//edgePairList.removeChan = make(chan int, 4) // should only be 2 concurrent requests against an edgepair list
				//edgePairList.restoreChan = make(chan int, 4)
				//edgePairList.responseChan = responseChan
				//go edgePairList.restoreTile()
				//go edgePairList.removeTile()
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
