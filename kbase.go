package kbase

import (
	"github.com/gl-works/kbase/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Kbase struct {
	lock     sync.RWMutex
	cluster  *cluster
	hostport string
	clients  map[string]rpc.KeyServiceClient
	kspaces  map[string]*Kspace
}

type Kspace struct {
	kbase      *Kbase
	name       string
	lockkvs    sync.RWMutex
	lockoffers sync.RWMutex
	kvs        map[string]*authValue //locally stored keys
	offers     map[string]map[string][]byte

	serializer   Serializer
	deserializer Deserializer
	chanScrape   chan interface{}
	chanSharding chan *sharding
}

type Key string
type Value interface{}
type Serializer func(Value) ([]byte, error)
type Deserializer func([]byte) (Value, error)
type KspaceDef struct {
	Name         string
	Serializer   Serializer
	Deserializer Deserializer
}

const cluster_consistent_safe_gap_in_minute = 5

func NewKbase(hostport string, seednodes []string, spaceDefs ...*KspaceDef) (kbase *Kbase) {
	kbase = &Kbase{
		hostport: hostport,
		clients:  make(map[string]rpc.KeyServiceClient),
		kspaces:  make(map[string]*Kspace, len(spaceDefs)),
	}
	if socket, err := net.Listen("tcp", rpcAddress(hostport)); err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		s := grpc.NewServer( /*grpc.MaxConcurrentStreams(10000)*/ )
		rpc.RegisterKeyServiceServer(s, kbase)
		go s.Serve(socket)
	}
	kbase.cluster = createCluster(hostport, seednodes, kbase)
	for _, def := range spaceDefs {
		kbase.kspaces[def.Name] = createSpace(kbase, def)
	}
	return
}

func (this *Kbase) SpaceOf(space string) *Kspace {
	return this.kspaces[space]
}

func (this *Kbase) ClusterSize() int {
	if n, err := this.cluster.ringpop.CountReachableMembers(); err != nil {
		return -1
	} else {
		return n
	}
}

// 发布key-value
func (this *Kspace) Offer(key string, value Value) (err error) {

	authority := this.kbase.authorityOf(key)

	this.lockkvs.Lock()
	this.kvs[key] = &authValue{value: value, authority: authority}
	this.lockkvs.Unlock()

	if bytes, err := this.serializer(value); err == nil {
		// 判断key对应的权威服务器
		if authority == this.kbase.hostport { // 本机为key的权威服务器
			this.putOffer(key, this.kbase.hostport, bytes)
		} else { // Offer 到远端权威服务器
			if _, err := this.kbase.rpcClient(authority).HandleOffer(context.Background(),
				&rpc.Offer{Space: this.name, Key: key, Server: this.kbase.hostport, Value: bytes}); err != nil {
				log.Errorf("%s", err)
			}
		}
		return nil
	} else {
		log.Errorf("serialize error %s", err)
		return err
	}
}

// 撤销key
func (this *Kspace) Revoke(key string) {
	authority := this.kbase.authorityOf(key)

	this.lockkvs.Lock()
	delete(this.kvs, key)
	this.lockkvs.Unlock()

	// 判断key对应的权威服务器
	if authority == this.kbase.hostport { // 本机为key的权威服务器
		this.dropOffer(key, this.kbase.hostport)
	} else { // Revoke 到远端权威服务器
		if _, err := this.kbase.rpcClient(authority).HandleRevoke(context.Background(),
			&rpc.Revoke{Space: this.name, Key: key, Server: this.kbase.hostport}); err != nil {
			log.Errorf("%s", err)
		}
	}
}

// 查询key
func (this *Kspace) Poll(key string) (results map[string]Value) {
	// 判断key对应的权威服务器
	var values map[string][]byte
	if server := this.kbase.authorityOf(key); server == this.kbase.hostport {
		// 本机为key的权威服务器
		this.lockoffers.RLock()
		values = this.offers[key]
		this.lockoffers.RUnlock()
	} else {
		// 向远端服务器发出 Poll
		if r, err := this.kbase.rpcClient(server).HandlePoll(
			context.Background(), &rpc.Poll{Space: this.name, Key: key}); err != nil {
			log.Errorf("poll rpc error: %s", err)
		} else {
			values = r.Values
		}
		if values == nil {
			values = emptyvalues
		}
	}

	results = make(map[string]Value, len(values))
	for key, bytes := range values {
		if v, err := this.deserializer(bytes); err != nil {
			log.Errorf("deserializ error %s", err)
		} else {
			results[key] = v
		}
	}
	return
}

func (this *Kspace) GetLocal(key string) (result Value) {
	this.lockkvs.RLock()
	if t, ok := this.kvs[key]; ok {
		result = t.value
	}
	this.lockkvs.RUnlock()
	return
}

var emptyvalues map[string][]byte = make(map[string][]byte)
