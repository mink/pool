package main

import (
	"fmt"
	"time"
)

type Pool struct {
	Config  *PoolConfig
	Factory func() (interface{}, error)
	Channel chan interface{}
}

type PoolConfig struct {
	Min int
	Max int
}

func NewPool(config *PoolConfig, factory func() (interface{}, error)) *Pool {
	return &Pool{
		Config:  config,
		Factory: factory,
		Channel: make(chan interface{}, config.Max),
	}
}

func (pool *Pool) Populate() {
	for i := len(pool.Channel); i < pool.Config.Min; i++ {
		conn, err := pool.Factory()
		if err != nil {
			panic(err)
		}
		pool.Channel <- conn
	}
}

func (pool *Pool) Get() (interface{}, error) {
	select {
	case conn := <-pool.Channel:
		return conn, nil
	case <-time.After(time.Second):
		return nil, fmt.Errorf("timeout getting connection")
	}
}
