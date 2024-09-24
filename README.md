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

## Assignment Checklist

- [x] Create an open-source (i.e public) GitHub repository to host your project.
- [x] Create a project that get on-chain statistics using the data between block 100 and 200.
- [x] Add tests for the functionality that you created.
- [x] Add to the README file
    - [x] Instructions to run the code.
    - [ ] What were the main technical decisions you made and why you made them.
    - [ ] Relevant comments about your project and how each of the steps were performed.
- [x] Implement the solution using GoLang.
- [ ] Create tests using Behaviour Driven Development
- [x] Add a GitHub Action to run a linter (i.e, golang-ci) and tests on pull-requests.
- [ ] Save the stats to sqlite or a csv files.

