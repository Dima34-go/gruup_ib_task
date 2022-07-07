package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	store := NewStorage()
	go store.StartWork()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			key := strings.TrimPrefix(r.URL.Path,"/")
			value :=r.URL.Query().Get("v")
			store.PushResources(key, value)
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			key := strings.TrimPrefix(r.URL.Path,"/")
			timeout,_ := strconv.Atoi(r.URL.Query().Get("timeout"))
			if timeout == 0 {
				if str, ok := store.PushRequestsWithoutTimeout(key); ok {
					fmt.Fprintf(w,"%s",str)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			} else {
				rq := newResourcesRequest(time.Now().Add(time.Duration(timeout)*time.Second))
				store.PushRequests(key,rq)

				select {
				case <-time.After(time.Duration(timeout)*time.Second):
					w.WriteHeader(http.StatusBadRequest)
					rq.successChan <- false
				case str := <-rq.infoChan:
					rq.successChan <- true
					fmt.Fprintf(w,"%s",str)
				}
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
		})
	http.ListenAndServe(":8080",mux)
}