package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var duration *time.Duration = flag.Duration("interval", 1000*time.Millisecond, "interval(unit: millisecond)")

type ThrottleWriter struct {
	sc     *bufio.Scanner
	writer *bufio.Writer
	ctlC   chan int
	rtnC   chan bool
}

const (
	flushCtl int = iota
)

func NewWriter(input io.Reader, output io.Writer) *ThrottleWriter {
	sc := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)

	w := &ThrottleWriter{
		sc:     sc,
		writer: writer,
		ctlC:   make(chan int),
		rtnC:   make(chan bool),
	}

	return w
}

func (w *ThrottleWriter) Run() {
	go func() {
		for w.sc.Scan() {
			if _, err := w.writer.Write(w.sc.Bytes()); err != nil {
				panic(err.Error())
			}
			if err := w.writer.WriteByte('\n'); err != nil {
				panic(err.Error())
			}
		}
	}()

	go w.loop()
}

func (w *ThrottleWriter) loop() {
OuterLoop:
	for {
		select {
		case ctl := <-w.ctlC:
			switch ctl {
			case flushCtl:
				// if err := w.writer.Flush(); err != nil {
				// 	fmt.Println(outW.String())
				// 	panic(err.Error())
				// }
				w.writer.Flush()
				w.Print()
				break OuterLoop
			}
		}
	}
	w.rtnC <- true
}

func (w *ThrottleWriter) Flush() {
	w.ctlC <- flushCtl
}

func (w *ThrottleWriter) Print() {
	w.writer.Flush()
}

func init() {
	flag.Parse()
}

func main() {
	// out := &writer{inputReader: os.Stdin, outWriter: os.Stdout, errWriter: os.Stderr}
	out := NewWriter(os.Stdin, os.Stdout)
	out.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	interval := time.Tick(*duration)

	go func() {
		for {
			<-interval
			out.Print()
		}
	}()

	<-c

	// if err := w.Flush(); err != nil {
	// 	panic(err.Error())
	// }
}
