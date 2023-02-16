package main

import (
	"peer-to-peer-app-v0.01/p2pnet"
	"peer-to-peer-app-v0.01/pubsub"
)

const topic = "rex/filegroup/1"
const service = "rex/service/test"

func main() {
	ctx, host := p2pnet.EstablishP2P()
	kad_dht := p2pnet.HandleDHT(ctx, host)
	pubsub.HandlePubSub(ctx, host, topic)
	go p2pnet.DiscoverPeers(ctx, host, kad_dht, service)
	var i = 0
	for i = 10; ; i++ {
	}
}
