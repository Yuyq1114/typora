条件表达式是流程自动化引擎（如 Flowable、Activiti）中，用于控制流程的流转和决策的核心机制之一。条件表达式通常用于流程网关（如排他网关、并行网关等）中，以根据业务逻辑动态决定流程的下一步操作。Flowable 支持使用 Java 表达式语言（JUEL）和 Spring EL 来编写这些条件表达式。

------

### **条件表达式的应用场景**

1. **流程分支控制**：基于流程变量的值或者某些业务条件，决定流程下一步应该走的路径。
2. **条件网关**：在流程中，条件网关允许根据表达式的结果（布尔值）来决定流程的流向。通常，这些表达式会比较某个流程变量或执行某种条件判断。
3. **任务执行条件**：在用户任务或服务任务中，表达式可以控制某个任务是否被执行。

------

### **常见条件表达式格式**

Flowable 使用 `JUEL（Java Unified Expression Language）` 来编写表达式。典型的条件表达式在 BPMN 文件中通过 `${expression}` 的形式表示。

1. **简单变量比较**：

   ```
   <conditionExpression xsi:type="tFormalExpression">${amount > 1000}</conditionExpression>
   ```

   这段表达式会判断流程变量 `amount` 是否大于 1000。如果为真，流程将沿着该分支执行。

2. **流程变量检查**：

   ```
   <conditionExpression xsi:type="tFormalExpression">${approved == true}</conditionExpression>
   ```

   判断 `approved` 变量是否为 `true`，适用于审批流程等场景。

3. **字符串比较**：

   ```
   <conditionExpression xsi:type="tFormalExpression">${status == 'completed'}</conditionExpression>
   ```

   判断流程变量 `status` 是否等于 `'completed'`。

4. **方法调用**：

   ```
   <conditionExpression xsi:type="tFormalExpression">${userService.isEligible(userId)}</conditionExpression>
   ```

   这段表达式会调用 Spring Bean `userService` 的方法 `isEligible`，传入 `userId` 作为参数，然后根据方法返回的布尔值来决定流程的走向。

------

### **条件表达式在 BPMN 中的应用**

#### 1. **排他网关（Exclusive Gateway）**

排他网关（Exclusive Gateway）是一种常用的 BPMN 元素，它允许在多个互斥条件下选择一个满足条件的流程分支。排他网关依赖于条件表达式来决定走哪条流程分支。

**排他网关的 BPMN 示例**：

```
<exclusiveGateway id="gateway1" name="Evaluate Conditions" />
<sequenceFlow id="flow1" sourceRef="gateway1" targetRef="task1">
    <conditionExpression xsi:type="tFormalExpression">${amount <= 1000}</conditionExpression>
</sequenceFlow>
<sequenceFlow id="flow2" sourceRef="gateway1" targetRef="task2">
    <conditionExpression xsi:type="tFormalExpression">${amount > 1000}</conditionExpression>
</sequenceFlow>
```

在这个例子中：

- 如果 `amount` 小于或等于 1000，流程将走向 `task1`。
- 如果 `amount` 大于 1000，流程将走向 `task2`。

#### 2. **并行网关（Parallel Gateway）**

并行网关（Parallel Gateway）允许流程同时分发到多个分支，而不需要条件判断。通常并行网关不使用条件表达式，但如果你想在并行执行的基础上加入条件，可以结合使用其他网关元素。

------

### **条件表达式的复杂应用**

#### 1. **多条件组合**

你可以在表达式中结合多个条件，使用逻辑运算符来控制流程的走向：

- **与条件**：

  ```
  <conditionExpression xsi:type="tFormalExpression">${amount > 1000 && user == 'admin'}</conditionExpression>
  ```

  该条件表达式要求 `amount` 大于 1000 且 `user` 等于 `'admin'` 时，流程才能沿该分支执行。

- **或条件**：

  ```
  <conditionExpression xsi:type="tFormalExpression">${status == 'approved' || status == 'pending'}</conditionExpression>
  ```

  该条件表达式允许 `status` 为 `'approved'` 或 `'pending'` 时，流程执行。

#### 2. **调用外部服务或 Java 方法**

在条件表达式中可以调用外部服务的 Java 方法，获取动态数据来做判断。例如，通过 Spring Bean 里的方法决定流程分支：

```
<conditionExpression xsi:type="tFormalExpression">${userService.isAdmin(userId)}</conditionExpression>
```

这里，`userService.isAdmin(userId)` 方法会返回一个布尔值。如果用户是管理员，流程会沿着这个分支执行。

#### 3. **字符串处理和正则匹配**

表达式支持基本的字符串处理和正则匹配：

```
<conditionExpression xsi:type="tFormalExpression">${username matches '[a-zA-Z0-9]{6,}'}</conditionExpression>
```

这段表达式判断 `username` 是否匹配指定的正则表达式，检查其是否由 6 位及以上的字母或数字组成。

------

### **条件表达式与流程变量**

在条件表达式中，流程变量是核心要素。Flowable 引擎通过流程变量来存储和传递业务数据。

1. **设置流程变量**： 流程变量可以在流程启动时设置，也可以在流程执行过程中通过任务、服务等设置。设置流程变量通常通过 `execution.setVariable()` 或 `taskService` 来完成。

   例如，在服务任务中可以通过以下代码来设置变量：

   ```
   execution.setVariable("amount", 1500);
   ```

2. **获取流程变量**： 在条件表达式中，你可以通过 `${variableName}` 的形式来获取流程变量值。

   ```
   <conditionExpression xsi:type="tFormalExpression">${amount > 1000}</conditionExpression>
   ```

3. **流程变量的生命周期**： 流程变量的生命周期通常与流程实例相关。流程变量会随着流程实例的创建、流转、完成而产生和销毁。

------

### **调试条件表达式**

调试条件表达式时，可能会遇到以下常见问题：

- **语法错误**：表达式编写不正确，如缺少括号或使用了不支持的语法。
- **变量未定义**：表达式中引用的流程变量不存在或未初始化。
- **类型错误**：流程变量类型与表达式中的期望类型不匹配。

**解决方法**：

1. **检查流程变量的定义**：确保变量在流程中已经初始化，且与表达式中的使用方式匹配。
2. **使用 Flowable Admin 工具**：查看流程实例的日志和错误信息，特别是表达式执行时抛出的异常。
3. **日志调试**：在表达式执行之前，使用 `System.out.println` 或日志框架输出变量值，确认它们的状态是否符合预期。