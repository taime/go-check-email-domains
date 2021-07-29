package main

import (
	"bufio"
	_ "embed"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	var data = [][]string{{"Email", "Status Code", "SSL"}}
	email, err := readLines(dir + "/emails.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	} else {
		for i, email := range email {
			fmt.Println(i)
			fmt.Println(email)
			url, success := emailToDomain(email)
			if success {
				var dmn, scode = checkUrl("http://" + url)
				fmt.Println(dmn)
				fmt.Println(scode)
				var newLine = []string{email, scode, "http"}
				data = append(data, newLine)

				var dmn2, scode2 = checkUrl("https://" + url)
				fmt.Println(dmn2)
				fmt.Println(scode2)
				var newLine2 = []string{email, scode2, "https"}
				data = append(data, newLine2)
			}
		}
	}

	writeToFileData(dir+"/result.csv", data)
}

func writeToFileData(filepath string, data [][]string) {
	file, err := os.Create(filepath)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
}

func checkUrl(url string) (string, string) {

	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return url, "None"
	}
	defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// fmt.Println(body)
	scode := strconv.Itoa(resp.StatusCode)
	return url, scode

}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func emailToDomain(email string) (string, bool) {

	at := strings.LastIndex(email, "@")
	if at >= 0 {
		username, domain := email[:at], email[at+1:]
		fmt.Printf("Username: %s, Domain: %s\n", username, domain)
		return (domain), true
	} else {
		fmt.Printf("Error: %s is an invalid email address\n", email)
		return (""), false
	}

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
