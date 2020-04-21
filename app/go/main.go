package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DataCovid struct {
	Fecha          string
	Ccaa           string
	Casos          int
	Hospitalizados int
	Uci            int
	Fallecidos     int
	Recuperados    int
	Activos        int
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("usage: bin url outputPath")
		os.Exit(1)
	}
	url := os.Args[1]
	filename := "casos.csv"
	path := os.Args[2]

	err := DownloadFile(url, filename)
	if err != nil {
		panic(err)
	}

	// read data from CSV file

	csvFile, err := os.Open("casos.csv")

	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dt := time.Now()
	// Create a file for writing
	f, _ := os.Create(path + "/covid19_" + dt.Format("01-02-2006_15:04:05") + ".txt")

	// Create a writer
	w := bufio.NewWriter(f)

	for _, each := range csvData {
		var oneRecord DataCovid
		oneRecord.Ccaa = each[0]
		oneRecord.Fecha = each[1] + " 12:00:00"
		oneRecord.Casos, _ = strconv.Atoi(each[2])
		oneRecord.Hospitalizados, _ = strconv.Atoi(each[3])
		oneRecord.Uci, _ = strconv.Atoi(each[4])
		oneRecord.Fallecidos, _ = strconv.Atoi(each[5])
		oneRecord.Recuperados, _ = strconv.Atoi(each[6])
		oneRecord.Activos = oneRecord.Casos - (oneRecord.Fallecidos + oneRecord.Hospitalizados)
		jsondata, _ := json.Marshal(oneRecord) // convert to JSON
		w.WriteString(string(jsondata) + "\n")
	}

	// Very important to invoke after writing a large number of lines
	w.Flush()

	os.Remove(filename)
}

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
