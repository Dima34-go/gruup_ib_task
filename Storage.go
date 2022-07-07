package main

import (
	"sync"
	"time"
)

type KeyStorage struct {
	resources   *Stack
	requests    *Queue
	mx          *sync.Mutex
}

func NewKeyStorage() *KeyStorage {
	return &KeyStorage{
		resources: NewStack(),
		requests:  NewQueue(),
		mx:        new(sync.Mutex),
	}
}

type Storage struct {
	s    map[string]*KeyStorage
}

func NewStorage() *Storage {
	return &Storage{
		s: make(map[string]*KeyStorage),
	}
}

func (s *Storage) PushResources(key, value string)  {
	if store, ok := s.s[key]; ok {
		store.resources.Push(value)
		return
	}
	s.s[key] = NewKeyStorage()
	s.s[key].resources.Push(value)
}

func (s *Storage) PushRequests(key string,rq *resourcesRequest) {
	if store, ok := s.s[key]; ok {
		store.requests.Push(rq)
		return
	}
	s.s[key] = NewKeyStorage()
	s.s[key].requests.Push(rq)
}

func (s *Storage) StartWork()  {
	for true {
		for _, keySt := range s.s {
			//if !keySt.requests.Empty() && !keySt.resources.Empty() {
			//	rq, _ := keySt.requests.GetFront()
			//	keySt.requests.Pop()
			//	str, _ := keySt.resources.GetFront()
			//	rq.infoChan <- str
			//	if <- rq.successChan {
			//		keySt.resources.Pop()
			//	}
			//}
			if !keySt.requests.Empty() {
				keySt.mx.Lock()
				rq, _ := keySt.requests.GetFront()
				if rq.timeEnd.Sub(time.Now()) < 0 {
					keySt.requests.Pop()
				}else if !keySt.resources.Empty() {
					keySt.requests.Pop()
					str, _ := keySt.resources.GetFront()
					rq.infoChan <- str
					if <- rq.successChan {
						keySt.resources.Pop()
					}
				}
				keySt.mx.Unlock()
			}
		}
	}
}
func (s *Storage) PushRequestsWithoutTimeout(key string) (string, bool) {
	if keySt, ok := s.s[key]; ok {
		keySt.mx.Lock()
		defer keySt.mx.Unlock()
		if !keySt.resources.Empty() {
			str, _ := keySt.resources.GetFront()
			keySt.resources.Pop()
			return  str, true
		}
		return "", false
	}
	return  "", false
}