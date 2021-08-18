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
	pool := NewPool(&PoolConfig{
		Min: 1,
		Max: 5,
	}, func() (interface{}, error) {
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

func TestPoolPopulation(t *testing.T) {
	pool := NewPool(&PoolConfig{
		Min: 2,
		Max: 5,
	}, func() (interface{}, error) {
		conn, err := net.Dial("tcp", "0.0.0.0:30000")
		if err != nil {
			return nil, err
		}
		return conn, nil
	})

	pool.Populate()

	if len(pool.Channel) != 2 {
		t.Error("Pool channel is empty")
	}
}
