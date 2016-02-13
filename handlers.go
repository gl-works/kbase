package kbase

import (
	"./rpc"
	"github.com/vektra/errors"
	"golang.org/x/net/context"
)

var (
	emptyresponse *rpc.KeyResponse = &rpc.KeyResponse{}

	ErrMissingKspace = errors.New("ErrMissingKspace")
)

func (this *Kbase) HandleOffer(ctx context.Context, offer *rpc.Offer) (response *rpc.KeyResponse, err error) {
	if kspace, ok := this.kspaces[offer.Space]; ok {
		kspace.putOffer(offer.Key, offer.Server, offer.Value)
		response = emptyresponse
	} else {
		err = ErrMissingKspace
	}
	return
}

func (this *Kbase) HandlePoll(ctx context.Context, poll *rpc.Poll) (response *rpc.KeyResponse, err error) {
	if kspace, ok := this.kspaces[poll.Space]; ok {
		kspace.lockoffers.RLock()
		values := kspace.offers[poll.Key]
		kspace.lockoffers.RUnlock()
		response = &rpc.KeyResponse{Values: values}
	} else {
		err = ErrMissingKspace
	}
	return
}

func (this *Kbase) HandleRevoke(ctx context.Context, revoke *rpc.Revoke) (response *rpc.KeyResponse, err error) {
	if kspace, ok := this.kspaces[revoke.Space]; ok {
		kspace.dropOffer(revoke.Key, revoke.Server)
		response = emptyresponse
	} else {
		err = ErrMissingKspace
	}
	return
}
