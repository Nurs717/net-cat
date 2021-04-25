package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Users struct {
	Conn net.Conn
	Name string
	Msg  []string
}

func PrintLogo(conn net.Conn) {
	file, err := os.Open("logo.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
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
