package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

var data []string
var connections int

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

func hub(ch <-chan string) {
	for {
		fmt.Printf("chanel prints %v\n", <-ch)
	}
}

func handleConnection(conn net.Conn, ch1 chan<- string) {
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

	// otherconn := user.Conn
	// fmt.Println(otherconn)
	// if otherconn != nil {
	// 	connmsg := name + "connected to the server..."
	// 	otherconn.Write([]byte(connmsg))
	// }
	// user.Conn = conn
	for {
		time := time.Now().Format("2006-01-02 15:04:05")
		terminal := "[" + time + "]" + "[" + name + "]" + ":"
		conn.Write([]byte(terminal))
		msg, _, err := bufio.NewReader(conn).ReadLine()
		if err != nil {
			connections = connections - 1
			fmt.Println("disconnected")
			break
		}
		outgoing := terminal + string(msg)
		// var msg2 []string
		// msg2 = append(msg2, outgoing)
		// e := &Users{
		// 	Conn: conn,
		// 	Name: name,
		// 	Msg:  msg2,
		// }
		// fmt.Println(e)
		data = append(data, outgoing)
		ch1 <- string(outgoing)
	}
}
