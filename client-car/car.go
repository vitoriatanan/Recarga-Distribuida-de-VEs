package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	brokerURL      = "tcp://mosquitto:1883"
	clientID       = "carroA"
	reservaServidorA = "http://servidor-a:8080/reserva"
	reservaServidorB = "http://servidor-b:8080/reserva"
	reservaServidorC = "http://servidor-c:8080/reserva"
)

type Trecho struct {
	Origem  string `json:"origem"`
	Destino string `json:"destino"`
}

type ReservaRequest struct {
	ClienteID       string `json:"cliente_id"`
	Trecho          Trecho `json:"trecho"`
	ProximoServidor string `json:"proximo_servidor,omitempty"`
}

type ReservaResponse struct {
	Status    string `json:"status"`
	ReservaID string `json:"reserva_id"`
	Mensagem  string `json:"mensagem"`
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

	func main() {
		rand.Seed(time.Now().UnixNano())
	
		client := connectMQTT(brokerURL, clientID)
		defer disconnectMQTT(client)
	
		subscribe(client, "car/recarga", defaultMessageHandler)
		startCarLoop(client)
	}

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

func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("📩 Mensagem recebida [%s]: %s\n", msg.Topic(), msg.Payload())
}

func generatePosition() (int, int) {
	return rand.Intn(1000), rand.Intn(1000)
}

func reservarTrechos() {
	// Define a sequência dos trechos com prioridade
	requestA := ReservaRequest{
		ClienteID: clientID,
		Trecho: Trecho{
			Origem:  "João Pessoa",
			Destino: "Maceió",
		},
		ProximoServidor: reservaServidorB,
	}

	requestB := ReservaRequest{
		ClienteID: clientID,
		Trecho: Trecho{
			Origem:  "Maceió",
			Destino: "Sergipe",
		},
		ProximoServidor: reservaServidorC,
	}

	requestC := ReservaRequest{
		ClienteID: clientID,
		Trecho: Trecho{
			Origem:  "Sergipe",
			Destino: "Feira de Santana",
		},
	}

	// Inicia a requisição em cadeia
	fmt.Println("🚗 Iniciando sequência de reserva de pontos...")
	if !fazerReserva(reservaServidorA, requestA) {
		fmt.Println("❌ Falha na reserva inicial (Empresa A)")
		return
	}
	if !fazerReserva(reservaServidorB, requestB) {
		fmt.Println("❌ Falha na reserva intermediária (Empresa B)")
		return
	}
	if !fazerReserva(reservaServidorC, requestC) {
		fmt.Println("❌ Falha na reserva final (Empresa C)")
		return
	}
	fmt.Println("✅ Todas as reservas realizadas com sucesso!")
}

func fazerReserva(url string, req ReservaRequest) bool {
	body, _ := json.Marshal(req)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("❌ Erro ao chamar %s: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()

	var res ReservaResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Printf("❌ Erro ao decodificar resposta: %v\n", err)
		return false
	}

	if res.Status != "ok" {
		log.Printf("❌ Falha na reserva: %s\n", res.Mensagem)
		return false
	}

	log.Printf("🔒 Reserva bem-sucedida: %s (%s)\n", res.ReservaID, res.Mensagem)
	return true
}

func startCarLoop(client mqtt.Client) {
	for {
		x, y := generatePosition()
		position := fmt.Sprintf("Carro A - posição x:%d, y:%d", x, y)

		publish(client, "car/position", position)
		fmt.Println("📤 Enviado para car/position:", position)

		publish(client, "car/recarga", position)
		fmt.Println("📤 Enviado para car/recarga:", position)

		time.Sleep(5 * time.Second)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	client := connectMQTT(brokerURL, clientID)
	defer disconnectMQTT(client)

	go startCarLoop(client)

	// Aguarda um pouco e tenta a reserva
	time.Sleep(3 * time.Second)
	reservarTrechos()

	// Mantém o programa vivo
	select {}
}
