package main

import (
	"fmt"
	"os"
)

// eh
var board Board
var tileSet TileSet
var highestProgress int

func main() {

	//argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) > 0 {
		for _, f := range argsWithoutProg {
			fmt.Println(f)
			width, height, tiles, err := readTilesFromFile(f)
			if err != nil {
				fmt.Println(err)
			} else {

				err := tileSet.setUpTileSet(width, height, tiles)
				if err != nil {
					fmt.Println(err)
				} else {
					// set up board

					err := board.createBoard(tileSet)
					if err != nil {
						fmt.Println(err)
					}
					//fmt.Println(board)
					//
					highestProgress = 0
					tiles[0].rotation = 0

					//
					// This is required to place 1st tile in top corner
					//
					loc := &board.loc[0][0]
					loc.edgePairList = loc.edgePairMap[calcEdgePairID(0, 0)]
					loc.edgePairList.needCount++
					tiles[0].placeTileOnBoard(loc, 1)
					//fmt.Println(tiles)
				}
			}

		}
	}

}
