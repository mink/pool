package main

type Pool struct {
	Factory func() (interface{}, error)
	Channel chan interface{}
}

func NewPool(factory func() (interface{}, error)) *Pool {
	return &Pool{Factory: factory}
}
