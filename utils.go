package kbase

import (
	"github.com/gl-works/kbase/rpc"
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"
	"strings"
	"time"
)

var (
	log *logrus.Logger = logrus.New()
)

type authValue struct {
	authority string
	value     Value
}

type sharding struct {
	key       string
	value     *authValue
	authority string
}

func (this *Kspace) shardingDaemon() {
	for sharding := range this.chanSharding {
		if bytes, err := this.serializer(sharding.value.value); err == nil {
			if _, err := this.kbase.rpcClient(sharding.authority).HandleOffer(context.Background(),
				&rpc.Offer{Key: sharding.key, Server: this.kbase.hostport, Value: bytes}); err != nil {
				log.Errorf("reshard key %s has failed (%s)", sharding.key, err)
			} else {
				sharding.value.authority = sharding.authority
			}
		} else {
			log.Errorf("serialize error %s", err)
		}
	}
}

func (this *Kbase) onClusterChanged() {
	for _, kspace := range this.kspaces {
		kspace.onClusterChanged()
	}
}

func (this *Kspace) onClusterChanged() {
	this.chanScrape <- "" //kickstart debounced scraping of offers

	this.lockkvs.RLock()
	for key, value := range this.kvs {
		authority, err := this.kbase.cluster.Lookup(key)
		if err != nil {
			log.Errorf("No authority for %s", err)
			break
		}
		if value.authority != authority {
			log.Infof("Authority changed for key %s from %s => %s", key, value.authority, authority)
			this.chanSharding <- &sharding{key: key, value: value, authority: authority}
		}
	}
	this.lockkvs.RUnlock()
}

func (this *Kspace) scrapeOffers(arg interface{}) {
	this.lockoffers.Lock()
	for key, _ := range this.offers {
		if authority, err := this.kbase.cluster.Lookup(key); err != nil && authority != this.kbase.hostport {
			log.Infof("Dropping authority on key %s. new authority being %s", key, authority)
			delete(this.offers, key)
		}
	}
	this.lockoffers.Unlock()
}

func (this *Kbase) authorityOf(key string) (server string) {
	server, err := this.cluster.Lookup(key)
	if err != nil {
		log.Errorf("No authority for %s. fallback locally", key)
		server = this.hostport
	}
	return server
}

func (this *Kbase) rpcClient(server string) (client rpc.KeyServiceClient) {
	this.lock.Lock()
	client, exists := this.clients[server]
	if !exists {
		conn, err := grpc.Dial(rpcAddress(server), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		client = rpc.NewKeyServiceClient(conn)
		this.clients[server] = client
	}
	this.lock.Unlock()
	return
}

func rpcAddress(server string) string {
	ss := strings.Split(server, ":")
	x, _ := strconv.Atoi(ss[1])
	return ss[0] + ":" + strconv.Itoa(x+1)
}

func debounce(interval time.Duration, debouncedFunction func(arg interface{})) chan interface{} {
	channel := make(chan interface{})
	backgroundFunction := func() {
		var (
			item    interface{}
			timeout <-chan time.Time
		)
		for {
			select {
			case item = <-channel:
				timeout = time.After(interval)
			case <-timeout:
				debouncedFunction(item)
			}
		}
	}
	go backgroundFunction()
	return channel
}

func createSpace(kbase *Kbase, def *KspaceDef) (kspace *Kspace) {
	kspace = &Kspace{
		name:         def.Name,
		kbase:        kbase,
		offers:       make(map[string]map[string][]byte),
		serializer:   def.Serializer,
		deserializer: def.Deserializer,
		kvs:          make(map[string]*authValue),
		chanSharding: make(chan *sharding, 1000), //TODO configurable
	}
	kspace.chanScrape = debounce(time.Minute*cluster_consistent_safe_gap_in_minute, kspace.scrapeOffers)
	go kspace.shardingDaemon()
	return
}

func (this *Kspace) putOffer(key string, server string, value []byte) {
	this.lockoffers.RLock()
	m := this.offers[key]
	this.lockoffers.RUnlock()
	if m == nil {
		m = make(map[string][]byte)
		this.lockoffers.Lock()
		this.offers[key] = m
		this.lockoffers.Unlock()
	}
	m[server] = value
}

func (this *Kspace) dropOffer(key string, server string) {
	this.lockoffers.Lock()
	if m := this.offers[key]; m != nil {
		delete(m, server)
		if len(m) == 0 {
			delete(this.offers, key)
		}
	}
	this.lockoffers.Unlock()
}
