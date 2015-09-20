package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	flag.Parse()
	rand.Seed(time.Now().Unix())
}

func main() {
	for i := 0; i < 1000; i++ {
		fmt.Printf("Welcome %d times\n", i)
		// sleep 0.1-0.3s
		time.Sleep((time.Duration)(rand.Intn(3)+1) * 100 * time.Millisecond)
	}
}
