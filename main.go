package main

import (
	"peer-to-peer-app-v0.01/p2pnet"
	"peer-to-peer-app-v0.01/pubsub"
	"peer-to-peer-app-v0.01/pubsub/msghandle"
)

const topic = "rex/filegroup/1"
const service = "rex/service/test"

func main() {
	ctx, host := p2pnet.EstablishP2P()
	kad_dht := p2pnet.HandleDHT(ctx, host)
	sub, topic := pubsub.HandlePubSub(ctx, host, topic)
	go p2pnet.DiscoverPeers(ctx, host, kad_dht, service)
	msghandle.HandlePubSubMessages(ctx, host, sub, topic)

}
