package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	add, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		panic(err)
	}
	con, err := net.DialUDP("udp", nil, add)
	if err != nil {
		panic(err)
	}
	defer con.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">\n")
		userMsg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}
		_, err = con.Write([]byte(userMsg))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}
