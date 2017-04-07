package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
	//	"tools"
)

type uplink string

var uplinked uplink

func (u uplink) init() {

	switch conf.protocol {

	case "tcp":

		if conf.server == true {
			log.Println("start tcptun server ")
			go u.startTCPServer()

		} else {
			log.Println("start tcptun client ,connect to  ", conf.leftaddr)
			go u.conntoserver()
		}

	case "udp":
		go u.startUDPServer()
		go u.startUDPClient()

	}

}

func (u uplink) conntoserver() {
	for {
		conn, err := net.Dial("tcp", conf.leftaddr)
		if err != nil {
			log.Println("connect to server  err: ", err)
			time.Sleep(time.Second * 30)
			continue
		}

		go func() {
			for {

				//发送数据
				//	_, err := conn.Write([]byte{0x00, 0x00, 0x00, 0x00})

				_, err := conn.Write(<-packetchanR)
				if err != nil {
					log.Println("conn err,close", conn.RemoteAddr())
					conn.Close()
					break
				}

			}
		}()

		log.Println("connected tcptun server :", conn.RemoteAddr().String())

		buf := make([]byte, 8192)
		//	tmpbuf := make([]byte, 0, 16384)
		for {

			lenght, err := conn.Read(buf)

			if err != nil {
				log.Println("server conn err ,closed to server connect !", conn.RemoteAddr().String(), err)
				conn.Close()

				break
			}

			packetchanS <- buf[:lenght]

			/*
				if lenght > 0 {
					tmpbuf = append(tmpbuf, buf[:lenght]...)
				} else {
					//	time.Sleep(time.Microsecond * 500)
					continue
				}

				for len(tmpbuf) >= 4 {
				}
				packetlenght := tools.BytestoInt(tmpbuf[:4])

				if packetlenght == 0 {
					//	fmt.Println("send heart 0000")
					tmpbuf = tmpbuf[4:]

				} else {

					if len(tmpbuf) >= packetlenght+4 {

						packetchan <- tmpbuf[:lenght+4] //送去转发

						tmpbuf = tmpbuf[packetlenght+4:]
					} else {
						break
					}
				}

			*/
		}

	}

	time.Sleep(time.Second * 15)

	log.Println("reconnecting tcptun server ......")

}

func (u uplink) startUDPClient() {

	addr, err := net.ResolveUDPAddr("udp", conf.rightaddr)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	for {
		_, err = conn.Write([]byte(<-packetchanR))
		if err != nil {
			fmt.Println("failed:", err)
			continue
		}
	}
}

func (u uplink) startUDPServer() {
	// 创建监听

	addr, err := net.ResolveUDPAddr("udp", conf.leftaddr)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}

	log.Println("UDP Listening:", conf.port)
	socket, err := net.ListenUDP("udp4", addr)
	// socket, err := net.ListenUDP("udp4", &net.UDPAddr{
	// 	IP:   net.IPv4(0, 0, 0, 0),
	// 	Port: StrtoInt(conf.port),
	// })
	if err != nil {
		fmt.Println("UDP Listening  error !", err)
		return
	}
	defer socket.Close()

	for {
		// 读取数据
		data := make([]byte, 4096)
		lenght, _, err := socket.ReadFromUDP(data)

		if err != nil {
			fmt.Println("UDP read data error!", err)
			continue
		}

		// if ip.IP.String() != conf.clientaddr {
		// 	continue
		// }

		packetchanS <- data[:lenght]

	}

}

//StartServer StartServer
func (u uplink) startTCPServer() {

	service := ":" + conf.port //strconv.Itoa(port);

	l, err := net.Listen("tcp", service)
	if err != nil {
		log.Println("tcptun server Listening error", err)
		os.Exit(1)
	}

	for {
		log.Println("tcptun server  Listening...:", conf.port)
		conn, err := l.Accept()
		if err != nil {
			log.Println("tcptun server  Listening error :", err)
			os.Exit(1)
		}

		go u.handler(conn)

	}

}

func (u *uplink) handler(conn net.Conn) {

	log.Println("tcptun client is connected from :", conn.RemoteAddr().String())
	buf := make([]byte, 8192)
	//	tmpbuf := make([]byte, 0, 2048)

	go func() {
		for {
			_, err := conn.Write(<-packetchanR)
			if err != nil {
				log.Println("tcptun server  send data err ,conn close  ")
				conn.Close()
				break
			}
		}
	}()

	for {

		lenght, err := conn.Read(buf)

		if err != nil {

			log.Println("slave read err ,closed conn :", conn.RemoteAddr().String(), err)

			conn.Close()
			break
		}

		//	log.Println("recive packet", lenght, " byte: ", buf[:lenght])

		packetchanS <- buf[:lenght]
		/*

			if lenght > 0 {
				tmpbuf = append(tmpbuf, buf[:lenght]...)
			} else {
				//	time.Sleep(time.Microsecond * 500)
				continue
			}

			for len(tmpbuf) >= 4 {

				if tools.BytestoInt(tmpbuf[:4]) == 0 {

					tmpbuf = tmpbuf[4:]
					packetchan <- tmpbuf
					//	conn.Write([]byte{0x00, 0x00, 0x00, 0x00})
				} else {
					tmpbuf = tmpbuf[4:]
					log.Println("tcptun client  data format err", tmpbuf)
					break
				}

			}

		*/

	}
}
