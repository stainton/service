# sidecar

## 问题记录与解决

### iptables 将 sidecar 转发的流量重新转发给 sidecar

iptables 按出口来转发流量，sidecar 的流量会重新转给 sidecar，导致死循环。

解决方法：通过 UID 来识别是否做转发——app要避免使用同一个用户

```bash
# 创建用户组和用户
RUN addgroup -g 1337 proxy && adduser -D -u 1337 -G proxy proxy
# 使用该用户组和用户组
USER 1337
```

kubernetes 的 API 对象中添加这些内容用于主要用于防止覆盖

```yaml
securityContext:
  runAsUser: 1337
  runAsGroup: 1337
  runAsNonRoot: true
```