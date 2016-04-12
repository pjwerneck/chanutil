package chanutil

import (
	"testing"
	"time"
)

func TestRecvNegativeTimeout(t *testing.T) {
	c := make(chan int)

	time.AfterFunc(time.Millisecond, func() { c <- 1 })

	v, err := RecvTimeout(c, -1)

	if err != nil {
		t.Fatal("expected err == nil, got", err)
	}

	if v != 1 {
		t.Fatal("expected v == 1, got", v)
	}
}

func TestRecvZeroTimeout(t *testing.T) {
	c := make(chan int, 1)
	c <- 1

	v, err := RecvTimeout(c, 0)

	if err != nil {
		t.Fatal("expected err == nil, got", err)
	}

	if v != 1 {
		t.Fatal("expected v == nil, got", v)
	}
}

func TestRecvZeroTimeoutEmpty(t *testing.T) {
	c := make(chan int)

	v, err := RecvTimeout(c, 0)

	if err != ErrTimeout {
		t.Fatal("expected err == ErrTimeout, got", err)
	}

	if v != nil {
		t.Fatal("expected v == nil, got", v)
	}
}

func TestRecvTimedOut(t *testing.T) {
	c := make(chan int)

	time.AfterFunc(time.Duration(10)*time.Millisecond, func() { c <- 1 })

	v, err := RecvTimeout(c, time.Millisecond)

	if err != ErrTimeout {
		t.Fatal("expected err == ErrTimeout, got", err)
	}

	if v != nil {
		t.Fatal("expected v == nil, got", v)
	}
}

func TestRecvNotTimedOut(t *testing.T) {
	c := make(chan int)

	time.AfterFunc(time.Millisecond, func() { c <- 1 })

	v, err := RecvTimeout(c, time.Duration(10)*time.Millisecond)

	if err != nil {
		t.Fatal("expected err == nil, got", err)
	}

	if v != 1 {
		t.Fatal("expected v == nil, got", v)
	}
}
