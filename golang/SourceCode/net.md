# 客户端结构

src\net\http\client.go

```
type Client struct {

    Transport RoundTripper//Transport 是一个接口，用于执行 HTTP 请求。RoundTripper 接口定义了一个方法 RoundTrip(*Request) (*Response, error)，用于执行 HTTP 请求并返回响应。Transport 可以是自定义的，也可以是标准库提供的，如 http.DefaultTransport。通过设置 Transport，你可以控制请求的实际传输过程，比如使用自定义的代理或修改请求/响应。

    CheckRedirect func(req *Request, via []*Request) error//CheckRedirect 是一个回调函数，用于处理 HTTP 重定向。当客户端接收到一个重定向响应时（状态码 3xx），CheckRedirect 会被调用。它接收当前请求 (req) 和一个包含重定向链的请求切片 (via) 作为参数。如果返回 nil，客户端将跟随重定向。如果返回一个错误，重定向将被停止，错误将被返回。

    Jar CookieJar//Jar 是一个接口，用于处理 HTTP cookies。CookieJar 接口定义了方法 SetCookies(u *url.URL, cookies []*http.Cookie) 和 Cookies(u *url.URL) []*http.Cookie。通过 Jar，你可以控制如何存储和发送 cookies。例如，http.CookieJar 可以用于管理 cookies 的持久性和跨请求的共享。

    Timeout time.Duration//Timeout 指定了 HTTP 请求的超时时间。如果请求超过这个时间没有完成，它将被取消。time.Duration 是一个表示时间间隔的类型，通常使用 time.Second、time.Millisecond 等单位。设置 Timeout 可以帮助你避免长时间等待响应，从而提高程序的可靠性和响应速度。
}
```





# 服务端结构

src\net\http\server.go

```
type Server struct {
    
    Addr string

    Handler Handler // Handler 是一个接口，定义了如何处理传入的 HTTP 请求。Handler 接口有一个方法 ServeHTTP(ResponseWriter, *Request)。如果 Handler 为 nil，服务器将使用默认的多路复用器（http.DefaultServeMux）来处理请求。
    
    DisableGeneralOptionsHandler bool//字段控制是否禁用对 OPTIONS 请求方法的处理。如果设置为 true，服务器将不再处理 OPTIONS 请求。

    TLSConfig *tls.Config//TLSConfig 字段允许你为服务器配置 TLS（Transport Layer Security）。tls.Config 包含了各种配置选项，如证书、密钥、加密协议等，用于启用 HTTPS。

    ReadTimeout time.Duration//设置了读取请求体的超时时间。即服务器在等待客户端发送请求的过程中，最大允许的时间。超过这个时间，服务器将关闭连接。

    ReadHeaderTimeout time.Duration//设置了读取请求头的超时时间。服务器在解析请求头时最大允许的时

    WriteTimeout time.Duration//设置了写入响应体的超时时间。

    IdleTimeout time.Duration//设置了连接的空闲超时时间。

    MaxHeaderBytes int//设置了请求头的最大字节数。

    TLSNextProto map[string]func(*Server, *tls.Conn, Handler)//是一个映射，键是协议名，值是处理程序函数。用于处理 TLS 连接的下一个协议（如 HTTP/2）。如果服务器启用了 HTTPS，这个字段允许你注册不同的协议处理程序。

    ConnState func(net.Conn, ConnState)//是一个回调函数，用于通知服务器连接状态的变化。

    ErrorLog *log.Logger

    BaseContext func(net.Listener) context.Context//是一个回调函数，用于为监听器生成基础上下文。

    ConnContext func(ctx context.Context, c net.Conn) context.Context//是一个回调函数，用于在连接建立时为每个连接生成上下文。

    HTTP2 *HTTP2Config//配置 HTTP/2 相关的选项。

    inShutdown atomic.Bool // 是一个原子布尔值，表示服务器是否处于关闭状态。

    disableKeepAlives atomic.Bool//表示是否禁用连接的 Keep-Alive 功能
    nextProtoOnce     sync.Once // 用于确保 HTTP/2 的初始化代码只运行一次，避免重复初始化。
    nextProtoErr      error     // 保存 http2.ConfigureServer 函数的错误结果，

    mu         sync.Mutex
    listeners  map[*net.Listener]struct{}//是一个映射，用于跟踪所有当前正在监听的网络连接。
    activeConn map[*conn]struct{}//用于跟踪所有当前活跃的连接。
    onShutdown []func()//是一个函数切片，在服务器关闭时需要执行的回调函数列表。

    listenerGroup sync.WaitGroup//用于等待所有监听器的关闭，确保在服务器完全关闭之前，所有的监听器都已经正确关闭。
}
```



# 响应结构

src\net\http\response.go

```
type Response struct {
    Status     string // e.g. "200 OK"
    StatusCode int    // e.g. 200
    Proto      string // e.g. "HTTP/1.0"表示使用的协议
    ProtoMajor int    // e.g. 1协议的主版本号
    ProtoMinor int    // e.g. 0协议的次版本号
    Header Header//Header 是一个映射类型 map[string][]string，包含了响应的所有头部字段及其对应的值。

    Body io.ReadCloser// io.ReadCloser 接口的值，表示响应的主体内容。

    ContentLength int64

    TransferEncoding []string//表示用于传输内容的编码方式。例如，"chunked" 表示分块传输编码。

    Close bool

    Uncompressed bool

    Trailer Header//一个 Header 类型的值，表示响应的尾部头部字段。它用于在响应体的末尾发送额外的头部字段。

    Request *Request//指向 http.Request 结构体的指针，表示引发该响应的请求。它提供了请求的详细信息，如请求头、方法、URL 等。

    TLS *tls.ConnectionState//指向 tls.ConnectionState 结构体的指针，表示与响应相关的 TLS 连接状态。它提供了有关 TLS 连接的详细信息，如证书、加密协议等。
}
```



# 请求结构

src\net\http\request.go

```
type Request struct {
    
    Method string//GET POST

    URL *url.URL

    Proto      string // "HTTP/1.0"
    ProtoMajor int    // 1
    ProtoMinor int    // 0

    Header Header

    Body io.ReadCloser

    GetBody func() (io.ReadCloser, error)一个函数类型，用于在需要时获取请求的主体内容。

    ContentLength int64

    TransferEncoding []string

    Close bool

    Host string

    Form url.Values//一个 url.Values 类型的值，表示解析后的表单数据。

    PostForm url.Values//专门用于处理 POST 请求中的表单数据。

    MultipartForm *multipart.Form//表示解析后的 multipart/form-data 表单数据。这通常用于处理文件上传。

    Trailer Header//

    RemoteAddr string//示发起请求的客户端的网络地址，例如 "192.168.1.1:12345"。

    RequestURI string//表示请求的原始 URI 部分，即请求行中的 URI。它包含了未经过处理的 URI 字符串。

    TLS *tls.ConnectionState

    Cancel <-chan struct{}

    Response *Response

    Pattern string//表示匹配请求的路由模式。它通常用于路由和处理请求时的模式匹配。

    ctx context.Context

    pat         *pattern          // 表示与请求相关的路由模式。
    matches     []string          // 一个字符串切片，存储与路由模式中的通配符匹配的值。
    otherValues map[string]string // for calls to SetPathValue that don't match a wildcard
}
```