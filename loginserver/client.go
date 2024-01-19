package main

import (
	"net"
	"fmt"
)

func main() {
	conn, err := net.Dial("tcp", ":10000")
	if err != nil {
		fmt.Println("dial err")
	}
	defer conn.Close()

	data := "hello"
	conn.Write([]byte(data))
	fmt.Println("write ok")

}
