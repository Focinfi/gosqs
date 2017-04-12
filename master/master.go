package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/Focinfi/sqs/agent"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/node"
	"github.com/Focinfi/sqs/storage"
	"github.com/Focinfi/sqs/util/urlutil"
)

const (
	nodesKey              = "sqs.nodes"
	getNodeStatsURLFormat = "%s/stats"
)

var (
	heartbeatPeriod = time.Second
)

type nodes map[string]node.Info

func (m nodes) nodeURLSlice() []string {
	nodes := make([]string, len(m))
	i := 0
	for node := range m {
		nodes[i] = node
		i++
	}

	return nodes
}

func (m nodes) statsSlice() InfoSlice {
	slice := make([]node.Info, len(m))
	i := 0
	for node := range m {
		slice[i] = m[node]
	}

	return slice
}

func nodeURLSliceToNodes(nodes []string) nodes {
	m := make(map[string]node.Info, len(nodes))
	for _, node := range nodes {
		m[node] = node.Info{}
	}

	return m
}

// Service for a master server
type Service struct {
	sync.RWMutex
	address string
	db      models.KV
	nodes   nodes
	agent   *agent.MasterAgent
}

func NewService(address string) *Service {
	service := &Service{
		address: address,
		db:      storage.EtcdKV,
	}

	service.agent = agent.NewMasterAgent(service, address)

	urlSlice, err := service.fetchNodes()
	if err != nil {
		panic(err)
	}
	service.nodes = nodeURLSliceToNodes(urlSlice)

	return service
}

func (s *Service) Start() {
	http.ListenAndServe(s.address, s.agent)
}

func (s *Service) AssignNode(userID int64, queueName string, squadName string) (string, error) {
	s.RLock()
	nodeStatsSlice := s.nodes.statsSlice()
	s.RUnlock()

	if nodeStatsSlice.Len() == 0 {
		return "", errors.NewInternalErrorf("service unavailable, empty node cluster")
	}

	sort.Sort(nodeStatsSlice)
	return nodeStatsSlice[0].Node, nil
}

func (s *Service) Join(stats node.Info) {
	s.Lock()
	defer s.Unlock()

	s.nodes[stats.Node] = stats
}

func (s *Service) fetchNodes() ([]string, error) {
	val, err := s.db.Get(nodesKey)
	if err != nil {
		return nil, err
	}

	nodes := []string{}
	if err := json.Unmarshal([]byte(val), nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (s *Service) heartbeat() {
	ticker := time.NewTicker(heartbeatPeriod)

	for {
		<-ticker.C
		s.RLock()
		nodesMap := s.nodes
		s.RUnlock()

		for node := range nodesMap {
			go func(n string) {
				stats, err := s.getNodeStat(n)
				if err != nil {
					log.Internal.Errorf("node[%s] can not connect\n", n)
					s.removeNode(n)
				}

				s.updateNode(stats)
			}(node)
		}
	}
}

func (s *Service) getNodeStat(node string) (node.Info, error) {
	url := fmt.Sprintf(getNodeStatsURLFormat, urlutil.MakeURL(node))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	stats := node.Info{}
	err = json.Unmarshal(respBytes, &stats)
	return stats, err
}

func (s *Service) removeNode(node string) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.nodes[node]; ok {
		delete(s.nodes, node)
	}
}

func (s *Service) updateNode(info node.Info) {
	s.Lock()
	defer s.Unlock()

	if old := s.nodes[info.Node]; old == info {
		return
	}

	s.nodes[info.Node] = info
}
