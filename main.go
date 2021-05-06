package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type message struct {
	from net.Conn
	body string
	info string
}

var (
	data        []string
	connections int
	allconn     map[string]net.Conn
)

func init() {
	allconn = make(map[string]net.Conn)
}

func main() {
	port := ":"
	args := os.Args[1:]

	if len(args) == 1 {
		port = port + args[0]
	} else if len(args) == 0 {
		port = ":8989"
	} else {
		log.Fatal("[USAGE]: ./TCPChat $port")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on the port %s\n", port)

	ch1 := make(chan message)

	// sending msg for other connections
	go hub(ch1)

	// accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		connections++
		if connections <= 10 {
			go handleConnection(conn, ch1)
		} else {
			conn.Write([]byte("Server is busy. Please try later.\n"))
			conn.Close()
		}
	}
}

func handleConnection(conn net.Conn, ch1 chan<- message) {
	defer conn.Close()

	printLogo(conn)

	//entering name
	name := enterName(conn)
	allconn[name] = conn

	// exist data for new users
	for _, message := range data {
		conn.Write([]byte(message))
		conn.Write([]byte("\n"))
	}

	onetime := time.Now().Format("2006-01-02 15:04:05")
	connMessage := message{info: "\n" + name + " has joined our chat...\n", body: "[" + onetime + "]", from: conn}
	ch1 <- connMessage

	for {
		time := time.Now().Format("2006-01-02 15:04:05")
		terminal := "[" + time + "]" + "[" + name + "]" + ":"
		conn.Write([]byte(terminal))
		msg, _, err := bufio.NewReader(conn).ReadLine()
		if err != nil {
			connections--
			fmt.Println(name + " disconnected")
			connMessage := message{info: "\n" + name + " has left our chat...\n", body: "[" + time + "]", from: conn}
			ch1 <- connMessage
			break
		}

		var connMessage message
		if string(msg) != "" {
			connMessage = message{body: "\n" + terminal + string(msg) + "\n" + "[" + time + "]", from: conn}
			ch1 <- connMessage
			data = append(data, terminal+string(msg))
		}
	}
}
