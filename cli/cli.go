package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"natneam.github.io/dfs-core/server"
)

func Start() (int, []string, error) {
	listenAddress := flag.Int("port", 0, "Listen address of the server")
	peers := flag.String("peers", "", "Comma-separated list of bootstrapped nodes url to connect to")

	flag.Parse()

	if *listenAddress <= 0 || *listenAddress > 65535 {
		return 0, []string{}, fmt.Errorf("invalid port")
	}

	nodes := []string{}
	if len(*peers) > 0 {
		nodes = strings.Split(*peers, ",")
	}

	return *listenAddress, nodes, nil
}

func InteractiveCli(s *server.FileServer) {
	for {
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		line := scanner.Text()
		parts := strings.Split(line, " ")
		cmd, args := parts[0], parts[1:]

		switch cmd {
		case "put":
			handlePutCommand(s, args)
		case "get":
			handleGetCommand(s, args)
		case "delete":
			handleDeleteCommand(s, args)
		case "clear":
			fmt.Print("\033[H\033[2J")
			fmt.Println(s.Transporter.RemoteAddr())
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Unknown command:", cmd)
		}
	}
}

func handlePutCommand(s *server.FileServer, args []string) {
	if len(args) != 2 {
		fmt.Println("Usage: put <local_file_path> <remote_filename>")
		return
	}
	localFilePath := args[0]
	remoteFileName := args[1]

	file, err := os.Open(localFilePath)
	if err != nil {
		fmt.Printf("Error opening local file: %+v\n", err)
		return
	}

	defer file.Close()

	if s.Store(remoteFileName, file); err != nil {
		fmt.Printf("Error storing file on the network: %+v\n", err)
		return
	}

	fmt.Printf("File '%s' successfully stored as '%s' on the network.\n", localFilePath, remoteFileName)

}

func handleGetCommand(s *server.FileServer, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: get <remote_filename>")
		return
	}
	remoteFileName := args[0]

	_, file, err := s.Get(remoteFileName)
	if err != nil {
		fmt.Printf("Error retrieving data from the network: %+v\n", err)
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading data from the network: %+v\n", err)
		return
	}

	println(string(data))
}

func handleDeleteCommand(s *server.FileServer, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: delete <remote_filename>")
		return
	}
	fileName := args[0]

	if err := s.Delete(fileName); err != nil {
		fmt.Printf("Error deleting data from the server: %+v\n", err)
		return
	}

	fmt.Printf("Data deleted Successfully from you local server.\n")
}
