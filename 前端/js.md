# **1. JavaScript 基础**

### **1.1 引入 JavaScript**

JS 可以通过以下方式添加到 HTML：

1. 内联 JS

   （不推荐）

   ```
   html
   
   
   复制编辑
   <button onclick="alert('你好，JS！')">点击我</button>
   ```

2. 内部 JS

   （写在 

   ```
   <script>
   ```

    标签内）

   ```
   html复制编辑<script>
       console.log("Hello, JavaScript!");
   </script>
   ```

3. 外部 JS

   （推荐，将代码写入 

   ```
   .js
   ```

    文件）

   ```
   html
   
   
   复制编辑
   <script src="script.js"></script>
   ```

------

# **2. 变量 & 数据类型**

### **2.1 变量声明**

- `var`（不推荐）
- `let`（推荐，可变变量）
- `const`（推荐，不可变变量）

```
js复制编辑let name = "张三";
const age = 25;
```

### **2.2 数据类型**

| 数据类型    | 示例               |
| ----------- | ------------------ |
| `String`    | `"Hello"`          |
| `Number`    | `100`，`3.14`      |
| `Boolean`   | `true`，`false`    |
| `Array`     | `[1, 2, 3]`        |
| `Object`    | `{ name: "张三" }` |
| `null`      | `null`             |
| `undefined` | `undefined`        |

```
js复制编辑let isStudent = true;
let scores = [90, 80, 100];
let person = { name: "张三", age: 25 };
```

------

# **3. 操作符**

### **3.1 算术运算符**

```
js复制编辑let a = 10, b = 3;
console.log(a + b);  // 加法
console.log(a - b);  // 减法
console.log(a * b);  // 乘法
console.log(a / b);  // 除法
console.log(a % b);  // 取余
console.log(a ** b); // 幂运算
```

### **3.2 比较运算符**

```
js复制编辑console.log(5 > 3);   // true
console.log(5 == "5");  // true（值相等，不考虑类型）
console.log(5 === "5"); // false（值和类型都要相等）
```

------

# **4. 条件语句**

```
js复制编辑let score = 85;
if (score >= 90) {
    console.log("优秀");
} else if (score >= 60) {
    console.log("及格");
} else {
    console.log("不及格");
}
```

**三元运算符**

```
js


复制编辑
let result = score >= 60 ? "及格" : "不及格";
```

------

# **5. 循环**

### **5.1 for 循环**

```
js复制编辑for (let i = 1; i <= 5; i++) {
    console.log(i);
}
```

### **5.2 while 循环**

```
js复制编辑let i = 1;
while (i <= 5) {
    console.log(i);
    i++;
}
```

------

# **6. 函数**

### **6.1 普通函数**

```
js复制编辑function greet(name) {
    return "你好，" + name;
}
console.log(greet("张三"));
```

### **6.2 箭头函数（ES6）**

```
js复制编辑const greet = (name) => `你好，${name}`;
console.log(greet("李四"));
```

------

# **7. 数组 & 对象**

### **7.1 数组**

```
js复制编辑let fruits = ["苹果", "香蕉", "橙子"];
console.log(fruits.length); // 获取长度
console.log(fruits[0]); // 访问元素
fruits.push("葡萄"); // 添加元素
console.log(fruits);
```

### **7.2 对象**

```
js复制编辑let person = {
    name: "张三",
    age: 25
};
console.log(person.name);
person.job = "程序员";
```

------

# **8. DOM 操作**

### **8.1 选择元素**

```
js复制编辑let title = document.getElementById("title");  // 通过 ID 选取
let paragraphs = document.querySelectorAll("p");  // 选取所有 p
```

### **8.2 修改内容**

```
js


复制编辑
title.innerText = "新标题";
```

### **8.3 监听事件**

```
js复制编辑document.getElementById("btn").addEventListener("click", function() {
    alert("按钮被点击了！");
});
```

------

# **9. 定时器**

### **9.1 `setTimeout`（延迟执行）**

```
js复制编辑setTimeout(() => {
    console.log("3 秒后执行");
}, 3000);
```

### **9.2 `setInterval`（循环执行）**

```
js复制编辑setInterval(() => {
    console.log("每秒执行一次");
}, 1000);
```

------

# **10. 异步 & Promise**

### **10.1 回调函数**

```
js复制编辑function fetchData(callback) {
    setTimeout(() => {
        callback("数据加载完成");
    }, 2000);
}
fetchData(console.log);
```

### **10.2 Promise**

```
js复制编辑let promise = new Promise((resolve, reject) => {
    setTimeout(() => resolve("成功"), 2000);
});
promise.then(console.log);
```

### **10.3 async/await**

```
js复制编辑async function fetchData() {
    let data = await new Promise(resolve => setTimeout(() => resolve("数据加载成功"), 2000));
    console.log(data);
}
fetchData();
```