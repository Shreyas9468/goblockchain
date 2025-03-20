# Goblockchain

Goblockchain is a lightweight blockchain implementation written in Go, designed to demonstrate core cryptocurrency concepts such as wallet creation, transaction signing, proof-of-work mining, and basic network consensus. This project runs two HTTP servers: one for blockchain operations (port `:8080`) and another for wallet management (port `:8081`). This README provides detailed instructions for setting up and running the project on Windows using Visual Studio Code (VS Code).

## Project Structure

- **`blockchain/`**: Core blockchain functionality.
  - `block.go`: Defines the `Block` struct and methods for creating and hashing blocks.
  - `transaction.go`: Manages the `Transaction` struct and ECDSA signature verification.
  - `blockchain.go`: Implements the `BlockChain` struct, handling transaction pooling, mining, proof-of-work, and balance calculation.
- **`wallet/`**: Wallet operations.
  - `wallet.go`: Generates ECDSA key pairs, derives addresses, and signs transactions.
- **`server/`**: HTTP server implementations.
  - `blockchain_server.go`: Runs on `:8080`, providing endpoints for transaction submission, mining, chain retrieval, balance checks, and consensus.
  - `wallet_server.go`: Runs on `:8081`, offering endpoints for wallet creation, balance queries, and transaction creation.
- **`network/`**: Network synchronization.
  - `network.go`: Implements the longest chain rule for consensus across peer nodes.
- **`main.go`**: Entry point that initializes the blockchain and starts both servers.
- **`.gitignore`**: Excludes `.json` files (e.g., `wallet1.json`, `chain.json`) from version control.
- **`go.mod`**: Defines the Go module (`goblockchain`) and dependencies.

## Prerequisites

- **Go**: Version 1.16 or higher installed on Windows.
  - Verify: Open PowerShell and run `go version`. Expected output: `go version go1.xx.x windows/amd64`.
  - Install from [golang.org/dl/](https://golang.org/dl/) if needed, and add `C:\Go\bin` to your system PATH.
- **Visual Studio Code (VS Code)**: Installed with the Go extension for enhanced development support.
  - Install the Go extension: `Ctrl+Shift+X`, search "Go" by Go Team at Google, and install.
- **Git**: Installed for version control.
  - Verify: `git --version` in PowerShell. Install from [git-scm.com](https://git-scm.com/downloads) if missing.

## Setup on Windows

1. **Clone or Download the Repository**:
   - Open PowerShell:
     ```powershell
     git clone https://github.com/Shreyas9468/goblockchain.git
     cd goblockchain
