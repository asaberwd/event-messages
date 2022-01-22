package event

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"context"
	s3lib "github.com/LF-Engineering/insights-datasource-shared/aws/s3"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Provider used in connecting to s3
type S3Provider interface {
	Save(payload []byte) error
	GetKeys() ([]string, error)
	Get(key string) ([]byte, error)
	Delete(key string) error
}

// Manager ...
type Manager struct {
}

// NewManager ...
func NewManager() *Manager {
	return &Manager{}
}

// Read ...
func (m *Manager) Read(topic string) ([]string, error) {
	messages := make([]string, 0)
	stage := os.Getenv("STAGE")
	bucketName := fmt.Sprintf("event-messages-datasource-%s-%s", strings.ToLower(topic), stage)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return messages, err
	}
	client := s3.NewFromConfig(cfg)
	_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{Bucket: &bucketName})
	if err != nil {
		return messages, fmt.Errorf("topic: %s does not exist", topic)
	}

	s3Provider := s3lib.NewManager(bucketName, os.Getenv("REGION"))
	keys, err := s3Provider.GetKeys()
	if err != nil {
		return messages, err
	}
	for _, key := range keys {
		result, err := s3Provider.Get(key)
		if err != nil {
			return messages, err
		}

		messages = append(messages, string(result))
		err = s3Provider.Delete(key)
		if err != nil {
			return messages, err
		}
	}
	return messages, nil
}

// Write ...
func (m *Manager) Write(data string, topic string) error {
	v, err := json.Marshal(data)
	if err != nil {
		return err
	}
	stage := os.Getenv("STAGE")
	bucketName := fmt.Sprintf("event-messages-datasource-%s-%s", strings.ToLower(topic), stage)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)
	_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{Bucket: &bucketName})
	if err != nil {
		_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{Bucket: &bucketName})
		if err != nil {
			return err
		}
	}
	s3Provider := s3lib.NewManager(bucketName, os.Getenv("REGION"))

	err = s3Provider.Save(v)
	if err != nil {
		return err
	}
	return nil
}
