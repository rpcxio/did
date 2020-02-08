# 分布式ID生成器服务

基于[twitter snowflake](https://blog.twitter.com/2010/announcing-snowflake)的算法实现的ID生成器高性能的健壮的可容错的网络服务。

它提供了:

- 可配置的节点bit数和最大序列数
- 可以**批量**获取ID
- 基于[rpcx](https://rpcx.io)提供网络服务，可以提供一组服务节点
- 基于rpcx,可以提供分布式的、容错的网络服务



snowflake算法的实现基于[bwmarrin/snowflake](https://github.com/bwmarrin/snowflake), 额外提供了批量获取ID的方法。


## 例子

[]()