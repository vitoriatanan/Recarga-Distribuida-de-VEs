package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://mosquitto:1883").SetClientID("server") // Altere para o endereço do broker MQTT
	client := mqtt.NewClient(opts)                                                            // Cria um novo cliente MQTT

	// Conecta ao broker MQTT
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	}
	fmt.Println("✅ Servidor A conectado ao broker MQTT!")

	defer client.Disconnect(250) // Desconecta do broker MQTT após o término

	// Inscreve-se no tópico "car/position" para receber mensagens do carro
	// Quando uma mensagem é recebida, ela é processada na função de callback
	if token := client.Subscribe("car/position", 0, func(client mqtt.Client, msg mqtt.Message) {
		position := string(msg.Payload())
		fmt.Println("📥 Posição recebida do carro:", position)

	// 	// Envia requisição POST para Empresa B
	// 	okB := sendReservationRequest("http://localhost:8081/reservar", position)
	// 	// Envia requisição POST para Empresa C
	// 	okC := sendReservationRequest("http://localhost:8082/reservar", position)

	// 	// Envia resposta de volta ao carro
	// 	var response string
	// 	if okB && okC {
	// 		response = "Reserva confirmada nas duas empresas"
	// 	} else {
	// 		response = "Falha na reserva"
	// 	}

	// 	client.Publish("car/response", 0, false, response)
	// 	fmt.Println("📤 Resposta enviada ao carro:", response)
	// }); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }

	select {} // Mantém o servidor em execução indefinidamente
}

func sendReservationRequest(url string, position string) bool {
	jsonData := []byte(fmt.Sprintf(`{"car_position":"%s"}`, position))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Erro ao enviar requisição para", url, "->", err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
