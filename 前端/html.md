### **1. HTML 基础**

#### **HTML 结构**

HTML 由标签（tags）组成，基本结构如下：

```
html复制编辑<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>我的第一个网页</title>
</head>
<body>
    <h1>你好，HTML！</h1>
    <p>这是一个段落。</p>
</body>
</html>
```

#### **常见标签**

- **标题**：`<h1>` ~ `<h6>`
- **段落**：`<p>`
- **换行**：`<br>`
- **分隔线**：`<hr>`
- **强调**：`<strong>`（加粗），`<em>`（斜体）

------

### **2. 列表**

#### **无序列表**

```
html复制编辑<ul>
    <li>苹果</li>
    <li>香蕉</li>
    <li>橙子</li>
</ul>
```

#### **有序列表**

```
html复制编辑<ol>
    <li>第一步</li>
    <li>第二步</li>
    <li>第三步</li>
</ol>
```

------

### **3. 链接和图片**

#### **超链接**

```
html


复制编辑
<a href="https://www.example.com">访问 Example</a>
```

- `_blank` 在新窗口打开：

```
html


复制编辑
<a href="https://www.example.com" target="_blank">新窗口打开</a>
```

#### **图片**

```
html


复制编辑
<img src="image.jpg" alt="描述文本" width="300">
```

------

### **4. 表格**

```
html复制编辑<table border="1">
    <tr>
        <th>姓名</th>
        <th>年龄</th>
    </tr>
    <tr>
        <td>张三</td>
        <td>25</td>
    </tr>
</table>
```

------

### **5. 表单**

```
html复制编辑<form action="submit.php" method="post">
    <label for="name">姓名：</label>
    <input type="text" id="name" name="name" required>
    
    <label for="email">邮箱：</label>
    <input type="email" id="email" name="email" required>

    <button type="submit">提交</button>
</form>
```