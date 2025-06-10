# 1.将本机的各种https流量通过kafka发送出去

## 1.导入CA证书

### **通过 mitmproxy 手动导出 CA 证书 (推荐)**

`mitmproxy` 自带命令可以直接生成证书文件，无需访问 `http://mitm.it/`。

### **步骤 1：启动 mitmproxy 并生成证书**

在终端中执行以下命令：

```
sh


复制编辑
mitmproxy --set block_global=false
```

- `--set block_global=false`：允许非本地连接（防止 mitmproxy 阻止外部请求）
- 启动后，`mitmproxy` 会自动在以下路径生成证书文件：

**Windows 路径**（默认）：

```
makefile


复制编辑
C:\Users\<你的用户名>\.mitmproxy\mitmproxy-ca-cert.pem
```

------

### **步骤 2：将证书导入 Windows 系统**

1. 打开 `certmgr.msc`（按 `Win + R`，输入 `certmgr.msc` 回车）
2. 点击左侧的 **“受信任的根证书颁发机构”**
3. 右键 **“证书”** → **“所有任务” → “导入...”**
4. 在向导中：
   - 点击 **“下一步”**
   - 点击 **“浏览...”** 并选择 `mitmproxy-ca-cert.pem`
   - 点击 **“下一步”**，确认存储位置为 **“受信任的根证书颁发机构”**
   - 点击 **“完成”**

## 2.使用mitmproxy代理

启动 mitmproxy：

```
mitmweb --mode reverse:http://example.com --listen-host 127.0.0.1 --listen-port 8080
mitmweb --listen-host 0.0.0.0 --listen-port 8080
mitmweb --mode upstream:http://127.0.0.1:7897 --listen-host 0.0.0.0 --listen-port 8080

mitmweb --listen-host 0.0.0.0 --listen-port 8080 --set stream_websocket=true

```

设置代理：

- Windows 设置 → 网络和 Internet → 代理 → 手动代理设置
- 服务器：`127.0.0.1`
- 端口：`8080`

访问 HTTPS 网站（如 `https://example.com`）

mitmproxy 控制台中应显示该网站的请求和响应

## 3.go监听代理

```
url := "wss://127.0.0.1:8081/?token=d1fb520ec81a32e4be3161088e8ac169"
注意go中的url使用wss才可以，而不是ws
```

## 4.使用python代理发送kafka

```
pip install mitmproxy kafka-python
pip install mitmproxy
```

