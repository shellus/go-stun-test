package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/ccding/go-stun/stun"
	"github.com/google/shlex"
	"net"
	"os"
	"time"
)

var remoteUdpAddr *net.UDPAddr

func main() {
	var test = flag.Bool("t", false, "Run the test suite")
	var serverAddr = flag.String("s", "stun.syncthing.net:3478", "STUN server address")
	var verboseLevel = flag.Int("v", 0, "Verbose level (0: none, 1: verbose, 2: double verbose, 3: triple verbose)")
	flag.Parse()

	// Validate verbose level
	if *verboseLevel < 0 || *verboseLevel > 3 {
		_, _ = fmt.Fprintln(os.Stderr, "Error: Invalid verbose level. Use -v with values 0, 1, 2, or 3.")
		os.Exit(1)
	}

	// Create Socket
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	printLocalAddr(conn.LocalAddr().(*net.UDPAddr).Port)
	// Discover the NAT
	if !*test {
		printExternalAddr(conn, *serverAddr, *verboseLevel)
	}

	go printMessage(conn)
	listenerInput(conn)

	_ = conn.Close()
}

func printLocalAddr(port int) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 检查是否是局域网地址
			if ip.IsLoopback() || !ip.IsGlobalUnicast() {
				continue
			}
			if ip.To4() != nil {
				fmt.Printf("%s: addr %s:%d\n", i.Name, ip.String(), port)
			} else {
				fmt.Printf("%s: addr [%s]:%d\n", i.Name, ip.String(), port)
			}
		}
	}
}
func printExternalAddr(conn *net.UDPConn, serverAddr string, verboseLevel int) {
	// Create a STUN client
	client := stun.NewClientWithConnection(conn)
	client.SetServerAddr(serverAddr)
	client.SetVerbose(verboseLevel >= 1)
	client.SetVVerbose(verboseLevel >= 2)

	nat, host, err := client.Discover()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	if host != nil {
		fmt.Printf("%s: %s\n", nat, host.TransportAddr())
	}
}
func printMessage(conn *net.UDPConn) {
	for {
		err := conn.SetReadDeadline(time.Now().Add(time.Duration(100) * time.Millisecond))
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		// 如果是超时，继续等待
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			continue
		}

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		// 将等下的接收人设置为刚给我们发送消息的人
		remoteUdpAddr = addr
		fmt.Printf("Received from %s: %s\n", addr.String(), string(buf[:n]))
	}
}
func listenerInput(conn *net.UDPConn) {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1] // 去除换行符

		// input以linux shell的方式解析成args，
		args, err := shlex.Split(input)
		if err != nil || len(args) == 0 {
			fmt.Println("Error parsing input:", err)
			continue
		}

		if args[0] == "addr" {
			if len(args) != 2 {
				_, _ = fmt.Fprintln(os.Stderr, "Error: invalid arguments")
				continue
			}
			remoteUdpAddr, err = net.ResolveUDPAddr("udp", args[1])
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			_, _ = fmt.Printf("Remote address set to %s\n", remoteUdpAddr.String())
		} else {
			if len(args) != 1 {
				_, _ = fmt.Fprintln(os.Stderr, "Error: invalid arguments")
				continue
			}
			// 否则就发送text到remoteUdpAddr
			if remoteUdpAddr == nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error: remote address is not set")
				continue
			}
			_, err = conn.WriteToUDP([]byte(args[0]), remoteUdpAddr)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			_, _ = fmt.Printf("Sent to %s: %s\n", remoteUdpAddr.String(), input)
		}
	}
}
