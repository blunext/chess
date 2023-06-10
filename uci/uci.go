package uci

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Start() {
	f, err := os.Create("~/Documents/projects/chess/data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		command := strings.TrimSpace(line)
		fmt.Printf("echo: %s\n", command)
		_, err = f.WriteString(fmt.Sprintf("c: %s\n", command))
		if err != nil {
			log.Fatal(err)
		}
		switch command {
		case "uci":
			respond(f, "id name MyChessEngine")
			respond(f, "id author MyName")
			respond(f, "uciok")
		}
		_, err := f.WriteString("\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", scanner.Err())
	}
}

func respond(f *os.File, response string) {
	fmt.Printf("%s\n", response)
	logUCI(f, response)
}

func logUCI(f *os.File, response string) {
	_, err := f.WriteString(fmt.Sprintf("r: %s\n", response))
	if err != nil {
		log.Fatal(err)
	}
}
