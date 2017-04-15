package storage

import (
	"encoding/json"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/models"
)

type Nodes struct {
	store *Storage
	db    models.KV
}

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

func (n *Nodes) UpdateAll(key string, nodes []string) error {
	nodesBytes, err := json.Marshal(nodes)
	if err != nil {
		return err
	}

	return n.db.Put(key, string(nodesBytes))
}
