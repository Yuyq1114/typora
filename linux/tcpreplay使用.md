## **tcpreplay 的核心功能**

### 1. **基础命令**

将 `.pcap` 文件流量重放到指定网卡：

```
bash


复制代码
sudo tcpreplay --intf1=eth0 traffic.pcap
```

- `--intf1=eth0`：指定发送流量的网卡为 `eth0`。
- `traffic.pcap`：要重放的 `.pcap` 文件。

### 2. **调整发送速度**

默认情况下，`tcpreplay` 会以原始捕获速度发送流量（即捕获时的时间间隔）。

- **倍速模式**（以指定倍数调整发送速度）：

  ```
  bash
  
  
  复制代码
  sudo tcpreplay --intf1=eth0 --multiplier=2 traffic.pcap
  ```

  （`--multiplier=2` 将流量发送速度提高到原来的 2 倍）

- **指定速率**（以特定速率发送流量）：

  ```
  bash
  
  
  复制代码
  sudo tcpreplay --intf1=eth0 --mbps=100 traffic.pcap
  ```

  （以 100 Mbps 的速率发送）

- **忽略时间间隔**（以最大速度发送流量）：

  ```
  bash
  
  
  复制代码
  sudo tcpreplay --intf1=eth0 --topspeed traffic.pcap
  ```

### 3. **双网卡模式**

使用双网卡，支持同时发送和接收：

```
bash


复制代码
sudo tcpreplay --intf1=eth0 --intf2=eth1 traffic.pcap
```

- `--intf1=eth0`：发送流量。
- `--intf2=eth1`：接收流量。

### 4. **发送部分包**

仅发送 `.pcap` 文件中的前 N 个包：

```
bash


复制代码
sudo tcpreplay --intf1=eth0 --limit=100 traffic.pcap
```

- `--limit=100`：发送前 100 个数据包。

### 5. **循环发送**

重复发送 `.pcap` 中的数据包：

```
bash


复制代码
sudo tcpreplay --intf1=eth0 --loop=5 traffic.pcap
```

- `--loop=5`：将数据包循环发送 5 次。
- `--loop=0`：无限循环，直到手动终止。

### 6. **反向流量**

将流量中的源 IP 和目标 IP 对调，模拟双向通信：

```
bash


复制代码
sudo tcpreplay --intf1=eth0 --reverse traffic.pcap
```

------

## **tcpreplay 的其他工具**

### **tcprewrite**

`tcprewrite` 是一个附带的工具，用于修改 `.pcap` 文件的内容，方便测试不同的场景。

#### 修改 MAC 地址

将 `.pcap` 文件中的目标 MAC 地址修改为指定值：

```
bash


复制代码
tcprewrite --infile=traffic.pcap --outfile=modified.pcap --enet-dmac=00:11:22:33:44:55
```

#### 修改 IP 地址

将源 IP 修改为 `192.168.1.100`，目标 IP 修改为 `192.168.1.200`：

```
bash


复制代码
tcprewrite --infile=traffic.pcap --outfile=modified.pcap --srcipmap=0.0.0.0/0:192.168.1.100 --dstipmap=0.0.0.0/0:192.168.1.200
```

#### 删除 VLAN 标签

从 `.pcap` 文件中移除 VLAN 标签：

```
bash


复制代码
tcprewrite --infile=traffic.pcap --outfile=modified.pcap --delete-vlan
```

------

### **tcpreplay-edit**

`tcpreplay-edit` 是 `tcpreplay` 的增强版，支持动态修改包内容后再重放。

例如，修改源 IP 并重放：

```
bash


复制代码
sudo tcpreplay-edit --intf1=eth0 --srcipmap=0.0.0.0/0:192.168.1.1 traffic.pcap
```

------

## **常用选项和参数**

| 参数           | 描述                                        |
| -------------- | ------------------------------------------- |
| `--intf1`      | 指定发送流量的网卡。    **-i veth0**        |
| `--intf2`      | 指定接收流量的网卡（双网卡模式）。          |
| `--multiplier` | 调整发送速度的倍数。                        |
| `--mbps`       | 指定发送速度（Mbps）。  **-M 30**           |
| `--topspeed`   | 以最大速度发送流量。                        |
| `--limit`      | 指定只发送前 N 个包。                       |
| `--loop`       | 循环发送流量。         **-l 20**            |
| `--reverse`    | 反向流量：交换源 IP 和目标 IP。             |
| `--unique-ip`  | 修改 `.pcap` 文件中的 IP 地址，使其唯一化。 |
| `--stats`      | 显示重放完成后的统计信息。                  |