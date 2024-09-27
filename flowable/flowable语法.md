Flowable 允许在业务流程中使用表达式（Expressions）来动态地控制流程的行为。表达式主要用于决策、动态任务分配、条件判断等场景。Flowable 支持 **JUEL（Java Unified Expression Language）** 表达式以及 **Spring EL（Spring Expression Language）** 表达式。

通过表达式，开发人员可以引用流程变量、执行复杂的条件判断、获取实体数据等，在流程模型中实现更灵活的业务逻辑。

------

### **表达式的应用场景**

1. **动态任务分配**：
   - 根据流程变量或其他条件动态分配任务给特定用户或用户组。
2. **条件网关**：
   - 在流程执行过程中，根据流程变量的值，决定流程的下一步流向（控制流程分支）。
3. **流程启动条件**：
   - 在启动流程时，依据条件来决定是否启动特定的子流程或进行某些操作。
4. **事件监听器**：
   - 在流程执行过程中，根据某些事件触发特定的业务逻辑。
5. **服务任务执行逻辑**：
   - 在执行服务任务时，使用表达式来调用 Java 方法、设置任务参数等。

------

### **常见的表达式类型**

1. **流程变量**：

   - 通过表达式可以读取或设置流程变量的值。流程变量是流程中存储数据的载体，可以在流程中传递和使用。

   例如：`${requesterName}` 可以获取名为 `requesterName` 的流程变量的值。

2. **属性访问**：

   - 表达式可以访问 Java 对象的属性或调用其方法。

   例如：`${employee.name}` 可以获取 `employee` 对象的 `name` 属性。

3. **调用方法**：

   - 表达式可以直接调用 Java 方法。

   例如：`${employee.getFullName()}` 可以调用 `employee` 对象的 `getFullName()` 方法。

4. **条件判断**：

   - 表达式可以用于条件网关或其他需要判断逻辑的场景。

   例如：`${amount > 1000}` 用于判断 `amount` 变量是否大于 1000。

5. **布尔值表达式**：

   - 通过返回 `true` 或 `false` 来控制流程的流转或任务分配。

   例如：`${approval == true}`。

6. **调用 Spring Bean**：

   - Flowable 与 Spring 集成后，可以通过表达式调用 Spring 容器中的 Bean。

   例如：`${myService.processRequest(request)}` 调用 `myService` Bean 的 `processRequest()` 方法。

------

### **常见的表达式用法**

#### 1. **任务动态分配**

Flowable 中的用户任务可以通过表达式动态分配给特定用户或用户组。可以通过表达式引用流程变量或方法来决定任务的负责人。

- **通过流程变量分配任务**：

  在用户任务中，配置 `assignee` 字段：

  ```
  <userTask id="task1" name="Approve Request" flowable:assignee="${approver}" />
  ```

  在上面的例子中，`approver` 是一个流程变量，任务将分配给 `approver` 变量的值所指定的用户。

- **通过方法分配任务**：

  如果要基于业务逻辑分配任务，可以调用方法来返回分配人：

  ```
  <userTask id="task2" name="Assign Task" flowable:assignee="${userService.getManager(requestId)}" />
  ```

  在这里，`userService.getManager(requestId)` 调用 `userService` Bean，返回特定任务的负责人。

#### 2. **条件网关**

条件网关用于控制流程的分支流转。在网关中可以使用表达式进行条件判断，决定下一步的执行路径。

- **简单条件判断**：

  ```
  <exclusiveGateway id="gateway1" />
  <sequenceFlow id="flow1" sourceRef="gateway1" targetRef="task1">
    <conditionExpression xsi:type="tFormalExpression">${amount <= 1000}</conditionExpression>
  </sequenceFlow>
  <sequenceFlow id="flow2" sourceRef="gateway1" targetRef="task2">
    <conditionExpression xsi:type="tFormalExpression">${amount > 1000}</conditionExpression>
  </sequenceFlow>
  ```

  上面的例子中，`amount` 是流程变量，根据 `amount` 的值选择 `task1` 或 `task2` 执行。

#### 3. **服务任务中的表达式**

服务任务是用于自动化执行的任务，例如调用外部服务或处理业务逻辑。可以使用表达式来调用 Java 方法或 Spring Bean。

- **调用 Java 方法**：

  ```
  <serviceTask id="serviceTask1" name="Process Order" flowable:expression="${orderService.process(orderId)}" />
  ```

  这里，`orderService` 是一个 Spring Bean，`process(orderId)` 是调用的业务方法。

- **设置输入参数和输出参数**：

  可以使用表达式为服务任务设置输入输出参数：

  ```
  <serviceTask id="serviceTask2" name="Calculate Discount">
    <extensionElements>
      <flowable:field name="discount">
        <flowable:expression>${discountService.calculateDiscount(order)}</flowable:expression>
      </flowable:field>
    </extensionElements>
  </serviceTask>
  ```

  在这个例子中，`discountService.calculateDiscount(order)` 调用服务来计算折扣，并将结果作为服务任务的输入参数。

#### 4. **事件监听器中的表达式**

Flowable 支持在流程中定义事件监听器，可以监听任务创建、完成等事件，并在事件触发时执行某些操作。表达式可以用于处理这些事件。

- **任务完成事件监听器**：

  ```
  <userTask id="approveTask" name="Approve Request">
    <extensionElements>
      <flowable:taskListener event="complete" expression="${notificationService.notifyApprover(approver)}" />
    </extensionElements>
  </userTask>
  ```

  在任务完成时，调用 `notificationService.notifyApprover(approver)` 方法，通知审批人任务已完成。

------

### **Flowable 支持的表达式类型**

- **${} 表达式**：表示普通的 Java 表达式或变量引用。
- **Spring EL 表达式**：当 Flowable 与 Spring 集成时，支持使用 Spring EL 表达式。

------

### **Flowable 表达式的配置位置**

1. **BPMN 2.0 流程定义**：
   - 可以直接在 BPMN 模型的 XML 文件中使用表达式。
2. **Flowable Modeler**：
   - 在 Flowable 的可视化流程建模工具中，可以在任务分配、网关条件、服务任务等元素中直接配置表达式。
3. **代码中设置表达式**：
   - 如果流程模型是通过代码动态生成的，可以通过 API 将表达式注入到流程定义中。

------

### **表达式调试**

在使用表达式时，可能会遇到以下常见问题：

- 表达式解析错误：如果表达式语法不正确或引用了不存在的变量，可能会导致流程执行失败。
- 类型转换错误：表达式中的变量类型与实际使用不匹配时，可能会抛出异常。

解决方法：

- 使用 Flowable Admin 工具查看流程实例的执行日志和错误信息。
- 检查表达式中使用的变量名和方法是否正确。
- 在开发环境中通过日志输出或断点调试表达式的执行结果。