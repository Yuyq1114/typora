------

## **1. CSS 基础**

### **1.1 CSS 语法**

CSS 由选择器（selector）、属性（property）和属性值（value）组成：

```
css复制编辑选择器 {
    属性: 值;
}
```

示例：

```
css复制编辑p {
    color: red;  /* 设置文本颜色为红色 */
    font-size: 16px;  /* 字体大小 */
}
```

------

### **1.2 引入 CSS**

有三种方式将 CSS 应用于 HTML：

1. 内联样式（Inline CSS）

   （写在 

   ```
   style
   ```

    属性中）

   ```
   html
   
   
   复制编辑
   <p style="color: blue;">这是蓝色文本</p>
   ```

2. 内部样式（Internal CSS）

   （写在 

   ```
   <style>
   ```

    标签内）

   ```
   html复制编辑<style>
       p { color: blue; }
   </style>
   ```

3. 外部样式（External CSS）

   （推荐，将样式写入 

   ```
   .css
   ```

    文件）

   ```
   html
   
   
   复制编辑
   <link rel="stylesheet" href="styles.css">
   ```

------

## **2. 选择器**

### **2.1 基础选择器**

- 标签选择器

  ```
  css
  
  
  复制编辑
  p { color: green; }
  ```

- 类选择器（`.class`）

  ```
  css
  
  
  复制编辑
  .highlight { color: red; }
  ```

- ID 选择器（`#id`，一个页面中唯一）

  ```
  css
  
  
  复制编辑
  #title { font-size: 24px; }
  ```

### **2.2 组合选择器**

- 后代选择器（空格）

   选择 

  ```
  div
  ```

   内的所有 

  ```
  p
  ```

  ```
  css
  
  
  复制编辑
  div p { color: blue; }
  ```

- 子选择器（`>`）

   选择 

  ```
  div
  ```

   直接子元素 

  ```
  p
  ```

  ```
  css
  
  
  复制编辑
  div > p { color: red; }
  ```

- 并集选择器（`,`）

   选择多个元素

  ```
  css
  
  
  复制编辑
  h1, h2, p { color: brown; }
  ```

------

## **3. 盒模型（Box Model）**

HTML 元素是一个盒子，包含：

1. **内容（content）**
2. **内边距（padding）**
3. **边框（border）**
4. **外边距（margin）**

```
css复制编辑.box {
    width: 200px;
    height: 100px;
    padding: 20px;
    border: 2px solid black;
    margin: 10px;
}
```

------

## **4. 布局**

### **4.1 浮动（Float，较老）**

```
css复制编辑img {
    float: left;
    margin-right: 10px;
}
```

### **4.2 Flexbox（现代布局，推荐）**

```
css复制编辑.container {
    display: flex;
    justify-content: space-between;
}
```

常见属性：

- `justify-content: center;`（居中）
- `justify-content: space-between;`（左右对齐）
- `align-items: center;`（垂直居中）

------

## **5. 颜色 & 背景**

```
css复制编辑body {
    background-color: #f0f0f0;
    color: #333;
}
```

- 背景图片

  ```
  css复制编辑body {
      background-image: url('background.jpg');
      background-size: cover;
  }
  ```

------

## **6. 文字 & 字体**

```
css复制编辑p {
    font-size: 16px;
    font-weight: bold;
    text-align: center;
    line-height: 1.5;
}
```

可以引入 **Google Fonts**：

```
html


复制编辑
<link href="https://fonts.googleapis.com/css2?family=Roboto&display=swap" rel="stylesheet">
css复制编辑body {
    font-family: 'Roboto', sans-serif;
}
```

------

## **7. 动画与过渡**

### **7.1 过渡**

```
css复制编辑button {
    background: blue;
    transition: background 0.3s ease-in-out;
}
button:hover {
    background: red;
}
```

### **7.2 动画**

```
css复制编辑@keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
}
div {
    animation: fadeIn 2s;
}
```