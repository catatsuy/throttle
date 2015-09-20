package main

import (
	"bufio"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var duration *time.Duration = flag.Duration("interval", 1000*time.Millisecond, "interval(unit: millisecond)")

func init() {
	flag.Parse()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	sc := bufio.NewScanner(os.Stdin)
	w := bufio.NewWriter(os.Stdout)

	interval := time.Tick(*duration)
	go func() {
		for sc.Scan() {
			if _, err := w.Write(sc.Bytes()); err != nil {
				panic(err.Error())
			}
			if err := w.WriteByte('\n'); err != nil {
				panic(err.Error())
			}
		}
	}()

	go func() {
		for {
			<-interval
			w.Flush()
		}
	}()

	<-c

	w.Flush()
}
