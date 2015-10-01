package main

import (
	"fmt"
	"os"
)

// eh

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
				var tileSet TileSet
				err := tileSet.setUpTileSet(width, height, tiles)
				if err != nil {
					fmt.Println(err)
				}
			}

		}
	}

}
