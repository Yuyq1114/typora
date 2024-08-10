# **Aircrack-ng** 

是一款强大的开源工具集，专门用于无线网络安全测试和破解。它广泛应用于无线网络的安全评估、密码破解、网络监控和数据包分析等领域。以下是对 Aircrack-ng 工具集的详细介绍，包括其主要工具、功能、使用方法和常见应用场景。

### 1. **Aircrack-ng 概述**

**Aircrack-ng** 包含一系列工具，用于以下任务：

- **无线网络监控**：捕获和分析无线网络数据包。
- **网络破解**：破解无线网络的加密密码。
- **网络注入**：向无线网络注入数据包以测试网络的安全性。
- **信息收集**：收集有关无线网络的信息，包括网络配置和设备信息。

### 2. **主要工具和功能**

#### 2.1. **airmon-ng**

- 功能

  ：

  - 用于启动和停止无线网卡的监控模式。监控模式使网卡能够捕获无线网络上的所有数据包，而不仅仅是发送到网卡的包。

- 常用命令

  ：

  - 启动监控模式：

    ```
    airmon-ng start wlan0
    ```

  - 停止监控模式：

    ```
    airmon-ng stop wlan0mon
    ```

#### 2.2. **airodump-ng**

- 功能

  ：

  - 用于捕获和记录无线网络的数据包，包括管理帧、数据帧和控制帧。
  - 主要用于发现和分析无线网络，收集网络流量。

- 常用命令

  ：

  - 捕获数据包并保存到文件：

    ```
    airodump-ng -w outputfile --bssid [目标BSSID] -c [频道] wlan0mon
    ```

  - 捕获所有网络数据包：

    ```
    airodump-ng wlan0mon
    ```

#### 2.3. **aireplay-ng**

- 功能

  ：

  - 用于注入数据包到无线网络中。可以用于执行各种攻击，如重放攻击、欺骗攻击等。

- 常用命令

  ：

  - 执行重放攻击：

    ```
    aireplay-ng --deauth [次数] -a [目标BSSID] wlan0mon
    ```

  - 注入伪造的认证数据包：

    ```
    aireplay-ng --auth [次数] -a [目标BSSID] -c [目标客户端MAC] wlan0mon
    ```

#### 2.4. **aircrack-ng**

- 功能

  ：

  - 用于破解捕获的数据包中包含的无线网络密码。支持多种加密算法，如 WEP 和 WPA/WPA2。

- 常用命令

  ：

  - 破解 WPA/WPA2 密码：

    ```
    aircrack-ng -w [密码字典文件] -b [目标BSSID] capturefile.cap
    ```

  - 破解 WEP 密码：

    ```
    aircrack-ng -b [目标BSSID] capturefile.cap
    ```

#### 2.5. **airdecap-ng**

- 功能

  ：

  - 用于解密捕获的加密数据包，适用于 WEP 和 WPA 密码破解之后的数据包分析。

- 常用命令

  ：

  - 解密 WEP 数据包：

    ```
    airdecap-ng -w [密码] capturefile.cap
    ```

#### 2.6. **airbase-ng**

- 功能

  ：

  - 用于创建一个虚假的接入点。可以用于测试和欺骗攻击。

- 常用命令

  ：

  - 创建一个虚假的接入点：

    ```
    airbase-ng -e [虚假SSID] -c [频道] wlan0mon
    ```

#### 2.7. **packetforge-ng**

- 功能

  ：

  - 用于创建自定义的数据包，可以用于注入和测试。

- 常用命令

  ：

  - 创建数据包：

    ```
    packetforge-ng -0 -a [目标BSSID] -h [目标客户端MAC] -k [虚假SSID] -w outputfile.cap
    ```

### 3. **使用 Aircrack-ng 的步骤**

#### 3.1. **准备**

1. **安装 Aircrack-ng**：

   - 在 Debian/Ubuntu 系统中：

     ```
     sudo apt-get install aircrack-ng
     ```

   - 在其他系统中，根据发行版和系统版本选择适当的安装方法。

2. **准备无线网卡**：

   - 确保你的无线网卡支持监控模式，并能与 Aircrack-ng 兼容。

#### 3.2. **启动监控模式**

1. 启动监控模式

   ：

   - 使用 

     ```
     airmon-ng
     ```

      启动监控模式：

     ```
     sudo airmon-ng start wlan0
     ```

#### 3.3. **捕获数据包**

1. 使用 airodump-ng 捕获数据包

   ：

   - 捕获数据包并保存到文件：

     ```
     sudo airodump-ng -w capturefile --bssid [目标BSSID] -c [频道] wlan0mon
     ```

#### 3.4. **注入数据包**

1. 使用 aireplay-ng 执行重放攻击

   ：

   - 例如，执行 deauth 攻击：

     ```
     sudo aireplay-ng --deauth 10 -a [目标BSSID] wlan0mon
     ```

#### 3.5. **破解密码**

1. **使用 aircrack-ng 破解 WPA/WPA2 密码**：

   - 例如，使用密码字典破解：

     ```
     sudo aircrack-ng -w [密码字典文件] -b [目标BSSID] capturefile.cap
     ```

2. **使用 aircrack-ng 破解 WEP 密码**：

   - 例如：

     ```
     sudo aircrack-ng -b [目标BSSID] capturefile.cap
     ```

### 4. **常见应用场景**

#### 4.1. **WEP 密码破解**

1. **捕获 WEP 数据包**：使用 `airodump-ng` 捕获大量数据包。
2. **使用 `aircrack-ng` 破解密码**：应用密码字典进行破解。

#### 4.2. **WPA/WPA2 密码破解**

1. **捕获握手数据包**：使用 `airodump-ng` 捕获 WPA/WPA2 握手数据包。
2. **使用 `aircrack-ng` 破解密码**：使用密码字典进行破解。

#### 4.3. **创建虚假接入点**

1. **使用 `airbase-ng` 创建虚假接入点**：测试网络安全性或进行欺骗攻击。

#### 4.4. **执行重放攻击**

1. **使用 `aireplay-ng` 执行重放攻击**：测试网络的抗攻击能力。



# 参数详解

```
aircrack-ng [options] [capturefile]
```

### **主要参数**

#### **1. 捕获文件**

- `[capturefile]`

  :

  - **说明**: 指定包含捕获数据包的文件。这些文件通常由 `airodump-ng` 捕获生成，包含网络流量和握手数据。

  - 示例

    :

    ```
    aircrack-ng capturefile.cap
    ```

#### **2. 破解 WEP 密码**

- **`-b [BSSID]`**:

  - **说明**: 指定目标无线网络的 BSSID。这个参数用于过滤目标网络。

  - 示例

    :

    ```
    aircrack-ng -b 00:11:22:33:44:55 capturefile.cap
    ```

- **`-w [wordlist]`**:

  - **说明**: 指定密码字典文件。`aircrack-ng` 会用这个字典尝试破解密码。

  - 示例

    :

    ```
    aircrack-ng -w /path/to/wordlist.txt capturefile.cap
    ```

- **`-e [SSID]`**:

  - **说明**: 指定网络的 SSID，通常与 BSSID 一起使用。

  - 示例

    :

    ```
    aircrack-ng -b 00:11:22:33:44:55 -e MyNetwork capturefile.cap
    ```

#### **3. 破解 WPA/WPA2 密码**

- **`-w [wordlist]`**:

  - **说明**: 同样用于指定破解 WPA/WPA2 密码的密码字典文件。

  - 示例

    :

    ```
    aircrack-ng -w /path/to/wordlist.txt capturefile.cap
    ```

- **`-b [BSSID]`**:

  - **说明**: 用于指定目标网络的 BSSID。特别是在 WPA/WPA2 中，此参数用于确定目标网络。

  - 示例

    :

    ```
    aircrack-ng -b 00:11:22:33:44:55 capturefile.cap
    ```

- **`-e [SSID]`**:

  - **说明**: 指定目标网络的 SSID，用于提高破解精度。

  - 示例

    :

    ```
    aircrack-ng -b 00:11:22:33:44:55 -e MyNetwork capturefile.cap
    ```

- **`-a`**:

  - **说明**: 指定破解算法，通常 `aircrack-ng` 自动选择最适合的算法。如果需要手动指定，可以使用 `-a 1` (WEP) 或 `-a 2` (WPA/WPA2)。

  - 示例

    :

    ```
    aircrack-ng -a 2 -w /path/to/wordlist.txt capturefile.cap
    ```

- **`-l [outputfile]`**:

  - **说明**: 指定将破解密码保存到的文件中。这对于记录破解过程和结果非常有用。

  - 示例

    :

    ```
    aircrack-ng -w /path/to/wordlist.txt -l cracked_password.txt capturefile.cap
    ```

#### **4. 其他选项**

- **`-n [keysize]`**:

  - **说明**: 用于指定 WEP 密码的长度（以比特为单位）。这在尝试破解 WEP 密码时很有用。

  - 示例

    :

    ```
    aircrack-ng -n 128 capturefile.cap
    ```

- **`-t [thread count]`**:

  - **说明**: 设置使用的线程数，以提高破解速度。

  - 示例

    :

    ```
    aircrack-ng -t 4 -w /path/to/wordlist.txt capturefile.cap
    ```

- **`-j [key]`**:

  - **说明**: 用于指定 WEP 密钥的候选值。这个选项主要用于 WEP 破解时提供初始猜测。

  - 示例

    :

    ```
    aircrack-ng -j 0123456789 capturefile.cap
    ```

- **`-h`**:

  - **说明**: 显示帮助信息，列出所有可用参数及其用法。

  - 示例

    :

    ```
    aircrack-ng -h
    ```

### **常见使用示例**

1. **破解 WEP 密码**：

   - 捕获 WEP 数据包并尝试破解：

     ```
     aircrack-ng -b 00:11:22:33:44:55 -w /path/to/wordlist.txt capturefile.cap
     ```

2. **破解 WPA/WPA2 密码**：

   - 捕获 WPA/WPA2 握手数据包并尝试破解：

     ```
     aircrack-ng -b 00:11:22:33:44:55 -e MyNetwork -w /path/to/wordlist.txt capturefile.cap
     ```

3. **指定输出文件保存破解结果**：

   - 将破解密码保存到指定文件：

     ```
     aircrack-ng -w /path/to/wordlist.txt -l cracked_password.txt capturefile.cap
     ```

4. **使用多线程提高破解速度**：

   - 使用 4 个线程进行破解：

     ```
     aircrack-ng -t 4 -w /path/to/wordlist.txt capturefile.cap
     ```

