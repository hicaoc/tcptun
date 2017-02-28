package main

import (
	"log"
	"net"
	"os"
	"time"
	//	"tools"
)

type uplink string

var uplinked uplink

func (u uplink) init() {

	if conf.server == true {
		log.Println("start tcptun server ")
		go u.StartServer()

	} else {
		log.Println("start tcptun client ,connect to  ", conf.serveraddr)
		go u.conntoserver()
	}

}

func (u uplink) conntoserver() {
	for {
		conn, err := net.Dial("tcp", conf.serveraddr)
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

//StartServer StartServer
func (u uplink) StartServer() {

	service := ":" + conf.tcpport //strconv.Itoa(port);

	l, err := net.Listen("tcp", service)
	if err != nil {
		log.Println("tcptun server Listening error", err)
		os.Exit(1)
	}

	for {
		log.Println("tcptun server  Listening...:", conf.tcpport)
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
