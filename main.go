package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ips = [8]string{"64.233", "66.102", "66.249", "72.14", "74.125", "209.85", "216.239", "66.184"}

func getLogs(pattern string) []string {
	fileList := []string{}

	files, _ := filepath.Glob(pattern)

	for _, f := range files {
		fileList = append(fileList, f)
	}

	return fileList
}

func checkuseragent(logfile string) []string {
	googlebothits := []string{}

	file, err := os.Open(logfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		splitLine := strings.Split(currentLine, " ")
		useragent := strings.Join(splitLine[11:], " ")
		googlebot := strings.Contains(useragent, "Googlebot")
		if googlebot {
			googlebothits = append(googlebothits, currentLine)
		}
	}

	return googlebothits
}

func checkip(googlehits []string) []string {
	validatedhits := []string{}

	for _, f := range googlehits {
		splitLine := strings.Split(f, " ")
		ipAddress := splitLine[0]
		for _, ip := range ips {
			validated := strings.Contains(ipAddress, ip)
			if validated {
				validatedhits = append(validatedhits, f)
				break
			}
		}
	}
	return validatedhits
}

func fixDate(date string) (string, string) {
	date = strings.Replace(date, "[", "", 1)
	dateTime := strings.Split(date, ":")
	date = dateTime[0]
	time := strings.Join(dateTime[1:], ":")
	return date, time
}

func formatOutput(line string) []string {
	splitLine := strings.Split(line, " ")
	ipAddress := splitLine[0]
	date, time := fixDate(splitLine[3])
	url := splitLine[6]
	status := splitLine[8]
	numbytes := splitLine[9]
	useragent := strings.Join(splitLine[11:], " ")
	csvRow := []string{}
	csvRow = append(csvRow, ipAddress, date, time, url, status, numbytes, useragent)
	return csvRow
}

func resultsToCsv(verifiedhits []string) {
	records := [][]string{}

	for _, f := range verifiedhits {
		result := formatOutput(f)
		records = append(records, result)
	}

	file, err := os.Create("logs.csv")
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range records {
		err := writer.Write(value)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	list := getLogs("*log")

	for _, f := range list {
		res := checkuseragent(f)
		res2 := checkip(res)
		resultsToCsv(res2)
	}
}
