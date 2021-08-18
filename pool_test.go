package main

import (
	"net"
	"testing"
)

func init() {
	listener, err := net.Listen("tcp", "0.0.0.0:30000")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			_, err := listener.Accept()
			if err != nil {
				panic(err)
			}
		}
	}()
}

func TestTCPConnection(t *testing.T) {
	pool := NewPool(func() (interface{}, error) {
		conn, err := net.Dial("tcp", "0.0.0.0:30000")
		if err != nil {
			return nil, err
		}
		return conn, nil
	})

	_, err := pool.Factory()
	if err != nil {
		t.Error(err)
	}
}
