package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	mplex "github.com/libp2p/go-libp2p-mplex"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	security := libp2p.Security(secio.ID, secio.New)

	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/tcp/0/ws",
	)

	// TODO: Configure libp2p to use a DHT with a libp2p.Routing option

	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		muxers,
		security,
	)
	if err != nil {
		panic(err)
	}

	host.SetStreamHandler(chatProtocol, chatHandler)

	for _, addr := range host.Addrs() {
		fmt.Println("Listening on", addr)
	}

	targetAddr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/59675/p2p/QmchtgByxcQPdVrwnq2KneE3ako5o65MyiWdYDrBjMpVDP")
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

	// TODO: Implement a simple read loop to read user input and send messages
	donec := make(chan struct{}, 1)
	go chatInputLoop(ctx, host, donec)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	case <-donec:
		host.Close()
	}
}
