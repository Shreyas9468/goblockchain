# Goblockchain

Goblockchain is a lightweight blockchain implementation written in Go, designed to demonstrate core cryptocurrency concepts such as wallet creation, transaction signing, proof-of-work mining, and basic network consensus. This project runs two HTTP servers: one for blockchain operations (`:8080`) and another for wallet management (`:8081`).

## Project Structure

```
├── blockchain/
│   ├── block.go          # Defines the Block struct and hashing methods
│   ├── transaction.go    # Manages Transactions and ECDSA verification
│   ├── blockchain.go     # Handles transaction pooling, mining, PoW, and balances
│
├── wallet/
│   ├── wallet.go         # Generates ECDSA key pairs, derives addresses, and signs transactions
│
├── server/
│   ├── blockchain_server.go  # Runs on :8080, provides blockchain-related endpoints
│   ├── wallet_server.go      # Runs on :8081, manages wallets and transactions
│
├── network/
│   ├── network.go        # Implements network synchronization and consensus
│
├── main.go               # Entry point to start blockchain and wallet servers
├── .gitignore            # Excludes JSON files (wallet data, blockchain data)
├── go.mod                # Go module dependencies
```

## Prerequisites

- **Go** (1.16 or higher)
  - Verify installation: `go version`
  - Download: [golang.org/dl](https://golang.org/dl/)
  - Ensure `C:\Go\bin` is added to system PATH.
- **Visual Studio Code (VS Code)** with Go extension.
  - Install the Go extension from the Extensions Marketplace (`Go` by Go Team at Google).
- **Git**
  - Verify installation: `git --version`
  - Download: [git-scm.com](https://git-scm.com/downloads)

## Setup Instructions

### 1. Clone the Repository

```sh
git clone https://github.com/Shreyas9468/goblockchain.git
cd goblockchain
```

Alternatively, download the ZIP from GitHub, extract it, and navigate to the folder.

### 2. Open in VS Code

1. Launch VS Code.
2. Open the project folder: **File > Open Folder** > Select `goblockchain`.

### 3. Install Go Tools (Recommended)

- In VS Code, press `Ctrl+Shift+P`.
- Type **Go: Install/Update Tools** and select it.
- Check all tools (`gopls`, `go-outline`, etc.) and install.

### 4. Verify Dependencies

```sh
go mod tidy
```

Ensures all dependencies are correctly installed.

## Running the Project

### Start the Servers

Ensure you're in the project root and run:

```sh
go run main.go
```

**Expected Output:**

```
Blockchain Server running on :8080
Wallet Server running on :8081
```

Keep the terminal open while running the servers.

## Testing the Blockchain

### 1. Create a Wallet

```sh
curl -X GET http://localhost:8081/wallet/new -o wallet1.json
cat wallet1.json
```

**Output:**

```json
{
  "address": "46d2dcdb...",
  "private_key": "143601dd...",
  "public_key": "041e8b62..."
}
```

### 2. Check Initial Balance

```sh
ADDRESS=$(jq -r '.address' wallet1.json)
curl -X GET "http://localhost:8081/wallet/balance?address=$ADDRESS"
```

**Output:** `{ "balance": 0 }`

### 3. Create a Second Wallet

```sh
curl -X GET http://localhost:8081/wallet/new -o wallet2.json
cat wallet2.json
```

### 4. Send a Transaction

```sh
PRIVATE_KEY=$(jq -r '.private_key' wallet1.json)
RECIPIENT=$(jq -r '.address' wallet2.json)

curl -X POST http://localhost:8081/wallet/transaction \
  -H "Content-Type: application/json" \
  -d "{ \"private_key\": \"$PRIVATE_KEY\", \"recipient\": \"$RECIPIENT\", \"value\": 0.5 }"
```

**Output:** `{ "message": "Transaction sent to blockchain" }`

_(Note: Auto-mining processes this within 10 seconds.)_

### 5. Mine a Block (Optional)

```sh
curl -X GET http://localhost:8080/mine
```

**Output:** `{ "message": "New block mined" }`

### 6. Check Balances

#### Wallet 1 (Sender)

```sh
curl -X GET "http://localhost:8081/wallet/balance?address=$ADDRESS"
```

**Expected Output:** `{ "balance": -0.5 }`

#### Wallet 2 (Recipient)

```sh
RECIPIENT=$(jq -r '.address' wallet2.json)
curl -X GET "http://localhost:8081/wallet/balance?address=$RECIPIENT"
```

**Expected Output:** `{ "balance": 0.5 }`

#### Miner (node1_address)

```sh
curl -X GET "http://localhost:8081/wallet/balance?address=node1_address"
```

**Expected Output:** `{ "balance": 1 }` (or higher if multiple blocks are mined).

### 7. View the Blockchain

```sh
curl -X GET http://localhost:8080/chain -o chain.json
cat chain.json
```

**Purpose:** Retrieves the full chain for inspection.

### 8. Stop the Servers

In the terminal running `go run main.go`, press `Ctrl+C`.

## Features

- **Wallet Creation**: Secure ECDSA key generation.
- **Transactions**: Signed transactions stored in a pool before mining.
- **Mining**: Proof-of-work with difficulty requiring 3 leading zeros; auto-mines transactions every 10 seconds.
- **Consensus**: Implements the longest chain rule for synchronization across nodes.
- **REST APIs**: Blockchain operations available on `:8080`, wallet operations on `:8081`.

---

### License
This project is licensed under the MIT License. See `LICENSE` for details.

### Author
Developed by [Shreyas](https://github.com/Shreyas9468).

### Contributions
Contributions are welcome! Feel free to submit issues or pull requests.

---

This `README.md` follows best practices, ensuring clarity and usability. 🚀

