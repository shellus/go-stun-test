# 使用公共STUN服务器进行NAT打洞的实验

## 步骤
1. 在两台NAT网络环境下的终端，分别运行本程序。我们将其称为A和B。
2. 当两处终端输出了NAT类型和公网地址，你需要借助另一个手段来交换他们的地址（是的，我用眼睛看，然后记在本子上来交换他们的公网地址）
3. 在A终端输入B终端的公网地址，在B终端输入A终端的公网地址。例如 `addr 100.100.100.100:50000`
4. 然后，你可以输入任意内容，输入的内容会通过UDP被发送到上一步输入的地址。
5. 观察另一端是否输出了你刚刚输入的内容，双方都需要至少发送一次数据包，才能打通NAT。

## NAT类型和打洞匹配

1. 完全圆锥型NAT（Full Cone NAT）
   > 内网主机同一个内网IP和端口映射出来的公网IP和端口，任何公网主机向映射的公网IP和端口发送数据包，数据包就可以到达内网主机。
2. 地址限制圆锥型NAT（Restricted Cone NAT）
   > 内网主机同一个内网IP和端口映射出来的公网IP和端口，只要内网主机发过数据包给公网主机，这个公网主机就可以向映射的公网IP和端口发送数据包，数据包就可以到达内网主机。
3. 端口限制圆锥型NAT（Port Restricted Cone NAT）
4. > 内网主机同一个内网IP和端口映射出来的公网IP和端口，只要内网主机发过数据包到给公网主机的特定端口，这个公网主机就可以通过特定端口向映射的公网IP和端口发送数据包，数据包就可以到达内网主机。
4. 对称型NAT（Symmetric NAT）
   > 内网主机给不同的公网主机和端口发送数据包，映射出来的公网IP和端口不一样，同时只有内网主机发过数据包到给公网主机的特定端口，这个公网主机才可以通过特定端口向特定端口所映射的公网IP和端口发送数据包，数据包才可以到达内网主机。


| NAT类型1    | NAT类型2    | 结果            |
|-------------|-------------|-----------------|
| 全锥型      | 全锥型    | ✔               |
| 全锥型      | 受限锥型    | ✔               |
| 全锥型      | 端口受限锥型| ✔               |
| 全锥型      | 对称型      | ✔               |
| 受限锥型    | 受限锥型    | ✔               |
| 受限锥型    | 端口受限锥型| ✔               |
| 受限锥型    | 对称型      | ✔               |
| 端口受限锥型| 端口受限锥型| ✔               |
| 端口受限锥型| 对称型      | ❌，无法打通     |
| 对称型      | 对称型      | ❌，无法打通    |
