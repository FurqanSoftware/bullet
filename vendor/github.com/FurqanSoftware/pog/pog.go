package pog

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Pogger struct {
	out      io.Writer
	logger   *log.Logger
	status   Status
	stopCh   chan struct{}
	m        sync.Mutex
	initOnce sync.Once
}

func NewPogger(out io.Writer, prefix string, flag int) *Pogger {
	pogger := Pogger{
		out:    out,
		logger: log.New(out, prefix, flag),
		stopCh: make(chan struct{}),
	}
	go pogger.loop()
	return &pogger
}

func (p *Pogger) SetStatus(status Status) {
	p.m.Lock()
	p.status = status
	p.m.Unlock()
}

func (p *Pogger) Debug(v ...any) {
	a := []any{color.WhiteString("[d] ")}
	a = append(a, v...)
	p.logger.Print(a...)
}

func (p *Pogger) Debugln(v ...any) {
	a := []any{color.WhiteString("[d]")}
	a = append(a, v...)
	p.logger.Println(a...)
}

func (p *Pogger) Debugf(format string, v ...any) {
	p.logger.Printf(color.WhiteString("[d] ")+format, v...)
}

func (p *Pogger) Info(v ...any) {
	a := []any{"[i] "}
	a = append(a, v...)
	p.logger.Print(a...)
}

func (p *Pogger) Infoln(v ...any) {
	a := []any{"[i]"}
	a = append(a, v...)
	p.logger.Println(a...)
}

func (p *Pogger) Infof(format string, v ...any) {
	p.logger.Printf("[i] "+format, v...)
}

func (p *Pogger) Warn(v ...any) {
	a := []any{color.YellowString("[w] ")}
	a = append(a, v...)
	p.logger.Print(a...)
}

func (p *Pogger) Warnln(v ...any) {
	a := []any{color.YellowString("[w]")}
	a = append(a, v...)
	p.logger.Println(a...)
}

func (p *Pogger) Warnf(format string, v ...any) {
	p.logger.Printf(color.YellowString("[w] ")+format, v...)
}

func (p *Pogger) Error(v ...any) {
	a := []any{color.RedString("[E] ")}
	a = append(a, v...)
	p.logger.Print(a...)
}

func (p *Pogger) Errorln(v ...any) {
	a := []any{color.RedString("[E]")}
	a = append(a, v...)
	p.logger.Println(a...)
}

func (p *Pogger) Errorf(format string, v ...any) {
	p.logger.Printf(color.RedString("[E] ")+format, v...)
}

func (p *Pogger) Fatal(v ...any) {
	p.Error(v...)
	os.Exit(1)
}

func (p *Pogger) Fatalln(v ...any) {
	p.Errorln(v...)
	os.Exit(1)
}

func (p *Pogger) Fatalf(format string, v ...any) {
	p.Errorf(format, v...)
	os.Exit(1)
}

func (p *Pogger) Stop() {
	close(p.stopCh)
}

func (p *Pogger) loop() {
	cur := ""
	pad := ""
	if p.logger.Flags() != 0 {
		pad = strings.Repeat(" ", 20)
	}
L:
	for i := 0; ; i = (i + 1) % 10 {
		var s string
		p.m.Lock()
		b := []byte{' '}
		if p.status != nil {
			if i < 5 || !p.status.Throb() {
				b[0] = p.status.Icon()
			}
			s = "[" + string(b) + "]"
			if color := p.status.Color(); color != nil {
				s = color.Sprint(s)
			}
			s += " " + p.status.Text()
		}
		p.m.Unlock()
		if s != cur {
			fmt.Fprintf(p.out, "\033[2K\r%s%s\r", pad, s)
			cur = s
		}
		select {
		case <-p.stopCh:
			break L
		case <-time.After(125 * time.Millisecond):
		}
	}
}

type Status interface {
	Icon() byte
	Text() string
	Color() *color.Color
	Throb() bool
}

var (
	defaultPogger *Pogger
)

func InitDefault() {
	defaultPogger = NewPogger(log.Default().Writer(), log.Default().Prefix(), log.Default().Flags())
}

func Default() *Pogger {
	return defaultPogger
}

func SetStatus(status Status)        { defaultPogger.SetStatus(status) }
func Debug(v ...any)                 { defaultPogger.Debug(v...) }
func Debugln(v ...any)               { defaultPogger.Debugln(v...) }
func Debugf(format string, v ...any) { defaultPogger.Debugf(format, v...) }
func Info(v ...any)                  { defaultPogger.Info(v...) }
func Infoln(v ...any)                { defaultPogger.Infoln(v...) }
func Infof(format string, v ...any)  { defaultPogger.Infof(format, v...) }
func Warn(v ...any)                  { defaultPogger.Warn(v...) }
func Warnln(v ...any)                { defaultPogger.Warnln(v...) }
func Warnf(format string, v ...any)  { defaultPogger.Warnf(format, v...) }
func Error(v ...any)                 { defaultPogger.Error(v...) }
func Errorln(v ...any)               { defaultPogger.Errorln(v...) }
func Errorf(format string, v ...any) { defaultPogger.Errorf(format, v...) }
func Fatal(v ...any)                 { defaultPogger.Fatal(v...) }
func Fatalln(v ...any)               { defaultPogger.Fatalln(v...) }
func Fatalf(format string, v ...any) { defaultPogger.Fatalf(format, v...) }
func Stop()                          { defaultPogger.Stop() }
