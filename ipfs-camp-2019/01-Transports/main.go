package main

import (
	"context"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"

	"github.com/libp2p/go-libp2p"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: add some libp2p.Transport options to this chain!
	transports := libp2p.ChainOptions(
		//answer
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
		)

	host, err := libp2p.New(ctx, transports)
	if err != nil {
		panic(err)
	}

	host.Close()
}
