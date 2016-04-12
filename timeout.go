package chanutil

import (
	"errors"
	"reflect"
	"time"
)

var ErrTimeout = errors.New("channel receive timed out")

// RecvTimeout takes a readable channel (chan or <-chan, but not chan<-) and
// attempts to receive a value within a positive duration d. Returns the value
// received as interface{} or an error. If d is 0 it will return immediately.
// If d is negative it will wait forever.
func RecvTimeout(ch interface{}, d time.Duration) (v interface{}, err error) {
	t := reflect.TypeOf(ch)
	if t.Kind() != reflect.Chan || t.ChanDir()&reflect.RecvDir == 0 {
		panic("must be a readable channel")
	}

	src := reflect.ValueOf(ch)

	switch {
	case d < 0:
		// timeout negative, wait forever
		recv, _ := src.Recv()
		v = recv.Interface()

	case d == 0:
		// timeout zero, return immediately
		recv, ok := src.TryRecv()
		if ok {
			v = recv.Interface()
		} else {
			err = ErrTimeout
		}

	case d > 0:
		// timeout positive, wait until timeout
		cases := make([]reflect.SelectCase, 2)
		cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: src}
		cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(time.After(d))}

		chosen, recv, _ := reflect.Select(cases)
		if chosen == 0 {
			v = recv.Interface()
		} else {
			err = ErrTimeout
		}

	}
	return v, err
}
