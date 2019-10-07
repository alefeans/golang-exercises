package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Problem struct stores the question and the related answer
type Problem struct {
	question string
	answer   string
}

func readArguments() (string, int) {
	filename := flag.String("file", "problems.csv", "CSV file with the quiz questions")
	timeLimit := flag.Int("time", 30, "Time limit to complete the Quiz")
	flag.Parse()
	return *filename, *timeLimit
}

func timeQuiz(timeLimit int, timer chan int) {
	time.Sleep(time.Duration(timeLimit) * time.Second)
	timer <- 1
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

func readInput(done chan string) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		done <- scanner.Text()
	}
}

func makeQuestion(i int, problem Problem, timer chan int, input chan string) (int, error) {

	fmt.Printf("Problem #%v: %v = ", i, problem.question)
	for {
		select {
		case <-timer:
			return 0, errors.New("timed out")
		case resp := <-input:
			if resp == problem.answer {
				return 1, nil
			}
			return 0, nil
		}
	}

}

func startQuiz(problems []Problem, timer chan int) {
	score := 0
	input := make(chan string)
	go readInput(input)

	for i, problem := range problems {
		result, err := makeQuestion(i, problem, timer, input)
		if err != nil {
			fmt.Printf("\nYou scored %v out of %v.", score, len(problems))
			os.Exit(0)
		}
		score += result
	}
}

func main() {
	filename, timeLimit := readArguments()
	csvfile := openCSV(filename)
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	problems := readCSVProblems(r)
	timer := make(chan int)
	go timeQuiz(timeLimit, timer)
	startQuiz(problems, timer)
}
