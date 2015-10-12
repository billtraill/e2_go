package main

import (
	"fmt"
	"log"
)

import "github.com/davidminor/uint128"

func cTileSide(a int, b int) side {

	cside := a<<8 + b

	return side(cside)
}

func cTileSideSwap(s side) side {

	a := s >> 8
	b := s & 0xFF
	cside := b<<8 + a

	return side(cside)
}

func (board *Board) createCompositTileBoard(tileSet TileSet, width int, height int, tileType byte) error {
	board.width = width
	board.height = height
	board.loc = make([][]BoardLocation, board.height)
	// loop over the rows allocating the slice for each row
	for y := range board.loc {
		board.loc[y] = make([]BoardLocation, board.width)
	}

	board.setTraversal()
	return nil
}

func createTileSideLists(tiles tileArray, tileType byte) tileEdgePairMap {
	m := make(tileEdgePairMap)
	//fmt.Println("Building List:")
	for _, tile := range tiles {
		for i, side := range tile.sides {
			if _, ok := m[edgePairID(side)]; !ok {
				var edgePairList tileEdgePairList
				edgePairList.edgePairID = edgePairID(side) // this is just to make debugging easier :-)
				edgePairList.tileType = tileType
				m[edgePairID(side)] = &edgePairList
			}
			var edgePairList = m[edgePairID(side)] // note the suttle difference in type from above, this is a ponter, the other was the struct
			if edgePairList.tileType != tile.tileType {
				log.Fatalln("Missmatching tile type and list type!")
			}
			edgePairList.addTile(tile, i) // note the index and rotation are one and the same!
		}
	}
	return m
}

func createCompositeTile(t0 *Tile, t1 *Tile, t2 *Tile, t3 *Tile, tileType byte) Tile {
	var cTile Tile
	cTile.cTiles[0] = t0
	cTile.cTiles[1] = t1
	cTile.cTiles[2] = t2
	cTile.cTiles[3] = t3
	cTile.cTileRotations[0] = t0.rotation
	cTile.cTileRotations[1] = t1.rotation
	cTile.cTileRotations[2] = t2.rotation
	cTile.cTileRotations[3] = t3.rotation

	cTile.sides[0] = cTileSide(int(t0.sides[(t0.rotation+0)%4]), int(t1.sides[(t1.rotation+0)%4]))
	cTile.sides[1] = cTileSide(int(t1.sides[(t1.rotation+1)%4]), int(t2.sides[(t2.rotation+1)%4]))
	cTile.sides[2] = cTileSide(int(t2.sides[(t2.rotation+2)%4]), int(t3.sides[(t3.rotation+2)%4]))
	cTile.sides[3] = cTileSide(int(t3.sides[(t3.rotation+3)%4]), int(t0.sides[(t0.rotation+3)%4]))
	//var used uint128.Uint128
	one := uint128.Uint128{}
	one.L = 1
	cTile.cTileUsed = cTile.cTileUsed.Add(one.ShiftLeft(uint(t0.tileNumber - 1)))
	cTile.cTileUsed = cTile.cTileUsed.Add(one.ShiftLeft(uint(t1.tileNumber - 1)))
	cTile.cTileUsed = cTile.cTileUsed.Add(one.ShiftLeft(uint(t2.tileNumber - 1)))
	cTile.cTileUsed = cTile.cTileUsed.Add(one.ShiftLeft(uint(t3.tileNumber - 1)))
	cTile.composite = true
	cTile.tileType = tileType
	cTile.tileNumber = 0 // maybe set this to something else in future

	cTile.setEdgePairs()
	return cTile
}

func generateCompositeTilesNormal() {
	size := 2
	fmt.Println("generateCompositeTilesNormal")

	//for i := range tileSet.normalTiles {
	//	fmt.Println("Index", i, "Tile no", tileSet.normalTiles[i].tileNumber)
	//}
	var tBoard Board
	tBoard.createCompositTileBoard(tileSet, size, size, 'N')
	//fmt.Println(tBoard)
	tileSidesMap := createTileSideLists(tileSet.normalTiles, 'N')
	//fmt.Println(tileSidesMap)
	// Place Tile on board 0..3
	//var compTile compTileType
	noCompositeTiles := 0
	for _, baseTile := range tileSet.normalTiles {

		for r := 0; r < 4; r++ {
			// four rotations of our starting tile!
			t0 := baseTile
			loc := &tBoard.loc[1][0]
			loc.tile = t0 // just do first tile in list
			loc.tile.rotation = r
			sidesList := tileSidesMap[edgePairID(t0.sides[(r+1)%4])]

			for t1i := 0; t1i < sidesList.availableNoTiles; t1i++ {
				t1 := sidesList.tiles[t1i].tile
				if t0 != t1 && t1.tileNumber > t0.tileNumber {
					tBoard.loc[0][0].tile = t1
					tBoard.loc[0][0].tile.rotation = (1 + sidesList.tiles[t1i].rotationForEdgePair) % 4

					//fmt.Println(tBoard)
					//fmt.Println("BOARD:", board.loc[0][0].tile.tileNumber, board.loc[1][0].tile.tileNumber)
					t1Side := t1.sides[(t1.rotation+2)%4]
					// fmt.Println("t1 side:", t1Side)
					sidesList2 := tileSidesMap[edgePairID(t1Side)]

					for t2i := 0; t2i < sidesList2.availableNoTiles; t2i++ {
						t2 := sidesList2.tiles[t2i].tile
						if t0 != t1 && t1 != t2 && t0 != t2 && t2.tileNumber > t0.tileNumber {
							tBoard.loc[0][1].tile = t2
							tBoard.loc[0][1].tile.rotation = (0 + sidesList2.tiles[t2i].rotationForEdgePair) % 4

							t2Side := t2.sides[(t2.rotation+3)%4]
							//fmt.Println("t2 side:", t2Side)
							sidesList3 := tileSidesMap[edgePairID(t2Side)]
							for t3i := 0; t3i < sidesList3.availableNoTiles; t3i++ {
								t3 := sidesList3.tiles[t3i].tile
								if t0 != t1 && t1 != t2 && t0 != t2 && t0 != t3 && t1 != t3 && t2 != t3 && t3.tileNumber > t0.tileNumber {
									tBoard.loc[1][1].tile = t3
									tBoard.loc[1][1].tile.rotation = (3 + sidesList3.tiles[t3i].rotationForEdgePair) % 4
									if t0.sides[(t0.rotation+2)%4] == t3.sides[(t3.rotation+0)%4] {
										//fmt.Println(tBoard)
										noCompositeTiles++
										cTile := createCompositeTile(t0, t1, t2, t3, 'N')
										tileSet.compositNormalTiles = append(tileSet.compositNormalTiles, &cTile)
										//fmt.Println(&cTile)
									} // t3 t0 validation
								} // t3 validation check
							} // t3
						} // t2 validation check
					} // t2
				} // t1 validation check
			} // t1
		} // t0 rotations

	} // t0
	fmt.Println("No of combinatons of composite middle tiles:", noCompositeTiles)
}

func generateCompositeTilesEdge() {
	size := 2
	fmt.Println("generateCompositeTilesEdge")

	//fmt.Println(&tileSet)
	var tBoard Board
	tBoard.createCompositTileBoard(tileSet, size, size, 'N')
	//fmt.Println(tBoard)
	edgeSidesMap := createTileSideLists(tileSet.edgeTiles, 'E')
	tileSidesMap := createTileSideLists(tileSet.normalTiles, 'N')
	//fmt.Println(edgeSidesMap)
	// Place Tile on board 0..3
	//var compTile compTileType
	noCompositeTiles := 0
	for _, t0 := range tileSet.edgeTiles {

		t0.rotation = 0
		loc := &tBoard.loc[1][0]
		loc.tile = t0 // just do first tile in list

		sidesList := edgeSidesMap[edgePairID(t0.sides[(1)%4])]
		//fmt.Println("BOARD [1,0]:", tBoard.loc[1][0].tile.tileNumber, tBoard.loc[1][0].tile.rotation)
		for t1i := 0; t1i < sidesList.availableNoTiles; t1i++ {
			t1 := sidesList.tiles[t1i].tile
			if t1 == t0 {
				continue
			}
			t1.rotation = (1 + sidesList.tiles[t1i].rotationForEdgePair) % 4

			if t0 != t1 && t1.rotation == 0 {
				tBoard.loc[0][0].tile = t1

				//fmt.Println(tBoard)
				//fmt.Println("BOARD:", tBoard.loc[0][0].tile.tileNumber, tBoard.loc[0][0].tile.rotation, tBoard.loc[1][0].tile.tileNumber, tBoard.loc[1][0].tile.rotation)

				t1Side := t1.sides[(t1.rotation+2)%4]
				// fmt.Println("t1 side:", t1Side)
				sidesList2 := tileSidesMap[edgePairID(t1Side)]

				for t2i := 0; t2i < sidesList2.availableNoTiles; t2i++ {
					t2 := sidesList2.tiles[t2i].tile

					tBoard.loc[0][1].tile = t2
					tBoard.loc[0][1].tile.rotation = (0 + sidesList2.tiles[t2i].rotationForEdgePair) % 4

					t2Side := t2.sides[(t2.rotation+3)%4]
					//fmt.Println("t2 side:", t2Side)
					sidesList3 := tileSidesMap[edgePairID(t2Side)]
					for t3i := 0; t3i < sidesList3.availableNoTiles; t3i++ {
						t3 := sidesList3.tiles[t3i].tile
						if t2 != t3 {
							tBoard.loc[1][1].tile = t3
							tBoard.loc[1][1].tile.rotation = (3 + sidesList3.tiles[t3i].rotationForEdgePair) % 4
							if t0.sides[(t0.rotation+2)%4] == t3.sides[(t3.rotation+0)%4] {
								//fmt.Println(tBoard)
								noCompositeTiles++
								cTile := createCompositeTile(t0, t1, t2, t3, 'E')
								tileSet.compositEdgeTiles = append(tileSet.compositEdgeTiles, &cTile)
							} // t3 t0 validation
						} // t3 validation check
					} // t3

				} // t2

			} // t1 validation check
		} // t1

	} // t0
	fmt.Println("No of combinatons of composite edge tiles:", noCompositeTiles)

}

func generateCompositeTilesCorner() {
	size := 2
	fmt.Println("generateCompositeTilesCorner")

	//fmt.Println(&tileSet)
	var tBoard Board
	tBoard.createCompositTileBoard(tileSet, size, size, 'N')
	//fmt.Println(tBoard)
	edgeSidesMap := createTileSideLists(tileSet.edgeTiles, 'E')
	tileSidesMap := createTileSideLists(tileSet.normalTiles, 'N')
	cornerSidesMap := createTileSideLists(tileSet.cornerTiles, 'C')
	//fmt.Println(edgeSidesMap)
	// Place Tile on board 0..3
	//var compTile compTileType
	noCompositeTiles := 0
	for _, t0 := range tileSet.edgeTiles {

		t0.rotation = 0
		loc := &tBoard.loc[1][0]
		loc.tile = t0 // just do first tile in list

		sidesList, ok := cornerSidesMap[edgePairID(t0.sides[(1)%4])]
		if ok {
			//fmt.Println("BOARD [1,0]:", tBoard.loc[1][0].tile.tileNumber, tBoard.loc[1][0].tile.rotation)
			for t1i := 0; t1i < sidesList.availableNoTiles; t1i++ {
				t1 := sidesList.tiles[t1i].tile
				t1.rotation = (1 + sidesList.tiles[t1i].rotationForEdgePair) % 4

				if t1.rotation == 0 {
					tBoard.loc[0][0].tile = t1

					t1Side := t1.sides[(t1.rotation+2)%4]
					// fmt.Println("t1 side:", t1Side)
					sidesList2 := edgeSidesMap[edgePairID(t1Side)]

					for t2i := 0; t2i < sidesList2.availableNoTiles; t2i++ {
						t2 := sidesList2.tiles[t2i].tile

						if t2 == t0 {
							continue
						}

						tBoard.loc[0][1].tile = t2
						tBoard.loc[0][1].tile.rotation = (0 + sidesList2.tiles[t2i].rotationForEdgePair) % 4

						if t2.rotation != 3 {
							continue
						}

						t2Side := t2.sides[(t2.rotation+3)%4]

						sidesList3 := tileSidesMap[edgePairID(t2Side)]
						for t3i := 0; t3i < sidesList3.availableNoTiles; t3i++ {
							t3 := sidesList3.tiles[t3i].tile
							if t2 != t3 {
								tBoard.loc[1][1].tile = t3
								tBoard.loc[1][1].tile.rotation = (3 + sidesList3.tiles[t3i].rotationForEdgePair) % 4
								if t0.sides[(t0.rotation+2)%4] == t3.sides[(t3.rotation+0)%4] {
									//fmt.Println(tBoard)
									noCompositeTiles++
									cTile := createCompositeTile(t0, t1, t2, t3, 'C')
									tileSet.compositCornerTiles = append(tileSet.compositCornerTiles, &cTile)
									//fmt.Println(&cTile)
								} // t3 t0 validation
							} // t3 validation check
						} // t3

					} // t2

				} // t1 validation check
			} // t1
		}
	} // t0
	fmt.Println("No of combinatons of composite corner tiles:", noCompositeTiles)

}

func generateTileCombinations() {
	generateCompositeTilesNormal()
	generateCompositeTilesEdge()
	generateCompositeTilesCorner()

	//fmt.Println(tileSet.compositCornerTiles)
	//fmt.Println(tileSet.compositEdgeTiles)
	//fmt.Println(tileSet.compositNormalTiles)
}
