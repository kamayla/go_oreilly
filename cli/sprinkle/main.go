package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherWord = "*"

var transforms = []string{}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	file, err := os.Open("wordlist.txt")
	if err != nil {
		log.Fatalln(err)
	}

	s := bufio.NewScanner(file)

	for s.Scan() {
		word := s.Text()

		transform := word + otherWord

		if rand.Intn(2) == 1 {
			transform = otherWord + word
		}

		transforms = append(transforms, transform)
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		// strings.Replaceは、第4引数が-1の場合、全ての一致する文字列を置換する
		fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
	}
}
