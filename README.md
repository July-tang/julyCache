# julyCache
  用Go实现的简易分布式缓存，原项目地址：https://github.com/geektutu/7days-golang。
  
  支持的特性有：
  1. 缓存淘汰策略采用lru算法。
  2. 使用一致性哈希选择节点，实现负载均衡，哈希算法默认为crc算法。
  3. 节点间采用http进行通信，采用protobuf协议传输。
  4. 实现singlefight防止缓存击穿。
  
  本人修改如下：
  - 新增缓存淘汰策略lfu可供参考。
  - 实现缓存过期时间设置，并实验过期淘汰策略为定时删除+惰性删除。
