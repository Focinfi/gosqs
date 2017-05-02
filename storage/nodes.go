package storage

import (
	"encoding/json"

	"github.com/Focinfi/gosqs/errors"
	"github.com/Focinfi/gosqs/models"
)

// Nodes for the cluster of server node
type Nodes struct {
	store *Storage
	db    models.KV
}

// All returns the all the joined nodes
func (n *Nodes) All(key string) ([]string, error) {
	val, err := n.db.Get(key)
	if err == errors.DataNotFound {
		return []string{}, nil
	}

	if err != nil {
		return nil, err
	}

	nodes := []string{}
	if err := json.Unmarshal([]byte(val), &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

// UpdateAll reset with the given nodes
func (n *Nodes) UpdateAll(key string, nodes []string) error {
	nodesBytes, err := json.Marshal(nodes)
	if err != nil {
		return err
	}

	return n.db.Put(key, string(nodesBytes))
}
