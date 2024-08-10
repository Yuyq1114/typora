### WebSocket 的基本原理

1. **建立连接**：
   - WebSocket 连接通过一个标准的HTTP请求开始。这是一个握手过程，客户端发送一个包含特殊头部的HTTP请求到服务器，以请求升级连接到WebSocket协议。
2. **协议升级**：
   - 服务器收到请求后，如果支持WebSocket，将返回一个101状态码，表示协议切换成功，然后连接升级为WebSocket协议。
3. **数据帧传输**：
   - 建立连接后，客户端和服务器可以互相发送数据帧。这些帧可以包含文本数据或二进制数据。
   - WebSocket 使用帧（frame）的形式传输数据，每个帧都包含一个标志位，用于指示消息的开始和结束。
4. **关闭连接**：
   - 任何一方都可以发送关闭帧来关闭连接。关闭帧可以包含一个状态码和关闭原因。

### WebSocket 的特点

- **全双工通信**：客户端和服务器都可以同时发送和接收消息。
- **低延迟**：由于连接是持久的，消除了传统HTTP请求的握手延迟。
- **高效**：相比HTTP请求，WebSocket帧的开销更小。
- **实时性**：适合需要实时更新的应用，如在线聊天、实时数据流、在线游戏等。

### WebSocket 和 HTTP 的对比

- **连接方式**：
  - HTTP：基于请求-响应模式，每次请求都会建立一个新的连接。
  - WebSocket：基于持久连接，一旦建立连接，通信可以持续进行。
- **数据传输**：
  - HTTP：每次请求都包含完整的头部信息，数据传输效率较低。
  - WebSocket：使用帧传输数据，头部开销小，数据传输效率高。
- **通信模式**：
  - HTTP：单向通信，客户端发起请求，服务器响应。
  - WebSocket：双向通信，客户端和服务器可以随时发送和接收消息。

### 握手过程

1. **客户端请求**：

   - 客户端发起一个标准的HTTP请求，带有升级到WebSocket协议的请求头。

   ```
   GET /chat HTTP/1.1
   Host: example.com:8000
   Upgrade: websocket
   Connection: Upgrade
   Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
   Sec-WebSocket-Version: 13
   ```

   关键头部字段：

   - `Upgrade: websocket`：表示请求升级到WebSocket协议。
   - `Connection: Upgrade`：表示需要升级连接。
   - `Sec-WebSocket-Key`：一个Base64编码的随机值，用于服务器验证。
   - `Sec-WebSocket-Version`：协议版本，当前为13。

2. **服务器响应**：

   - 服务器返回一个101状态码的HTTP响应，表示协议切换成功。

   ```
   HTTP/1.1 101 Switching Protocols
   Upgrade: websocket
   Connection: Upgrade
   Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
   ```

   关键头部字段：

   - `Sec-WebSocket-Accept`：通过 `Sec-WebSocket-Key` 计算得出的值，用于验证握手请求的合法性。

### 数据帧格式

WebSocket 数据以帧（frame）的形式传输。每个帧都有一个固定的头部结构，后面跟随可变长度的数据。数据帧的基本格式如下：

```
0               1               2               3               
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| OpCode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |           (16/64)             |
|N|V|V|V|       |S|             | (if payload len==126/127)     |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                     Payload Data continued ...                |
+---------------------------------------------------------------+
```

### 详细字段解释

1. **FIN（1 bit）**：
   - 表示是否为消息的最后一帧。1表示这是消息的最后一帧，0表示还有后续帧。
2. **RSV1, RSV2, RSV3（各1 bit）**：
   - 保留位，通常为0，除非双方协商使用某个扩展。
3. **OpCode（4 bits）**：
   - 表示帧类型：
     - 0x0：延续帧（continuation frame）
     - 0x1：文本帧（text frame）
     - 0x2：二进制帧（binary frame）
     - 0x8：关闭连接帧（connection close frame）
     - 0x9：ping帧（ping frame）
     - 0xA：pong帧（pong frame）
     - 其它值为保留。
4. **Mask（1 bit）**：
   - 表示是否对Payload Data进行掩码处理。客户端到服务器的帧必须设置为1，服务器到客户端的帧通常为0。
5. **Payload length（7 bits, 7+16 bits, or 7+64 bits）**：
   - 表示负载数据的长度：
     - 0-125：表示实际数据长度（0-125字节）。
     - 126：接下来的2个字节表示实际数据长度（16位无符号整数）。
     - 127：接下来的8个字节表示实际数据长度（64位无符号整数）。
6. **Masking-key（0 or 4 bytes）**：
   - 如果Mask位为1，这里有4个字节的掩码键，用于对负载数据进行掩码处理。
7. **Payload Data（x+y bytes）**：
   - 实际的负载数据。长度由Payload length字段确定。如果Mask位为1，数据需要使用掩码键进行解码。

### 控制帧

- **关闭连接帧（Close Frame）**：
  - 用于关闭WebSocket连接。包含状态码和可选的关闭原因。
- **Ping帧**：
  - 用于检查连接的可用性。接收到Ping帧时，必须回复一个Pong帧。
- **Pong帧**：
  - 是对Ping帧的响应。可以包含Ping帧中的负载数据。

### WebSocket 数据掩码

客户端发送给服务器的数据帧必须进行掩码处理。掩码处理使用一个32位的掩码键进行。掩码处理的方式如下：

```
transformed-octet-i = original-octet-i XOR masking-key-octet-(i MOD 4)
```

这种掩码处理可以防止代理缓存和其他潜在的安全问题。





### 简单示例

#### 服务器端（Go语言）

使用`github.com/gorilla/websocket`包来实现WebSocket服务器。

```
package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Client Connected")

    go func() {
        for {
            // 定时发送消息
            if err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from server")); err != nil {
                fmt.Println(err)
                return
            }
        }
    }()

    for {
        // 读取消息
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            fmt.Println(err)
            return
        }
        // 将收到的消息回送给客户端
        if err := conn.WriteMessage(messageType, p); err != nil {
            fmt.Println(err)
            return
        }
    }
}

func setupRoutes() {
    http.HandleFunc("/ws", wsEndpoint)
}

func main() {
    fmt.Println("Starting WebSocket server on :8080")
    setupRoutes()
    http.ListenAndServe(":8080", nil)
}
```

#### 客户端（HTML + JavaScript）

使用HTML和JavaScript与WebSocket服务器通信并动态显示消息。

```
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WebSocket Client</title>
</head>
<body>
    <h1>WebSocket Client</h1>
    <div id="messages"></div>
    <script>
        let socket = new WebSocket("ws://localhost:8080/ws");

        socket.onopen = function(event) {
            console.log("Connected to WebSocket server");
        };

        socket.onmessage = function(event) {
            let messagesDiv = document.getElementById("messages");
            let newMessage = document.createElement("div");
            newMessage.textContent = event.data;
            messagesDiv.appendChild(newMessage);
        };

        socket.onclose = function(event) {
            console.log("Disconnected from WebSocket server");
        };

        socket.onerror = function(error) {
            console.log("WebSocket error:", error);
        };
    </script>
</body>
</html>
```