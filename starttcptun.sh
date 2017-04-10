#!/bin/sh 

killall tcptun
sleep 2
nohup ./tcptun &
sleep 2
ifconfig tun10 10.5.0.1 netmask 255.255.255.0 up
ip route add 192.168.0.0/16  dev tun10
