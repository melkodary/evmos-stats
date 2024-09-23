# On-Chain Statistics Project

This project is designed to fetch and analyze on-chain statistics from the Evmos blockchain. It specifically addresses the following requirements:

1. List all the smart contracts that were used between block 100 and 200 and sort them by the amount of interactions.
2. Of all the wallets that interacted with the network, sort them by balance to find the richest user.

## Prerequisites

- Go 1.21 or later
- An Evmos node running and accessible at `http://localhost:8545`

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/melkodary/evmos-stats.git
    cd evmos-stats
    ```

2. Run the project:
    ```sh
    go run main.go
    ```