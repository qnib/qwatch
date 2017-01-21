package utils

import (
    "fmt"

    "github.com/qnib/algoimpl/go/graph"
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

// PrintGraph outputs the graph to stdout
func (qg *QGraph) PrintGraph() {
    for k,n := range qg.Collectors {
        fmt.Printf("%s has downstream: %s\n", k, qg.Graph.ReachableNeighbors(n))
    }
}

// AddCollector adds a collector node to the graph
func (qg *QGraph) AddCollector(in string) {
    qg.Collectors[in] = qg.Graph.MakeNamedNode(in)
}

// AddOutput adds a handler to the graph, which either listens to
// - a collector directly if the name is found in the collectors list
// - a filter of some sort.
// if parents are empty it will be added without a parent
func (qg *QGraph) AddOutput(out string, parents []string) {
    qg.Outputs[out] = qg.Graph.MakeNamedNode(out)
    for _, p := range parents {
        if pn, ok := qg.Collectors[p]; ok {
            qg.Graph.MakeEdgeWeight(pn, qg.Outputs[out], 1)
        } else if pn, ok = qg.Filters[p]; ok {
            qg.Graph.MakeEdgeWeight(pn, qg.Outputs[out], 1)
        }
    }
}

// AddFilter adds a filter to the graph, which either listens to
// - a collector directly if the name is found in the collectors list
// - a filter of some sort.
// if parents are empty it will be added without a parent
func (qg *QGraph) AddFilter(f string, parents []string) {
    qg.Filters[f] = qg.Graph.MakeNamedNode(f)
    for _, p := range parents {
        if pn, ok := qg.Collectors[p]; ok {
            qg.Graph.MakeEdgeWeight(pn, qg.Filters[f], 1)
        } else if pn, ok = qg.Filters[p]; ok {
            qg.Graph.MakeEdgeWeight(pn, qg.Filters[f], 1)
        }
    }
}

// ConnectFilters looks for filters w/o inputs and will connect them to all collectors
func (qg *QGraph) ConnectFilters() {
    _ = 1
}

// ConnectOutputs looks for outputs w/o inputs (collectors/filters) and will connect them to all leafs of the collectors/filter graph.
func (qg *QGraph) ConnectOutputs() {
    _ = 1
}
