package Lexer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Printer struct {
	msgChan chan string

	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

func NewPrinter() *Printer {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Printer{
		ctx:    ctx,
		cancel: cancel,
	}

	return p
}

func (p *Printer) Start() {
	p.msgChan = make(chan string)

	p.wg.Add(1)

	go p.msgListener()
}

func (p *Printer) Close() {
	select {
	case <-p.ctx.Done():
		// Do nothing
	default:
		p.cancel()

		close(p.msgChan)

		p.wg.Wait()

		p.msgChan = nil
	}
}

func (p *Printer) msgListener() {
	defer p.wg.Done()

	for msg := range p.msgChan {
		ok := strings.HasSuffix(msg, "\n")

		if ok {
			fmt.Print(msg)
		} else {
			fmt.Println(msg)
		}

		fmt.Println() // Add a new line
	}
}

func (p *Printer) Print(a ...interface{}) {
	select {
	case <-p.ctx.Done():
		// Do nothing
	default:
		p.msgChan <- fmt.Sprint(a...)
	}
}

func (p *Printer) Printf(format string, a ...interface{}) {
	select {
	case <-p.ctx.Done():
		// Do nothing
	default:
		p.msgChan <- fmt.Sprintf(format, a...)
	}
}

type Verbose struct {
	isActive bool
	printer  *Printer
}

func NewVerbose(active bool) *Verbose {
	p := NewPrinter()

	v := &Verbose{
		isActive: active,
		printer:  p,
	}

	p.Start()

	return v
}

func (v *Verbose) Close() {
	v.printer.Close()
}

func (v *Verbose) DoIf(doFunc func(p *Printer)) {
	if !v.isActive {
		return
	}

	doFunc(v.printer)
}
