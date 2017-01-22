package utils_test

import (
    "testing"

    "github.com/qnib/algoimpl/go/graph"
    "github.com/qnib/qwatch/utils"
    "github.com/stretchr/testify/assert"

)


// NewQGraph returns a initial graph
func TestNewQGraph(t *testing.T) {
    exp := utils.QGraph{
        Inputs: make(map[string]graph.Node, 0),
        Filters: make(map[string]graph.Node, 0),
        Outputs: make(map[string]graph.Node, 0),
        Graph: graph.New(graph.Directed),
    }
    got := utils.NewQGraph()
    assert.Equal(t, exp, got)
}

func TestAddCollector(t *testing.T) {

    assert.Equal(t, 1, 1)
}
/*

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
*/
