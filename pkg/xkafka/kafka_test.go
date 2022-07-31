package xkafka_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/comeonjy/go-kit/pkg/xkafka"
)

func TestNew(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*50000)
	cli, err := xkafka.New(ctx, `{"addrs":["kafka.tool:9092"],"client_id":"go-kit"}`)
	if err != nil {
		t.Error(err)
		return
	}
	defer cli.Close()
	topic := "Test1"
	topics := "Test1,Test2"

	t.Run("demo", func(t *testing.T) {
		strings, err := cli.Client.Topics()
		if err != nil {
			return
		}
		log.Println(strings)
	})

	t.Run("producer", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if err := cli.SendSyncMessage(&sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(fmt.Sprintf("%d", i)),
				Value: sarama.StringEncoder(fmt.Sprintf("Helllo-%d", i)),
			}); err != nil {
				t.Error(err)
				return
			}
		}
		//for i := 0; i < 10; i++ {
		//	if err := cli.SendSyncMessage(&sarama.ProducerMessage{
		//		Topic: "Test2",
		//		Key:   sarama.StringEncoder(fmt.Sprintf("%d", i)),
		//		Value: sarama.StringEncoder(fmt.Sprintf("Helllo-%d", i)),
		//	}); err != nil {
		//		t.Error(err)
		//		return
		//	}
		//}
	})
	t.Run("Consumes", func(t *testing.T) {
		// 两个消费者，消费两遍
		if err := cli.Consumes(topic, nil); err != nil {
			t.Error(err)
			return
		}
		if err := cli.Consumes(topic, DefaultHandler); err != nil {
			t.Error(err)
			return
		}
		<-ctx.Done()
	})
	t.Run("Consumes1", func(t *testing.T) {
		if err := cli.Consumes(topic, nil); err != nil {
			t.Error(err)
			return
		}
		<-ctx.Done()
	})
	t.Run("ConsumeGroup", func(t *testing.T) {
		if err := cli.ConsumerGroup("group_1", topics, nil); err != nil {
			t.Error(err)
			return
		}
		<-ctx.Done()
	})
	t.Run("ConsumeGroup", func(t *testing.T) {
		if err := cli.ConsumerGroup("group_1", "Test1", nil); err != nil {
			t.Error(err)
			return
		}
		<-ctx.Done()
	})

}

func DefaultHandler(msg *sarama.ConsumerMessage) {
	log.Printf("Partition:%d Topic:%s Key:%s Value:%s Offset:%d \n", msg.Partition, msg.Topic, msg.Key, msg.Value, msg.Offset)
}
