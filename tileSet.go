package main

import "fmt"

// TileSet : Holds the set of all tiles to be used
// gets initialised once and no further changes done to it.
// Also gets validated for types, numbers etc
type TileSet struct {
	width       int
	height      int
	cornerTiles []*Tile
	edgeTiles   []*Tile
	normalTiles []*Tile
}

/*
// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}


// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}
*/

// setUpTileSet takes width and height of a board and a set of tiles to cover the board and populates
// the given tileSet
func (tileSet *TileSet) setUpTileSet(width int, height int, tiles []Tile) error {
	tileSet.width = width
	tileSet.height = height

	for i := range tiles {
		switch tiles[i].tileType {
		case 'C':
			tileSet.cornerTiles = append(tileSet.cornerTiles, &tiles[i])
		case 'E':
			tileSet.edgeTiles = append(tileSet.edgeTiles, &tiles[i])
		case 'N':
			tileSet.normalTiles = append(tileSet.normalTiles, &tiles[i])
		}
	}

	// check we have the correct number of tiles for the shape of the board
	numberOfTiles := len(tiles)
	if numberOfTiles != (width * height) {
		return fmt.Errorf("Number of tiles:%v does not match width:%v height:%v ", numberOfTiles, width, height)
	}
	// check we have the correct number of corners/edges/normal tiles
	numberOfCorners := len(tileSet.cornerTiles)
	if numberOfCorners != 4 {
		return fmt.Errorf("Only %v corners. Should be 4", numberOfCorners)
	}
	// check we have the correct number of edge tiles
	numberOfEdges := len(tileSet.edgeTiles)
	requiredNumberOfEdges := 2 * ((width - 2) + (height - 2))
	if numberOfEdges != requiredNumberOfEdges {
		return fmt.Errorf("Only %v edges. Should be %v", numberOfEdges, requiredNumberOfEdges)
	}
	return nil
}
