/*
dump packets sent to TUN device to stdout

example:
-- Packet Received --

IP Version: 4
Length: 84
Protocol: 1 (1=ICMP, 6=TCP, 17=UDP)
Source IP: 10.12.0.2
Destination IP: 10.12.0.1
Payload: [08001467cc79000156bcc24d0004130d08090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f3031323334353637]
*/
package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	//	"os"
	//	"log"
	"time"

	"github.com/jamescun/tuntap"
)

var packetchanR = make(chan []byte, 10)
var packetchanS = make(chan []byte, 10)

func main() {

	conf.init()
	uplinked.init()
	tcptun()

}

func tcptun() {
	// if len(os.Args) < 2 {
	// 	fmt.Println("help: icmp.go <interface name|device path>")
	// 	return
	// }

	tun, err := tuntap.Tun(conf.interfacename)
	if err != nil {
		fmt.Println("error: tun:", err)
		return
	}
	defer tun.Close()

	buf := make([]byte, 8192)

	go func() {
		for {
			pp := <-packetchanS
			_, err := tun.Write(pp)
			//		log.Println("packet:", pp)
			if err != nil {

			}
		}

	}()

	for {
		n, err := tun.Read(buf)
		if err == tuntap.ErrNotReady {
			fmt.Println("warning: tun: interface not ready, waiting 1s...")
			time.Sleep(1 * time.Second)
			continue
		} else if err != nil {
			fmt.Println("error: read:", err)
			break
		}

		//	fmt.Print("\n-- Packet Received --\n\n")

		packetchanR <- buf[:n]

		continue

		switch buf[0] & 0xF0 {
		case 0x40:
			fmt.Println("IP Version: 4")

			fmt.Printf("Length: %d\n", binary.BigEndian.Uint16(buf[2:4]))
			fmt.Printf("Protocol: %d (1=ICMP, 6=TCP, 17=UDP)\n", buf[9])
			fmt.Printf("Source IP: %s\n", net.IP(buf[12:16]))
			fmt.Printf("Destination IP: %s\n", net.IP(buf[16:20]))

			ihl := (buf[0] & 0x0F) * 4
			fmt.Printf("Payload: [%x]\n", buf[ihl:n])

		case 0x60:
			fmt.Println("IP Version: 6")

			fmt.Printf("Length: %d\n", binary.BigEndian.Uint16(buf[4:6]))
			fmt.Printf("Protocol: %d (1=ICMP, 6=TCP, 17=UDP)\n", buf[7])
			fmt.Printf("Source IP: %s\n", net.IP(buf[8:24]))
			fmt.Printf("Destination IP: %s\n", net.IP(buf[24:40]))

			fmt.Printf("Payload: [%x]\n", buf[40:n])
		}

	}
}

//StrtoInt strings to int
func StrtoInt(a string) int {

	b, error := strconv.Atoi(a)
	if error != nil {
		log.Println("string to int fail")

	}
	return b
}
