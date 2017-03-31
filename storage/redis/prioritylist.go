package redis

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"gopkg.in/redis.v5"
)

var index = 0

// PriorityList for priority consumers using redis
type PriorityList struct {
	plKey         string
	setKey        string
	pushedChanMap map[int64]chan bool
	db            *redis.Client
}

// Add adds c into pl
func (pl *PriorityList) Add(c models.Consumer) error {
	key := c.Client().Key()
	pushedRes := pl.db.SIsMember(pl.setKey, key)

	if err := pushedRes.Err(); err != nil {
		return nil
	}

	log.Biz.Infof("ADD: key: %s, member: %s, %#v\n", pl.setKey, key, pushedRes.Val())

	if pushedRes.Val() {
		return nil
	}

	// SAdd add c into set
	saddRes := pl.db.SAdd(pl.setKey, c.Client().Key())
	if err := saddRes.Err(); err != nil {
		return err
	}

	return pl.Push(c)
}

// Push pushes the c into the pl
func (pl *PriorityList) Push(c models.Consumer) error {
	// ZAdd add c into pl
	if err := pl.ZAdd(c); err != nil {
		return err
	}

	// broadcast push notification
	go pl.broadcastPushEvent()

	return nil
}

func (pl *PriorityList) broadcastPushEvent() {
	for k := range pl.pushedChanMap {
		go func(i int64) {
			pl.pushedChanMap[i] <- true
		}(k)
	}
}

// Pop returns the consumer with the highest Score
func (pl *PriorityList) Pop() (models.Consumer, error) {
	now := time.Now().UnixNano()
	pushedChan := make(chan bool)
	pl.pushedChanMap[now] = pushedChan
	defer delete(pl.pushedChanMap, now)

	for {
		if pl.db.ZCard(pl.plKey).Val() == 0 {
			log.Biz.Infoln("TRY POP")
			select {
			case <-time.After(time.Second):
			case <-pushedChan:
			}
		}

		consumer, err := pl.pop()
		// pl is empty
		if err == errors.NoConsumer {
			continue
		}

		// db error
		if err != nil {
			return nil, err
		}

		// poped
		return consumer, nil
	}
}

func (pl *PriorityList) pop() (models.Consumer, error) {
	consumer, err := pl.ZHighest()
	if err != nil {
		return nil, err
	}

	key, err := json.Marshal(consumer)
	if err != nil {
		return nil, errors.NewInternalErrorf(err.Error())
	}

	if res := pl.db.ZRem(pl.plKey, string(key)); res.Val() > 0 {
		return consumer, nil
	}
	return nil, nil
}

// Remove removes the c
func (pl *PriorityList) Remove(c models.Consumer) error {
	return pl.db.SRem(pl.setKey, c.Client().Key()).Err()
}

// ZAdd add c using redis.ZAdd
func (pl *PriorityList) ZAdd(c models.Consumer) error {
	b, err := json.Marshal(c)
	if err != nil {
		return errors.NewInternalErrorf(err.Error())
	}

	res := pl.db.ZAdd(pl.plKey, redis.Z{Member: string(b), Score: float64(c.Priority())})
	if err := res.Err(); err != nil {
		log.DB.Error(err)
		return err
	}

	return nil
}

// ZHighest returns the highest-score consumer
func (pl *PriorityList) ZHighest() (models.Consumer, error) {
	return pl.zTop(-1)
}

// ZTop returns the highest-score(top=-1) or lowest-score(top=0)
func (pl *PriorityList) zTop(top int64) (models.Consumer, error) {
	res := pl.db.ZRange(pl.plKey, top, top)
	if err := res.Err(); err != nil {
		return nil, err
	}

	if len(res.Val()) == 0 {
		return nil, errors.NoConsumer
	}

	consumer := &Consumer{}
	err := json.Unmarshal([]byte(res.Val()[0]), consumer)
	if err != nil {
		return nil, errors.DataBroken(pl.plKey, err)
	}

	return consumer, nil
}

// New returns a new PriorityList
func New() (*PriorityList, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Config().RedisAdrr,
		Password: config.Config().RedisPwd,
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil || pong != "PONG" {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Int() + index + 1

	pl := &PriorityList{
		db:            client,
		plKey:         fmt.Sprintf("sqs.pl.%d", id),
		setKey:        fmt.Sprintf("sqs.set.%d", id),
		pushedChanMap: make(map[int64]chan bool),
	}

	log.Biz.Infof("%#v\n", pl)
	return pl, nil
}
