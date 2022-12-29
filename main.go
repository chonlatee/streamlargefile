package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
)

type node struct {
	conn net.Conn
}

type server struct {
	nodes map[string]*node
}

func NewServer() *server {

	conn1, err := net.Dial("tcp", ":3333")
	if err != nil {
		log.Fatal("can't connect to :3333", err.Error())
	}

	n1 := &node{
		conn: conn1,
	}

	conn2, err := net.Dial("tcp", ":3334")
	if err != nil {
		log.Fatal("can't connect :3334", err.Error())
	}

	n2 := &node{
		conn: conn2,
	}

	s := &server{
		nodes: make(map[string]*node),
	}

	s.nodes["3333"] = n1
	s.nodes["3334"] = n2

	return s

}

func main() {

	s := NewServer()

	s.broadcast()

	select {}

}

func (s *server) broadcast() {

	var wr []io.Writer

	f, err := os.Open("img.jpg")
	if err != nil {
		log.Fatal("can't open file ", err.Error())
	}
	fi, err := f.Stat()
	if err != nil {
		log.Fatal("can't stat file ", err.Error())
	}

	size := fi.Size()

	for _, v := range s.nodes {
		binary.Write(v.conn, binary.LittleEndian, int64(size))
		if err != nil {
			log.Fatal("can't send file meta data")
		}
		wr = append(wr, v.conn)
	}

	mw := io.MultiWriter(wr...)
	_, err = io.CopyN(mw, f, int64(size))
	if err != nil {
		log.Fatalln("can't copy file ", err.Error())
	}
}
