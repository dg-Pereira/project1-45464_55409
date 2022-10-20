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
	Target     string
	ToChildren []chan *Msg
	ToParents  chan *Msg
	ParentNum  int
}

type MsgType = int

const (
	BuildSuccess MsgType = iota
	BuildError
)

type Msg struct {
	Type      MsgType
	Timestamp time.Time
}

func NewGraph() map[string][]*Node {
	return make(map[string][]*Node)
}

func Add(target string, newNode *Node, graph map[string][]*Node) map[string][]*Node {

	if _, ok := graph[target]; !ok {
		graph[target] = make([]*Node, 0)
	}
	graph[target] = append(graph[target], newNode)

	return graph
}

func IsLeaf(node *Node, graph map[string][]*Node) bool {
	return len(graph[node.Target]) == 0
}
