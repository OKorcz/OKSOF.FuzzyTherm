package mqtt

import (
	"fmt"
	"os"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

const (
	defaultBrokerURL = "tcp://localhost:1883"
	connectTimeout   = 10 * time.Second
)

// ConnectMQTT connects to the configured MQTT broker and returns the live client.
// Set MQTT_BROKER, MQTT_CLIENT_ID, MQTT_USERNAME, and MQTT_PASSWORD as needed.
func ConnectMQTT() (paho.Client, error) {
	options := paho.NewClientOptions().
		AddBroker(envOrDefault("MQTT_BROKER", defaultBrokerURL)).
		SetClientID(envOrDefault("MQTT_CLIENT_ID", "fuzzy-therm"))

	if username := os.Getenv("MQTT_USERNAME"); username != "" {
		options.SetUsername(username)
	}
	if password := os.Getenv("MQTT_PASSWORD"); password != "" {
		options.SetPassword(password)
	}

	client := paho.NewClient(options)
	token := client.Connect()
	if !token.WaitTimeout(connectTimeout) {
		return nil, fmt.Errorf("MQTT connection timed out after %s", connectTimeout)
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("connect to MQTT broker: %w", err)
	}

	return client, nil
}

func envOrDefault(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}
