package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

//PrintLogo Prints logo gopher
func printLogo(conn net.Conn) {
	file, err := ioutil.ReadFile("logo.txt")
	if err != nil {
		fmt.Println(err)
	}
	conn.Write(file)
}

//EnterName asking to enter name
func enterName(conn net.Conn) string {
	bname, _, err := bufio.NewReader(conn).ReadLine()
	if err != nil {
		log.Fatal(err)
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
	fmt.Println(name + " connected")
	return name
}

// sending msg for other connections
func hub(ch <-chan message) {
	for {
		msg := <-ch
		for name, conn := range allconn {
			if conn == msg.from {
				continue
			}
			if msg.info == "" {
				conn.Write([]byte(msg.body))
				bname := "[" + name + "]" + ":"
				conn.Write([]byte(bname))
			} else {
				conn.Write([]byte(msg.info))
				aname := msg.body + "[" + name + "]" + ":"
				conn.Write([]byte(aname))
			}
		}
		msg.info = ""
	}
}
