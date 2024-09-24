# On-Chain Statistics Project

This project is designed to fetch and analyze on-chain statistics from the Evmos blockchain. It specifically addresses the following requirements:

1. List all the smart contracts that were used between block 100 and 200 and sort them by the amount of interactions.
2. Of all the wallets that interacted with the network, sort them by balance to find the richest user.

## Project Structure
The project is structured as follows:
- `main.go`: The entry point of the application.
- `evmos_client.go`: Contains the client to interact with the Evmos node.
- `service.go`: Contains the service to fetch and analyze on-chain statistics.

#### Support several endpoints:

- **/**: For health check
- **accounts**: Returns the list of accounts found in a local node of evmos. Not really utilized. Just there for testing purposes.
- **balance**: Returns the balance of a specific account at a specific block (Default latest).
- **blocknumber**: Returns the block number of the latest block.
- **block**: Returns the block information of a specific block number.
- **transactiontrace**: Returns the transaction trace of a specific transaction hash.
- **smartcontracts**: Retrieves the interactions of smart contracts used between block 100 and 200.
- **richestusers**: Calculates the richest users based on their wallet balances at block 200.

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

## Technical Decisions
1. **Concurrency with Goroutines**: Utilized goroutines to fetch wallet balances concurrently, reducing the overall execution time.
2. **Mocked Data**: Evmos endpoint for blocks, always returned an empty transaction list. To test the application, 
I created a mock data with transactions between blocks 100 and 200, and assumed the response of `transcation_tracer`.
3. **Save stats to csv \& BDD**: I did not have enough time to implement them.


## Assignment Checklist

- [x] Create an open-source (i.e public) GitHub repository to host your project.
- [x] Create a project that get on-chain statistics using the data between block 100 and 200.
- [x] Add tests for the functionality that you created.
- [x] Add to the README file
    - [x] Instructions to run the code.
    - [x] What were the main technical decisions you made and why you made them.
    - [x] Relevant comments about your project and how each of the steps were performed.
- [x] Implement the solution using GoLang.
- [x] Add a GitHub Action to run a linter (i.e, golang-ci) and tests on pull-requests.
- [ ] Create tests using Behaviour Driven Development
- [ ] Save the stats to sqlite or a csv files.
