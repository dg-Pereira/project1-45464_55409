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
	go func() {
		<-reqCh
		graph, nodes := makeGraph(file)
		printGraph(graph)

		for node := range nodes {
			go build(nodes[node], graph)
		}

		//wait for message from root
		<-nodes[file.Rules[0].Object].ToParents

		reqCh <- &dependency_graph.Msg{Type: dependency_graph.BuildSuccess}
	}()

	//make leaves start waiting from 1 message, to make the leave message start the build
	return reqCh
}

func buildNode(node *dependency_graph.Node) {
	_, err := utils.Build(node.Target)
	if err != nil {
		for i := 0; i < node.ParentNum; i++ {
			node.ToParents <- &dependency_graph.Msg{Type: dependency_graph.BuildError}
		}
	} else {
		for i := 0; i < node.ParentNum; i++ {
			node.ToParents <- &dependency_graph.Msg{Type: dependency_graph.BuildSuccess}
		}
	}
}

func build(node *dependency_graph.Node, graph map[string][]*dependency_graph.Node) {

	// wait for messages from all children to build
	for _, child := range node.ToChildren {
		m := <-child
		if m.Type == dependency_graph.BuildError {
			node.ToParents <- m
			return
		}
	}

	if dependency_graph.IsLeaf(node, graph) {
		_, err := utils.Status(node.Target)
		if err == nil { //don't build
			node.ToParents <- &dependency_graph.Msg{Type: dependency_graph.BuildSuccess}
		} else {
			buildNode(node)
		}
	} else {
		buildNode(node)
	}
}

func printGraph(graph map[string][]*dependency_graph.Node) {
	println("===== Graph =====")
	for target, deps := range graph {
		println(target)
		for _, dep := range deps {
			println("\t", dep.Target, dep.ParentNum)
		}
	}
	println("=================")
}

func makeGraph(file *parser.DepFile) (map[string][]*dependency_graph.Node, map[string]*dependency_graph.Node) {
	graph := dependency_graph.NewGraph()
	nodes := make(map[string]*dependency_graph.Node)

	for _, rule := range file.Rules {
		target := rule.Object

		if _, ok := nodes[target]; !ok {
			nodes[target] = &dependency_graph.Node{Target: target, ToChildren: make([]chan *dependency_graph.Msg, 0), ParentNum: 1}
		}

		for _, dep := range rule.Deps {
			if _, ok := nodes[dep]; !ok {
				newNode := &dependency_graph.Node{Target: dep, ToChildren: make([]chan *dependency_graph.Msg, 0), ParentNum: 1}
				graph = dependency_graph.Add(target, newNode, graph)
				nodes[dep] = newNode
			} else {
				newNode := nodes[dep]
				graph = dependency_graph.Add(target, newNode, graph)
			}
		}
	}

	//add communication channels between nodes
	nodes[file.Rules[0].Object].ToParents = make(chan *dependency_graph.Msg)
	for target := range graph {
		for _, node := range graph[target] {
			if node.ToParents != nil {
				nodes[target].ToChildren = append(nodes[target].ToChildren, node.ToParents)
				node.ParentNum++
			} else {
				toChild := make(chan *dependency_graph.Msg)
				node.ToParents = toChild
				nodes[target].ToChildren = append(nodes[target].ToChildren, toChild)
			}
		}
	}

	println(nodes["root"].Target, len(nodes["root"].ToChildren), nodes["root"].ToParents)

	return graph, nodes

}
