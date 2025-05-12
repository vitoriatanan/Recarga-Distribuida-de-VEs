package main

import (
	"encoding/json"
	"fmt"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Define o formato da mensagm q wue o servidor vai receber
type CarroMsg struct {
	CarID   string  `json:"car_id"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Bateria int     `json:"bateria"`
}

// Configura o servidor MQTT e começa a escutar mensagens dos carros
func InitMQTT(serverID string) {
	areaMin := 0.0
	areaMax := 0.0
	// Definição das áreas de cada empresa
	switch serverID {
	case "A":
		areaMin = 0
		areaMax = 100
	case "B":
		areaMin = 100
		areaMax = 200
	case "C":
		areaMin = 200
		areaMax = 300
	default:
		fmt.Println("Servidor com ID inválido")
		return
	}

	// Conecta ao broker MQTT
	opts := MQTT.NewClientOptions().AddBroker("tcp://mqtt:1883")
	opts.SetClientID("server-" + serverID)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

    fmt.Printf("✅ Servidor %s conectado ao broker MQTT!\n", serverID)
    
	// Assina o tópico e escuta o canal onde os carros publicam suas mensagens
	topic := "veiculo/recarga"
	client.Subscribe(topic, 0, func(client MQTT.Client, msg MQTT.Message) {
		var carro CarroMsg
		err := json.Unmarshal(msg.Payload(), &carro)
		if err != nil {
			fmt.Println("Erro ao decodificar mensagem:", err)
			return
		}

		// Verifica se o carro está na área da empresa
		if carro.X >= areaMin && carro.X < areaMax {
			fmt.Printf("[Servidor %s] Carro %s está na área da empresa %s (x=%.2f)\n", serverID, carro.CarID, serverID, carro.X)
		}
	})
}

// Recebe mensagaens de todos os carros. Só responde e imprime no terminal se o carro estiver dentro de sua área
func main() {
	serverID := os.Getenv("SERVER_ID")
	fmt.Printf("Servidor %s iniciado\n", serverID)
	InitMQTT(serverID)

	select {} // Mantém o servidor em execução
}