package main

import (
	"fmt"
	"sync"
	"time"
)

type (

	// **Event** saved the data
	Event struct {
		data int
	}

	// **Observer** defined the NotifyCallback function
	Observer interface {
		NotifyCallback(Event)
	}

	// **Subject** interface
	Subject interface {
		AddListener(Observer)
		RemoveListener(Observer)
		Notify(Event)
	}

	eventObserver struct {
		id   int
		time time.Time
	}

	eventSubject struct {
		observers sync.Map
	}
)

// the **eventObserver is implemented all the functions
// of the type **Observer** interface
func (e *eventObserver) NotifyCallback(event Event) {
	fmt.Printf("Observer: %d, Received: %d, After %v\n",
		e.id, event.data, time.Since(e.time))
}

// the **eventSubject** is implemented all the
// functions of the type **Subject** inteface
func (s *eventSubject) AddListener(obs Observer) {
	s.observers.Store(obs, struct{}{})
}

func (s *eventSubject) RemoveListener(obs Observer) {
	s.observers.Delete(obs)
}

func (s *eventSubject) Notify(event Event) {
	s.observers.Range(func(key interface{}, value interface{}) bool {
		if key == nil || value == nil {
			return false
		}

		key.(Observer).NotifyCallback(event)
		return true
	})
}

func fib(n int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i, j := 0, 1; i < n; i, j = i+j, i {
			out <- i
		}
	}()
	return out
}

// 0, 1, 1, 2, 3, 5, 8, 13, 21, 34 ...
func main() {
	// for x := range fib(10000000) {
	// 	fmt.Println(x)
	// }
	n := eventSubject{
		observers: sync.Map{},
	}

	var obs1 = eventObserver{id: 1, time: time.Now()}
	var obs2 = eventObserver{id: 2, time: time.Now()}

	// access the pointer of the **eventObserver** struct
	// &struct{} -> interface
	n.AddListener(&obs1)
	n.AddListener(&obs2)

	go func() {
		select {
		case <-time.After(time.Microsecond * 20):
			n.RemoveListener(&obs2)

		}
	}()

	for x := range fib(10000000) {
		n.Notify(Event{data: x})
	}
}
