package utils

import (
    "github.com/twmb/algoimpl/go/graph"
)

// QGraph provides a DAG to keep the relations between input,filters and output
type QGraph struct {
    Collectors  map[string]graph.Node
    Filters     map[string]graph.Node
    Outputs     map[string]graph.Node
    Graph       *graph.Graph
}

// NewQGraph returns a initial graph
func NewQGraph() QGraph {
    qg := QGraph{
        Collectors: make(map[string]graph.Node, 0),
        Filters: make(map[string]graph.Node, 0),
        Outputs: make(map[string]graph.Node, 0),
        Graph: graph.New(graph.Directed),
    }
    return qg
}

// AddCollector adds a collector node to the graph
func (qg *QGraph) AddCollector(in string) {
    qg.Collectors[in] = qg.Graph.MakeNode()
}

// AddOutput adds a handler to the graph, which either listens to
func (qg *QGraph) AddOutput(out string, parents []string) {
    qg.Outputs[out] = qg.Graph.MakeNode()
    for _, p := range parents {
        if pn, ok := qg.Collectors[p]; ok {
            qg.Graph.MakeEdgeWeight(pn, qg.Outputs[out], 1)
        } else if pn, ok = qg.Filters[p]; ok {
            qg.Graph.MakeEdgeWeight(pn, qg.Outputs[out], 1)
        }
    }
}
