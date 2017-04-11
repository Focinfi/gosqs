package master

import "github.com/hashicorp/consul/api"

// Service for a master server
type Service struct {
	address string
	api.Agent
}

func (s *Service) AssignNode(userID int64, queueName string, squadName string) (string, error) {
	_, err := s.Members(true)
	if err != nil {
		return "", err
	}

	return "", nil
}
