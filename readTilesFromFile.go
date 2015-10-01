package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readTilesFromFile(fileName string) (width int, height int, tiles []Tile, err error) {

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
			for i, v := range fields {
				newTile.sides[i], err = strconv.Atoi(v)
				if err != nil {
					return
				}
			}

			newTile.tileNumber = tileNumber

			//fmt.Println(newTile)
			// TODO: tileType, edgePairs etc
			newTile.setTileProperties()
			tiles = append(tiles, newTile)
			tileNumber++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(tiles)
	// TODO: check validity of tile set , no of tiles = width*height, no of sides are even etc....
	return
}
