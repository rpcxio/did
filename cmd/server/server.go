package main

import (
	"flag"

	"github.com/rpcxio/did/snowflake"
	"github.com/smallnest/rpcx/server"
)

var (
	serverID = flag.Int64("serverID", 1, "unique server ID")
	epoch    = flag.Int64("epoch", 1580601600000, "epoch time for base timeline")
	nodeBits = flag.Int64("nodeBits", 8, "the number of bits to use for Node")
	addr     = flag.String("addr", ":8972", "server listened address")
)

// 这是一个无注册中心的ID生成器的服务,客户端通过IP直接访问这些服务节点.
// 你可以配置etcd/zookeeper/consul/nacos作为注册中心，实现高可用的服务.
func main() {
	flag.Parse()

	snowFlake := snowflake.NewSnowFlake(*serverID, *epoch, uint8(*nodeBits), 22-uint8(*nodeBits))

	s := server.NewServer()

	// 配置其它注册中心

	// 配置插件，比如trace等

	// 注册SnowFlake服务
	s.RegisterName("snowflake", snowFlake, "")

	s.Serve("tcp", *addr)
}
