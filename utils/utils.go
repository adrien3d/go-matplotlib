package utils

import (
	"encoding/csv"
	"github.com/snwfdhmp/errlog"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func CheckErr(e error) {
	if e != nil {
		errlog.Debug(e)
		panic(e)
	}
}

func OpenCSV(fileName string) (datasetNames []string, ret [][]float64) {
	csvfile, err := os.Open(fileName + ".csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)
	firstLine := true
	for {
		line, err := r.Read()
		var lineData []float64
		if err == io.EOF {
			break
		}
		if firstLine {
			for i := 0; i < len(line); i++ {
				datasetNames = append(datasetNames, strings.Trim(line[i], " "))
			}
			firstLine = false
		} else {
			for j := 0; j < len(line); j++ {
				val, err := strconv.ParseFloat(line[j], 64)
				CheckErr(err)
				lineData = append(lineData, val)
			}
			ret = append(ret, lineData)
		}
	}
	return datasetNames, ret
}
