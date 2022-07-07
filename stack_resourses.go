package main

import (
	"container/list"
	"fmt"
	"sync"
)

type Stack struct {
	stack *list.List
	mx    *sync.Mutex
}

func NewStack() *Stack {
	return &Stack{
		stack: list.New(),
		mx:    new(sync.Mutex),
	}
}

func (c *Stack) Push(value string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.stack.PushFront(value)
}

func (c *Stack) GetFront() (string, error) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if c.stack.Len() > 0 {
		if val, ok := c.stack.Front().Value.(string); ok {
			return val, nil
		}
		return "", fmt.Errorf("Peep Error: Stack Datatype is incorrect")
	}
	return "", fmt.Errorf("Peep Error: Stack is empty")
}

func (c *Stack) Empty() bool {
	return c.stack.Len() == 0
}

func (c *Stack) Pop() error {
	if c.stack.Len() > 0 {
		ele := c.stack.Front()
		c.stack.Remove(ele)
	}
	return fmt.Errorf("Pop Error: Stack is empty")
}