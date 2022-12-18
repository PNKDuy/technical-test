package main

import (
	"net"
	"technical-test/msg"
	"technical-test/node"
	"time"
)

func main() {

	node1 := node.Node{
		Host:      "127.0.0.1",
		Port:      "8000",
		KnowPeers: []node.Node{},
	}

	node2 := node.Node{
		Host:      "127.0.0.1",
		Port:      "8080",
		KnowPeers: []node.Node{},
	}

	node3 := node.Node{
		Host:      "127.0.0.1",
		Port:      "5050",
		KnowPeers: []node.Node{},
	}

	node1.KnowPeers = append(node1.KnowPeers, node2, node3)
	node2.KnowPeers = append(node2.KnowPeers, node1, node3)
	node3.KnowPeers = append(node1.KnowPeers, node1, node2)

	node1.ListenAndServe()
	node2.ListenAndServe()
	node3.ListenAndServe()

	node1Conns, _ := node1.DiscoverPeers()
	node2Conns, _ := node2.DiscoverPeers()
	node3Conns, _ := node3.DiscoverPeers()

	go func() {
		for _, c := range node1Conns {
			sendMsgEvery5s(c)
		}
	}()

	go func() {
		for _, c := range node2Conns {
			sendMsgEvery5s(c)
		}
	}()

	go func() {
		for _, c := range node3Conns {
			sendMsgEvery5s(c)
		}
	}()

	// wait forever
	select {}
}

func sendMsgEvery5s(conn net.Conn) {
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)
	msg := msg.Message{
		Msg:          "this is a message",
		OwnerAddress: conn.LocalAddr().String(),
	}
	msgBytes, err := msg.ParseIntoByte()
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				conn.Write(msgBytes)
			}
		}
	}()

	time.Sleep(2 * time.Minute)
	ticker.Stop()
	done <- true
}
