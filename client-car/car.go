package main

import (

	"fmt"
	"log"
	"math/rand"
	"time"


	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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
	}
}

