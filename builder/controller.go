package builder

import (
	"cpl_go_proj22/dependency_graph"
	"cpl_go_proj22/parser"
	"cpl_go_proj22/utils"
)

func MakeController(file *parser.DepFile) chan *dependency_graph.Msg {
	reqCh := make(chan *dependency_graph.Msg)
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
	// go func() {
	// 	<-reqCh
	// 	graph, root := makeGraph(file)
	// 	printGraph(graph)
	// 	ch := make(chan *dependency_graph.Msg)
	// 	for obj := range graph {
	// 		//go build(obj, graph, root, ch)
	// 	}
	// 	<-ch
	// 	reqCh <- &dependency_graph.Msg{Type: dependency_graph.BuildSuccess}
	// }()

	graph, _ := makeGraph(file)
	printGraph(graph)
	return reqCh
}

func build(obj string, graph map[string][]dependency_graph.Node, root string, toDependents chan *dependency_graph.Msg) {
	//wait until can build

	_, err := utils.Build(obj)
	if err == nil {
		toDependents <- &dependency_graph.Msg{Type: dependency_graph.BuildSuccess}
	} else {
		toDependents <- &dependency_graph.Msg{Type: dependency_graph.BuildError}
	}
}

func printGraph(graph map[string][]dependency_graph.Node) {
	println("===== Graph =====")
	for target, deps := range graph {
		println(target)
		for _, dep := range deps {
			println("\t", dep.Target)
		}
	}
	println("=================")
}

func makeGraph(file *parser.DepFile) (map[string][]dependency_graph.Node, string) {
	graph := dependency_graph.NewGraph()
	for _, rule := range file.Rules {
		toDependents := make(chan *dependency_graph.Msg)
		for _, dep := range rule.Deps {
			graph = dependency_graph.Add(dep, rule.Object, graph, toDependents)
		}
	}
	root := file.Rules[0].Object
	return graph, root
}
