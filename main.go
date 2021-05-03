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
	info string
}

var (
	data        []string
	connections int
	allconn     []net.Conn
	allname     []string
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

	ch1 := make(chan message)

	// printing from chanel
	go hub(ch1)

	// accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		connections = connections + 1
		if connections < 10 {
			allconn = append(allconn, conn)
			go handleConnection(conn, ch1)
		} else {
			conn.Write([]byte("Server is busy. Please try later.\n"))
			conn.Close()
		}
	}
}

func hub(ch <-chan message) {
	for {
		msg := <-ch
		for i, conn := range allconn {
			if conn == msg.from {
				continue
			}
			if msg.info == "" {
				conn.Write([]byte(msg.body))
				bname := "[" + allname[i] + "]" + ":"
				conn.Write([]byte(bname))
			} else {
				conn.Write([]byte(msg.info))
				aname := msg.body + "[" + allname[i] + "]" + ":"
				conn.Write([]byte(aname))
			}
		}
		msg.info = ""
	}
}

func handleConnection(conn net.Conn, ch1 chan<- message) {
	defer conn.Close()

	PrintLogo(conn)

	//entering name
	name := EnterName(conn)
	allname = append(allname, name)

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
			connections = connections - 1
			fmt.Println(name + " disconnected")
			connMessage := message{info: "\n" + name + " has left our chat...\n", body: "[" + time + "]", from: conn}
			ch1 <- connMessage
			break
		}

		var connMessage message
		if string(msg) != "" {
			connMessage = message{body: "\n" + terminal + string(msg) + "\n" + "[" + time + "]", from: conn}
		}
		ch1 <- connMessage

		data = append(data, terminal+string(msg))
	}
}
