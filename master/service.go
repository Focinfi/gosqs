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
	"github.com/Focinfi/sqs/util/fmtutil"
	"github.com/Focinfi/sqs/util/urlutil"
)

const (
	getNodeStatsURLFormat = "%s/stats"
	logPrefix             = "[sqs.master]"
)

var (
	heartbeatPeriod = time.Second
	format          = fmtutil.NewFormat("sqs.master")
)

// Service for a master server
type Service struct {
	sync.RWMutex
	address string
	db      *database
	nodes   nodes
	agent   *agent.MasterAgent
}

// NewService allocates a new Service
func NewService(address string) *Service {
	service := &Service{
		address: address,
		db:      db,
	}

	service.agent = agent.NewMasterAgent(service, address)

	urlSlice, err := service.db.fetchNodes()
	if err != nil {
		panic(err)
	}
	service.nodes = nodeURLSliceToNodes(urlSlice)
	log.Internal.Infoln(format.Sprintln("init nodes:", service.nodes))

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
	log.Internal.Infoln(format.Sprintln("ready to be assigned nodes:", s.nodes))
	return nodeStatsSlice[0].Addr, nil
}

// AddNode joins a node to the ready-to-serve nodes list
func (s *Service) AddNode(info models.NodeInfo) {
	s.Lock()
	defer s.Unlock()

	log.Biz.Infoln(format.Sprintln("node to be joined:", info.Addr))
	s.nodes[info.Addr] = info
	if err := s.db.updateNodes(s.nodes.nodeURLSlice()); err != nil {
		log.DB.Errorln(format.Sprintln("failed to update nodes into db"))
	}
}

func (s *Service) heartbeat() {
	ticker := time.NewTicker(heartbeatPeriod)

	for {
		<-ticker.C
		s.RLock()
		nodesMap := s.nodes
		s.RUnlock()

		log.Cluster.Infof(format.Sprintf("get stats from [%d] nodes", len(nodesMap)))
		for node := range nodesMap {
			go func(n string) {
				stats, err := s.getNodeStat(n)
				if err != nil {
					log.Cluster.Errorln(format.Sprintf("node[%s] can not be connected, err: %v\n", n, err))
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
		log.Cluster.Infoln(format.Sprintf("%s node[%s] removed\n", logPrefix, node))
		// TODO: alert for handle failed node
		s.db.updateNodes(s.nodes.nodeURLSlice())
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
