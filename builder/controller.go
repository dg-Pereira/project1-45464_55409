package builder

import (
	"cpl_go_proj22/dependency_graph"
	"cpl_go_proj22/parser"
	"cpl_go_proj22/utils"
	"time"
)

type MsgType = int

const (
	BuildSuccess MsgType = iota
	BuildError
)

type Msg struct {
	Type      MsgType
	object    string
	timestamp time.Time
	//TODO: May add more fields here.
}

func MakeController(file *parser.DepFile) chan *Msg {
	reqCh := make(chan *Msg)
	//rootCh := make(chan *Msg)
	//var leafCh chan interface{} //TODO: You may want to change this type
	// TODO: Startup system that emits outcome of build on rootCh, triggered by an output on leafCh
	//go func() {
	//for{
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
	//}
	go func() {
		<-reqCh
		graph, root := makeGraph(file)
		ch := make(chan *Msg)
		go build(graph, &root, ch)
		<-ch
		reqCh <- &Msg{Type: BuildSuccess, object: root.Object}
	}()
	return reqCh
}

func makeGraph(file *parser.DepFile) (map[string][]*dependency_graph.Node, dependency_graph.Node) {
	graph := dependency_graph.NewGraph()
	for _, rule := range file.Rules {
		for _, dep := range rule.Deps {
			graph = dependency_graph.AddToGraph(graph, rule.Object, dep)
		}
	}
	printGraph(graph)
	return graph, dependency_graph.MakeNode(file.Rules[0].Object)
}

func printGraph(graph map[string][]*dependency_graph.Node) {
	println("===== Graph =====")
	for target, deps := range graph {
		println(target)
		for _, dep := range deps {
			println("\t", dep.Object)
		}
	}
	println("=================")
}

func build(graph map[string][]*dependency_graph.Node, node *dependency_graph.Node, ch chan *Msg) {
	//println("Building", node.Object)
	deps, isLeaf := dependency_graph.GetDeps(graph, *node)
	currMaxTimestamp := time.Time{}
	if !isLeaf {
		i := 0
		loc_ch := make(chan *Msg)
		for _, dep := range deps {
			//println(dep)
			//println(dep.Object)
			//println("Building ", node.Object, "launching ", dep.Object)
			go build(graph, dep, loc_ch)
			i++
		}

		//wait for launched goroutines to finish
		for j := i; j > 0; j-- {
			m := <-loc_ch
			if m.Type == BuildError {
				ch <- m
				return
			}
			if m.timestamp.After(currMaxTimestamp) {
				currMaxTimestamp = m.timestamp
			}
		}
	} else { // if isLeaf, only build the object if object file does not exist
		_, err := utils.Status(node.Object)
		if err == nil {
			ch <- &Msg{Type: BuildSuccess, object: node.Object, timestamp: node.Timestamp}
			return
		} // else continue to build
	}

	println(currMaxTimestamp.UnixMicro(), "Building", node.Object)

	//only build the object if it is out of date
	//that is, if any of the dependencies are more recent than the object

	if !currMaxTimestamp.Before(node.Timestamp) {
		timeNow := time.Now()
		_, err := utils.Build(node.Object)
		if err == nil {
			node.Timestamp = timeNow
			ch <- &Msg{Type: BuildSuccess, object: node.Object, timestamp: timeNow}
		} else {
			ch <- &Msg{Type: BuildError, object: node.Object, timestamp: timeNow}
		}
	} else {
		ch <- &Msg{Type: BuildSuccess, object: node.Object, timestamp: node.Timestamp}
	}
}
