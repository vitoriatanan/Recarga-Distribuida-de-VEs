package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Conetando ao broker MQTT
	opts := mqtt.NewClientOptions().AddBroker("tcp://mosquitto:1883").SetClientID("server")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	}

	fmt.Println("✅ Servidor conectado ao broker MQTT!")

	defer client.Disconnect(250)

	// Assina no tópico car/position
	if token := client.Subscribe("car/position", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("📥 Recebido: %s\n", msg.Payload())
	}); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	select {} // mantém o programa rodando
}
