package bipartitegraph

import (
	"github.com/onsi/gomega/matchers/support/goraph/edge"
	"slices"
	"testing"
)

func buildEdgesArr(l, r []interface{}, edges edge.EdgeSet) []string {
	unpackArr := func(in []interface{}) []string {
		result := make([]string, 0, len(in))
		for _, el := range in {
			result = append(result, el.(string))
		}
		return result
	}

	vertexes := unpackArr(append(l, r...))

	result := make([]string, 0)
	for _, currEdge := range edges {
		result = append(result, vertexes[currEdge.Node1]+"-"+vertexes[currEdge.Node2])
	}
	return result
}

func expectedContains(t *testing.T, expected string, edges []string) {
	idx := slices.IndexFunc(edges, func(c string) bool { return c == expected })
	if idx == -1 {
		t.Fatalf("edges %v not contains expected: %s", edges, expected)
	}
}

func TestMaximumCardinalityMatch(t *testing.T) {
	edgesFunc := func(l, r interface{}) (bool, error) {
		ll := l.(string)
		rr := r.(string)

		type currEdge struct {
			l string
			r string
		}
		knownEdges := []currEdge{
			{"1", "A"},
			{"1", "B"},
			{"1", "C"},
			{"1", "D"},
			{"1", "E"},
			{"2", "A"},
			{"2", "D"},
			{"3", "B"},
			{"3", "D"},
			{"4", "B"},
			{"4", "D"},
			{"4", "E"},
			{"5", "A"},
		}

		for _, el := range knownEdges {
			if el.l == ll && el.r == rr {
				return true, nil
			}
		}
		return false, nil
	}

	leftPart := []interface{}{"1", "2", "3", "4", "5"}
	rightPart := []interface{}{"A", "B", "C", "D", "E"}

	bipartiteGraph, err := NewBipartiteGraph(
		leftPart,
		rightPart,
		edgesFunc,
	)
	if err != nil {
		t.Fatalf("NewBipartiteGraph returned error: %v", err)
	}
	if err != nil {
		t.Fatal(err)
	}
	edgeSet := bipartiteGraph.LargestMatching()
	if len(edgeSet) != 5 {
		t.Fatalf("bipartiteGraph.LargestMatching() returned not 5 elements: %v", edgeSet)
	}
	edges := buildEdgesArr(leftPart, rightPart, edgeSet)
	expectedContains(t, "1-C", edges)
	expectedContains(t, "2-D", edges)
	expectedContains(t, "3-B", edges)
	expectedContains(t, "4-E", edges)
	expectedContains(t, "5-A", edges)
}
