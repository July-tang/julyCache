# julyCache
  用Go实现的简易分布式缓存。
  
  1. 缓存淘汰策略采用lru算法。代码实现了lfu可供参考。
  2. 实现了一致性哈希算法，哈希算法默认为crc算法。
  3. 节点间采用http进行通信，采用protobuf协议传输。
  4. 实现singlefight防止缓存击穿。
