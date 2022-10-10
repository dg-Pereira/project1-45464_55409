package builder

import (
	"cpl_go_proj22/dependency_graph"
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
	//rootCh := make(chan *Msg)
	//var leafCh chan interface{} //TODO: You may want to change this type
	// TODO: Startup system that emits outcome of build on rootCh, triggered by an output on leafCh
	//go func() {
	//	<-reqCh
	//	leafCh <- nil
	//	m := <-rootCh
	//	switch m.Type {
	//	case BuildSuccess:
	//		reqCh <- m
	//		break
	//	case BuildError:
	//		reqCh <- m
	//		break
	//	}
	//}()
	makeGraph(file)
	return reqCh
}

func makeGraph(file *parser.DepFile) map[string][]string {
	graph := dependency_graph.NewGraph()
	for _, rule := range file.Rules {
		for _, dep := range rule.Deps {
			graph = dependency_graph.AddToGraph(graph, rule.Object, dep)
		}
	}
	return graph
}
