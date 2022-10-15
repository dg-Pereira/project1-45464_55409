package dependency_graph

type Graph interface {
	AddToGraph() map[string][]string
	GetDeps() []string
	GetTargets() []string
	ClearGraph()
	GetGraphSize() int
}

type Node struct {
	Target       string
	ToDependents chan *Msg
}

type MsgType = int

const (
	BuildSuccess MsgType = iota
	BuildError
)

// need to replicate Msg to not have circular dependency
type Msg struct {
	Type MsgType
	//TODO: May add more fields here.
}

func NewGraph() map[string][]Node {
	return make(map[string][]Node)
}

func Add(dep string, target string, graph map[string][]Node, toDependents chan *Msg) map[string][]Node {
	if _, ok := graph[dep]; !ok {
		graph[dep] = []Node{}
	}
	graph[dep] = append(graph[dep], Node{Target: target, ToDependents: toDependents})
	return graph
}
