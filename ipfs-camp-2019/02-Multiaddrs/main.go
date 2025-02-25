package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

	// TODO: add some listen addresses with the libp2p.ListenAddrs or
	// libp2p.ListenAddrStrings configuration options.
	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/tcp/0/ws",
	)



	host, err := libp2p.New(ctx, transports,listenAddrs)
	if err != nil {
		panic(err)
	}

	// TODO: with our host made, let's connect to our bootstrap peer
	targetAddr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/63785/p2p/QmWjz6xb8v9K4KnYEwP5Yk75k5mMBCehzWFLCvvQpYxF3d")
	if err != nil {
		panic(err)
	}

	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		panic(err)
	}

	err = host.Connect(ctx, *targetInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to", targetInfo.ID)


	host.Close()
}
