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
	server        bool
	protocol      string
	port          string
	leftaddr      string
	rightaddr     string
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
protocol=tcp
port=9999
leftaddr=10.0.0.1:9999
rightaddr=22.33.44.55:9999
interfacename=tun100
	
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
		case "port":
			c.port = s[1]

		case "protocol":
			c.protocol = s[1]

		case "interfacename":
			c.interfacename = s[1]

		case "leftaddr":
			c.leftaddr = s[1]

		case "clientaddr":
			c.rightaddr = s[1]

		}

	}
	log.Println("Read tcptun conf file ok ", c)
}
