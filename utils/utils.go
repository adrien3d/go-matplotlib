package utils

import (
	"encoding/csv"
	"github.com/snwfdhmp/errlog"
	"io"
	"log"
	"os"
	"strconv"
)

func CheckErr(e error) {
	if e != nil {
		errlog.Debug(e)
		panic(e)
	}
}

func OpenCSV(fileName string) (ret [][]float64) {
	csvfile, err := os.Open(fileName + ".csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)
	for {
		line, err := r.Read()
		var lineData []float64
		if err == io.EOF {
			break
		}
		CheckErr(err)
		for j := 0; j < len(line); j++ {
			val, err := strconv.ParseFloat(line[j], 64)
			CheckErr(err)
			lineData = append(lineData, val)
		}
		ret = append(ret, lineData)
	}
	return ret
}
