package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	var addr, dirName string
	flag.StringVar(&addr, "addr", "", "-addr=:3000")
	flag.StringVar(&dirName, "dirName", "", "-dirName=3000")
	flag.Parse()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("can't listen on %v err %v:\n", addr, err.Error())
	}

	err = os.Mkdir(fmt.Sprintf("./client/"+dirName), os.ModePerm)
	if err != nil {
		log.Fatalln("can't create folder: ", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("accept err: ", err.Error())
		}

		go handleConn(conn, dirName)
	}

}

func handleConn(conn net.Conn, dir string) {
	var size int64
	err := binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		log.Fatalln("can't read file size")
	}

	path := fmt.Sprintf("./client/%s/img.jpg", dir)
	f, err := os.Create(path)
	if err != nil {
		log.Printf("can't create file %v err: %v", path, err)
		return
	}

	defer f.Close()

	n, err := io.CopyN(f, conn, int64(size))

	if err != nil {
		log.Printf("copy file err: %v", err)
		return
	}

	log.Printf("written %v bytes\n", n)
}
