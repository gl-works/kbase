package kbase

import (
	"github.com/gl-works/ringpop-go"
	"github.com/gl-works/ringpop-go/events"
	"github.com/gl-works/ringpop-go/swim"
	"github.com/uber-common/bark"
	"github.com/uber/tchannel-go"
	//"hash/crc32"
	//"sort"
)

type clusterApp interface {
	onClusterChanged()
}

type cluster struct {
	//ring 	Ring
	ringpop *ringpop.Ringpop
	channel *tchannel.Channel
	app     clusterApp
}

func createCluster(hostport string, seednodes []string, app clusterApp) *cluster {
	channel, err := tchannel.NewChannel("kbase", nil)
	if err != nil {
		log.Fatalf("channel did not create successfully: %v", err)
	}

	ringpop, err := ringpop.New("kbase-app",
		ringpop.Channel(channel),
		ringpop.Identity(hostport),
		ringpop.Logger(bark.NewLoggerFromLogrus(log)),
	)
	if err != nil {
		log.Fatalf("Unable to create Ringpop: %v", err)
	}

	if err := channel.ListenAndServe(hostport); err != nil {
		log.Fatalf("could not listen on given hostport: %v", err)
	}

	opts := new(swim.BootstrapOptions)
	if len(seednodes) > 0 {
		opts.Hosts = seednodes
	}

	if _, err := ringpop.Bootstrap(opts); err != nil {
		log.Fatalf("ringpop bootstrap failed: %v", err)
	}

	cluster := &cluster{
		/*ring: Ring{
			hcodes: []int{},
			servers: make(map[int]string),
		},*/
		ringpop: ringpop,
		channel: channel,
		app:     app}

	ringpop.RegisterListener(cluster)

	return cluster
}

func (this *cluster) Lookup(key string) (string, error) {
	return this.ringpop.Lookup(key)
}

func (this *cluster) HandleEvent(e events.Event) {
	switch e.(type) {
	case events.RingChangedEvent:
		//this.handleRingChange(e)
		this.app.onClusterChanged()
	case swim.MemberlistChangesReceivedEvent:
		/*log.Printf("swim.MemberlistChangesReceivedEvent: %#v", e)
		a, b := this.ringpop.GetReachableMembers()
		log.Printf("__ring__ %#v, %s", a, b)*/
	}
}

/*
func (this *cluster) handleRingChange(e events.RingChangedEvent) {
	log.Debugf("events.RingChangedEvent: %#v", e)
	affected := []string{}
	for _, server := range e.ServersAdded {
		affected = append(affected, this.ring.Add(server)...)
	}
	for _, server := range e.ServersRemoved {
		affected = append(affected, this.ring.Remove(server)...)
	}
	this.app.onClusterChanged()
}


var hashFunc func(data []byte) uint32 = crc32.ChecksumIEEE

type Ring struct {
	hcodes  []int // Sorted
	servers map[int]string
}

func (this *Ring) Add(server string) (affected []string){
	const servers = this.servers
	const hash = int(hashFunc([]byte(server)))
	hcodes := this.hcodes
	affected = []string{}
	const N = len(hcodes)
	if idx := sort.Search(N, func(i int) bool { return hcodes[i] >= hash }); idx == N {
		this.hcodes = append(hcodes, hash)
		if N > 0 {
			affected = []string{servers[hcodes[0]]}
		}
	} else {
		if exists, ok := servers[hash]; ok{
			log.Errorf("server hashcode conflicts! %s vs %s", server, exists)
		} else {
			affected = []string{servers[hcodes[idx]]}
			t := append(hcodes[0:idx], hash)
			this.hcodes = append(t, hcodes[idx+1:])
			servers[hash] = server
		}
	}
	return
}

func (this *Ring) Remove(server string) (affected []string){
	const servers = this.servers
	const hash = int(hashFunc([]byte(server)))
	hcodes := this.hcodes
	affected = []string{server}
	const N = len(hcodes)
	if s, _ := servers[hash]; s == server {
		if idx := sort.Search(N, func(i int) bool { return hcodes[i] == hash }); idx == N {
			panic("bug")
		} else {
			hcodes = append(hcodes[0:idx], hcodes[idx+1:]...)
			if N := len(hcodes); N > 0 {
				if idx == N - 1 {
					affected = append(affected, servers[hcodes[0]])
				} else {
					affected = append(affected, servers[hcodes[idx]])
				}
			}
			this.hcodes = hcodes
		}
	} else {
		log.Errorf("Server %s not in ring", server)
	}
	return
}

// Gets the closest item in the hash to the provided key.
func (this *Ring) Lookup(key string) string {
	if N := len(this.hcodes); N == 0 {
		return ""
	} else {
		hash := int(hashFunc([]byte(key)))
		if idx := sort.Search(N, func(i int) bool { return this.hcodes[i] >= hash }); idx == N {
			return this.servers[this.hcodes[0]]
		} else {
			return this.servers[this.hcodes[idx]]
		}
	}
}
*/
