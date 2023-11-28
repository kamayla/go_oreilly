package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	// Vowelは母音を表す
	duplicateVowel = true
	removeVowel    = false
)

func randBool() bool {
	return rand.Intn(2) == 0
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := []byte(s.Text())
		if randBool() {
			var vI int = -1
			for i, char := range word {
				switch char {
				case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
					if randBool() {
						vI = i
					}
				}
			}
			if vI >= 0 {
				switch randBool() {
				case duplicateVowel:
					/*
						スライスの[x:y]指定は半開区間なのでこの場合は
						最初から母音の箇所を含む文字列と、母音の箇所以降の文字列を結合している。
						つまりchatという文字列の場合はchaatという文字列になる。
					*/
					word = append(word[:vI+1], word[vI:]...)
				case removeVowel:
					/*
						chatの場合はchとtが結合されてchtとなる。
					*/
					word = append(word[:vI], word[vI+1:]...)
				}
			}
		}
		fmt.Println(string(word))
	}
}
