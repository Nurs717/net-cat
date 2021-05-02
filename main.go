package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

type message struct {
	from net.Conn
	body string
}

var data []string
var connections int
var allconn []net.Conn

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

	ch1 := make(chan message)
	ch2 := make(chan net.Conn, 2)

	go othconn(ch2)

	// printing from chanel
	go hub(ch1)

	// accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		connections = connections + 1
		fmt.Println(connections)
		ch2 <- conn
		if connections < 4 {
			allconn = append(allconn, conn)
			go handleConnection(conn, ch1)
		} else {
			conn.Write([]byte("Server is busy. Please try later.\n"))
			conn.Close()
		}
	}
}

func othconn(ch <-chan net.Conn) {
	for {
		fmt.Printf("chanel 2 prints %v\n", <-ch)
	}
}

func hub(ch <-chan message) {
	for {
		msg := <-ch
		for _, conn := range allconn {
			if conn == msg.from {
				continue
			}

			conn.Write([]byte(msg.body))
		}
	}
}

func handleConnection(conn net.Conn, ch1 chan<- message) {
	defer conn.Close()

	PrintLogo(conn)

	//entering name
	name := EnterName(conn)
	// user := &Users{}

	// exist data for new users
	for _, message := range data {
		conn.Write([]byte(message))
		conn.Write([]byte("\n"))
	}

	connMessage := message{body: "\n" + name + " has joined our chat...", from: conn}
	ch1 <- connMessage

	for {
		time := time.Now().Format("2006-01-02 15:04:05")
		terminal := "[" + time + "]" + "[" + name + "]" + ":"
		conn.Write([]byte(terminal))
		msg, _, err := bufio.NewReader(conn).ReadLine()
		if err != nil {
			connections = connections - 1
			fmt.Println("disconnected")
			connMessage := message{body: "\n" + name + " has left our chat...", from: conn}
			ch1 <- connMessage
			break
		}
		connMessage := message{body: "\n" + terminal + string(msg), from: conn}
		ch1 <- connMessage

		data = append(data, terminal+string(msg))
	}
}
