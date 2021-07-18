package kafka_impl

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	pdf_creator "github.com/ag-students/support-service/pkg/pdf-creator"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

type DataForDocument struct {
	surname    string
	name       string
	patronymic string
}

func GetKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func handleMessage(message kafka.Message) error {
	if string(message.Key) == "DataForDocument" {
		data := DataForDocument{}
		buf := &bytes.Buffer{}
		buf.Write(message.Value)
		err := binary.Read(buf, binary.BigEndian, &data)
		if err != nil {
			log.Fatalf("Failed to convert message '%s' to struct", message.Key)
			return err
		} else {
			log.Printf("Key: %s\nValue: %s\n", message.Key, data)
			pdf_creator.CreatePDF(data.surname, data.name, data.patronymic)
			return nil
		}
	}
	return nil
}

// Listen TODO: добавить обработчик пришедшего запроса
func Listen(ctx context.Context, reader *kafka.Reader) {
	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		err = handleMessage(m)
		if err != nil {
			log.Fatalf("Probably, value of message '%s' isn't correct", m.Key)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
