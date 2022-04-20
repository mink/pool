package main

import (
	"fmt"
	"time"
)

type Pool[T any] struct {
	Config  *PoolConfig
	Factory func() (T, error)
	Channel chan T
}

type PoolConfig struct {
	Min            int
	Max            int
	FactoryTimeout time.Duration
}

func NewPool[T any](config *PoolConfig, factory func() (T, error)) *Pool[T] {
	return &Pool[T]{
		Config:  config,
		Factory: factory,
		Channel: make(chan T, config.Max),
	}
}

func (pool *Pool[T]) Populate() {
	for i := len(pool.Channel); i < pool.Config.Min; i++ {
		conn, err := pool.Factory()
		if err != nil {
			panic(err)
		}
		pool.Channel <- conn
	}
}

func (pool *Pool[T]) Get() (interface{}, error) {
	go pool.Populate()
	select {
	case conn := <-pool.Channel:
		return conn, nil
	case <-time.After(pool.Config.FactoryTimeout):
		return nil, fmt.Errorf("pool factory timeout")
	}
}
