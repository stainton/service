🔹 表 & 链（最重要）
|表|作用|必须|
|:---|:---|:---|
|nat|重定向|✅|
|filter	|防火墙	|❌|
|mangle	|打标	|⚠️|

|链	|用途|
|:---|:---|
|OUTPUT	|出站透明代理|
|PREROUTING	|入站透明代理|


🔹 match（你现在就该熟）

|flag|	为什么|
|:---|:---|
|-p tcp	|mesh 只管 TCP|
|--dport / --sport	|排除端口|
|-d CIDR	|Service IP|
|-m owner --uid-owner	|核心中的核心|
|-m mark	|高级|
|-m conntrack --ctstate	|入站必须|


🔹 target（你必须会）

|target	|用途|
|:---|:---|
|REDIRECT	|出站劫持|
|DNAT	|入站|
|RETURN	|放行|
|MARK	|标记|