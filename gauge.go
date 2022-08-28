package main

import "errors"

// RequestGauge is a struct that uses a buffered channel for it's queue implementation
// to limit the number of simultaneous in the application
type RequestGauge struct {
	channel chan struct{}
}

func (g *RequestGauge) run(callback func()) error {
	select {

	// successful reads from the channel indicates theres a spot in our queue.
	// reading hangs when the buffer is empty(simultaneous limit reached) causing the default case to run
	case <-g.channel:

		// runs given callback function
		callback()

		// write to the channel, enabling a read from next goroutine.
		g.channel <- struct{}{}
		return nil

	default:
		return errors.New("request gauge capacity exceeded")
	}
}

func NewRequestGauge(capacity int) *RequestGauge {
	return &RequestGauge{
		channel: getFilledChan(struct{}{}, capacity),
	}
}

func getFilledChan[T any](token T, capacity int) chan T {
	c := make(chan T, capacity)

	for i := 0; i < capacity; i++ {
		c <- token
	}

	return c
}
