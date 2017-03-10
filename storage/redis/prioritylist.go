package redis

import (
	"fmt"
	"sync"

	"encoding/json"

	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/errors"
	"github.com/Focinfi/sqs/log"
	"github.com/Focinfi/sqs/models"
	"gopkg.in/redis.v5"
)

var mux sync.RWMutex
var index = 0

// PriorityList for priority consumers useing redis
type PriorityList struct {
	key    string
	locked map[string]bool
	db     *redis.Client
}

// Push pushes the c into the pl
func (pl *PriorityList) Push(c models.Consumer) error {
	if err := pl.ZAdd(c); err != nil {
		return err
	}

	// release lock
	mux.Lock()
	defer mux.Unlock()
	key := c.Client().Key()
	if _, ok := pl.locked[key]; ok {
		delete(pl.locked, key)
	}

	return nil
}

// Pop returns the consumer with the highest Score
// if locked it will try the next-highest-score consumer
func (pl *PriorityList) Pop() (models.Consumer, error) {
	//  get the highest-score consumer
	log.Biz.Debugln("POP-redis")
	heighest, err := pl.ZHeighest()
	if err != nil {
		return nil, err
	}

	// lock the map
	mux.RLock()
	// check if locked
	if _, ok := pl.locked[heighest.Client().Key()]; ok {
		// unlock read
		mux.RUnlock()
		return nil, errors.NoConsumer
	}
	// unlock read
	mux.RUnlock()

	// get the lowest-score consumer
	lowest, err := pl.ZLowest()
	if err != nil {
		return nil, err
	}

	// set the heighest score smaller than the lowest to simulate the removing
	heighest.IncPriority(-1 * (heighest.Priority() - lowest.Priority()))
	if err := pl.ZAdd(heighest); err != nil {
		return nil, err
	}

	// lock the heighest
	mux.Lock()
	defer mux.Unlock()
	pl.locked[heighest.Client().Key()] = true
	return heighest, nil
}

// ZAdd add c using redis.ZAdd
func (pl *PriorityList) ZAdd(c models.Consumer) error {
	b, err := json.Marshal(c)
	if err != nil {
		return errors.NewInternalErr(err.Error())
	}

	res := pl.db.ZAdd(pl.key, redis.Z{Member: string(b), Score: float64(c.Priority())})
	if err := res.Err(); err != nil {
		log.DB.Error(err)
		return err
	}

	return nil
}

// ZHeighest returns the heighest-score consmuer
func (pl *PriorityList) ZHeighest() (models.Consumer, error) {
	return pl.zTop(-1)
}

// ZLowest returns the lowest-score consumer
func (pl *PriorityList) ZLowest() (models.Consumer, error) {
	return pl.zTop(-1)
}

// ZTop returns the highest-score(top=-1) or lowest-score(top=0)
func (pl *PriorityList) zTop(top int64) (models.Consumer, error) {
	res := pl.db.ZRange(pl.key, top, top)
	if err := res.Err(); err != nil {
		return nil, err
	}

	if len(res.Val()) == 0 {
		return nil, errors.NoConsumer
	}

	consumer := &Consumer{}
	err := json.Unmarshal([]byte(res.Val()[0]), consumer)
	if err != nil {
		return nil, errors.DataBroken(pl.key, err)
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

	mux.Lock()
	defer mux.Unlock()

	key := fmt.Sprintf("sqs.pl.%d", index+1)
	return &PriorityList{db: client, key: key, locked: map[string]bool{}}, nil
}
