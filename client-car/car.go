package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	brokerURL = "tcp://mosquitto:1883" // URL do broker MQTT
)

var clientID int // ID único do cliente (carro), gerado aleatoriamente

/**
 * Conecta ao broker MQTT utilizando o endereço e o ID fornecidos.
 *
 * @param broker  o endereço do broker MQTT
 * @param clientID o ID do cliente
 * @return mqtt.Client o cliente MQTT conectado
 */
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

/**
 * Desconecta o cliente MQTT do broker.
 *
 * @param client o cliente MQTT a ser desconectado
 */
func disconnectMQTT(client mqtt.Client) {
	client.Disconnect(250)
	fmt.Println("🚪 Cliente desconectado.")
}

/**
 * Publica uma mensagem em um tópico MQTT com QoS 0.
 *
 * @param client  o cliente MQTT
 * @param topic   o tópico para onde publicar
 * @param payload o conteúdo da mensagem
 */
func publish(client mqtt.Client, topic string, payload string) {
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("⚠️ Erro ao publicar no tópico %s: %v", topic, err)
		return
	}
	fmt.Printf("📄 Publicado no tópico %s\n", topic)
}

/**
 * Inscreve o cliente em um tópico MQTT, associando um handler de mensagens.
 *
 * @param client   o cliente MQTT
 * @param topic    o tópico a se inscrever
 * @param callback o handler para mensagens recebidas
 */
func subscribe(client mqtt.Client, topic string, callback mqtt.MessageHandler) {
	token := client.Subscribe(topic, 0, callback)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Fatalf("❌ Erro ao se inscrever no tópico %s: %v", topic, err)
	}
	fmt.Printf("📡 Inscrito no tópico %s\n", topic)
}

/**
 * Handler padrão para mensagens recebidas via MQTT.
 *
 * @param client o cliente MQTT
 * @param msg    a mensagem recebida
 */
func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("📩 Mensagem recebida [%s]: %s\n", msg.Topic(), msg.Payload())
}

/**
 * Handler específico para mensagens de confirmação de reserva.
 *
 * @param client o cliente MQTT
 * @param msg    a mensagem de confirmação
 */
func reservationHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("✅ Confirmação de reserva recebida para o carro %d: %s\n", clientID, string(msg.Payload()))
}

/**
 * Gera uma posição aleatória com coordenadas de origem e destino.
 *
 * @return quatro inteiros representando originX, originY, destinationX, destinationY
 */
func generatePosition() (int, int, int, int) {
	originX := rand.Intn(1000)
	originY := rand.Intn(1000)
	destinationX := rand.Intn(1000)
	destinationY := rand.Intn(1000)

	return originX, originY, destinationX, destinationY
}

/**
 * Inicia o loop de simulação do carro.
 * Publica coordenadas aleatórias de posição e destino nos tópicos apropriados a cada 2 segundos.
 *
 * @param client o cliente MQTT
 */
func startCarLoop(client mqtt.Client) {
	for {
		origX, origY, destX, destY := generatePosition()
		route := fmt.Sprintf("%d, %d, %d, %d, %d", clientID, origX, origY, destX, destY)

		publish(client, "car/position", route)
		fmt.Println("📤 Enviado para car/position:", route)

		publish(client, "car/recarga", route)
		fmt.Println("📤 Enviado para car/recarga:", route)

		time.Sleep(2 * time.Second)
	}
}

/**
 * Função principal da aplicação.
 * Gera um ID aleatório para o carro, conecta ao broker MQTT, se inscreve nos tópicos necessários
 * e inicia a simulação de movimentação do carro.
 */
func main() {
	rand.Seed(time.Now().UnixNano())

	clientID = rand.Intn(1000)
	id := fmt.Sprintf("%d", clientID)

	client := connectMQTT(brokerURL, id)
	defer disconnectMQTT(client)

	subscribe(client, "car/recarga", defaultMessageHandler)
	reservationTopic := fmt.Sprintf("car/%d/reservation", clientID)
	subscribe(client, reservationTopic, reservationHandler)

	startCarLoop(client)
}
