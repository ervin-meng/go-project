package utils

import "net"

func GetPort() int {
	addr, _ := net.ResolveTCPAddr("tcp", "localhost:0")
	lis, _ := net.ListenTCP("tcp", addr)
	defer lis.Close()
	return lis.Addr().(*net.TCPAddr).Port
}
