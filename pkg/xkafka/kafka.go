package xkafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

type Kafka struct {
	ctx           context.Context
	Client        sarama.Client
	Admin         sarama.ClusterAdmin
	SyncProducer  sarama.SyncProducer
	AsyncProducer sarama.AsyncProducer
	OffsetManager sarama.OffsetManager
	GroupMap      sync.Map
}

type _config struct {
	Addrs    []string `json:"addrs"`
	ClientID string   `json:"client_id"`
}

func New(ctx context.Context, confStr string) (*Kafka, error) {
	conf := _config{}
	if err := json.Unmarshal([]byte(confStr), &conf); err != nil {
		return nil, err
	}
	configs := sarama.NewConfig()
	configs.ClientID = conf.ClientID
	configs.Producer.Return.Successes = true
	client, err := sarama.NewClient(conf.Addrs, configs)
	if err != nil {
		return nil, err
	}
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return nil, err
	}
	syncProducer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	asyncProducer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &Kafka{
		ctx:           ctx,
		Client:        client,
		Admin:         admin,
		SyncProducer:  syncProducer,
		AsyncProducer: asyncProducer,
	}, nil
}

// Close 关闭连接
func (k *Kafka) Close() {
	k.SyncProducer.Close()
	k.AsyncProducer.Close()
	k.Client.Close()
}

// SendSyncMessage 发送同步消息
func (k *Kafka) SendSyncMessage(message *sarama.ProducerMessage) error {
	partition, i, err := k.SyncProducer.SendMessage(message)
	if err != nil {
		return err
	}
	log.Printf("SendSyncMessage success partition:%d offset:%d %v \n", partition, i, *message)
	return nil
}

// Consumes
// 消费者
// 为topic每个partition启动一个协程处理
// 用于分布式环境时，需注意消息重复消费
func (k *Kafka) Consumes(topic string, handle func(*sarama.ConsumerMessage)) error {
	if handle == nil {
		handle = defaultHandler
	}
	consumer, err := sarama.NewConsumerFromClient(k.Client)
	if err != nil {
		return err
	}
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}
	for _, v := range partitions {
		pc, err := consumer.ConsumePartition(topic, v, sarama.OffsetNewest)
		if err != nil {
			return err
		}
		go func(partitionConsumer sarama.PartitionConsumer, p int32) {
			defer partitionConsumer.Close()
			log.Printf("Partition:%d Topic:%s Start... \n", p, topic)
			for msg := range partitionConsumer.Messages() {
				handle(msg)
			}
			log.Printf("Partition:%d Topic:%s Exit... \n", p, topic)
		}(pc, v)
	}
	return nil
}

func defaultHandler(msg *sarama.ConsumerMessage) {
	log.Printf("Partition:%d Topic:%s Key:%s Value:%s Offset:%d \n", msg.Partition, msg.Topic, msg.Key, msg.Value, msg.Offset)
}

// ConsumerGroup
// 消费者组
// 分布式环境下、需注意分区数量与服务数量的关系
func (k *Kafka) ConsumerGroup(groupID string, topics string, handler sarama.ConsumerGroupHandler) error {
	if handler == nil {
		handler = &DefaultConsumer{}
	}
	client, err := sarama.NewConsumerGroupFromClient(groupID, k.Client)
	if err != nil {
		return err
	}
	go func() {
		for {
			if err := client.Consume(k.ctx, strings.Split(topics, ","), handler); err != nil {
				log.Println("ConsumeGroup err:", err.Error())
				return
			}
			log.Println("Rebalanced ...")
		}
	}()
	log.Printf("ConsumerGroup GroupID:%s Topics:%s Start... \n", groupID, topics)
	return nil
}

type DefaultConsumer struct{}

func (consumer *DefaultConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *DefaultConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *DefaultConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for m := range claim.Messages() {
		log.Printf("Message claimed: partition = %d, key = %s, value = %s, timestamp = %v, topic = %s", m.Partition, string(m.Key), string(m.Value), m.Timestamp, m.Topic)
		session.MarkMessage(m, "")
	}
	return nil
}
