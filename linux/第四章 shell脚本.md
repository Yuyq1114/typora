Linux 的 Shell 脚本是一种用于自动化操作的脚本语言，可以用来执行一系列的命令。Shell 是用户与操作系统之间的接口，它允许用户输入命令并执行。通过编写 Shell 脚本，可以将多个命令组合在一起，以完成复杂的任务。

### **基本概念**

1. **Shell 类型**：

   - **Bash (Bourne Again SHell)**：最常用的 Shell，提供了强大的功能。
   - **Sh (Bourne Shell)**：早期的 Shell，功能较少。
   - **Csh (C Shell)**：与 C 语言语法相似的 Shell。
   - **Zsh**：功能强大的 Shell，支持更多的特性。

2. **文件扩展名**：

   - Shell 脚本通常使用 `.sh` 作为文件扩展名，但实际上没有强制要求。

3. **执行权限**：

   - 在执行 Shell 脚本之前，需要确保脚本文件具有可执行权限，可以使用以下命令：

     ```
     chmod +x script.sh
     ```

### **脚本结构**

一个基本的 Shell 脚本结构如下：

```
#!/bin/bash
# 这是一个注释
echo "Hello, World!"  # 输出 "Hello, World!"
```

- `#!/bin/bash`：指定使用的 Shell 解释器，称为 shebang。
- 注释行以 `#` 开头。

### **基本命令**

1. **变量**：

   - 声明变量时，不需要空格：

     ```
     my_var="Hello"
     echo $my_var  # 输出 Hello
     ```

2. **`pwd` 命令**：
   `pwd`（print working directory）是一个Shell命令，它用于显示当前工作目录的完整路径。例如，如果当前工作目录是 `/home/user/project`，运行 `pwd` 将输出 `/home/user/project`。

   **反引号（```）**：
   反引号是Shell中的一种命令替换（Command Substitution）方式，它会执行其中的命令，并将命令的输出结果返回。也可以用 `$()` 代替反引号，例如：`root=$(pwd)，这两者效果相同。

   **变量赋值**：
   将 `pwd` 命令的输出结果（即当前工作目录的路径）赋值给 `root` 变量。

   ### 例子：

   ```
   #!/bin/bash
   root=`pwd`
   echo "当前目录是: $root"
   ```

3. **条件语句**：

   - 使用 

     ```
     if
     ```

      语句进行条件判断：

     ```
     if [ "$my_var" = "Hello" ]; then
         echo "Variable is Hello"
     else
         echo "Variable is not Hello"
     fi
     ```

4. **循环**：

   - 使用 `for` 循环：

     ```
     for i in {1..5}; do
         echo "Iteration $i"
     done
     ```

   - 使用 `while` 循环：

     ```
     count=1
     while [ $count -le 5 ]; do
         echo "Count is $count"
         ((count++))
     done
     ```

5. **函数**：

   - 定义和调用函数：

     ```
     my_function() {
         echo "This is my function"
     }
     
     my_function  # 调用函数
     ```

### **使用示例**

一个简单的 Shell 脚本示例，演示如何创建一个备份文件：

```
#!/bin/bash

# 备份目录
backup_dir="/path/to/backup"
# 要备份的文件
file_to_backup="/path/to/file"

# 创建备份
cp $file_to_backup $backup_dir

# 输出备份完成消息
echo "Backup of $file_to_backup completed in $backup_dir"
```

### **脚本调试**

在调试 Shell 脚本时，可以使用 `-x` 选项执行脚本，以便查看执行的每个命令：

```

bash -x script.sh
```

### **脚本执行**

要执行 Shell 脚本，可以使用以下两种方法之一：

1. 通过直接调用：

   ```
   
   ./script.sh
   ```

2. 通过 Shell 命令调用：

   ```
   
   bash script.sh
   ```

在Shell脚本中，`set -e` 是一种用于控制脚本执行行为的命令，具体作用是 **如果脚本中的任何命令返回非零退出状态（即命令失败），脚本将立即退出**，而不会继续执行后续的命令。

###  **`$@` 代表所有参数**：

- 在Shell脚本中，`$@` 是一个特殊变量，它表示传递给脚本或函数的所有参数。
- 使用 `$@` 传递参数时，它会保留每个参数的独立性，比如参数中如果有空格或特殊字符，它们会被正确处理。

在Shell脚本中，`$#` 是一个特殊变量，表示传递给脚本或函数的 **参数数量**。

### 详细解释：

- **脚本级别**：如果你在脚本中使用 `$#`，它会返回运行该脚本时传递的参数的数量。
- **函数级别**：如果你在函数内部使用 `$#`，它会返回调用该函数时传递的参数数量。

- 
- `-e`：**多点编辑**。允许在同一命令中使用多个编辑脚本。
- `-i`：**直接编辑文件**。将编辑结果直接写回文件，而不是输出到标准输出。
- `-r`：启用扩展的正则表达式（ERE，Extended Regular Expression）。
- `-f`：**从文件读取脚本**。可以将复杂的 `sed` 脚本放在文件中，并通过 `-f` 选项执行。

### 常用命令

#### 1. 替换（`s`）

`sed` 中最常用的功能是替换，即用新内容替换匹配的文本。其基本格式如下：

```
sed 's/模式/替换文本/标志' 文件
```

- **模式**：要匹配的文本或正则表达式。

- **替换文本**：用于替换的内容。

- 标志

  ：

  - `g`：全局替换，即替换行中所有匹配的内容。
  - `p`：打印替换后的结果（通常与 `-n` 一起使用）。
  - `1`：仅替换每行第一个匹配项。
  - `n`：第 `n` 个匹配项。

#### 示例 1：将 `apple` 替换为 `orange`

```
echo "apple is red" | sed 's/apple/orange/'
# 输出：orange is red
```

#### 示例 2：替换每行中所有的 `apple`

```
echo "apple apple apple" | sed 's/apple/orange/g'
# 输出：orange orange orange
```

#### 示例 3：只替换每行中第二个 `apple`

```
echo "apple apple apple" | sed 's/apple/orange/2'
# 输出：apple orange apple
```

#### 示例 4：带正则表达式的替换

将所有以 `a` 开头的单词替换为 `fruit`：

```
echo "apple banana cherry" | sed 's/\ba\w*/fruit/g'
# 输出：fruit banana cherry
```

#### 2. 删除（`d`）

`sed` 的 `d` 命令用于删除指定的行。

#### 示例 1：删除文件中的第 2 行

```
sed '2d' filename
```

#### 示例 2：删除第 3 到第 5 行

```
sed '3,5d' filename
```

#### 示例 3：删除包含某个模式的行

删除所有包含 `apple` 的行：

```
sed '/apple/d' filename
```

#### 3. 打印（`p`）

`p` 命令用于打印指定的行。通常与 `-n` 选项配合使用，只打印匹配到的行。

#### 示例 1：只打印文件中的第 2 行

```
sed -n '2p' filename
```

#### 示例 2：打印包含某个模式的行

打印包含 `apple` 的行：

```
sed -n '/apple/p' filename
```

#### 4. 插入和追加（`i` 和 `a`）

- **`i`**：在匹配的行前面插入文本。
- **`a`**：在匹配的行后面追加文本。

#### 示例 1：在第 2 行前插入一行 "Hello"

```
sed '2i\Hello' filename
```

#### 示例 2：在第 3 行后追加一行 "World"

```
sed '3a\World' filename
```

#### 示例 3：在匹配模式的行前插入

在所有包含 `apple` 的行前插入 "Fruit":

```
sed '/apple/i\Fruit' filename
```

#### 5. 修改（`c`）

`c` 命令用于将指定行替换为给定的文本。

#### 示例 1：将第 3 行替换为 "This is new line"

```
sed '3c\This is new line' filename
```

#### 示例 2：将包含 `apple` 的行替换为 "This is an apple line"

```
sed '/apple/c\This is an apple line' filename
```

#### 6. 多命令（`;` 和 `{}`）

`sed` 可以在同一个命令中执行多个操作，使用 `;` 或 `{}` 作为分隔符。

#### 示例 1：同时删除第 1 行和第 3 行

```
sed '1d; 3d' filename
```

#### 示例 2：将第 2 行替换内容，并在第 4 行前插入一行

```
sed '2s/apple/orange/; 4i\This is a new line' filename
```

#### 7. 替换的高级用法（替换带分隔符的字符串）

有时你可能会替换路径或其他包含 `/` 字符的字符串。在这种情况下，可以使用其他分隔符来避免混淆，如 `|`：

#### 示例 1：替换 `/path/to/old` 为 `/path/to/new`

```
sed 's|/path/to/old|/path/to/new|g' filename
```

### 高级用法

#### 1. 结合正则表达式

`sed` 支持使用正则表达式进行模式匹配和替换。默认使用基础正则表达式（BRE），可以通过 `-r` 启用扩展正则表达式（ERE）。

#### 示例：匹配以数字开头的行，并替换数字

```
echo "123abc" | sed 's/^[0-9]*/[number]/'
# 输出：[number]abc
```

#### 2. 使用 `-i` 修改文件内容

`-i` 选项可以直接修改文件内容，而不输出到标准输出。

#### 示例：将文件中的 `apple` 替换为 `orange`，并将结果直接保存到原文件

```

sed -i 's/apple/orange/g' filename
```

#### 3. 在指定范围内执行命令

你可以指定行号范围，或者使用模式匹配来限制 `sed` 的作用范围。

#### 示例 1：替换第 2 到第 4 行中的 `apple`

```
sed '2,4s/apple/orange/' filename
```

#### 示例 2：在匹配到 `pattern` 的行到文件末尾，删除所有行

```
sed '/pattern/,$d' filename
```
