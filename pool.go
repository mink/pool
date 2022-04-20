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
		item, err := pool.Factory()
		if err != nil {
			panic(err)
		}
		pool.InPool++
		pool.Channel <- item
	}
}

func (pool *Pool[T]) Pop() (T, error) {
	go pool.Populate()
	select {
	case item := <-pool.Channel:
		pool.InPool--
		pool.InUse++
		return item, nil
	case <-time.After(pool.Config.FactoryTimeout):
		return *new(T), fmt.Errorf("pool factory timeout")
	}
}

func (pool *Pool[T]) Push(item T) {
	pool.InUse--
	pool.InPool++
	pool.Channel <- item
}
