package main

import (
	"github.com/sh3rp/gabel/core"
)

func main() {
	t1 := core.Loopback{}
	t2 := core.Loopback{}
	t1state := core.NewState()
	t2state := core.NewState()

	t1.AddListener(t1state)
	t2.AddListener(t2state)

}
