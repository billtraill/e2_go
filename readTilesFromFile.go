package readTilesFromFile

import  (
	"os"
	"fmt"
	"bufio"
	"log"
	"strconv"
	"strings"
)


func readTilesFromFile(fileName string) error {
	var width, height int

	file, err := os.Open(fileName)
	if err != nil {
	    log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	first := true

	for scanner.Scan()  {
		result := strings.Fields(scanner.Text())
		if first {
			first = false
			width, err  = strconv.Atoi(result[0])
			height, err = strconv.Atoi(result[1])
			if err != nil {
					return  err
			}
			fmt.Println(width,height)
		} else {
			e1,err := strconv.Atoi(result[0])
			e2,err := strconv.Atoi(result[1])
			e3,err := strconv.Atoi(result[2])
			e4,err := strconv.Atoi(result[3])
			if err != nil {
					return  err
			}
			fmt.Println(e1,e2,e3,e4,e4)
		}
	}
	if err := scanner.Err(); err != nil {
	    log.Fatal(err)
	}

	return err
}
