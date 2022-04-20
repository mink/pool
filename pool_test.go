package main

import (
	"net"
	"testing"
	"time"
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
		Min:            1,
		Max:            5,
		FactoryTimeout: time.Second,
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
		Min:            2,
		Max:            5,
		FactoryTimeout: time.Second,
	}, func() (interface{}, error) {
		conn, err := net.Dial("tcp", "0.0.0.0:30000")
		if err != nil {
			return nil, err
		}
		return conn, nil
	})

	pool.Populate()

	if len(pool.Channel) != pool.Config.Min {
		t.Errorf("Pool channel size = %d, expected %d", len(pool.Channel), pool.Config.Min)
	}
}

func TestGetConnection(t *testing.T) {
	pool := NewPool(&PoolConfig{
		Min:            2,
		Max:            5,
		FactoryTimeout: time.Second,
	}, func() (interface{}, error) {
		conn, err := net.Dial("tcp", "0.0.0.0:30000")
		if err != nil {
			return nil, err
		}
		return conn, nil
	})

	pool.Populate()

	conn, err := pool.Get()
	if err != nil {
		t.Error(err)
	}

	if len(pool.Channel) != (pool.Config.Min - 1) {
		t.Errorf("Pool channel size = %d, expected %d", len(pool.Channel), pool.Config.Min-1)
	}

	if _, ok := conn.(*net.TCPConn); !ok {
		t.Error("Connection received is of incorrect type")
	}
}

func TestFactoryTimeout(t *testing.T) {
	config := &PoolConfig{
		Min:            1,
		Max:            2,
		FactoryTimeout: time.Millisecond * 50,
	}
	pool := NewPool(config, func() (string, error) {
		time.Sleep(time.Second)
		return "hello", nil
	})

	item, err := pool.Get()

	if item != nil {
		t.Errorf("item should be null, got %s", item)
	}

	if err == nil {
		t.Error("error expected")
	}

	if err.Error() != "pool factory timeout" {
		t.Errorf("expected error \"pool factory timeout\", got error \"%s\"\n", err)
	}
}
