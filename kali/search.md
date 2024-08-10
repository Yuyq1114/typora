# google hacking

### **基本搜索运算符**

1. **`site:`**

   - **说明**: 限制搜索结果为特定网站或域名。

   - 示例

     : 查找某网站中的所有 PDF 文件：

     ```
     site:example.com filetype:pdf
     ```

2. **`filetype:`**

   - **说明**: 查找特定类型的文件。

   - 示例

     : 查找所有 

     ```
     .doc
     ```

      文件：

     ```
     filetype:doc
     ```

3. **`intitle:`**

   - **说明**: 查找标题中包含指定关键词的网页。

   - 示例

     : 查找标题中包含 "login" 的页面：

     ```
     intitle:login
     ```

4. **`inurl:`**

   - **说明**: 查找 URL 中包含指定关键词的页面。

   - 示例

     : 查找 URL 中包含 "admin" 的页面：

     ```
     inurl:admin
     ```

5. **`intext:`**

   - **说明**: 查找页面内容中包含指定关键词的页面。

   - 示例

     : 查找页面内容中包含 "password" 的页面：

     ```
     intext:password
     ```

6. **`allintext:`**

   - **说明**: 查找页面内容中包含所有指定关键词的页面。

   - 示例

     : 查找页面内容中同时包含 "username" 和 "password" 的页面：

     ```
     allintext:username password
     ```

7. **`allintitle:`**

   - **说明**: 查找标题中包含所有指定关键词的页面。

   - 示例

     : 查找标题中同时包含 "admin" 和 "panel" 的页面：

     ```
     allintitle:admin panel
     ```

8. **`allinurl:`**

   - **说明**: 查找 URL 中包含所有指定关键词的页面。

   - 示例

     : 查找 URL 中同时包含 "login" 和 "admin" 的页面：

     ```
     allinurl:login admin
     ```

9. **`cache:`**

   - **说明**: 查看 Google 缓存中的网页内容。

   - 示例

     : 查看 

     ```
     example.com
     ```

      的缓存：

     ```
     cache:example.com
     ```

10. **`related:`**

    - **说明**: 查找与指定网站相关的其他网站。

    - 示例

      : 查找与 

      ```
      example.com
      ```

       相关的其他网站：

      ```
      related:example.com
      ```

### **进阶技巧**

1. **查找敏感文件和目录**

   - 查找暴露的配置文件：

     ```
     site:example.com filetype:conf
     ```

   - 查找暴露的备份文件：

     ```
     site:example.com filetype:bak
     ```

   - 查找暴露的数据库文件：

     ```
     site:example.com filetype:sql
     ```

2. **发现管理面板**

   - 查找常见的管理面板登录页面：

     ```
     intitle:"Admin Login" OR intitle:"Login" inurl:admin
     ```

3. **查找公开的文档和数据**

   - 查找公开的 Excel 文档：

     ```
     filetype:xls inurl:financial
     ```

   - 查找公开的用户数据：

     ```
     filetype:csv inurl:users
     ```

4. **寻找网页漏洞**

   - 查找含有 "403 Forbidden" 或 "404 Not Found" 的页面（可能存在隐藏的内容）：

     ```
     "403 Forbidden" OR "404 Not Found"
     ```

   - 查找可能的文件上传漏洞：

     ```
     intitle:"File Upload" inurl:upload
     ```

5. **暴露的网络摄像头**

   - 查找公开的网络摄像头：

     ```
     intitle:"Live View / - AXIS" OR intitle:"IP Camera" OR intitle:"Webcam"
     ```

6. **寻找公开的企业数据**

   - 查找包含企业敏感信息的文件：

     ```
     site:example.com "confidential" OR "internal use only"
     ```



# shodan hacking

#### **1. 访问 Shodan**

- **访问网址**：
  - 打开 [Shodan 官方网站](https://www.shodan.io/)。
- **创建账户**：
  - 注册一个 Shodan 账户以获得更多功能和使用权限。

#### **2. 使用基本搜索**

- 进行基本搜索
  - 在搜索框中输入关键词，Shodan 将返回相关设备和服务的结果。
- 查看设备信息
  - 点击搜索结果中的设备，可以查看该设备的详细信息，包括 IP 地址、端口、服务、地理位置等。

### **常见搜索技巧**

#### **1. 基本搜索**

- **搜索特定设备**：

  - 查找特定类型的设备，例如：

    ```
    "webcam"
    ```

- **搜索特定服务**：

  - 查找特定服务的暴露，例如：

    ```
    "apache"
    ```

#### **2. 使用过滤器**

- **搜索特定端口**：

  - 查找特定端口开放的设备，例如：

    ```
    port:80
    ```

- **搜索特定国家的设备**：

  - 查找特定国家的设备，例如：

    ```
    country:"US"
    ```

- **搜索特定城市的设备**：

  - 查找特定城市的设备，例如：

    ```
    city:"San Francisco"
    ```

- **搜索特定操作系统**：

  - 查找运行特定操作系统的设备，例如：

    ```
    os:"Windows"
    ```

- **搜索特定组织的设备**：

  - 查找属于特定组织的设备，例如：

    ```
    org:"Google"
    ```

#### **3. 使用高级搜索**

- **组合多个搜索条件**：

  - 使用布尔运算符组合多个搜索条件，例如：

    ```
    port:22 AND country:"US"
    ```

- **搜索特定版本的服务**：

  - 查找特定版本的服务，例如：

    ```
    "Apache/2.4.41"
    ```

- **搜索特定响应内容**：

  - 查找包含特定响应内容的设备，例如：

    ```
    "Welcome to my website"
    ```

### **实用策略**

#### **1. 发现暴露的设备**

- **查找未授权的设备**：
  - 查找暴露在互联网上的未授权设备，例如网络摄像头、路由器、打印机等。
- **检查默认配置**：
  - 查找使用默认配置或密码的设备，这可能表明安全配置不当。

#### **2. 漏洞扫描**

- **查找已知漏洞的设备**：
  - 查找已知存在漏洞的设备和服务，例如特定版本的服务器或应用程序。
- **分析公开的管理界面**：
  - 查找暴露的管理界面，分析其安全性和配置错误。

#### **3. 网络监控**

- **监控网络上的设备**：
  - 使用 Shodan 监控网络中的设备状态，检测配置变更和新设备的加入。
- **审计公开服务**：
  - 审计网络中公开的服务，确保没有未授权的服务暴露在外。

### **示例应用**

#### **1. 查找暴露的摄像头**

- 查找公开的 IP 摄像头

  ```
  "webcam" OR "camera"
  ```

#### **2. 查找暴露的数据库**

- 查找暴露的数据库服务

  ```
  port:3306 OR port:5432
  ```

#### **3. 查找特定版本的服务**

- 查找运行特定版本 Apache 的设备

  ```
  "Apache/2.4.41"
  ```

1. - 