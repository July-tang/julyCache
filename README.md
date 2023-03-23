# julyCache
用Go实现的简易分布式缓存, 采用lru算法, 实现了一致性哈希算法，哈希算法默认为crc算法，节点间采用http进行通信，使用protobuf传输。
