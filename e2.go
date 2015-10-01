package main

import  (
	"os"
	"fmt"
  "readTilesFromFile"
)


func main() {

	//argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) > 0 {
		for _,f := range argsWithoutProg {
			fmt.Println(f)
			readTilesFromFile(f)
		}
	}

}
