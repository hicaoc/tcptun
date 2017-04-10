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
var udpaddr *net.UDPAddr

func (u uplink) init() {

	if conf.server == true {
		log.Println("start tcptun server ")
		if conf.protocol == "tcp" {
			go u.startTCPServer()
		} else {
			go u.startUDPServer()
		}

	} else {
		log.Println("start tcptun client ,connect to  ", conf.serveraddr)
		if conf.protocol == "tcp" {
			go u.conntoserver()
		} else {
			go u.startUDPClient()
		}
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

		log.Println("connected tcptun TCP server :", conn.RemoteAddr().String())

		buf := make([]byte, 4096)
		//	tmpbuf := make([]byte, 0, 16384)
		for {

			lenght, err := conn.Read(buf)

			if err != nil {
				log.Println("server conn err ,closed to server connect !", conn.RemoteAddr().String(), err)
				conn.Close()

				break
			}

			packetchanS <- buf[:lenght]

		}

	}

	time.Sleep(time.Second * 15)

	log.Println("reconnecting tcptun server ......")

}

func (u uplink) startUDPClient() {

	addr, err := net.ResolveUDPAddr("udp", conf.serveraddr)
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

	log.Println("connected tcptun UDP server :", conn.RemoteAddr().String())

	go func() {

		data := make([]byte, 1500)
		for {
			lenght, _ := conn.Read(data)
			// if err != nil {
			// 	fmt.Println("failed to read UDP msg because of ", err)
			// 	continue
			// }

			packetchanS <- data[:lenght]
		}

	}()

	for {
		_, _ = conn.Write([]byte(<-packetchanR))
		// if err != nil {
		// 	fmt.Println("failed:", err)
		// 	continue
		// }
	}
}

func (u uplink) startUDPServer() {
	// 创建监听

	addr, err := net.ResolveUDPAddr("udp", conf.serveraddr)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}

	log.Println("UDP Listening:", conf.serveraddr)
	socket, err := net.ListenUDP("udp4", addr)
	// socket, err := net.ListenUDP("udp4", &net.UDPAddr{
	// 	IP:   net.IPv4(0, 0, 0, 0),
	// 	Port: StrtoInt(conf.port),
	// })
	if err != nil {
		fmt.Println("UDP Listening  error !", err)
		os.Exit(1)
	}
	defer socket.Close()

	var lenght int
	//var udperr error

	go func() {
		for {

			_, _ = socket.WriteToUDP([]byte(<-packetchanR), udpaddr)
			// if err != nil {
			// 	fmt.Println("send udp packet to client failed:", err)
			// 	continue
			// }
		}
	}()

	for {
		// 读取数据
		data := make([]byte, 8192)
		lenght, udpaddr, _ = socket.ReadFromUDP(data)

		// if udperr != nil {
		// 	fmt.Println("UDP read data error!", udperr)
		// 	continue
		// }

		// if ip.IP.String() != conf.clientaddr {
		// 	continue
		// }

		packetchanS <- data[:lenght]

	}

}

//StartServer StartServer
func (u uplink) startTCPServer() {

	//	service := ":" + conf.port //strconv.Itoa(port);

	l, err := net.Listen("tcp", conf.serveraddr)
	if err != nil {
		log.Println("tcptun server Listening error", err)
		os.Exit(1)
	}

	for {
		log.Println("tcptun server TCP  Listening...:", conf.serveraddr)
		conn, err := l.Accept()
		if err != nil {
			log.Println("tcptun server TCP Listening error :", err)
			os.Exit(1)
		}

		go u.handler(conn)

	}

}

func (u *uplink) handler(conn net.Conn) {

	log.Println("tcptun client is connected from whth TCP:", conn.RemoteAddr().String())
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
