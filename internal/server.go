package internal

import (
	"errors"
	"github.com/mborders/logmatic"
	"net"
	"strconv"
)

const BUFFER_SIZE = 1024

type Server struct {
	portConfig     []string
	config         Config
	listenSocket   net.Listener
	availablePorts chan string
}

type Link struct {
	PortConfig      string
	ClientCon       net.Conn
	NodeConn        net.Conn
	NodeWriteChan   chan []byte
	NodeReadChan    chan []byte
	ClientWriteChan chan []byte
	ClientReadChan  chan []byte
	Quit            bool
}

/**
Available port algorithm: Sorted map in Server struct contains all ports. A flag on each port determines whether it is
available. The ports collection contains a pointer to the next available port. The pointer will always point to the
lowest port available.
*/

type Node struct {
	Port      int
	Available bool
}

type NodeManager struct {
	Nodes []Node
}

func (mgr *NodeManager) NextNode() (*Node, error) {
	for _, node := range mgr.Nodes {
		if node.Available {
			node.Available = false
			return &node, nil
		}
	}

	return nil, errors.New("no nodes available")
}

func (server *Server) Serve(config Config) {

	logger := logmatic.NewLogger()

	server.config = config
	server.portConfig = config.GetNodes()
	server.availablePorts = make(chan string, len(server.portConfig))

	for _, port := range server.portConfig {
		server.availablePorts <- port
	}

	logger.Info("Serving customers on %s:%d", config.Service.Host, config.Service.Port)
	listenSocket, err := net.Listen("tcp", config.Service.Host+":"+strconv.Itoa(config.Service.Port))

	if err != nil {
		logger.Error("Error while establishing server socket: %s", err.Error())
	}

	defer listenSocket.Close()

	server.listenSocket = listenSocket

	for {
		clientConn, err := listenSocket.Accept()

		if err != nil {
			logger.Error("Error while accepting connection: %s", err.Error())
			clientConn.Close()
			continue
		} else {
			logger.Info("Accepted connection from client.")
		}

		tcpAddr, err := net.ResolveTCPAddr("tcp", server.portConfig[0])
		if err != nil {
			logger.Error("ResolveTCPAddr failed: %s", err.Error())
			continue
		}

		nodeConn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			logger.Error("Dial failed: %s", err.Error())
			clientConn.Close()
			continue
		}

		logger.Info("Made connection to node: %s", server.portConfig[0])

		link := Link{ClientCon: clientConn, NodeConn: nodeConn, PortConfig: server.portConfig[0], Quit: false}
		link.ClientReadChan = make(chan []byte, 50)
		link.ClientWriteChan = make(chan []byte, 50)
		link.NodeReadChan = make(chan []byte, 50)
		link.NodeWriteChan = make(chan []byte, 50)

		go processConnection(link)
	}

}

func processConnection(link Link) {

	link.ClientCon.Write([]byte("Connecting you..."))

	go handleClientWrite(link)
	go handleNodeWrite(link)
	go handleNodeRead(link)
	go handleClientRead(link)
}

/*
*
A not very clean way to handle closing connections....
*/
func handleQuit(link Link) {
	link.Quit = true
	link.ClientCon.Close()
	link.NodeConn.Close()
}

/*
*
Reads data from the node socket and places it on the client write channel to be written
*/
func handleNodeRead(link Link) {
	logger := logmatic.NewLogger()

	for link.Quit == false {
		buffer := make([]byte, BUFFER_SIZE)

		mLen, err := link.NodeConn.Read(buffer)

		if err != nil {
			logger.Error("Error reading node: %s", err.Error())
			link.Quit = true
			go handleQuit(link)
		} else {
			link.ClientWriteChan <- buffer[:mLen]
		}
	}
}

func handleNodeWrite(link Link) {
	logger := logmatic.NewLogger()

	for link.Quit == false {
		buffer := <-link.NodeWriteChan

		_, err := link.NodeConn.Write(buffer)

		if err != nil {
			logger.Error("Error writing data to node: %s", err.Error())
			link.Quit = true
			go handleQuit(link)
		}
	}
}

func handleClientRead(link Link) {
	logger := logmatic.NewLogger()

	for link.Quit == false {
		buffer := make([]byte, BUFFER_SIZE)
		mLen, err := link.ClientCon.Read(buffer)

		if err != nil {
			logger.Error("Error reading client: %s", err.Error())
			link.Quit = true
			go handleQuit(link)
		} else {
			link.NodeWriteChan <- buffer[:mLen]
		}
	}
}

func handleClientWrite(link Link) {
	logger := logmatic.NewLogger()

	for link.Quit == false {
		buffer := <-link.ClientWriteChan

		logger.Debug("Client Write: %s", string(buffer))

		_, err := link.ClientCon.Write(buffer)

		if err != nil {
			logger.Error("Error writing data to client: %s", err.Error())
			link.Quit = true
			go handleQuit(link)
		}
	}
}
