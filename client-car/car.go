package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	brokerURL = "tcp://mosquitto:1883"
	clientID  = "carroA"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	client := connectMQTT(brokerURL, clientID)
	defer disconnectMQTT(client)

	subscribe(client, "car/recarga", defaultMessageHandler)
	startCarLoop(client)
}

func connectMQTT(broker, clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetDefaultPublishHandler(defaultMessageHandler)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Erro ao conectar ao broker: %v", token.Error())
	}

	fmt.Println("✅ Conectado ao broker MQTT!")
	return client
}

func disconnectMQTT(client mqtt.Client) {
	client.Disconnect(250)
	fmt.Println("🚪 Cliente desconectado.")
}

func publish(client mqtt.Client, topic string, payload string) {
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("⚠️ Erro ao publicar no tópico %s: %v", topic, err)
	}
	fmt.Printf("📄 Publicado no tópico %s\n", topic)
}

func subscribe(client mqtt.Client, topic string, callback mqtt.MessageHandler) {
	token := client.Subscribe(topic, 0, callback)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Fatalf("❌ Erro ao se inscrever no tópico %s: %v", topic, err)
	}
	fmt.Printf("📡 Inscrito no tópico %s\n", topic)
}

func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("📩 Mensagem recebida [%s]: %s\n", msg.Topic(), msg.Payload())
}

func generatePosition() (int, int) {
	return rand.Intn(1000), rand.Intn(1000)
}

func startCarLoop(client mqtt.Client) {
	for {
		x, y := generatePosition()
		position := fmt.Sprintf("Carro A - posição x:%d, y:%d", x, y)

		// Envia para ambos os tópicos
		publish(client, "car/position", position)
		fmt.Println("📤 Enviado para car/position:", position)

		publish(client, "car/recarga", position)
		fmt.Println("📤 Enviado para car/recarga:", position)

		time.Sleep(2 * time.Second)
	}
}
