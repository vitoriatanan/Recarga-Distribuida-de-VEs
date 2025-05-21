package functions

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// ======== FUNÇÕES UTILITÁRIAS ========

/**
 * Obtém o nome do host atual da máquina ou da variável de ambiente INSTANCE_NAME.
 *
 * @return string nome do host ou da instância.
 */
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Erro ao obter o nome do host: %v", err)
	}
	if hostname == "" {
		hostname = os.Getenv("INSTANCE_NAME")
	}
	return hostname
}

/**
 * Define os limites da localização do servidor com base no nome da empresa.
 * Os limites são representados por um slice com os valores [xMin, xMax, yMin, yMax].
 *
 * @param serverName nome do servidor (empresa-a, empresa-b, etc.)
 * @param serverLocation slice que será sobrescrito com os novos limites
 * @return []int limites da área do servidor
 */
func SetServerLocation(serverName string, serverLocation []int) []int {
	switch serverName {
	case "empresa-a":
		serverLocation = []int{0, 500, 0, 500}
	case "empresa-b":
		serverLocation = []int{0, 500, 501, 1000}
	case "empresa-c":
		serverLocation = []int{501, 1000, 0, 500}
	case "empresa-d":
		serverLocation = []int{501, 1000, 501, 1000}
	default:
		serverLocation = []int{-1, -1}
	}
	return serverLocation
}

/**
 * Verifica se uma posição (x, y) está dentro dos limites de uma empresa.
 *
 * @param x coordenada X da posição
 * @param y coordenada Y da posição
 * @param serverLocation slice com os limites da empresa [xMin, xMax, yMin, yMax]
 * @return bool true se estiver dentro dos limites, false caso contrário
 */
func IsPositionInCompanyLimits(x, y int, serverLocation []int) bool {
	if (x >= serverLocation[0] && x <= serverLocation[1]) && (y >= serverLocation[2] && y <= serverLocation[3]) {
		fmt.Println("🚗 O carro está dentro dos limites da empresa.")
		return true
	}
	fmt.Println("🚗 O carro está fora dos limites da empresa.")
	return false
}

/**
 * Tenta reservar um ponto de recarga para um carro.
 * Percorre o mapa de estações e reserva a primeira que estiver livre (valor 0).
 *
 * @param carID ID do carro que solicita a reserva
 * @param stationsSpots mapa de estações com seus respectivos IDs de reserva (0 = livre)
 * @return string nome da estação reservada, ou string vazia se nenhuma estiver disponível
 */
func StationReservation(carID int, stationsSpots map[string]int) string {
	for station, reservation := range stationsSpots {
		if reservation == 0 {
			stationsSpots[station] = carID
			return station
		}
	}
	return ""
}

/**
 * Envia a posição atual do carro (x, y) para os demais servidores da rede,
 * exceto para o servidor local.
 *
 * @param x coordenada X da posição do carro
 * @param y coordenada Y da posição do carro
 * @param serverName nome do servidor atual (para evitar autoenvio)
 */
func SendPositionToServers(x, y int, serverName string) {
	servers_ip := []string{"8081", "8082", "8083", "8084"}
	servers := []string{"empresa-a", "empresa-b", "empresa-c", "empresa-d"}

	for i := range servers {
		name := servers[i]
		ip := servers_ip[i]

		// Evita envio para o próprio servidor
		if name == serverName {
			continue
		}

		url := fmt.Sprintf("http://%s:%s/server/forward", name, ip)
		jsonStr := fmt.Sprintf(`{"x":%d, "y":%d}`, x, y)

		resp, err := http.Post(url, "application/json", strings.NewReader(jsonStr))
		if err != nil {
			log.Printf("❌ Erro ao enviar para %s: %v\n", name, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("📨 Posição (%d, %d) enviada para %s\n", x, y, name)
		} else {
			fmt.Printf("⚠️ Resposta do servidor %s: %s\n", name, resp.Status)
		}
	}
}
