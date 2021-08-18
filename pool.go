package main

type Pool struct {
	Factory func() (interface{}, error)
	Channel chan interface{}
}
