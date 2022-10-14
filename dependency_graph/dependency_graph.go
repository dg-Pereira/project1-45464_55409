package dependency_graph

import "time"

type Graph interface {
	AddToGraph() map[string][]string
	GetDeps() []string
	GetTargets() []string
	ClearGraph()
	GetGraphSize() int
}

type Node struct {
	Object    string
	Visited   bool
	Timestamp time.Time
}

func MakeNode(object string) Node {
	return Node{object, false, time.Time{}}
}

func NewGraph() map[string][]*Node {
	return make(map[string][]*Node)
}

func AddToGraph(graph map[string][]*Node, target string, dep string) map[string][]*Node {
	newNode := &Node{dep, false, time.Time{}}
	graph[target] = append(graph[target], newNode)
	return graph
}

func GetDeps(graph map[string][]*Node, target Node) ([]*Node, bool) {
	deps, ok := graph[target.Object]
	return deps, !ok
}

func GetTargets(graph map[string][]*Node) []string {
	var targets []string
	for target := range graph {
		targets = append(targets, target)
	}
	return targets
}

func GetGraphSize(graph map[string][]*Node) int {
	return len(graph)
}
