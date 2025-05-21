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
)

var clientID int

// Conecta ao broker MQTT
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

// Desconecta do broker MQTT
func disconnectMQTT(client mqtt.Client) {
	client.Disconnect(250)
	fmt.Println("🚪 Cliente desconectado.")
}

// Publica uma mensagem em um tópico MQTT
func publish(client mqtt.Client, topic string, payload string) {
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("⚠️ Erro ao publicar no tópico %s: %v", topic, err)
		return
	}
	fmt.Printf("📄 Publicado no tópico %s\n", topic)
}

// Inscreve em um tópico MQTT
func subscribe(client mqtt.Client, topic string, callback mqtt.MessageHandler) {
	token := client.Subscribe(topic, 0, callback)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Fatalf("❌ Erro ao se inscrever no tópico %s: %v", topic, err)
	}
	fmt.Printf("📡 Inscrito no tópico %s\n", topic)
}

// Tratador padrão de mensagens MQTT recebidas
func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("📩 Mensagem recebida [%s]: %s\n", msg.Topic(), msg.Payload())
}

// Um callback para tratar a resposta
func reservationHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("✅ Confirmação de reserva recebida para o carro %d: %s\n", clientID, string(msg.Payload()))
}

// Gera uma posição aleatória simulada
func generatePosition() (int, int, int, int) {

	originX := rand.Intn(1000)      // coordenada X entre 0 e 999
	originY := rand.Intn(1000)      // coordenada Y entre 0 e 999
	destinationX := rand.Intn(1000) // coordenada X entre 0 e 999
	destinationY := rand.Intn(1000) // coordenada Y entre 0 e 999

	return originX, originY, destinationX, destinationY
}

// Loop que simula a movimentação do carro publicando posição periodicamente
func startCarLoop(client mqtt.Client) {
	for {
		origX, origY, destX, destY := generatePosition()
		route := fmt.Sprintf("%d, %d, %d, %d, %d", clientID, origX, origY, destX, destY)

		// Envia para ambos os tópicos
		publish(client, "car/position", route)
		fmt.Println("📤 Enviado para car/position:", route)

		publish(client, "car/recarga", route)
		fmt.Println("📤 Enviado para car/recarga:", route)

		time.Sleep(2 * time.Second)
	}
}

// Função principal
func main() {
	rand.Seed(time.Now().UnixNano())

	// Gera um ID único para o carro
	clientID = rand.Intn(1000)
	id := fmt.Sprintf("%d", clientID)

	client := connectMQTT(brokerURL, id)
	defer disconnectMQTT(client)

	subscribe(client, "car/recarga", defaultMessageHandler)
	reservationTopic := fmt.Sprintf("car/%d/reservation", clientID)
	subscribe(client, reservationTopic, reservationHandler)

	startCarLoop(client)
}
