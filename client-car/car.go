package main

import (

	"fmt"
	"log"


	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Criando um novo cliente MQTT
	opts := mqtt.NewClientOptions().AddBroker("tcp://mosquitto:1883").SetClientID("car")
	client := mqtt.NewClient(opts)

	// Conecta ao broker MQTT
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	}

	fmt.Println("✅ Carro conectado ao broker MQTT!")

	defer client.Disconnect(250)

	// Subscribe to the topic "car/position"
	// if token := client.Subscribe("car/position", 0, func(client mqtt.Client, msg mqtt.Message) {
	// 	println(string(msg.Payload()))
	// }); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }

	// select {}
}

