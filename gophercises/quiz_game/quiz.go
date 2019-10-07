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
	"strings"
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

func setQuizTimer(timeLimit int, timer chan int) {
	time.Sleep(time.Duration(timeLimit) * time.Second)
	timer <- 1
}

func parseCSV(filename string, problem chan Problem) {
	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)

	go func(r *csv.Reader, problem chan Problem) {
		defer csvfile.Close()

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			problem <- Problem{question: record[0], answer: record[1]}
		}
		close(problem)
	}(r, problem)
}

func readInput(done chan string) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		done <- strings.TrimSpace(scanner.Text())
	}
}

func makeQuestion(problem Problem, timer chan int, input chan string) (int, error) {
	fmt.Printf("Problem: %v = ", problem.question)
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

func finalScore(score int) {
	fmt.Printf("\nYou scored %v.", score)
}

func startQuiz(problem chan Problem, timer chan int) {
	score := 0
	input := make(chan string)
	go readInput(input)

	for {
		p, ok := <-problem
		if ok == false {
			finalScore(score)
			break
		}

		result, err := makeQuestion(p, timer, input)
		if err != nil {
			finalScore(score)
			os.Exit(0)
		}
		score += result
	}
}

func main() {
	filename, timeLimit := readArguments()

	problem := make(chan Problem)
	go parseCSV(filename, problem)

	timer := make(chan int)
	go setQuizTimer(timeLimit, timer)
	startQuiz(problem, timer)
}
