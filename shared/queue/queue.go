package queue

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/adjust/redismq"
	"github.com/gauravbansal74/solution/config"
	"github.com/gauravbansal74/solution/distance"
	"github.com/gauravbansal74/solution/shared/logger"
)

var (
	con *QueueConfig
)

type QueueConfig struct {
	Queue    *redismq.Queue
	Consumer *redismq.Consumer
}

type RMessage struct {
}

func Init(conf config.Config) *QueueConfig {
	con, err := createQueue(conf)
	if err != nil {
		logger.Fatal("server", err, "Redis Queue init failed")
		return con
	}
	return con
}

func ReadQueue() *QueueConfig {
	return con
}

func createQueue(conf config.Config) (*QueueConfig, error) {
	que := redismq.CreateQueue(conf.RedisHost, strconv.Itoa(conf.RedisPort), "", 9, conf.RedisQueue)
	con = &QueueConfig{
		Queue: que,
	}
	consumer, err := addConsumer(con)
	if err != nil {
		logger.Fatal("server", err, "Redis Create Queue failed")
		return con, err
	}
	con.Consumer = consumer
	return con, nil
}

func addConsumer(c *QueueConfig) (*redismq.Consumer, error) {
	consumer, err := c.Queue.AddConsumer(string(time.Now().UnixNano()))
	if err != nil {
		logger.Fatal("server", err, "Redis Add Consumer failed")
		return consumer, err
	}
	return consumer, nil
}

func (c *QueueConfig) PutPayload(payload string) error {
	err := c.Queue.Put(payload)
	if err != nil {
		logger.Error("server", err, "Queue Error while put payload", nil)
	}
	return err
}

func (c *QueueConfig) GetPayload() {
	for {
		payloadObject, err := c.Consumer.Get()
		if err != nil {
			logger.Error("server", err, "Queue Error while reading payload", nil)
		}
		logger.Info("server", payloadObject.Payload, nil)
		err = payloadObject.Ack()
		if err != nil {
			payloadObject.Fail()
			logger.Error("server", err, "Queue Error while ack payload", nil)
		}
		if payloadObject.Payload != "" {
			output, err := distance.CheckForShortestDistance(payloadObject.Payload)
			if err != nil {
				logger.Error("redis-getpayload", err, "Error while process distance payload", logger.Fields{
					"Id":   output.ID,
					"Path": output.Path,
				})
				payloadObject.Fail()
			}
			if output.Status == 0 {
				payloadObject.Fail()
			}
			bytesse, err := json.Marshal(output)
			if err != nil {
				logger.Error("redis-getpayload", err, "Error while process distance payload", logger.Fields{
					"Id":   output.ID,
					"Path": output.Path,
				})
			}
			redisapi := ReadRedisClient()
			err = redisapi.Client.Set(output.ID, bytesse, 0).Err()
			if err != nil {
				logger.Error("redis-getpayload", err, "Error while process distance payload", logger.Fields{
					"Id":   output.ID,
					"Path": output.Path,
				})
			}
		}
	}
}
