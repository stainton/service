# 下一步，非“加功能”，而是加“维度”

目标路径为 `sidecar mesh -> ambinet mesh`，下一步要感知 L7 的语义。

# 优先引入“最小 L7 语义”，不是全 L7

一上来就做 HTTP/gRPC 会陷入协议细节泥潭。现在着重需要`让 proxy “知道一次请求在哪里开始、在哪里结束”`

要刻意训练的能力——设计边界：
* 哪些失败可以 retry
* 哪些失败必须 RST
* retry 次数 / backoff

真正理解`Envoy retry policy`是干什么的。

# 引入“控制面概念”

## 引入“弱控制面”

不需要复杂，哪怕一个 YAML 文件也行。控制面下发给 sidecar：
* upstream 实例列表
* weight
* 健康状态（哪怕定时 ping）

sidecar 开始区分：
```
dataplane ≠ controlplane
```

重试 ≠ 重连  
负载均衡 ≠ 随机挑 IP

# 刻意“去 sidecar 化”

## 把 sidecar 能力拆成两块

你现在的 sidecar 其实干了两件事：

    流量拦截（iptables / socket 层）
    连接 / 协议处理（proxy 逻辑）

Ambient mesh 的本质是：

每个 Pod 不再有 sidecar，但流量依然被代理，所以你下一步要问的不是：

    “sidecar 还能做什么”

而是：

    “这部分能力，能不能挪到节点级？”

## 你应该做的练习（非常关键）

在 不改 proxy 代码 的前提下：

把 proxy 从：

    per-pod

改成：

    per-node（DaemonSet / ztunnel 形态）

你会立刻撞上三个硬问题：

    多租户连接隔离
    源身份（这个是谁的流量）
    性能边界
