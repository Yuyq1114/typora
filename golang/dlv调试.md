# 本地安装go 安装dlv

# 远程安装go和dlv

## 安装go

wget https://go.dev/dl/go1.19.5.linux-amd64.tar.gz（可选）

tar -C /usr/local -xzf go1.22.linux-amd64.tar.gz

echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc

echo "export GOPATH=\$HOME/go" >> ~/.bashrc

echo "export GOROOT=/usr/local/go" >> ~/.bashrc

source ~/.bashrc

go version

 

## 安装dlv

go install [github.com/go-delve/delve/cmd/dlv@latest](mailto:github.com/go-delve/delve/cmd/dlv@latest)

echo $(go env GOPATH)/bin

echo "export PATH=$PATH:$(go env GOPATH)/bin" >> ~/.bashrc

source ~/.bashrc

 

## 调试

sudo firewall-cmd --zone=public --add-port=2345/tcp --permanent

sudo firewall-cmd --reload

go build -o myprogram .

dlv exec ./myprogram --headless --listen=:2345 --api-version=2 --accept-multiclient

 

## 本地

dlv connect 远程ip:2345

**a. break** **或 b**

在程序的特定位置设置断点（可以是函数或特定的行号）。

- 设置在某个函数的断点：

(dlv) break main.main

- 设置在特定行号的断点： 

(dlv) break main.go:10

**b. continue** **或 c**

继续执行程序，直到下一个断点或程序结束：

(dlv) continue

**c. next** **或 n**

执行下一行代码，跳过函数调用：

(dlv) next

**d. step** **或 s**

进入当前行的函数内部（逐步执行）：

(dlv) step

**f. print** **或 p**

打印变量的值：

(dlv) print someVariable

- 可以通过 print 命令查看变量的值，甚至可以打印整个结构体或数组。
- 也可以打印表达式，例如：print someVariable + 10。

**g. list** **或 l**

显示代码的当前行及其上下文：

(dlv) list

这将显示当前停止位置的周围几行代码。

 

## 通过goland调试

## 1.配置调试

· **打开 GoLand 的配置界面**：

- 在 GoLand 中，点击     Run -> Edit Configurations... 来设置运行/调试配置。

· **添加远程调试配置**：

- 点击左上角的 +，选择     Go Remote 来添加一个新的远程调试配置。

· **配置远程调试**： 在弹出的配置窗口中填写以下字段：

- **Name**：设置为一个有意义的名称，比如 Remote Debugging.
- **Host**：填写远程服务器的 IP 地址或主机名。例如 192.168.1.100。
- **Port**：设置为 Delve 调试器监听的端口。例如，2345。
- **Remote GoPath**：指定远程服务器上 Go 项目的根目录路径。例如 /home/user/project。
- **Go executable**：如果 Go 没有安装在标准路径，可以指定 Go 的可执行文件路径。

· **调试选项**：

- **Use remote SDK**：选择远程的 Go SDK。
- **Allow     parallel run**：如果你要允许多任务并行调试，可以勾选这个选项。

配置完成后点击 Apply 和 OK。

![img](image\editRemote.png)

 

使用第一种时：可以不用编译，直接运行

使用第二种时：先编译，后面的路径换成项目路径，然后运行下面的

 

## 2.配置代码上传

**1.** **配置远程开发环境（通过 SSH）**

1. **打开 GoLand 设置**：
   - 打开 GoLand，点击菜单      File -> Settings (对于 Mac，选择      GoLand -> Preferences)。
2. **配置远程服务器**：
   - 在左侧菜单栏中，选择 Build, Execution, Deployment -> Deployment。
   - 点击右上角的 + 按钮来添加新的服务器。
   - 在弹出的对话框中，选择 SFTP（通过 SSH 上传）。
   - 配置如下：
     - **Name**: 你可以给服务器起个名字。
     - **Type**: 选择 SFTP。
     - **Host**: 远程服务器的 IP 地址或主机名。
     - **Port**: 通常是 22，除非你配置了其他端口。
     - **Root Path**: 设置项目在远程服务器上的根目录路径。
     - **User Name**: 填写你的 SSH 用户名。
     - **Password/Key       Pair**: 填写密码或者使用 SSH 密钥对认证。
3. **测试连接**：
   - 在配置完后，可以点击 Test Connection 来确保 GoLand 可以成功连接到远程服务器。
4. **设置映射路径**：
   - 在同一对话框的 Mappings      标签中，设置本地项目文件与远程服务器文件之间的映射关系。通常你需要映射的是本地的项目根目录与远程服务器上的对应目录。
5. **保存配置**：
   - 配置完成后，点击 OK      保存。

**2.** **通过 GoLand 上传代码**

1. **打开你的项目**：
   - 在 GoLand 中打开本地项目。
2. **手动上传**：
   - 在你本地的文件编辑完成后，右键点击项目文件，选择 Deployment -> Upload to，然后选择你之前配置的远程服务器。
   - GoLand 会自动将你所选文件或文件夹上传到远程服务器。
3. **自动上传**：
   - 在 GoLand 设置中启用自动上传功能：选择      Settings -> Build, Execution, Deployment -> Deployment      -> 选择你设置的远程服务器配置 -> 勾选 Automatic Upload。
   - 这样，GoLand 会在你本地文件保存时自动将文件上传到远程服务器。



 ![img](image\upload.png)

 ![img](image\upload2.png)





 

 

 

## 问题

 

每次修改代码都要重新运行远程的命令，结束调试时不可以再次调试

 

 

 

 

 