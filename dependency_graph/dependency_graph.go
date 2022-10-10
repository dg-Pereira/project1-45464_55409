package dependency_graph

type Graph interface {
	AddToGraph() map[string][]string
	GetDeps() []string
	GetTargets() []string
	ClearGraph()
	GetGraphSize() int
}

func NewGraph() map[string][]string {
	return make(map[string][]string)
}

func AddToGraph(graph map[string][]string, target string, dep string) map[string][]string {
	graph[target] = append(graph[target], dep)
	return graph
}

func GetDeps(graph map[string][]string, target string) []string {
	return graph[target]
}

func GetTargets(graph map[string][]string) []string {
	var targets []string
	for target := range graph {
		targets = append(targets, target)
	}
	return targets
}

func GetGraphSize(graph map[string][]string) int {
	return len(graph)
}
