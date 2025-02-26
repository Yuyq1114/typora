# **1. React 的基本概念**

### **1.1 组件**

React 主要基于组件（Component）来构建应用。组件是 UI 的基本构建块，可以是类组件或函数组件。每个组件包含自己的状态（state）和行为（事件处理）。

- **函数组件**：最简单的组件类型，使用 JavaScript 函数定义。

  ```
  js复制编辑function Welcome(props) {
    return <h1>Hello, {props.name}</h1>;
  }
  ```

- **类组件**：通过继承 `React.Component` 类来定义组件，适用于需要使用生命周期方法和本地状态的情况。

  ```
  js复制编辑class Welcome extends React.Component {
    render() {
      return <h1>Hello, {this.props.name}</h1>;
    }
  }
  ```

### **1.2 JSX**

JSX 是 JavaScript 的一种语法扩展，它允许你在 JavaScript 代码中写 HTML 结构。React 使用 JSX 来定义组件的渲染内容。

```
js


复制编辑
const element = <h1>Hello, world!</h1>;
```

JSX 会被 Babel 转换成 `React.createElement()` 函数调用。

### **1.3 Props（属性）**

Props（属性）是从父组件传递给子组件的数据。Props 是只读的，子组件不能直接修改它们，只能使用它们来渲染内容或传递给其他组件。

```
js复制编辑function Greeting(props) {
  return <h1>Hello, {props.name}</h1>;
}

function App() {
  return <Greeting name="John" />;
}
```

### **1.4 State（状态）**

State 是组件内部的数据或状态，是一个动态的对象。当 state 改变时，React 会重新渲染组件。你可以通过 `setState()` 更新组件的 state。

```
js复制编辑class Counter extends React.Component {
  constructor(props) {
    super(props);
    this.state = { count: 0 };
  }

  increment() {
    this.setState({ count: this.state.count + 1 });
  }

  render() {
    return (
      <div>
        <p>{this.state.count}</p>
        <button onClick={() => this.increment()}>Increment</button>
      </div>
    );
  }
}
```

### **1.5 事件处理**

React 通过 `on<EventName>` 来处理用户事件，例如点击、输入等。事件处理函数通常是组件的方法。

```
js复制编辑class Button extends React.Component {
  handleClick() {
    alert('Button clicked!');
  }

  render() {
    return <button onClick={this.handleClick}>Click me</button>;
  }
}
```

------

# **2. React 的核心特性**

### **2.1 单向数据流**

React 采用单向数据流，意味着数据从父组件流向子组件。父组件通过 props 将数据传递给子组件，子组件只能通过调用父组件传递的回调函数来修改数据。

### **2.2 虚拟 DOM**

React 使用虚拟 DOM 来优化渲染过程。每次组件的 state 或 props 发生变化时，React 会首先在内存中更新虚拟 DOM，然后与实际的 DOM 进行比较，最后仅更新发生变化的部分。这种优化提高了性能。

### **2.3 组件生命周期**

React 组件有一个生命周期，从组件的创建、更新到销毁，React 提供了一些生命周期方法来在这些阶段进行操作。生命周期方法包括：

- `componentDidMount()`：组件首次渲染后调用，适合进行数据获取等操作。
- `shouldComponentUpdate()`：决定组件是否需要更新，返回 `true` 或 `false`。
- `componentDidUpdate()`：组件更新后调用。
- `componentWillUnmount()`：组件销毁前调用，用于清理资源。

对于函数组件，React 提供了 **Hooks**，使得你可以在函数组件中使用类似生命周期的功能。

------

# **3. React Hooks**

React 16.8 引入了 Hooks，允许你在函数组件中使用 state 和其他 React 特性，避免了类组件的复杂性。常用的 Hooks 包括：

### **3.1 `useState`**

`useState` 用于在函数组件中声明状态变量。

```
js复制编辑import { useState } from 'react';

function Counter() {
  const [count, setCount] = useState(0);

  return (
    <div>
      <p>{count}</p>
      <button onClick={() => setCount(count + 1)}>Increment</button>
    </div>
  );
}
```

### **3.2 `useEffect`**

`useEffect` 用于处理副作用（side effects），例如数据获取、订阅等。它相当于类组件中的 `componentDidMount`、`componentDidUpdate` 和 `componentWillUnmount`。

```
js复制编辑import { useState, useEffect } from 'react';

function Timer() {
  const [seconds, setSeconds] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => setSeconds(seconds => seconds + 1), 1000);

    return () => clearInterval(interval);  // 清理副作用
  }, []);  // 依赖数组为空，表示只会在组件挂载时运行

  return <p>Elapsed time: {seconds}s</p>;
}
```

### **3.3 `useContext`**

`useContext` 用于访问 React 的 Context，用来在组件树中传递数据，而不需要手动传递 props。

```
js复制编辑import { useContext } from 'react';

const MyContext = React.createContext();

function ChildComponent() {
  const value = useContext(MyContext);
  return <p>{value}</p>;
}

function App() {
  return (
    <MyContext.Provider value="Hello from Context">
      <ChildComponent />
    </MyContext.Provider>
  );
}
```

### **3.4 `useReducer`**

`useReducer` 是 `useState` 的一个替代方案，用于管理复杂的状态逻辑，类似于 Redux。

```
js复制编辑import { useReducer } from 'react';

const initialState = { count: 0 };

function reducer(state, action) {
  switch (action.type) {
    case 'increment':
      return { count: state.count + 1 };
    case 'decrement':
      return { count: state.count - 1 };
    default:
      throw new Error();
  }
}

function Counter() {
  const [state, dispatch] = useReducer(reducer, initialState);

  return (
    <div>
      <p>{state.count}</p>
      <button onClick={() => dispatch({ type: 'increment' })}>Increment</button>
      <button onClick={() => dispatch({ type: 'decrement' })}>Decrement</button>
    </div>
  );
}
```

------

# **4. 路由与 React**

在 React 中，我们可以使用 **React Router** 来处理单页面应用（SPA）中的路由导航。

### **4.1 安装 React Router**

```
bash


复制编辑
npm install react-router-dom
```

### **4.2 设置路由**

```
js复制编辑import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';

function Home() {
  return <h2>Home Page</h2>;
}

function About() {
  return <h2>About Page</h2>;
}

function App() {
  return (
    <Router>
      <Switch>
        <Route path="/home" component={Home} />
        <Route path="/about" component={About} />
      </Switch>
    </Router>
  );
}
```

------

# **总结**

React 通过组件化、虚拟 DOM 和 Hooks 等特性，让开发者能够高效地构建用户界面。它强调可重用性、性能优化和灵活性，是现代前端开发中非常流行的框架。