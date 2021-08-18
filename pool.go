package main

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
	for len(pool.Channel) < pool.Config.Min {
		conn, err := pool.Factory()
		if err != nil {
			panic(err)
		}
		pool.Channel <- conn
	}
}
