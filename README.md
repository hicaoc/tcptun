# tcptun
a tunnel with  tcp/udp

# 跨区服务器内网通信

- gihtup网址:https://github.com/hicaoc/tcptun

- 内网网段计算网址:https://tool.520101.com/wangluo/ipjisuan/

- 计算图片方法: https://imgchr.com/i/NWsY9O

  ![NWsY9O.png](https://s1.ax1x.com/2020/06/29/NWsY9O.png)

## 1.master

```
1.  编译
  ####go install 
  go build -o tcptun

2. 授权    #master
  chmod 755 tcptun
    ./tcptun    
3.编辑  
  vi tcptun.ini
    
    server=true   # server端
    protocol=udp    #通信方式
    serveraddr=172.31.0.54:9999   #master端内网ip
    interfacename=tun100   # 网卡名字注意先查看是否有相同的名字  ifconfig

4. #新增网卡   
  和上面编辑的网卡名字一样tun100    网段 10.0.0.1   掩码  255.255.255.240 
ifconfig tun100 10.0.0.1 netmask 255.255.255.240  
ifconfig tun100 up/down   #启用  停止
    # 同时操作后先ping slave的10.0.0.2是否ping通
5. 添加路由
ip route add 172.31.32.0/20 via 10.0.0.2         #172.31.32.0/20为slave端的eth0网段以及slave的tun100的ip 
## 注意,不知道网段的可以用上述网址计算

6. 启动
./tcptun   
# 匿名启动
nohup ./tcptun &
7. ping 一下对方内网ip看是否通


#由于重启进程会消失,所以写成脚本方便启动
8.最后修改启动脚本starttcptun.sh
vi starttcptun.sh

#!/bin/sh 

killall tcptun
sleep 2
nohup ./tcptun &
sleep 2
ifconfig tun100 10.0.0.1 netmask 255.255.255.240 up
ip route add 172.31.32.0/20 via 10.0.0.2

# 授权
chmod 755 starttcptun.sh
./starttcptun.sh

```

## 2. slave

```
####  slave

1.  编译
 go build -o tcptun
2. 授权    #slave 
3. 编辑
  vi tcptun.ini
    
    #server=true   # slave端
    protocol=udp    #通信方式
    serveraddr=47.105.197.53:9999   #slave端内网ip
    interfacename=tun100   # 网卡名字注意先查看是否有相同的名字  ifconfig


4. #新增网卡   
  和上面编辑的网卡名字一样tun100    网段 10.0.0.2   掩码  255.255.255.240 
ifconfig tun100 10.0.0.2 netmask 255.255.255.240
ifconfig tun100 up/down   #启用  停止
# 同时操作后先ping slave的10.0.0.2是否ping通

5. 添加路由
ip route add   172.31.0.0/20 via 10.0.0.1   # 172.31.0.0/20为master的eth0的网段和master的tun100的ip


6. 启动
./tcptun   
# 匿名启动
nohup ./tcptun &
7. ping 一下对方内网ip看是否通

#由于重启进程会消失,所以写成脚本方便启动
8.最后修改启动脚本starttcptun.sh
vi starttcptun.sh

#!/bin/sh 

killall tcptun
sleep 2
nohup ./tcptun &
sleep 2
ifconfig tun100 10.0.0.3 netmask 255.255.255.240 up 
#ifconfig tun10 10.5.0.1 netmask 255.255.255.0 up
#ip route add 192.168.0.0/16  dev tun10
ip route add   172.31.0.0/20 via 10.0.0.1


# 授权
chmod 755 starttcptun.sh
./starttcptun.sh

```


