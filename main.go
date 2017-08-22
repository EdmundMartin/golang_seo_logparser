package main

import (
	"bufio"
	"fmt"
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

func main() {
	list := getLogs("*log")
	fmt.Print(list)

	for _, f := range list {
		res := checkuseragent(f)
		fmt.Print(res)

	}
}
