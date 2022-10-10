package builder

import (
	"cpl_go_proj22/parser"
)

type MsgType = int

const (
	BuildSuccess MsgType = iota
	BuildError
)

type Msg struct {
	Type MsgType
	//TODO: May add more fields here.
}

func MakeController(file *parser.DepFile) chan *Msg {
	reqCh := make(chan *Msg)
	rootCh := make(chan *Msg)
	var leafCh chan interface{} //TODO: You may want to change this type
	// TODO: Startup system that emits outcome of build on rootCh, triggered by an output on leafCh
	go func() {
		<-reqCh
		leafCh <- nil
		m := <-rootCh
		switch m.Type {
		case BuildSuccess:
			reqCh <- m
			break
		case BuildError:
			reqCh <- m
			break
		}
	}()
	return reqCh
}
