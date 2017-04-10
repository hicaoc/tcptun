package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type config struct {
	server   bool
	protocol string
	//	port          string
	serveraddr string
	//	clientaddr    string
	interfacename string
}

var conf = &config{}

func (c *config) init() {

	conf.readconffile()
	//	go c.cronread()

}

func (c *config) readconffile() {

	log.Println("read config file tcptun.ini ......")

	f, err := os.Open("./tcptun.ini")
	if err != nil {
		log.Println("open tcptun.ini file err:", err)
		fmt.Println(`
server=true 
protocol=udp
serveraddr=10.140.0.2:9999
interfacename=tun10

	
`)
		os.Exit(1)
	}
	defer f.Close()

	rd := bufio.NewReader(f)

	for {

		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err || line == ".\n" {
			//	log.Println("read basinfo file error :", err)
			break
		}

		ss := strings.Split(strings.TrimSuffix(line, "\n"), " ")

		s := strings.Split(strings.TrimSpace(ss[0]), "=")

		switch s[0] {

		case "server":
			if s[1] == "true" {
				c.server = true
			}

		case "protocol":
			c.protocol = s[1]

		case "interfacename":
			c.interfacename = s[1]

		case "serveraddr":
			c.serveraddr = s[1]

		}

	}
	log.Println("Read tcptun conf file ok ", c)
}
