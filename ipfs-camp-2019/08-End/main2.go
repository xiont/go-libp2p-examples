package main
//
//import (
//	"context"
//	"crypto/rand"
//	"fmt"
//	circuit "github.com/libp2p/go-libp2p-circuit"
//
//	//autonat "github.com/libp2p/go-libp2p-autonat"
//	"github.com/multiformats/go-multiaddr"
//
//	//"github.com/libp2p/go-libp2p"
//	//"github.com/libp2p/go-libp2p-core/crypto"
//	//"github.com/libp2p/go-libp2p-core/host"
//	//"github.com/libp2p/go-libp2p-core/network"
//	//"github.com/libp2p/go-libp2p-core/peer"
//
//	ds "github.com/ipfs/go-datastore"
//	dsync "github.com/ipfs/go-datastore/sync"
//	"github.com/libp2p/go-libp2p"
//	"github.com/libp2p/go-libp2p-core/host"
//	"github.com/libp2p/go-libp2p-core/peer"
//	"github.com/libp2p/go-libp2p-core/crypto"
//	dht "github.com/libp2p/go-libp2p-kad-dht"
//	mplex "github.com/libp2p/go-libp2p-mplex"
//	yamux "github.com/libp2p/go-libp2p-yamux"
//
//	//mplex "github.com/libp2p/go-libp2p-mplex"
//	pubsub "github.com/libp2p/go-libp2p-pubsub"
//	//secio "github.com/libp2p/go-libp2p-secio"
//	//yamux "github.com/libp2p/go-libp2p-yamux"
//	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
//	tcp "github.com/libp2p/go-tcp-transport"
//	ws "github.com/libp2p/go-ws-transport"
//	//ma "github.com/multiformats/go-multiaddr"
//	"io"
//	"log"
//	//mrand "math/rand"
//	"os"
//	"os/signal"
//	"syscall"
//)
//
//type mdnsNotifee struct {
//	h   host.Host
//	ctx context.Context
//}
//
//func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
//	m.h.Connect(m.ctx, pi)
//}
//
//// makeRoutedHost creates a LibP2P host with a random peer ID listening on the
//// given multiaddress. It will use secio if secio is true. It will bootstrap using the
//// provided PeerInfo
//func makeRoutedHost( bootstrapPeers []peer.AddrInfo) (host.Host, error) {
//
//	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
//	// deterministic randomness source to make generated keys stay the same
//	// across multiple runs
//	var r io.Reader
//
//	r = rand.Reader
//
//
//	// Generate a key pair for this host. We will use it at least
//	// to obtain a valid host ID.
//	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
//	if err != nil {
//		return nil, err
//	}
//
//	transports := libp2p.ChainOptions(
//		libp2p.Transport(tcp.NewTCPTransport),
//		libp2p.Transport(ws.New),
//	)
//
//	listenAddrs := libp2p.ListenAddrStrings(
//		"/ip4/0.0.0.0/tcp/0",
//		"/ip4/0.0.0.0/tcp/0/ws",
//	)
//
//	muxers := libp2p.ChainOptions(
//		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
//		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
//	)
//	opts := []libp2p.Option{
//		listenAddrs,
//		//libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", 9999)),
//		libp2p.Identity(priv),
//		transports,
//		//libp2p.DefaultTransports,
//		muxers,
//		//libp2p.DefaultMuxers,
//		libp2p.DefaultSecurity,
//		libp2p.NATPortMap(),
//		libp2p.EnableRelay(circuit.OptHop),
//	}
//
//	ctx := context.Background()
//
//	basicHost, err := libp2p.New(ctx, opts...)
//	if err != nil {
//		return nil, err
//	}
//
//	// Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
//	dstore := dsync.MutexWrap(ds.NewMapDatastore())
//
//	// Make the DHT
//	dht := dht.NewDHT(ctx, basicHost, dstore)
//
//	// Make the routed host
//	routedHost := rhost.Wrap(basicHost, dht)
//
//	// connect to the chosen ipfs nodes
//	err = bootstrapConnect(ctx, routedHost, bootstrapPeers)
//	if err != nil {
//		return nil, err
//	}
//
//	// Bootstrap the host
//	err = dht.Bootstrap(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	// Build host multiaddress
//	//hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))
//
//	// Now we can build a full multiaddress to reach this host
//	// by encapsulating both addresses:
//	// addr := routedHost.Addrs()[0]
//	//addrs := routedHost.Addrs()
//	//log.Println("I can be reached at:")
//	//for _, addr := range addrs {
//	//	log.Println(addr.Encapsulate(hostAddr))
//	//}
//
//	return routedHost, nil
//}
//
//func main() {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	var bootstrapPeers []peer.AddrInfo
//	//bootstrapPeers = IPFS_PEERS
//	bootstrapPeers = getLocalPeerInfo()
//
//	routedHost, err := makeRoutedHost(bootstrapPeers)
//
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//
//	for _, addr := range routedHost.Addrs() {
//		fmt.Println("Listening on", addr)
//	}
//
//	ps, err := pubsub.NewGossipSub(ctx, routedHost)
//	if err != nil {
//		panic(err)
//	}
//	sub, err := ps.Subscribe(pubsubTopic)
//	if err != nil {
//		panic(err)
//	}
//	go pubsubHandler(ctx, sub)
//
//	targetAddr, err := multiaddr.NewMultiaddr("/ip4/192.168.3.17/tcp/10000/p2p/QmSH3Ajg3JtgYyQ1XSz55o4VqZfqFxK8YpZsFVfFTDyt2V")
//	if err != nil {
//		panic(err)
//	}
//
//	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
//	if err != nil {
//		panic(err)
//	}
//
//	err = routedHost.Connect(ctx, *targetInfo)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println("Connected to", targetInfo.ID)
//
//	//mdns, err := discovery.NewMdnsService(ctx, routedHost, time.Second*10, "")
//	//if err != nil {
//	//	panic(err)
//	//}
//	//mdns.RegisterNotifee(&mdnsNotifee{h: routedHost, ctx: ctx})
//
//	//err = dht.Bootstrap(ctx)
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	donec := make(chan struct{}, 1)
//	go chatInputLoop(ctx, routedHost, ps, donec)
//
//	stop := make(chan os.Signal, 1)
//	signal.Notify(stop, syscall.SIGINT)
//
//	select {
//	case <-stop:
//		routedHost.Close()
//		os.Exit(0)
//	case <-donec:
//		routedHost.Close()
//	}
//}
//
//
//
