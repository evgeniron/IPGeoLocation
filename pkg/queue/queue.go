package queue

import (
	"github.com/go-redis/redis/v8"
)

type Queue struct {
	redisClient  *redis.Client
}

/* Create new task queue */
func NewQueue() (*Queue, error) {
	opts, err = redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return err, nil
	}
	redisClient := redis.NewClient(opts)

	return &Queue{redis.Client: redisClient}, nil
}

/* Pop new task from queue, block if empty */
func (q *Queue) PopTask() (string ,error) {
	task, err := q.redisClient.BLpop(0, "ipQuery").Result()
	if err != nil {
		return "", nil
	}

	return task[1], nil
}

/* Push new task into queue*/
func (q *Queue) PushTask (query string) error {
	_, err := q.redisClient.RPush("ipQuery", query).Result()
	return err
}
