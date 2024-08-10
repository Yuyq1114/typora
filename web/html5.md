### 1. **新语义元素**

HTML5 引入了一些新的语义化标签，用于更清晰地表示网页的结构和内容：

- `<header>`：定义网页的头部，通常包含导航、标题等。
- `<footer>`：定义网页的底部，通常包含版权信息、联系信息等。
- `<article>`：定义独立的内容块，通常是博客文章、新闻条目等。
- `<section>`：定义文档中的节，用于分组内容。
- `<nav>`：定义导航链接部分。
- `<aside>`：定义与主内容略相关的附加内容。
- `<figure>` 和 `<figcaption>`：定义图片或其他媒体内容及其说明。

### 2. **多媒体支持**

HTML5 内置了对多媒体的支持，无需使用插件即可嵌入音频和视频：

- **`<audio>`**：用于在网页中嵌入音频。支持多种音频格式，如 MP3、Ogg 和 WAV。

  ```
  <audio controls>
    <source src="audio.mp3" type="audio/mpeg">
    Your browser does not support the audio element.
  </audio>
  ```

- **`<video>`**：用于在网页中嵌入视频。支持多种视频格式，如 MP4、WebM 和 Ogg。

  ```
  <video width="320" height="240" controls>
    <source src="video.mp4" type="video/mp4">
    Your browser does not support the video tag.
  </video>
  ```

### 3. **表单控件**

HTML5 增强了表单控件，增加了新的输入类型和属性，使表单更具交互性和验证功能：

- **新输入类型**：如 `date`、`email`、`url`、`tel` 和 `range`，提供了更多输入方式和内置验证。

  ```
  <input type="date" name="birthdate">
  <input type="email" name="email">
  ```

- **表单验证**：新增了一些属性，如 `required`、`pattern` 和 `min`/`max`，用于在提交表单前验证用户输入。

  ```
  <input type="text" name="username" required pattern="[A-Za-z]{3,}">
  ```

### 4. **Canvas 元素**

HTML5 引入了 `<canvas>` 元素，用于在网页中绘制图形和动画。通过 JavaScript 访问 `<canvas>` 元素的绘图上下文，可以创建动态图形、游戏和其他可视效果。

```
<canvas id="myCanvas" width="200" height="100"></canvas>
<script>
  var canvas = document.getElementById('myCanvas');
  var ctx = canvas.getContext('2d');
  ctx.fillStyle = 'red';
  ctx.fillRect(10, 10, 150, 75);
</script>
```

### 5. **本地存储**

HTML5 提供了本地存储和会话存储，允许在用户的浏览器中存储数据，而不依赖于服务器。两者分别为 `localStorage` 和 `sessionStorage` 对象。

- **`localStorage`**：用于存储持久性的数据，数据在浏览器关闭后仍然存在。

  ```
  localStorage.setItem('key', 'value');
  console.log(localStorage.getItem('key'));
  ```

- **`sessionStorage`**：用于存储会话数据，数据在浏览器标签页关闭后会被删除。

  ```
  sessionStorage.setItem('key', 'value');
  console.log(sessionStorage.getItem('key'));
  ```

### 6. **地理定位**

HTML5 提供了 Geolocation API，用于获取用户的地理位置信息。这对于实现基于位置的服务和功能非常有用。

```
navigator.geolocation.getCurrentPosition(function(position) {
  console.log('Latitude: ' + position.coords.latitude);
  console.log('Longitude: ' + position.coords.longitude);
});
```

### 7. **Web 存储**

除了传统的 Cookies 外，HTML5 引入了 Web 存储，提供了更大容量的数据存储能力：

- **Local Storage**：存储的数据在关闭浏览器后仍然存在。
- **Session Storage**：存储的数据在浏览器会话期间有效，关闭标签页或窗口后数据会被清除。

### 8. **WebSocket**

HTML5 引入了 WebSocket API，允许与服务器进行持久性、双向通信，适用于实时数据交换，如在线聊天和实时更新。

```
var socket = new WebSocket('ws://example.com/socket');
socket.onopen = function() {
  socket.send('Hello Server!');
};
socket.onmessage = function(event) {
  console.log('Message from server:', event.data);
};
```

### 9. **Web Workers**

Web Workers 允许在后台线程中运行 JavaScript，以进行计算密集型任务而不阻塞主线程。这对于提高应用的性能非常有用。

```
// worker.js
self.onmessage = function(e) {
  self.postMessage('Received: ' + e.data);
};

// main.js
var worker = new Worker('worker.js');
worker.onmessage = function(e) {
  console.log('Message from worker:', e.data);
};
worker.postMessage('Hello Worker!');
```