package main

import "net"

func CreateServerUDP() {
	port := ":" + "3000"

	s, err := net.ResolveUDPAddr("udp4", port)
	if err != nil {
		panic(err)
	}
	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

}
