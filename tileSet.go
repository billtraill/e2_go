package main

import "fmt"

// TileSet : Holds the set of all tiles to be used
// gets initialised once and no further changes done to it.
// Also gets validated for types, numbers etc
type TileSet struct {
	width                 int
	height                int
	cornerTiles           []*Tile
	edgeTiles             []*Tile
	normalTiles           []*Tile
	allTiles              []*Tile // we probably just need this one and its associated map
	numberOfEdgeColours   int
	numberOfNormalColours int
	//tilesEdgePairMap      tileEdgePairMap // maybe we can just getaway with one map? as boarder colours don't apear on normal tiles!
	cornerTilesEdgePairsMap tileEdgePairMap
	edgeTilesEdgePairsMap   tileEdgePairMap
	normalTilesEdgePairsMap tileEdgePairMap
	//responseChan            chan int // used to respond to completion notices....
}

// setUpTileSet takes width and height of a board and a set of tiles to cover the board and populates
// the given tileSet
func (tileSet *TileSet) setUpTileSet(width int, height int, tiles tileArray) (err error) {
	tileSet.width = width
	tileSet.height = height
	tileSet.cornerTiles = nil
	tileSet.edgeTiles = nil
	tileSet.normalTiles = nil
	tileSet.allTiles = nil
	tileSet.numberOfEdgeColours = 0
	tileSet.numberOfNormalColours = 0

	for i := range tiles {
		tileSet.allTiles = append(tileSet.allTiles, tiles[i])
		// the rest might not be needed !
		switch tiles[i].tileType {
		case 'C':
			tileSet.cornerTiles = append(tileSet.cornerTiles, tiles[i])
		case 'E':
			tileSet.edgeTiles = append(tileSet.edgeTiles, tiles[i])
		case 'N':
			tileSet.normalTiles = append(tileSet.normalTiles, tiles[i])
		}
	}
	err = tileSet.checkTileSetIntegrety()
	if err != nil {
		return err
	}
	//tileSet.responseChan = make(chan int, 4) // we expect only 4 responses at a time (1 for each edgepair list)
	// Setup edgePairsmaps
	//tileSet.tilesEdgePairMap = createEdgePairLists(tileSet.allTiles)
	//fmt.Println(tileSet.tilesEdgePairMap)
	tileSet.cornerTilesEdgePairsMap = createEdgePairLists(tileSet.cornerTiles, 'C')
	//fmt.Println("Corner list:")
	//fmt.Println(tileSet.cornerTilesEdgePairsMap)
	tileSet.edgeTilesEdgePairsMap = createEdgePairLists(tileSet.edgeTiles, 'E')
	//fmt.Println("Edge list:")
	//fmt.Println(tileSet.edgeTilesEdgePairsMap)
	tileSet.normalTilesEdgePairsMap = createEdgePairLists(tileSet.normalTiles, 'N')
	//fmt.Println("Normal list:")
	//fmt.Println(tileSet.normalTilesEdgePairsMap)
	//fmt.Println(tiles)
	// HERE TODO need to add functions to add/remove entries from the lists. Use array as small linked list! swaping used elements to end of list!
	// to implement forward reservation of EP entries and constrain on count of available!
	return nil
}

func (tileSet *TileSet) checkTileSetIntegrety() (err error) {
	// check we have the correct number of tiles for the shape of the board
	numberOfTiles := len(tileSet.cornerTiles) + len(tileSet.edgeTiles) + len(tileSet.normalTiles)
	if numberOfTiles != (tileSet.width * tileSet.height) {
		return fmt.Errorf("Number of tiles:%v does not match width:%v height:%v ", numberOfTiles, tileSet.width, tileSet.height)
	}
	// check we have the correct number of corners/edges/normal tiles
	numberOfCorners := len(tileSet.cornerTiles)
	if numberOfCorners != 4 {
		return fmt.Errorf("Only %v corners. Should be 4", numberOfCorners)
	}
	// check we have the correct number of edge tiles
	numberOfEdges := len(tileSet.edgeTiles)
	requiredNumberOfEdges := 2 * ((tileSet.width - 2) + (tileSet.height - 2))
	if numberOfEdges != requiredNumberOfEdges {
		return fmt.Errorf("Only %v edges. Should be %v", numberOfEdges, requiredNumberOfEdges)
	}
	// check edge colours
	err = tileSet.checkEdgeColours()
	if err != nil {
		return err
	}
	err = tileSet.checkNormalColours()
	if err != nil {
		return err
	}

	return nil
}

func (tileSet *TileSet) checkEdgeColours() error {
	var colours [32 + 1]int
	borderEdgeIndexes := []int{1, 3} // this is the border edge indexes. must be in these positions are all tiles sides are normalised
	for _, tile := range tileSet.edgeTiles {
		for _, i := range borderEdgeIndexes {
			v := tile.sides[i]
			if v > 32 {
				return fmt.Errorf("Max edge colour exceeded:%v Max:%v", v, 32)
			}
			colours[v]++
		}
	}
	cornerEdgeIndexes := []int{2, 3} // this is the border corner indexes. must be in these positions are all tiles sides are normalised
	for _, tile := range tileSet.cornerTiles {
		for _, i := range cornerEdgeIndexes {
			v := tile.sides[i]
			if v > 32 {
				return fmt.Errorf("Max edge colour exceeded:%v Max:%v", v, 32)
			}
			colours[v]++
		}
	}
	// Check that each colour is even number
	maxColour := 0
	for i, v := range colours {
		if v > 0 {
			maxColour = i
			if v%2 != 0 {
				return fmt.Errorf("Missmatching edge colour count for colour:%v", i)
			}
		}
	}
	tileSet.numberOfEdgeColours = maxColour
	return nil
}

func (tileSet *TileSet) checkNormalColours() error {
	var colours [32 + 1]int
	borderEdgeIndexes := []int{2} // this is the border edge indexes. must be in these positions are all tiles sides are normalised
	for _, tile := range tileSet.edgeTiles {
		for _, i := range borderEdgeIndexes {
			v := tile.sides[i]
			if v > 32 {
				return fmt.Errorf("Max edge colour exceeded:%v Max:%v", v, 32)
			}
			colours[v]++
		}
	}
	// Now check rest of normal tiles
	for _, tile := range tileSet.normalTiles {
		for _, v := range tile.sides {
			if v > 32 {
				return fmt.Errorf("Max edge colour exceeded:%v Max:%v", v, 32)
			}
			colours[v]++
		}
	}
	// Check that each colour is even number
	maxColour := 0
	for i, v := range colours {
		if v > 0 {
			maxColour = i
			if v%2 != 0 {
				return fmt.Errorf("Missmatching centre colour count for colour:%v", i)
			}
		}
	}
	tileSet.numberOfNormalColours = maxColour
	return nil
}
