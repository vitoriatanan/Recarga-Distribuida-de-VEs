package main

import (

	"fmt"
	"log"
	"math/rand"
	"time"


	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connectAndCommunicate(broker string, topic string, clientID string, playload string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Mensagem recebida no tópico %s: %s\n", msg.Topic(), msg.Payload())
	})

	client := mqtt.NewClient(opts)
	if toker := client.Connect(); toker.Wait() && toker.Error() != nil {
		log.Fatalf("Erro de conexão: %v", toker.Error())
	}

	// Publica a mensagem no tópico
	token := client.Publish(topic, 0, false, playload)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao publicar mensagem: %v", token.Error())
	}

	// Se inscreve no tópico
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao se inscrever no tópico: %v", token.Error())
		return
	}

	fmt.Printf("Mensagem publicada no tópico %s: %s\n", topic, playload)

	// Aguarda a mensagem ser recebida		
	// Aguarda 10 segundos para receber mensagens
	time.Sleep(10 * time.Second)

	// Desconecta o cliente
	client.Disconnect(250)
	fmt.Println("Cliente desconectado do broker MQTT.")	

}

func main() {
	// Semente para gerar posições diferentes
	rand.Seed(time.Now().UnixNano())

	// Criando um novo cliente MQTT
	opts := mqtt.NewClientOptions().AddBroker("tcp://mosquitto:1883").SetClientID("car")
	client := mqtt.NewClient(opts)

	// Conecta ao broker MQTT
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	}

	fmt.Println("✅ Carro conectado ao broker MQTT!")

	defer client.Disconnect(250)

	// Gera coordenadas aleatórias e envia para o tópico car/position
	for {
		x := rand.Intn(1000)  // coordenada X entre 0 e 999
		y := rand.Intn(1000)  // coordenada Y entre 0 e 999

		position := fmt.Sprintf("Carro A - posição x:%d, y:%d", x, y)
		token := client.Publish("car/position", 0, false, position)
		token.Wait()

		fmt.Println("📤 Enviado:", position)
		time.Sleep(2 * time.Second)

		broker := "tcp://mosquitto:1883"
		topic := "car/recarga"
		clientID := "clienteCarroA"
		playload := "Carro A - posição x:100, y:200"

		connectAndCommunicate(broker, topic, clientID, playload)
	}

}

