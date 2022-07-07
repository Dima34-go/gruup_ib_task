package main

import (
	"container/list"
	"fmt"
	"sync"
)

type Queue struct {
	queue *list.List
	mx    *sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		queue: list.New(),
		mx:    new(sync.Mutex),
	}
}

func (c *Queue) Push(rq *resourcesRequest) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.queue.PushBack(rq)
}

func (c *Queue) GetFront() (*resourcesRequest, error) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if c.queue.Len() > 0 {
		if val, ok := c.queue.Front().Value.(*resourcesRequest); ok {
			return val, nil
		}
		return nil, fmt.Errorf("Peep Error: Stack Datatype is incorrect")
	}
	return nil, fmt.Errorf("Peep Error: Queue is empty")
}

func (c *Queue) Empty() bool {
	return c.queue.Len() == 0
}

func (c *Queue) Pop() error {
	if c.queue.Len() > 0 {
		ele := c.queue.Front()
		c.queue.Remove(ele)
	}
	return fmt.Errorf("Pop Error: Stack is empty")
}