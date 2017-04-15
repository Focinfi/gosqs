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
	"github.com/Focinfi/sqs/storage"
	"github.com/Focinfi/sqs/util/urlutil"
)

const (
	nodesKey              = "sqs.nodes"
	getNodeStatsURLFormat = "%s/stats"
	logPrefix             = "[sqs.master]"
)

var (
	heartbeatPeriod = time.Second
)

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

func (m nodes) statsSlice() models.InfoSlice {
	slice := make([]models.NodeInfo, len(m))
	i := 0
	for node := range m {
		slice[i] = m[node]
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

// Service for a master server
type Service struct {
	sync.RWMutex
	address string
	db      models.KV
	nodes   nodes
	agent   *agent.MasterAgent
}

// NewService allocates a new Service
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
	log.DB.Infoln(logPrefix, "init nodes:", service.nodes)

	return service
}

// Start starts the service
func (s *Service) Start() {
	go s.heartbeat()
	http.ListenAndServe(s.address, s.agent)
}

// AssignNode assigns a node to serve one client
func (s *Service) AssignNode(userID int64, queueName string, squadName string) (string, error) {
	s.RLock()
	nodeStatsSlice := s.nodes.statsSlice()
	s.RUnlock()

	if nodeStatsSlice.Len() == 0 {
		return "", errors.NewInternalErrorf("service unavailable, empty node cluster")
	}

	sort.Sort(nodeStatsSlice)
	log.DB.Infoln(logPrefix, "nodes:", s.nodes)
	return nodeStatsSlice[0].Addr, nil
}

// Join joins a node to the ready-to-serve nodes list
func (s *Service) Join(info models.NodeInfo) {
	s.Lock()
	defer s.Unlock()

	log.Biz.Infoln(logPrefix, "to join:", info.Addr)
	s.nodes[info.Addr] = info
	if err := s.updateNodes(s.nodes.nodeURLSlice()); err != nil {
		log.DB.Errorln(logPrefix, "failed to update nodes into db")
	}
}

func (s *Service) fetchNodes() ([]string, error) {
	val, err := s.db.Get(nodesKey)
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

func (s *Service) updateNodes(nodes []string) error {
	nodesBytes, err := json.Marshal(nodes)
	if err != nil {
		return err
	}

	return s.db.Put(nodesKey, string(nodesBytes))
}

func (s *Service) heartbeat() {
	ticker := time.NewTicker(heartbeatPeriod)

	for {
		<-ticker.C
		s.RLock()
		nodesMap := s.nodes
		s.RUnlock()

		log.DB.Infof("%s get stats from %d nodes", logPrefix, len(nodesMap))
		for node := range nodesMap {
			go func(n string) {
				stats, err := s.getNodeStat(n)
				if err != nil {
					log.Internal.Errorf("node[%s] can not be connected\n", n)
					s.removeNode(n)
					return
				}

				s.updateNode(*stats)
			}(node)
		}
	}
}

func (s *Service) getNodeStat(node string) (*models.NodeInfo, error) {
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

	stats := &struct {
		models.HTTPStatusMeta
		Data models.NodeInfo
	}{}
	err = json.Unmarshal(respBytes, stats)
	return &stats.Data, err
}

func (s *Service) removeNode(node string) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.nodes[node]; ok {
		delete(s.nodes, node)
		log.DB.Errorf("%s node[%s] removed\n", logPrefix, node)
		// TODO: alert for handle failed node
	}
}

func (s *Service) updateNode(info models.NodeInfo) {
	s.Lock()
	defer s.Unlock()

	if old := s.nodes[info.Addr]; old == info {
		return
	}

	s.nodes[info.Addr] = info
}
