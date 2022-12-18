package node

import (
	"encoding/json"
	"fmt"
	"net"
	"technical-test/msg"
)

type Node struct {
	Host      string
	Port      string
	KnowPeers []Node
}

func (n *Node) ListenAndServe() (err error) {
	listener, err := net.Listen("tcp", net.JoinHostPort(n.Host, n.Port))
	if err != nil {
		return err
	}

	fmt.Println("Listening on address: ", listener.Addr().String())

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go n.processConn(conn)
		}
	}()

	return nil
}

func (n *Node) processConn(conn net.Conn) {
	var msg msg.Message
	// create buffer
	buf := make([]byte, 1024)

	fmt.Println("Conn established at: ", conn.LocalAddr())
	go func() {
		for {
			l, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Read error", err)
				break
			}

			// msgBytes := string(buf[:l])
			err = json.Unmarshal(buf[:l], &msg)
			if err != nil {
				return
			}

			if n.isReceivedMsgValidated(&msg) {
				return
			}
			fmt.Println(conn.LocalAddr().String(), " : ", msg.Msg, " from ", msg.OwnerAddress)
		}
	}()
}

func (n *Node) DiscoverPeers() ([]net.Conn, error) {
	var conns []net.Conn
	for _, node := range n.KnowPeers {
		conn, err := net.Dial("tcp", net.JoinHostPort(node.Host, node.Port))
		if err != nil {
			return nil, err
		}
		conns = append(conns, conn)
	}

	return conns, nil
}

func (n *Node) SendMsg(conn net.Conn, msg msg.Message) error {
	msgBytes, err := msg.ParseIntoByte()
	if err != nil {
		return err
	}

	conn.Write(msgBytes)

	return nil
}

func (n *Node) isReceivedMsgValidated(msg *msg.Message) bool {
	for _, v := range n.KnowPeers {
		if net.JoinHostPort(v.Host, v.Port) == msg.OwnerAddress {
			return true
		}
	}

	return false
}
