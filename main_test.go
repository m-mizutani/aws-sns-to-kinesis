package main_test

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/uuid"
	gp "github.com/m-mizutani/generalprobe"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	Region        string `json:"region"`
	StackName     string `json:"stack_name"`
	SnsTopic      string
	KinesisStream string
}

func loadTestConfig() testConfig {
	stackConfig := os.Getenv("STACK_CONFIG")
	if stackConfig == "" {
		log.Fatal("Environment Variable 'STACK_CONFIG' is required fro test")
	}

	fd, err := os.Open(os.Getenv("STACK_CONFIG"))
	if err != nil {
		log.Fatal("Can not open test.json", err)
	}
	defer fd.Close()

	cfg := testConfig{}

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		sep := strings.Split(line, "=")
		if len(sep) != 2 {
			continue
		}

		switch sep[0] {
		case "StackName":
			cfg.StackName = sep[1]
		case "SnsTopicArn":
			cfg.SnsTopic = sep[1]
		case "KinesisStreamArn":
			cfg.KinesisStream = sep[1]
		}
	}

	out, err := exec.Command("aws", "configure", "get", "region").Output()
	if err != nil {
		log.Fatal("Fail to get default region by aws command", err)
	}

	cfg.Region = strings.TrimSpace(string(out))
	return cfg
}

func Test(t *testing.T) {
	cfg := loadTestConfig()
	log.WithField("config", cfg).Info("Start")
	testID := uuid.New().String()

	type testData struct {
		ID string `json:"id"`
	}

	g := gp.New(cfg.Region, cfg.StackName)
	g.AddScenes([]gp.Scene{
		// Send request
		gp.PublishSnsMessage(g.Arn(cfg.SnsTopic), []byte(`{"id":"`+testID+`"}`)),

		// Recv result
		gp.GetKinesisStreamRecord(g.Arn(cfg.KinesisStream), func(data []byte) bool {
			var td testData
			log.WithField("received", string(data)).Info("Recv Kinesis record")
			err := json.Unmarshal(data, &td)
			require.NoError(t, err)
			assert.Equal(t, testID, td.ID)
			return true
		}),
	})

	g.Act()
	log.Info("Exit")
}
