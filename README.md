# Godis-基于golang实现的redis服务器

Godis 是一个用 Go 语言实现的 Redis 服务器。

主要特性：

1. 数据结构支持：Godis支持多种数据结构，包括字符串（string）、列表（list）、哈希（hash）、集合（set）和有序集合（sorted set）

2. 自动过期功能(TTL)：Godis支持设置键的过期时间，通过设置键的生存时间（TTL），可以让键在一段时间后自动过期并被删除。这对于缓存和临时数据非常有用。

3. AOF持久化及AOF重写：Godis支持AOF（Append-Only File）持久化，将写命令追加到磁盘上的AOF文件中，以便在服务器重启后恢复数据。此外，Godis还支持AOF重写，可以对AOF文件进行压缩和优化，减小文件大小并提高性能。

4. 内置集群模式：Godis内置了集群模式，可以通过配置启动多个Godis实例，并使用一致性哈希算法将不同的命令路由到不同的节点。这样可以实现数据的分布式存储和负载均衡，提高系统的可扩展性和容错性。

5. 分布式命令执行：在集群模式下，Godis支持的一些命令（如exists、type、set、setnx、get、getset、ping、rename、renamenx、flushdb、del、Select）将分布在不同的节点上执行，以提高并发性能和减轻单个节点的负载。