package master

import "github.com/Focinfi/gosqs/models"

type nodes map[string]models.NodeInfo

func (m nodes) nodeURLSlice() []string {
	nodes := make([]string, len(m))
	i := 0
	for node := range m {
		nodes[i] = node
		i++
	}

	return nodes
}

func (m nodes) statsSlice() models.NodeSlice {
	slice := make([]models.NodeInfo, len(m))
	i := 0
	for node := range m {
		slice[i] = m[node]
		i++
	}

	return slice
}

func nodeURLSliceToNodes(nodes []string) nodes {
	m := make(map[string]models.NodeInfo, len(nodes))
	for _, node := range nodes {
		m[node] = models.NodeInfo{}
	}

	return m
}
