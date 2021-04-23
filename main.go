package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	port := ":"
	args := os.Args[1:]
	if len(args) == 1 {
		port = port + os.Args[1]
	} else {
		port = ":8989"
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on the port %s\n", port)

	ch1 := make(chan string)
	// printing from chanel
	go hub(ch1)
	// accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn, ch1)
	}
}

func hub(ch <-chan string) {
	for {
		fmt.Printf("chanel prints %v\n", <-ch)
	}
}

func printLogo(conn net.Conn) {
	file, err := os.Open("logo.txt")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	bd := bufio.NewReader(file)
	for {
		line, err := bd.ReadString('\n')
		if err == io.EOF {
			conn.Write([]byte(line))
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte(line))
	}
}

func handleConnection(conn net.Conn, ch1 chan<- string) {
	defer conn.Close()

	printLogo(conn)

	//entering name
	bname, _, err := bufio.NewReader(conn).ReadLine()
	if err != nil {
		fmt.Println(err)
		return
	}
	name := string(bname)
	for {
		if name == "" {
			conn.Write([]byte("[ENTER YOUR NAME]:"))
			bname, _, err = bufio.NewReader(conn).ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			name = string(bname)
		} else {
			break
		}
	}
	fmt.Println(name)

	// typing
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "Exit" {
			conn.Write([]byte("Bye\n\r"))
			fmt.Println(name, "disconnected")
			break
		} else if text != "" {
			ch1 <- text
			fmt.Println(name, "enters", text)
			conn.Write([]byte("You enter " + text + "\n\r"))
		}
	}
	// close(ch1)
}
