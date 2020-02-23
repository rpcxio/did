package snowflake

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"

	mc "github.com/rpcxio/gomemcached"
)

// MemcachedServer implements memcached protocol.
// You can use get,gets to get IDs.
type MemcachedServer struct {
	ms        *mc.Server
	ln        net.Listener
	snowFlake *SnowFlake
}

// NewMemcachedServer returns a new MemcachedServer.
func NewMemcachedServer(ln net.Listener, snowFlake *SnowFlake) *MemcachedServer {
	return &MemcachedServer{
		ln:        ln,
		snowFlake: snowFlake,
	}
}

// Serve serves memcached server.
func (ms *MemcachedServer) Serve() {
	ms.ms = mc.NewServer("")
	ms.ms.RegisterFunc("get", ms.get)
	ms.ms.RegisterFunc("gets", ms.gets)
	ms.ms.Serve(ms.ln)
}

// get id
func (ms *MemcachedServer) get(ctx context.Context, req *mc.Request, res *mc.Response) error {
	if len(req.Keys) != 1 || req.Keys[0] != "id" {
		return errors.New("key must be id")
	}

	id := strconv.FormatInt(ms.snowFlake.node.Generate(), 10)
	res.Values = append(res.Values, mc.Value{req.Keys[0], "0", []byte(id), ""})
	res.Response = mc.RespEnd
	return nil
}

// gets count
func (ms *MemcachedServer) gets(ctx context.Context, req *mc.Request, res *mc.Response) error {
	if len(req.Keys) != 1 {
		return errors.New("must set count")
	}

	if req.Keys[0] == "id" {
		return ms.get(ctx, req, res)
	}

	count, err := strconv.Atoi(req.Keys[0])
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("count must be greater than zero")
	}
	if count > 65535 {
		return errors.New("count must not be greater than 65535")
	}

	vs := make([]string, count)
	ints := ms.snowFlake.node.GenerateBatch(uint16(count))
	for i, v := range ints {
		vs[i] = strconv.FormatInt(v, 10)
	}
	res.Values = append(res.Values, mc.Value{req.Keys[0], "0", []byte(strings.Join(vs, ",")), ""})

	res.Response = mc.RespEnd
	return nil
}
