package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

type items []string

var tlds items

// itemsをVarインターフェースに適合させる
func (t *items) String() string {
	return strings.Join(*t, ",")
}

// itemsをVarインターフェースに適合させる
// Setはflagパッケージによって、コマンドライン引数がパースされた後に呼び出される
func (t *items) Set(value string) error {
	values := strings.Split(value, ",")
	for _, v := range values {
		*t = append(*t, v)
	}
	return nil
}

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789-_"

func init() {
	flag.Var(&tlds, "tld", "Top level domains")
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())

	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		text := strings.ToLower(s.Text())
		var newText []rune
		for _, r := range text {
			if unicode.IsSpace(r) {
				r = '-'
			}
			if !strings.ContainsRune(allowedChars, r) {
				continue
			}
			newText = append(newText, r)
		}
		fmt.Println(string(newText) + "." + tlds[rand.Intn(len(tlds))])
	}
}
