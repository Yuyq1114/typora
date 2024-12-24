### **Liquibase Changelog 的原理**

Changelog 是 Liquibase 的核心文件，用于定义数据库的所有变更。它支持多种格式，包括：

- **XML**（推荐格式）
- **YAML**
- **JSON**
- **SQL**

这些文件中包含 **ChangeSet** 的定义，每个 ChangeSet 代表一个独立的、可回滚的数据库变更操作。

------

#### **Changelog 的执行流程**

1. **初始化并记录状态**：

   - Liquibase 在执行 Changelog 之前，会检查数据库中是否有 `DATABASECHANGELOG` 和 `DATABASECHANGELOGLOCK` 两张表。
   - 如果没有，则自动创建：
     - **`DATABASECHANGELOG`**：记录所有执行过的 ChangeSet。
     - **`DATABASECHANGELOGLOCK`**：用于防止多个进程同时操作 Changelog。

2. **解析 Changelog 文件**：

   - Liquibase 解析 Changelog 文件中的 ChangeSet。
   - 每个 ChangeSet 包含一个或多个数据库变更操作（例如创建表、添加列、插入数据等）。

3. **检查 ChangeSet 是否已执行**：

   - 每个 ChangeSet 都有一个唯一标识（通常由 `id`、`author` 和 `file` 组成）。

   - Liquibase 会检查 

     ```
     DATABASECHANGELOG
     ```

      表，看该 ChangeSet 是否已经执行过。

     - **已执行**：跳过。
     - **未执行**：执行 ChangeSet 并记录到 `DATABASECHANGELOG`。

4. **执行 ChangeSet**：

   - 按顺序执行未执行的 ChangeSet。
   - 每次执行后，将其标记为已执行，并将相关信息写入 `DATABASECHANGELOG`。

5. **锁机制防止并发问题**：

   - 在执行过程中，Liquibase 会锁定 `DATABASECHANGELOGLOCK` 表，防止其他 Liquibase 实例同时操作。
   - 执行完成后，释放锁。

------

#### **Changelog 文件的结构**

以下是一个典型的 XML 格式的 Changelog 文件示例：

```
xml复制代码<databaseChangeLog
    xmlns="http://www.liquibase.org/xml/ns/dbchangelog"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://www.liquibase.org/xml/ns/dbchangelog
                        http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-4.3.xsd">

    <changeSet id="1" author="user1">
        <createTable tableName="person">
            <column name="id" type="int" autoIncrement="true" primaryKey="true"/>
            <column name="first_name" type="varchar(50)"/>
            <column name="last_name" type="varchar(50)"/>
            <column name="birth_date" type="date"/>
        </createTable>
    </changeSet>

    <changeSet id="2" author="user1">
        <addColumn tableName="person">
            <column name="email" type="varchar(100)"/>
        </addColumn>
    </changeSet>

</databaseChangeLog>
```

------

### **关键机制**

1. **ChangeSet 的唯一性**：
   - 每个 ChangeSet 必须有唯一的 `id` 和 `author`。
   - 文件路径也会影响 ChangeSet 的唯一性。
2. **ChangeSet 的顺序执行**：
   - 按文件中的定义顺序或依赖关系顺序执行。
   - 可通过 `runOrder` 或文件结构调整顺序。
3. **幂等性**：
   - Liquibase 会检查 `DATABASECHANGELOG` 表，确保每个 ChangeSet 只执行一次，避免重复变更。
4. **支持回滚**：
   - Liquibase 支持为每个 ChangeSet 定义回滚逻辑（`rollback` 标签），在需要时可以撤销变更。
5. **支持多种数据库变更操作**：
   - `createTable`、`addColumn`、`dropTable`、`update`、`delete` 等等。
   - 也支持运行自定义 SQL 脚本。

------

### **优点**

1. **版本化管理**：
   - 每个 ChangeSet 都可以看作数据库变更的一个版本。
   - 可以方便地将数据库状态同步到代码版本。
2. **幂等与一致性**：
   - 保证每次执行结果一致，不会因重复执行而产生错误。
3. **回滚支持**：
   - 可以根据需求快速回滚到特定版本，适合开发与测试场景。
4. **多数据库支持**：
   - Liquibase 可以兼容多种数据库（MySQL、PostgreSQL、Oracle、SQL Server 等）。
5. **集成与自动化**：
   - 可以与 CI/CD 管道集成，实现数据库变更的自动化部署。

------

### **典型问题与解决方法**

1. **问题：Changelog 冲突**：
   - 多个开发人员同时修改 Changelog，可能引发冲突。
   - **解决**：使用分布式版本控制工具（如 Git）协作，定期合并和验证。
2. **问题：数据库状态与代码不一致**：
   - 数据库可能被手动修改，导致状态与 Changelog 不符。
   - **解决**：严格规范操作，只通过 Liquibase 修改数据库。
3. **问题：执行失败**：
   - 可能是由于语法错误、权限不足或数据库状态异常。
   - **解决**：检查 ChangeSet 的定义，并修复问题后重新执行。