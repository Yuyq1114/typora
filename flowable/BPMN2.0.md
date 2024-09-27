BPMN 2.0（**Business Process Model and Notation 2.0**）是由 **OMG（Object Management Group）** 推出的业务流程建模语言标准，旨在为业务人员和技术人员提供一种可视化的、易于理解且标准化的流程设计语言。BPMN 2.0 广泛应用于业务流程管理（BPM）领域，支持从流程建模到执行和监控的全生命周期管理。

BPMN 2.0 的主要目标是提供一种标准化的方法来定义和沟通业务流程的各个方面，帮助业务分析师、开发人员、流程设计者等不同角色协同工作。它通过图形化符号表示流程逻辑，标准化了流程的建模和执行。

### 1. **BPMN 2.0 基础概念**

BPMN 2.0 的核心是 **流程图**，它由各种类型的 **图形元素** 组成，这些元素用于描述业务活动的执行流。BPMN 2.0 规范主要包括以下几个部分：

- **流程元素**：用于建模业务活动的基本符号。
- **流程类型**：描述不同类型的业务流程。
- **语法规则**：定义如何组合和连接这些流程元素。

#### BPMN 2.0 提供了四大类别的图形元素：

1. **流程对象（Flow Objects）**
2. **连接对象（Connecting Objects）**
3. **泳道（Swimlanes）**
4. **构件（Artifacts）**

------

### 2. **BPMN 2.0 的关键元素**

#### 2.1 流程对象（Flow Objects）

流程对象是 BPMN 2.0 的基本构建块，用于描述业务流程中的活动、事件和网关。主要包括：

##### a. **事件（Event）**

事件是流程的触发点或结果。它们是圆形的符号，有三种主要类型：

- **开始事件（Start Event）**：标志流程的起点。它可以是消息事件、定时器事件或信号事件。

![](image\start.png)

- **中间事件（Intermediate Event）**：表示流程中的中断或事件。它可以挂在流程线上，表示某个过程中会发生的某种事件，如消息接收或定时触发。

![](image\intermediate.png)

- **结束事件（End Event）**：标志流程的终点，可能表示流程的完成或终止。

![](image\end.png)

**事件类型示例**：

| 事件类型   | 符号     | 描述                               |
| ---------- | -------- | ---------------------------------- |
| 开始事件   | 空心圆圈 | 触发流程的开始                     |
| 结束事件   | 实心圆圈 | 标志流程结束，通常不再触发其他动作 |
| 定时器事件 | 时钟图标 | 按计划或延迟时间触发事件           |
| 消息事件   | 信封图标 | 发送或接收消息时触发               |
| 错误事件   | 闪电图标 | 发生错误时触发                     |

##### b. **活动（Activity）**

![](image\activity.png)

活动是流程中的任务或工作单元，矩形表示：

- **任务（Task）**：表示流程中的单一工作。常见类型包括用户任务、服务任务、脚本任务等。

![](image\task.png)

- **子流程（Sub-process）**：可以将流程逻辑封装为子流程，支持嵌套的层次结构。

![](image\subprocess.png)

- **事务（Transaction）**：用于表示需要满足所有步骤都完成才能提交的流程。

**常见任务类型**：

| 任务类型     | 描述                                      |
| ------------ | ----------------------------------------- |
| 用户任务     | 需要用户手动执行的任务，例如表单填写      |
| 服务任务     | 由系统自动执行的任务，例如 API 调用       |
| 脚本任务     | 执行嵌入式脚本（如 JavaScript 或 Groovy） |
| 业务规则任务 | 调用业务规则引擎来进行判断和决策          |
| 发送任务     | 向外部系统发送消息                        |

##### c. **网关（Gateway）**

![](image\gateway.png)

网关用于控制流程的分支、合并或同步。常见的网关有：

![](image\gateways.png)

- **排他网关（Exclusive Gateway）**：根据条件选择一条路径，只有一条路径可以执行。
- **并行网关（Parallel Gateway）**：所有分支同时执行，不依赖条件。
- **事件网关（Event-Based Gateway）**：等待事件的发生，然后选择路径。
- **包容性网关（Inclusive Gateway）**：可以选择多条路径同时执行，所有路径必须完成才能汇合。

------

#### 2.2 连接对象（Connecting Objects）

连接对象用于连接不同的流程元素，显示它们之间的关系。主要包括：

- **顺序流（Sequence Flow）**：表示活动之间的执行顺序，用实线表示。

![](image\sequence.png)

- **消息流（Message Flow）**：用于表示两个参与者之间的通信或消息交换，通常用于不同的泳道之间。

![](image\message.png)

- **关联（Association）**：将注释或数据与流程对象关联。

![](image\association.png)

**连接对象示例**：

| 连接类型 | 描述                                             |
| -------- | ------------------------------------------------ |
| 顺序流   | 定义任务或事件的执行顺序，常用实线箭头表示       |
| 消息流   | 表示不同泳道或参与者之间的消息传递，虚线箭头表示 |
| 关联     | 连接附加数据、注释或文本到特定任务或事件         |

------

#### 2.3 泳道（Swimlanes）

泳道用于表示不同的角色或部门，它们负责某些流程活动。泳道可以帮助清晰地定义流程责任和任务归属。

- **泳道（Lane）**：泳道是水平或垂直的区域，表示不同的角色或参与者。

![](image\lane.png)

- **池（Pool）**：泳道可以组合成池，表示一个完整的组织或系统。

![](image\pool.png)

------

#### 2.4 构件（Artifacts）

构件用于补充业务流程图的信息，不影响流程执行。常见的构件包括：

- **数据对象（Data Object）**：表示任务需要的数据或产生的输出。
- **注释（Annotation）**：用于添加说明性的文本，帮助解释流程。

![](image\annotation.png)

- **组（Group）**：将流程中的相关元素分组。

![](image\group.png)

------

### 3. **BPMN 2.0 的流程类型**

BPMN 2.0 提供了多种流程类型，用于应对不同的业务场景：

- **顺序流程（Sequential Flow）**：表示任务按固定的顺序执行。
- **并行流程（Parallel Flow）**：多个任务可以同时执行，通常通过并行网关实现。
- **子流程（Sub-process）**：封装为子流程，可以在主流程中复用。
- **事件驱动流程（Event-driven Flow）**：流程的执行依赖于事件的触发，例如消息或定时器。

------

### 4. **BPMN 2.0 与业务自动化**

BPMN 2.0 不仅是业务流程的建模工具，还能够直接驱动业务自动化系统。在业务自动化工具中（如 Flowable、Camunda 等），BPMN 2.0 定义的流程可以被直接执行，支持：

- **任务调度和分配**：将任务分配给具体的用户或系统角色。
- **业务规则引擎**：自动根据规则或条件判断流程路径。
- **事件处理**：监听事件触发的流程活动。
- **流程监控和优化**：通过 BPMN 2.0 模型监控流程执行，优化流程性能。

### 5. **BPMN 2.0 XML 语法示例**

以下是一个简单的 BPMN 2.0 XML 语法示例，表示一个用户审批流程：

```
<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL"
             xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
             xsi:schemaLocation="http://www.omg.org/spec/BPMN/20100524/MODEL BPMN20.xsd"
             targetNamespace="http://www.example.org/bpmn">

  <process id="approvalProcess" name="Approval Process">
    <!-- 开始事件 -->
    <startEvent id="start" name="Start" />

    <!-- 用户任务 -->
    <userTask id="submitRequest" name="Submit Request" flowable:assignee="employee" />

    <!-- 排他网关 -->
    <exclusiveGateway id="gateway" />

    <!-- 审批任务 -->
    <userTask id="approveRequest" name="Approve Request" flowable:assignee="manager" />

    <!-- 结束事件 -->
    <endEvent id="end" name="End" />

    <!-- 顺序流 -->
    <sequenceFlow id="flow1" sourceRef="start" targetRef="submitRequest" />
    <sequenceFlow id="flow2" sourceRef="submitRequest" targetRef="gateway" />
    <sequenceFlow id="flow3" sourceRef="gateway" targetRef="approveRequest" />
    <sequenceFlow id="flow4" sourceRef="approveRequest" targetRef="end" />

  </process>
</definitions>
```

在这个例子中：

- **StartEvent**：流程开始。
- **UserTask（Submit Request）**：员工提交请求。
- **Exclusive Gateway**：判断请求是否通过。
- **UserTask（Approve Request）**：经理审批请求。
- **EndEvent**：流程结束。

------

### 6. **BPMN 2.0 与若依平台的关系**

在若依平台中，流程管理模块借助了 BPMN 2.0 来定义和管理业务流程。通过 BPMN 2.0，若依能够提供流程可视化、动态审批和流程自动化功能，确保业务流程的透明性和灵活性。