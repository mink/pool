package main

import (
	"fmt"
	"time"
)

type Pool[T any] struct {
	Config  *PoolConfig
	Factory func() (T, error)
	Channel chan T
	InPool  int
	InUse   int
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

func (pool *Pool[T]) Count() int {
	return pool.InPool + pool.InUse
}

func (pool *Pool[T]) Populate() {
	for i := pool.Count(); i < pool.Config.Min; i++ {
		conn, err := pool.Factory()
		if err != nil {
			panic(err)
		}
		pool.InPool++
		pool.Channel <- conn
	}
}

func (pool *Pool[T]) Get() (T, error) {
	go pool.Populate()
	select {
	case conn := <-pool.Channel:
		pool.InPool--
		pool.InUse++
		return conn, nil
	case <-time.After(pool.Config.FactoryTimeout):
		return *new(T), fmt.Errorf("pool factory timeout")
	}
}
