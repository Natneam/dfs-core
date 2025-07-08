# Distributed File System

This project is a simple implementation of a distributed file storage system in Go. It creates a peer-to-peer network where files can be stored, retrieved, and deleted. The system is designed to be fault-tolerant and scalable, with data encrypted both in transit and at rest.

## Table of Contents

- [Architecture](#architecture)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Compilation](#compilation)
  - [Running the Application](#running-the-application)
- [Usage](#usage)
  - [Interactive CLI](#interactive-cli)
- [Testing](#testing)
- [Project Structure](#project-structure)

## Architecture

The system is built as a peer-to-peer network of nodes. Each node in the network is a server that can communicate with other nodes.

- **Peer-to-Peer Network**: Nodes connect to each other to form a network. When a file is uploaded to one node, it is broadcasted and replicated across other nodes in the network. When a file is requested, the network is searched to find and serve the file.
- **TCP Transport**: Communication between nodes is handled over TCP. Each node listens on a specific port for incoming connections from other peers.
- **File Storage**: Files are not stored with their original names. Instead, a key is used. The key is hashed, and this hash is used to determine the storage path and filename on disk. This provides a uniform way of addressing files across the network.
- **Encryption**: All files are encrypted before being written to disk using AES encryption. The same key is used for decryption when a file is retrieved. This ensures that the file contents are secure.

## Features

- **Distributed Storage**: Files are replicated across multiple nodes in the network.
- **Content-Addressable Storage**: Files are stored and retrieved using a key, which is hashed to create a unique address.
- **Data Encryption**: Files are encrypted using AES to ensure data privacy.
- **Command-Line Interface**: An interactive CLI is provided to interact with the file system.
- **Fault Tolerance**: The distributed nature of the system provides a level of fault tolerance. If one node goes down, files can still be retrieved from other nodes in the network.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) version 1.23 or higher.

### Compilation

You can compile the project using the provided `Makefile`.

```bash
make build
```

This will create an executable binary at `bin/fs`.

### Running the Application

To run the application, you can use the `run` target in the `Makefile` or execute the binary directly. The application accepts the following command-line flags:

- `-port`: The port for the server to listen on (default: `3000`).
- `-nodes`: A comma-separated list of bootstrap nodes to connect to.

#### Running a Single Node

To start the first node in the network, run the following command:

```bash
./bin/fs -port 3000
```

This will start a server listening on port 3000.

#### Running Multiple Nodes

To create a network, you can start more nodes and connect them to the first node. Open a new terminal and run:

```bash
./bin/fs -port 4000 -nodes localhost:3000
```

This will start a second server on port 4000 and connect it to the node running on port 3000. You can connect more nodes by specifying the address of any existing node across the internet.

## Usage

The application provides an interactive command-line interface to store, retrieve, and delete files.

### Interactive CLI

Once a node is running, you can use the following commands:

- **Connect to a node:**
  ```
  > connect <node_address>
  ```
  This command connects the current node to another node in the network. The address should be in the format `host:port`.

- **List connected nodes:**
  ```
  > peers
  ```
  This command lists all nodes currently connected to the network.

- **Store a file:**
  ```
  > put <local_file_path> <remote_filename>
  ```
  This command reads a local file and stores it on the network with the specified remote filename as the key.

- **Retrieve a file:**
  ```
  > get <remote_filename>
  ```
  This command retrieves the file associated with the given key from the network and prints its content to the console.

- **Delete a file:**
  ```
  > delete <remote_filename>
  ```
  This command deletes the file associated with the given key from the local storage of the node.

- **Clear the console:**
  ```
  > clear
  ```

- **Exit the application:**
  ```
  > exit
  ```

## Testing

To run the test suite for the project, use the `test` target in the `Makefile`.

```bash
make test
```

This will run all `_test.go` files in the project.

## Project Structure

```
.
├── Makefile          # Makefile for building, running, and testing.
├── README.md         # This file.
├── bin/              # Compiled binaries.
├── cipher/           # Cryptographic functions (encryption/decryption).
├── cli/              # Command-line interface logic.
├── network/          # Network transport and communication logic.
├── server/           # File server implementation.
└── store/            # File storage logic.
```
