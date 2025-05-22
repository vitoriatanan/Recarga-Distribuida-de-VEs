# Recarga-Distribuida-de-VEs

## Descrição do Projeto
Este projeto implementa um sistema distribuído para gerenciar a recarga de veículos elétricos (VEs). Ele utiliza contêineres Docker para simular diferentes componentes, como empresas e clientes, que se comunicam através do protocolo MQTT usando o broker Mosquitto e também comunicação API, utilizando Gin. Cada empresa é representada como um serviço independente, e o cliente (client-car) interage com essas empresas para realizar operações relacionadas à recarga.

## Objetivo
O objetivo principal do projeto é simular um ambiente distribuído para gerenciar a recarga de veículos elétricos, permitindo a comunicação entre diferentes entidades (empresas e clientes) de forma eficiente e escalável.

## Estrutura do Projeto
- **Mosquitto**: Broker MQTT para comunicação entre os serviços.
- **Client-Car**: Cliente que interage com as empresas para realizar operações.
- **Empresas (A, B, C, D)**: Serviços que representam diferentes empresas no sistema, cada uma com sua própria configuração e porta.

## Tecnologias Utilizadas
- **Linguagem**: Go (Golang)
- **Frameworks**:
  - [Gin](https://github.com/gin-gonic/gin): Para criação de APIs REST.
  - [Paho MQTT](https://github.com/eclipse/paho.mqtt.golang): Para comunicação MQTT.
- **Contêineres**: Docker e Docker Compose.
- **Broker MQTT**: Eclipse Mosquitto.

## Pré-requisitos
- Docker e Docker Compose instalados no sistema.

## Como Executar o Projeto
1. Clone o repositório:
   ```bash
   git clone <URL_DO_REPOSITORIO>
   cd Recarga-Distribuida-de-VEs

2. Execute o Docker Compose para iniciar os serviços:
   ```bash
   docker-compose up --build
   ```


3. Para parar os serviços, use:
   ```bash
   docker-compose down
   ```

## Configurações

- O broker Mosquitto está configurado para escutar na porta 1884 (mapeada para 1883 no contêiner).
- Cada empresa possui uma porta específica:
    - Empresa A: 8081
    - Empresa B: 8082
    - Empresa C: 8083
    - Empresa D: 8084

- O cliente (client-car) está configurado para se conectar ao broker Mosquitto na porta 1884 e às empresas nas portas correspondentes.

## APIs REST
Os servidores expõem as seguintes rotas:

- `GET /server/position`: Retorna os limites da área coberta pelo servidor.
- `GET /server/route`: Retorna a última rota processada.
- `GET /server/stations`: Retorna o status de todas as estações de recarga.
- `GET /server/stations/:station`: Retorna o status de uma estação específica.
- `POST /server/forward`: Recebe posição de outro servidor para tentativa de reserva.

## Autores
- Kauan Caio de Arruda Farias
- Nathielle Cerqueira Alves
- Vitória Tanan dos Santos
