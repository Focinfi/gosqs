package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Focinfi/sqs/agent"
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"github.com/Focinfi/sqs/util/fmtutil"
	"github.com/Focinfi/sqs/util/psutil"
	"github.com/Focinfi/sqs/util/urlutil"
)

const (
	joinURLFormat  = "%s/join"
	jsonHTTPHeader = "application/json"
)

var (
	defaultPullMessageCount = config.Config.PullMessageCount
	format                  = fmtutil.NewFormat("sqs.node")
)

// Service for one user info
type Service struct {
	addr       string
	masterAddr string
	*database
	agent *agent.QueueAgent
	info  *models.NodeInfo
}

// New allocates a new Service
func New(addr string, masterAddr string) *Service {
	service := &Service{
		addr:       addr,
		database:   db,
		masterAddr: masterAddr,
		info:       &models.NodeInfo{Addr: addr},
	}

	service.agent = agent.NewQueueAgent(service, addr)

	return service
}

// Start starts services
func (s *Service) Start() {
	if err := s.join(); err != nil {
		panic(err)
	}

	go s.updateInfo()
	log.Biz.Fatal(http.ListenAndServe(s.addr, s.agent))
}

// PullMessages pulls message from database, create it if the squad is not found
func (s *Service) PullMessages(userID int64, queueName, squadName string, length int) ([]models.Message, error) {
	squad, err := s.Squad.One(userID, queueName, squadName)

	log.Biz.Infoln(format.Sprintln("[PullMessages]", squad, err))

	if err == errors.DataNotFound {
		maxMessageID, err := s.database.Storage.Queue.MessageMaxID(userID, queueName)
		if err != nil {
			return nil, err
		}

		squad = &models.Squad{
			Name:              squadName,
			UserID:            userID,
			QueueName:         queueName,
			ReceivedMessageID: maxMessageID,
		}

		if err := s.Storage.Squad.Add(*squad); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	maxMessageID, err := s.database.Queue.MessageMaxID(userID, queueName)
	if err != nil {
		return nil, err
	}

	return s.Message.Nextn(userID, queueName, squad.ReceivedMessageID, maxMessageID, defaultPullMessageCount)
}

// ReportMaxReceivedMessageID reports the max recieved message id to mark forward of the squad process
func (s *Service) ReportMaxReceivedMessageID(userID int64, queueName, squadName string, messageID int64) error {
	squad, err := s.database.Squad.One(userID, queueName, squadName)
	if err != nil {
		return err
	}

	if squad.ReceivedMessageID >= messageID {
		return nil
	}

	squad.ReceivedMessageID = messageID
	return s.Squad.Update(squad)
}

// PushMessage receives message
func (s *Service) PushMessage(userID int64, queueName, content string, index int64) error {
	maxID, err := s.Queue.MessageMaxID(userID, queueName)
	log.Biz.Infoln(format.Sprintln("[PushMessage]", index, maxID))
	// create it if not found
	if err == errors.DataNotFound {
		if err := s.Queue.Add(models.NewQueue(userID, queueName)); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	if index > maxID {
		return errors.MessageIndexOutOfRange
	}

	return s.database.PushMessage(userID, queueName, content, index)
}

// ApplyMessageIDRange tries to apply a range a free message id
func (s *Service) ApplyMessageIDRange(userID int64, queueName string, size int) (maxID int64, err error) {
	if size > config.Config.MaxMessageIDRangeSize {
		return -1, errors.ApplyMessageIDRangeOversize
	}

	return s.Queue.ApplyMessageIDRange(userID, queueName, size)
}

// Info returns the info of current service
func (s *Service) Info() models.NodeInfo {
	return *s.info
}

func (s *Service) updateInfo() {
	ticker := time.NewTicker(time.Millisecond * 500)
	for {
		<-ticker.C

		cpuPercent, err := psutil.CPUUsedPercent()
		if err != nil {
			log.Biz.Error(format.Sprintln("failed to get the CPU used percent"))
		} else {
			s.info.CPU = cpuPercent
		}

		s.info.Memory = psutil.MemoryUsedPercent()
	}
}

func (s *Service) join() error {
	infoBytes, err := json.Marshal(s.info)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(joinURLFormat, urlutil.MakeURL(s.masterAddr))
	resp, err := http.Post(url, jsonHTTPHeader, bytes.NewReader(infoBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(format.Sprintln("can not join into master"))
	}
	return nil
}
