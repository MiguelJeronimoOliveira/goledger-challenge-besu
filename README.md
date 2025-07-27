# GoLedger Challenge - Besu Edition

This is the solution to the **GoLedger Challenge**, consisting of a REST API built in Go that interacts with a **Hyperledger Besu** blockchain network and a **PostgreSQL** database.

## Table of Contents

- [Project Architecture](#project-architecture)
- [Prerequisites](#prerequisites)
- [Setting up the Environment](#setting-up-the-environment)
  - [1. Start the Besu Network](#1-start-the-besu-network)
  - [2. Configure the API](#2-configure-the-api)
  - [3. Start the PostgreSQL Database](#3-start-the-postgresql-database)
  - [4. Install Go Dependencies](#4-install-go-dependencies)
  - [5. Run the Application](#5-run-the-application)
- [API Endpoints](#api-endpoints)
  - [`GET /get`](#get-get)
  - [`POST /set`](#post-set)
  - [`POST /sync`](#post-sync)
  - [`GET /check`](#get-check)


## Project Architecture

This project implements a Go application that interacts with a Hyperledger Besu blockchain network and a PostgreSQL database. The application exposes a REST API to manage a simple storage smart contract deployed on the blockchain and synchronize its value with the external database.

### Architecture Overview

#### 1. Go Application (Backend)

The core of the solution is a Go application organized into several packages, each with a specific responsibility:

- **main.go**: Entry point of the application. Loads environment variables, initializes the blockchain client, connects to the PostgreSQL database, and sets up the HTTP router.

- **config/config.go**: Loads configuration parameters from environment variables, such as the Besu node URL, smart contract address, private key for transactions, ABI file path, and PostgreSQL DSN.

- **contract/contract.go**: Provides the client for interacting with the Hyperledger Besu blockchain using the go-ethereum library to:

  - Connect to the Ethereum node via RPC.
  - Load the smart contract's ABI.
  - Initialize a bound contract instance for interacting with the SimpleStorage smart contract.
  - Implement methods to get the stored value and set a new value by sending a transaction.

- **db/postgresDB.go**: Manages the PostgreSQL database using `database/sql` and `github.com/lib/pq` to:

  - Open a connection using the provided DSN.
  - Ensure the `storage` table exists (creating it if necessary).
  - Initialize the table with a default value of `0` if empty.

- **handler/handler.go**: Contains the HTTP handlers for the REST API. Coordinates between incoming requests, the blockchain client, and the database to implement business logic.

- **router/router.go**: Sets up the API routes using `github.com/gin-gonic/gin`. Maps endpoints to their corresponding handlers.

## Prerequisites

Ensure the following tools are installed on your system:

- [Node.js and NPX](https://www.npmjs.com/get-npm)
- [Docker and Docker Compose](https://www.docker.com/)
- [Hyperledger Besu](https://besu.hyperledger.org/private-networks/get-started/install/binary-distribution)
- [Go (Golang)](https://golang.org/dl/)
- `jq`
	```	
	sudo apt-get install jq
	```
## Setting up the Environment

### 1. Start the Besu Network

Navigate to the `/besu` folder and run the following script:

```bash
cd besu
./startDev.sh
```

If the script exits the terminal before confirming the deployment, replace:

```bash
npx hardhat ignition deploy ./ignition/modules/deploy.js --network besu << EOF
y
EOF
```

with:

```bash
yes | npx hardhat ignition deploy ./ignition/modules/deploy.js --network besu
```

At the end of the execution, your **smart contract address** will be displayed in the terminal. Copy this address — you’ll need it in the next step.

### 2. Configure the API

- Navigate to the `/goledger-challenge` folder:

```bash
cd goledger-challenge
```

- Open the `.env` file and set the `CONTRACT_ADDRESS` variable:

```
CONTRACT_ADDRESS=your_contract_address_here
```

### 3. Start the PostgreSQL Database

```bash
docker-compose up -d
```

### 4. Install Go Dependencies

```bash
go mod tidy
```

### 5. Run the Application

```bash
go run .
```

## API Endpoints

the server is run in `http://localhost:8080`

### `GET /get`

Returns the current value stored in the smart contract.

**Response:**

```json
{
  "value": "<SMART_CONTRACT_VALUE>"
}
```

### `POST /set`

Sets a new value in the smart contract.

**Request body:**

```json
{
  "value": "<VALUE>"
}
```

### `POST /sync`

Fetches the smart contract value and persists it to the PostgreSQL database.

**Response:**

```json
{
  "synced_value": "<SMART_CONTRACT_VALUE>"
}
```

### `GET /check`

Checks if the smart contract value is equal to the value stored in the database.

**Response:**

```json
{
  "equal": true or false
}
```

