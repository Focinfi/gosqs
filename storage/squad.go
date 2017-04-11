package storage

import (
	"encoding/json"

	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/models"
)

// Squad is for squad storage
type Squad struct {
	store *Storage
	db    models.KV
}

// One gets a squad
func (s *Squad) One(userID int64, queueName string, squadName string) (*models.Squad, error) {
	key := models.SquadKey(userID, queueName, squadName)
	val, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}

	squad := &models.Squad{}
	if err := json.Unmarshal([]byte(val), squad); err != nil {
		return nil, errors.NewInternalWrap(err)
	}

	return squad, nil
}

// Add adds a squad
func (s *Squad) Add(squad models.Squad) error {
	_, err := s.One(squad.UserID, squad.QueueName, squad.Name)
	if err == nil {
		return errors.DuplicateSquad
	}

	if err != errors.DataNotFound {
		return err
	}

	squadBytes, err := json.Marshal(squad)
	if err != nil {
		return errors.NewInternalWrap(err)
	}

	return s.db.Put(squad.Key(), string(squadBytes))
}

// Update updates a squad
func (s *Squad) Update(squad *models.Squad) error {
	old, err := s.One(squad.UserID, squad.QueueName, squad.Name)
	if err != nil {
		return err
	}

	if squad.ReceivedMessageID <= old.ReceivedMessageID {
		return nil
	}

	squadBytes, err := json.Marshal(squad)
	if err != nil {
		return errors.NewInternalWrap(err)
	}

	s.db.Put(squad.Key(), string(squadBytes))
	return nil
}
