Flowable 是一个 **轻量级的开源 BPM（Business Process Management，业务流程管理）平台**，支持 **BPMN 2.0、CMMN（Case Management Model and Notation）** 和 **DMN（Decision Model and Notation）** 三大业务流程、案例管理和决策建模标准。Flowable 提供了丰富的功能，帮助开发人员和业务团队定义、自动化和执行业务流程。

Flowable 是从更早的 Activiti 项目分支出来的，经过独立开发后形成了一个强大且独立的工作流管理和业务流程自动化解决方案。它具有较高的扩展性和灵活性，支持从简单的工作流到复杂的企业级业务流程的设计、执行和监控。

------

### **Flowable 的主要组件**

Flowable 平台包含几个核心模块：

1. **Flowable Engine**：Flowable 的核心组件，用于执行 BPMN 2.0 规范定义的业务流程。它是一个 Java 库，能够嵌入到任何 Java 应用程序中或作为独立的服务使用。
2. **Flowable Modeler**：Web 界面的流程建模工具，用于设计业务流程（基于 BPMN 2.0），并通过拖拽操作定义复杂的业务流程。
3. **Flowable Task**：一个基于 Web 的任务管理应用程序，允许用户管理和处理 Flowable 引擎执行过程中产生的任务。
4. **Flowable Admin**：一个管理员工具，用于监控和管理流程引擎、流程实例和系统配置。
5. **Flowable REST API**：提供与 Flowable 引擎的交互接口，通过 REST API 来启动流程、查询流程状态和管理任务。
6. **Flowable CMMN Engine**：处理 CMMN（Case Management Model and Notation），适用于需要动态决策的场景，如案件管理系统。
7. **Flowable DMN Engine**：执行 DMN（Decision Model and Notation），用于业务决策的自动化，常见于规则引擎的实现。

------

### **Flowable 的核心特性**

1. **BPMN 2.0 支持**：
   - Flowable 完全支持 BPMN 2.0 标准，允许用户使用 BPMN 2.0 的符号和元素来设计流程。
   - 提供了对各种 BPMN 元素的支持，如任务、事件、网关、子流程等，涵盖所有业务场景。
2. **流程设计**：
   - Flowable 提供了基于浏览器的 **Flowable Modeler**，用于设计和发布业务流程模型。
   - 流程模型可以导出为 XML 格式，符合 BPMN 2.0 标准。
3. **流程执行**：
   - **Flowable Engine** 执行 BPMN 2.0 流程，管理流程实例和任务的生命周期。
   - 支持用户任务的分配和管理，自动任务的执行，异步处理任务等。
4. **流程监控与管理**：
   - **Flowable Admin** 提供对流程实例的监控、任务管理、流程统计和异常处理。
   - 支持流程的暂停、恢复、终止和手动干预。
5. **任务管理**：
   - Flowable 支持通过 **Flowable Task** 进行任务分配、通知、审批等操作。
   - 用户可以通过 Web 应用查看并处理分配的任务，支持任务的优先级和角色分配。
6. **集成与扩展性**：
   - **REST API**：Flowable 提供全面的 REST API，可以集成到任何系统中，实现流程启动、任务查询、事件触发等功能。
   - **消息事件**：支持通过消息队列或事件驱动系统进行流程的异步执行和集成。
   - **多数据库支持**：Flowable 支持多种数据库，如 MySQL、PostgreSQL、Oracle、SQL Server 等。
7. **CMMN 与 DMN 支持**：
   - **CMMN**：用于基于事件或动态决策的案例管理，适合没有明确流程但需要灵活处理的场景，如客户支持、项目管理等。
   - **DMN**：用于执行业务规则和决策表，支持复杂的业务规则引擎，适合在流程中嵌入自动决策。
8. **集群与高可用**：
   - Flowable 支持分布式部署，能够在集群环境中实现高可用性和负载均衡，适合大规模企业应用场景。

------

### **Flowable 的流程开发与执行**

#### 1. **BPMN 流程模型设计**

Flowable 使用 BPMN 2.0 的流程模型来定义业务流程。常见的 BPMN 元素如下：

- **开始事件（Start Event）**：流程的起点。
- **用户任务（User Task）**：分配给用户手动处理的任务。
- **服务任务（Service Task）**：由系统自动执行的任务，如调用 Web 服务、执行数据库操作等。
- **排他网关（Exclusive Gateway）**：用于分支和合并流程流。
- **结束事件（End Event）**：流程的终点。

在 Flowable Modeler 中，业务人员可以通过拖拽和配置这些元素，快速完成业务流程的设计。

#### 2. **Flowable 引擎执行流程**

设计好的 BPMN 2.0 模型可以部署到 Flowable 引擎上。流程的执行过程如下：

- **启动流程**：通过 REST API 或手动触发流程实例，流程按照模型定义的顺序开始执行。
- **处理任务**：当流程到达用户任务时，会创建任务并分配给指定用户，用户通过 Web 界面完成任务处理。
- **自动任务执行**：如果是服务任务，Flowable 引擎会自动执行预定义的操作，如调用外部系统或执行内部逻辑。
- **流程结束**：当所有任务完成，流程会到达结束事件并终止，记录流程的执行日志。

#### 3. **使用 REST API 操作流程**

Flowable 的 REST API 允许外部应用程序与 Flowable 引擎交互。常见的 API 操作包括：

- 启动流程实例
- 查询流程状态和任务
- 完成任务
- 取消或终止流程实例

------

### **Flowable 语法示例**

以下是一个简单的 BPMN 2.0 流程的 XML 示例，定义了一个简单的用户审批流程：

```
<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL"
             xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
             xsi:schemaLocation="http://www.omg.org/spec/BPMN/20100524/MODEL BPMN20.xsd"
             targetNamespace="http://www.example.org/bpmn">
  <process id="approvalProcess" name="Approval Process">
    <startEvent id="start" name="Start"/>
    <userTask id="submitRequest" name="Submit Request" flowable:assignee="employee"/>
    <exclusiveGateway id="gateway"/>
    <userTask id="approveRequest" name="Approve Request" flowable:assignee="manager"/>
    <endEvent id="end" name="End"/>
    <sequenceFlow id="flow1" sourceRef="start" targetRef="submitRequest"/>
    <sequenceFlow id="flow2" sourceRef="submitRequest" targetRef="gateway"/>
    <sequenceFlow id="flow3" sourceRef="gateway" targetRef="approveRequest"/>
    <sequenceFlow id="flow4" sourceRef="approveRequest" targetRef="end"/>
  </process>
</definitions>
```

------

### **Flowable 的实际应用场景**

1. **审批工作流**：如员工请假、报销流程，支持动态的审批规则。
2. **订单管理**：自动化处理订单状态的更新、通知和审批。
3. **客户支持**：基于案例的支持系统，动态处理客户的请求。
4. **自动化任务调度**：集成多个系统自动执行任务，减少人为操作。