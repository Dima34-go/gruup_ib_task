package main

import "time"


type resourcesRequest struct {
	timeEnd      time.Time
	infoChan     chan string
	successChan  chan bool
}

func newResourcesRequest(tm time.Time) *resourcesRequest {
	return &resourcesRequest{
		timeEnd:                         tm,
		infoChan:           make(chan string, 2),
		successChan:         make(chan bool, 2),
	}
}