# 分布式ID生成器服务

基于[twitter snowflake](https://blog.twitter.com/2010/announcing-snowflake)的算法实现的ID生成器高性能的健壮的可容错的网络服务。

它提供了:

- 可配置的节点bit数和最大序列数
- 可以**批量**获取ID
- 基于[rpcx](https://rpcx.io)提供网络服务，可以提供一组服务节点
- 基于rpcx,可以提供分布式的、容错的网络服务



snowflake算法的实现基于[bwmarrin/snowflake](https://github.com/bwmarrin/snowflake), 额外提供了批量获取ID的方法。


## 例子

### server

[server](https://github.com/rpcxio/did/tree/master/cmd/server) 提供了一个简单的服务节点，你可以提供多个服务节点以便提供高可用性。

你还可以配置插件，使用ZooKeeper、Etcd、Consul、Nacos等作为注册中心，配置插件进行trace监控，设置黑白名单等操作，这些微服务的特性通过[rpcx](https://rpcx.io)来实现。


[client](https://github.com/rpcxio/did/tree/master/cmd/client)实现远程调用ID生成器的服务，这个例子采用服务地址硬编码的方式实现。如果你的服务使用了etcd、nacos等服务中心，你可以配置client使用注册中心自动获取服务节点。

客户端的例子演示了获取单个ID和获取批量ID的方法。

## 性能

单个节点每秒可以产生**400万**左右的ID。

```sh
BenchmarkGenerate-4                 	 4918761	       244 ns/op	   4097826 ids/s	     115 B/op	       1 allocs/op
BenchmarkGenerateBatch-4            	   49309	     24472 ns/op	   4086058 ids/s	   10811 B/op	     104 allocs/op
BenchmarkGenerate_Parallel-4        	 4903029	       244 ns/op	   4093939 ids/s	      88 B/op	       0 allocs/op
BenchmarkGenerateBatch_Parallel-4   	   49328	     24402 ns/op	   4097864 ids/s	   12904 B/op	     126 allocs/op
```

### 使用网络服务获取ID

**1**、256个client并发，每次只获取1个ID, ID的产生速度是 **12万/秒**。

```sh
./bclient -addr 10.41.15.225:8972 -n 100000
total IDs: 25600000, duration: 3m31.581592489s, id/s: 120993
```

**2**、如果采用批量获取，尽量1减少网络消耗，256个client并发，每次只获取100个ID, ID的产生速度是 **297万/秒**。
```sh
./bclient -addr 10.41.15.225:8972 -n 1000000 -b 100
total IDs: 256000000, duration: 1m26.178942509s, id/s: 2970563
```