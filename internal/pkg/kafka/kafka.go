package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type Client struct {
	conn sarama.Client
}

func New(kafkaAddr []string) (*Client, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	client, err := sarama.NewClient(kafkaAddr, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create Kafka client: %w", err)
	}

	_, err = client.Topics()
	if err != nil {
		return nil, fmt.Errorf("unable to get topics: %w", err)
	}

	return &Client{conn: client}, nil
}

func (c *Client) NewProducer() (sarama.SyncProducer, error) {
	return sarama.NewSyncProducerFromClient(c.conn)
}

func (c *Client) NewConsumerGroup(groupID string) (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroupFromClient(groupID, c.conn)
}

func (c *Client) Close() {
	c.conn.Close()
}
