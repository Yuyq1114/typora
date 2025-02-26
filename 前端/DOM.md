# **1. 获取元素**

### **1.1 `getElementById()`**

通过 ID 获取一个元素，返回一个元素对象。

```
js


复制编辑
let element = document.getElementById("myElement");
```

**HTML**：

```
html


复制编辑
<div id="myElement">这是一个元素</div>
```

### **1.2 `getElementsByClassName()`**

通过类名获取元素，返回一个“类数组”对象。

```
js


复制编辑
let elements = document.getElementsByClassName("myClass");
```

### **1.3 `getElementsByTagName()`**

通过标签名获取元素，返回一个“类数组”对象。

```
js


复制编辑
let paragraphs = document.getElementsByTagName("p");
```

### **1.4 `querySelector()`**

返回匹配的第一个元素，可以使用 CSS 选择器。

```
js


复制编辑
let element = document.querySelector(".myClass");
```

### **1.5 `querySelectorAll()`**

返回所有匹配的元素，返回一个 NodeList 对象。

```
js


复制编辑
let elements = document.querySelectorAll("p.myClass");
```

------

# **2. 操作元素内容**

### **2.1 `innerHTML`**

获取或设置元素的 HTML 内容。

```
js复制编辑let content = element.innerHTML;  // 获取
element.innerHTML = "<p>新内容</p>";  // 设置
```

### **2.2 `innerText`**

获取或设置元素的文本内容，通常用于纯文本内容。

```
js复制编辑let text = element.innerText;  // 获取
element.innerText = "新文本内容";  // 设置
```

### **2.3 `textContent`**

与 `innerText` 类似，但 `textContent` 不会考虑 CSS 样式，始终返回纯文本。

```
js


复制编辑
let content = element.textContent;
```

------

# **3. 修改样式**

### **3.1 `style`**

直接修改元素的样式，通过 `.style` 访问内联样式。

```
js复制编辑element.style.color = "red";  // 修改文本颜色
element.style.fontSize = "20px";  // 修改字体大小
```

### **3.2 `classList`**

添加、删除或切换类名。

- `add()`：添加类
- `remove()`：删除类
- `toggle()`：切换类
- `contains()`：检查类是否存在

```
js复制编辑element.classList.add("active");
element.classList.remove("inactive");
element.classList.toggle("hidden");
```

------

# **4. 操作属性**

### **4.1 `getAttribute()` 和 `setAttribute()`**

获取和设置元素的属性值。

```
js复制编辑let href = element.getAttribute("href");  // 获取属性
element.setAttribute("href", "https://www.example.com");  // 设置属性
```

### **4.2 `removeAttribute()`**

删除元素的属性。

```
js


复制编辑
element.removeAttribute("href");
```

------

# **5. 创建和删除元素**

### **5.1 创建新元素**

```
js复制编辑let newElement = document.createElement("div");
newElement.innerHTML = "这是新元素";
document.body.appendChild(newElement);  // 添加到页面中
```

### **5.2 删除元素**

```
js


复制编辑
document.body.removeChild(element);  // 从页面中移除
```

------

# **6. 操作事件**

### **6.1 `addEventListener()`**

为元素添加事件监听器，处理用户的交互操作。

```
js复制编辑element.addEventListener("click", function() {
    alert("元素被点击了！");
});
```

可以监听常见的事件：

- `click`：点击
- `mouseover`：鼠标悬停
- `keydown`：键盘按键
- `submit`：表单提交

### **6.2 `removeEventListener()`**

移除事件监听器。

```
js复制编辑function handleClick() {
    alert("元素被点击了！");
}
element.addEventListener("click", handleClick);
element.removeEventListener("click", handleClick);
```

------

# **7. DOM 与 JSON**

JSON 格式数据与 DOM 元素的交互，例如使用 `fetch` 从服务器加载数据并动态更新页面：

```
js复制编辑fetch("https://api.example.com/data")
    .then(response => response.json())
    .then(data => {
        document.getElementById("title").innerText = data.title;
        document.getElementById("description").innerText = data.description;
    });
```