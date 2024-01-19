package main

import (
	"fmt"
	"net"
	"os"
	"gopkg.in/ini.v1"
)


type ConnList struct {
	ConnRoot *ConnNode
	ConnTail *ConnNode
	Mutex sync.Mutex
}

type ConnNode struct {
	Next *ConnNode
	ConnId int

}


func main() {
	fmt.Println("Start LoginServer...")

	// load servercfg
	config, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("Fail to read config: %v", err)
		os.Exit(1)
	}
	port := config.Section("connection").Key("port").String()
	maxConn, err := config.Section("connection").Key("maxConn").Int()
	if err != nil {
		fmt.Println("fail to get maxConn config: %v", err)
	}

	// check IsAlreadyRun("LoginServer")


	// startNetwork
	l, err:= net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("listen error")
		return
	}
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error")
			break
		}
		go handleConn(c, maxConn)
	}

	// register message

	// init web command manager

	// gift code manager init

	fmt.Println("Start LoginServer success")
}

func handleConn(c net.Conn, maxConn int) {
	defer c.Close()
	for {
		var buf = make([]byte, maxConn)
		n, err := c.Read(buf)
		if err != nil {
			fmt.Println("conn read error")
			return
		}
		fmt.Printf("conn got %d bytes, content is %s\n", n, string(buf[:n]))
	}
}
