package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readTilesFromFile(fileName string) (width int, height int, tiles tileArray, err error) {

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	first := true
	tileNumber := 1
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text()) // split the line up into fields
		if first {
			first = false
			width, err = strconv.Atoi(fields[0])
			height, err = strconv.Atoi(fields[1])
			if err != nil {
				return
			}
			fmt.Println(width, height)

		} else {
			var newTile Tile
			newTilep := &newTile
			for i, v := range fields {
				var s int
				s, err = strconv.Atoi(v)   // convert to integer
				newTile.sides[i] = side(s) // Make side typesafe
				if err != nil {
					return
				}
			}

			newTile.tileNumber = tileNumber

			//fmt.Println(newTile)
			// TODO: tileType, edgePairs etc
			newTile.setTileProperties()
			tiles = append(tiles, newTilep)
			tileNumber++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(tiles)
	//for _, v := range tiles {
	//	fmt.Println(v)
	//}

	return
}
