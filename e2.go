package main

import (
	"flag"
	"fmt"
)

// eh
var board Board
var tileSet TileSet
var highestProgress int

func main() {

	// boolPtr := flag.Bool("fork", false, "a bool")
	methodPtr := flag.String("method", "iter", "Method for traversing the board, default is iter - iterative, or recursive - recursive")
	composite := flag.Bool("comp", false, "if set uses compisite tiles ")

	flag.Parse()
	//argsWithProg := os.Args
	argsWithoutProg := flag.Args()
	fmt.Println("Files:", argsWithoutProg)
	if len(argsWithoutProg) > 0 {
		for _, f := range argsWithoutProg {
			fmt.Println(f)
			width, height, tiles, err := readTilesFromFile(f)
			if err != nil {
				fmt.Println(err)
			} else {

				err := tileSet.setUpTileSet(width, height, tiles, *composite)
				if err != nil {
					fmt.Println(err)
				} else {
					// set up board

					if *composite {

						//fmt.Println(tileSet.cornerTilesEdgePairsMap)
						//fmt.Println(tileSet.edgeTilesEdgePairsMap)
						//fmt.Println(tileSet.normalTilesEdgePairsMap)
						if width%2 == 1 || height%2 == 1 {
							fmt.Println("Composite tiles can only be used on even sized boards")
							return
						}
						err := board.createBoard(tileSet, width/2, height/2)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println(board)
						traverseCompositeBoard()
					} else {
						err := board.createBoard(tileSet, width, height)
						if err != nil {
							fmt.Println(err)
						}
						if *methodPtr == "iter" {
							traverseBoard()
						} else if *methodPtr == "recursive" {
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
							loc.tile = tiles[0]
							loc.tile.rotation = 0
							loc.placeTileOnBoard(1)
							//fmt.Println(tiles)
						}
					}
				}
			}

		}
	}

}
