package main

import (
	"context"
	//"flag"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	//"github.com/multiformats/go-multiaddr"
	//"io"
	//"log"
	mrand "math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	//"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	//pubsub "github.com/libp2p/go-libp2p-pubsub"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	tcp "github.com/libp2p/go-tcp-transport"
	//ws "github.com/libp2p/go-ws-transport"
)

type mdnsNotifee struct {
	h   host.Host
	ctx context.Context
}

func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	m.h.Connect(m.ctx, pi)
}

func main() {
	//var target string
	//flag.StringVar(&target,"d","","target peer to dial")
	//
	//flag.Parse()
	//
	//if target == "" {
	//	log.Println("please enter bootstrap node!")
	//	return
	//}


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()




	//bootstrapPeers = getLocalPeerInfo()

	r := mrand.New(mrand.NewSource(int64(0)))

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		//libp2p.Transport(ws.New),
	)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	security := libp2p.Security(secio.ID, secio.New)

	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/4002",
		//"/ip4/0.0.0.0/tcp/0/ws",
	)

	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(ctx, h)
		return dht, err
	}
	routing := libp2p.Routing(newDHT)

	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		libp2p.Identity(prvKey),
		muxers,
		security,
		routing,
		libp2p.NATPortMap(),
		libp2p.DefaultEnableRelay,
		libp2p.DefaultPeerstore,
		//libp2p.AddrsFactory(newAddrsFactory(bootstrapAddrs)),
	)



	if err != nil {
		panic(err)
	}

	//ps, err := pubsub.NewGossipSub(ctx, host)
	//if err != nil {
	//	panic(err)
	//}
	//sub, err := ps.Subscribe(pubsubTopic)
	//if err != nil {
	//	panic(err)
	//}
	//go pubsubHandler(ctx, sub)


	fmt.Printf("addr: %s\n", host.ID())
	for _, addr := range host.Addrs() {
		fmt.Println("Listening on", addr)
	}
	//targetAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/144.34.183.16/tcp/4001/p2p/%s",target))
	//if err != nil {
	//	panic(err)
	//}
	//
	//targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = host.Connect(ctx, *targetInfo)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Connected to", targetInfo.ID)

	mdns, err := discovery.NewMdnsService(ctx, host, time.Second*10, "")
	if err != nil {
		panic(err)
	}
	mdns.RegisterNotifee(&mdnsNotifee{h: host, ctx: ctx})

	err = dht.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	donec := make(chan struct{}, 1)
	//go chatInputLoop(ctx, host, ps, donec)

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
