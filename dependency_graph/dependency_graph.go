package dependency_graph

var graph map[string][]string

func init() {
	graph = make(map[string][]string)
}

func AddToGraph(target string, dep string) {
	graph[target] = append(graph[target], dep)
}

func GetDeps(target string) []string {
	return graph[target]
}

func GetTargets() []string {
	var targets []string
	for target := range graph {
		targets = append(targets, target)
	}
	return targets
}

func ClearGraph() {
	graph = make(map[string][]string)
}

func GetGraphSize() int {
	return len(graph)
}
