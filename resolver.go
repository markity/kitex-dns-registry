package registry

import (
	"context"
	"errors"
	"net"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

type dnsResolver struct {
}

// error is always nil
func NewDnsResolver() (discovery.Resolver, error) {
	return &dnsResolver{}, nil
}

func (r *dnsResolver) Target(ctx context.Context, target rpcinfo.EndpointInfo) (description string) {
	return target.ServiceName()
}

func (r *dnsResolver) Resolve(ctx context.Context, desc string) (discovery.Result, error) {
	result := discovery.Result{}

	names, err := net.LookupAddr(desc)
	if err != nil {
		return result, errors.New("no instcance remains for: " + desc)
	}
	for _, v := range names {
		ins := discovery.NewInstance("tcp", v+":8888", discovery.DefaultWeight, nil)
		result.Instances = append(result.Instances, ins)
	}

	if len(result.Instances) == 0 {
		return result, errors.New("no instance remains for: " + desc)
	}

	result.CacheKey = desc
	result.Cacheable = true
	return result, nil
}

func (r *dnsResolver) Diff(cacheKey string, prev, next discovery.Result) (discovery.Change, bool) {
	return discovery.DefaultDiff(cacheKey, prev, next)
}

// Name implements the Resolver interface.
func (r *dnsResolver) Name() string {
	return "dns"
}
