package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/sofyan48/svc_order/src/app/v1/entity"
	"github.com/sofyan48/svc_order/src/app/v1/event"
	"github.com/sofyan48/svc_order/src/utils/kafka"
	"github.com/sofyan48/svc_order/src/utils/logger"
)

// V1OrderEvents ...
type V1OrderEvents struct {
	Kafka  kafka.KafkaLibraryInterface
	Event  event.UserEventInterface
	Logger logger.LoggerInterface
}

// V1OrderEventsHandler ...
func V1OrderEventsHandler() *V1OrderEvents {
	return &V1OrderEvents{
		Kafka:  kafka.KafkaLibraryHandler(),
		Event:  event.OrderEventHandler(),
		Logger: logger.LoggerHandler(),
	}
}

// V1OrderEventsInterface ...
type V1OrderEventsInterface interface {
	Consume(topics []string, signals chan os.Signal)
}

// Consume ...
func (consumer *V1OrderEvents) Consume(topics []string, signals chan os.Signal) {
	// StateFullData := consumer.Kafka.GetStateFull()
	chanMessage := make(chan *sarama.ConsumerMessage, 256)
	csm, err := consumer.Kafka.InitConsumer()
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}
	for _, topic := range topics {
		partitionList, err := csm.Partitions(topic)
		if err != nil {
			log.Println("Unable to get partition got error ", err)
			continue
		}
		for _, partition := range partitionList {
			go consumeMessage(csm, topic, partition, chanMessage)
		}
	}
	log.Println("Event is Started....")

ConsumerLoop:
	for {
		select {
		case msg := <-chanMessage:
			eventData := &entity.StateFullFormatKafka{}
			json.Unmarshal(msg.Value, eventData)
			switch eventData.Action {
			case "users":
				consumer.userLoad(eventData)
			default:
				fmt.Println("OK")
			}

		case sig := <-signals:
			if sig == os.Interrupt {
				break ConsumerLoop
			}
		}
	}
}

func consumeMessage(consumer sarama.Consumer, topic string, partition int32, c chan *sarama.ConsumerMessage) {
	msg, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Println("Unable to consume partition got error ", partition, err)
		return
	}
	defer func() {
		if err := msg.Close(); err != nil {
			log.Println("Unable to close partition : ", partition, err)
		}
	}()
	for {
		msg := <-msg.Messages()
		c <- msg
	}

}

func (consumer *V1OrderEvents) userLoad(dataUser *entity.StateFullFormatKafka) {
	result, err := consumer.Event.InsertDatabase(dataUser)
	if err != nil {
		loggerData := map[string]interface{}{
			"code":  "400",
			"error": err,
		}
		consumer.Logger.Save(dataUser.UUID, "failed", loggerData)
	}
	loggerData := map[string]interface{}{
		"code":   "200",
		"result": result,
	}
	data, err := consumer.Logger.Save(dataUser.UUID, "success", loggerData)
	fmt.Println(data, err)
}
