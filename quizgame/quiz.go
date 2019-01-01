package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type QuizQuestion struct {
	question string
	answer   int
}

func readInput(input chan<- int) {
	for {
		var u int
		fmt.Scanf("%d\n", &u)
		input <- u
	}
}

func main() {
	var filename string
	if len(os.Args) == 1 {
		filename = "problems"
	} else {
		filename = os.Args[1]
	}
	csvFile, _ := os.Open(filename + ".csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	score := 0
	fmt.Println("Ready to start?")
	fmt.Scanln()
	userInput := make(chan int)
	go readInput(userInput)
problemLoop:
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		answer, _ := strconv.Atoi(line[1])
		question := line[0]
		fmt.Printf("%#v\n", question)
		select {
		case userAnswer := <-userInput:
			if userAnswer == answer {
				fmt.Println("Correct answer:", userAnswer)
				score++
			} else {
				fmt.Println("Wrong answer")
			}
		case <-time.After(5 * time.Second):
			fmt.Println("\n Time is over!")
			break problemLoop
		}
	}
	fmt.Printf("%#v\n", score)
}
