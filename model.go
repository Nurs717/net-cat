package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

type Users struct {
	Conn net.Conn
	Name string
	Msg  []string
}

func PrintLogo(conn net.Conn) {
	file, err := ioutil.ReadFile("logo.txt")
	if err != nil {
		fmt.Println(err)
	}
	conn.Write(file)
}

func EnterName(conn net.Conn) string {
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
	fmt.Println(name)
	return name
}
