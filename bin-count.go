package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

var binNamefileName = "binName.csv"
var binCountfileName = "binCount.csv"
var binNameCountfileName = "binNameCount.csv"

func main() {
	binNamefile, err := os.Open(binNamefileName)
	checkErrorq("Cannot open file binName", err)
	defer binNamefile.Close()
	rbinName := csv.NewReader(bufio.NewReader(binNamefile))
	wbinNameCountfile, err := os.Create(binNameCountfileName)
	checkErrorq("Cannot create file", err)
	defer wbinNameCountfile.Close()
	wbinNameCount := csv.NewWriter(wbinNameCountfile)
	defer wbinNameCount.Flush()
	wbinNameCount.Write([]string{"bin_no.","issuer", "name of variant", "count"})
	i := 0
	countFailed := 0
	for {
		data, err := rbinName.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if i == 0{
			i = 1
			continue
		}
		//checkErrorq("Not able to read from binName",err)
		binno, err := strconv.Atoi(data[0])
		checkErrorq("Conversion to int not possible", err)
		issuer := data[1]
		name := data[2]
		//log.Printf("Finding count for ")
		//binno = 486269
		count := FindCount(binno, name, issuer)
		if count == 0 {
			countFailed += 1
		}
		wbinNameCount.Write([]string{data[0], issuer, name, strconv.Itoa(count)})
	}
	log.Println(countFailed)
}

func FindCount(binno int, name, issuer string) int {
	binCountfile, err := os.Open(binCountfileName)
	checkErrorq("Cannot open file binCount", err)
	defer binCountfile.Close()
	rbinCount := csv.NewReader(bufio.NewReader(binCountfile))
	i := 0
	for {
		data, err := rbinCount.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if i == 0{
			i = 1
			continue
		}
		//checkErrorq("Not able to raed binCount", err)
		binno2,_ := strconv.Atoi(data[1])
		if binno == binno2{
			count, err := strconv.Atoi(data[0])
			checkErrorq("Conversion to int not possible", err)
			return count
		}
	}
	log.Printf("Can't find the count of bin no. = %d ; bank name = %s; variant name = %s", binno, issuer, name)
	return 0
}

func checkErrorq(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
