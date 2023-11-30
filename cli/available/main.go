package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func exists(domain string) (bool, error) {
	const whoisServer = "com.whois-servers.net"
	conn, err := net.Dial("tcp", whoisServer+":43")

	if err != nil {
		return false, err
	}

	defer conn.Close()

	conn.Write([]byte(domain + "\r\n"))
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), "no match") {
			return false, nil
		}
	}
	return true, nil
}

func main() {
	var marks = map[bool]string{true: "✅", false: "❌"}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		domain := s.Text()
		isExist, err := exists(domain)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s\n", marks[!isExist], domain)
		time.Sleep(1 * time.Second)
	}
}
