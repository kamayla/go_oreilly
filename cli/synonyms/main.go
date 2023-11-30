package main

import (
	"bufio"
	"cli/theaurus"
	"fmt"
	"os"
)

func main() {
	apiKey := "f8dc31de58dcf29533a5499f8ca18081"

	thesaurus := &theaurus.BigHuge{APIKey: apiKey}

	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)

		if err != nil {
			panic(err)
		}

		if len(syns) == 0 {
			panic("類義語がありませんでした")
		}

		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
