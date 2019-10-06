package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

// Problem struct stores the question and the related answer
type Problem struct {
	question string
	answer   string
}

func readArguments() string {
	filename := flag.String("file", "problems.csv", "CSV file with the quiz questions")
	flag.Parse()
	return *filename
}

func openCSV(filename string) *os.File {
	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	return csvfile
}

func readCSVProblems(r *csv.Reader) []Problem {
	var problems []Problem
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		problems = append(problems, Problem{question: record[0], answer: record[1]})
	}
	return problems
}

func readInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func startQuiz(problems []Problem) {
	var total int
	for i, problem := range problems {
		fmt.Printf("Problem #%v: %v = ", i+1, problem.question)
		input := readInput()
		if input == problem.answer {
			total++
		}
	}
	fmt.Printf("You scored %v out of %v.", total, len(problems))
}

func main() {
	filename := readArguments()
	csvfile := openCSV(filename)
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	problems := readCSVProblems(r)

	startQuiz(problems)
}
