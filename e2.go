package main

import (
	"fmt"
	"os"
)

// eh
var board Board
var tileSet TileSet
var highest_progress int

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
					highest_progress = 0
					tiles[0].placeTileOnBoard(BoardPosition{0, 0}, 0, 1)
					//fmt.Println(tiles)
				}
			}

		}
	}

}
