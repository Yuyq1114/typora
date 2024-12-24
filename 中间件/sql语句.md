## **2. SQL 的功能分类**

### **1. 数据定义语言 (DDL: Data Definition Language)**

**定义**：
DDL 用于定义、修改或删除数据库的结构，包括表、索引、视图等。

#### 常用命令：

1. **CREATE**：创建数据库或数据库对象（如表、索引、视图）。
2. **ALTER**：修改数据库对象的结构。
3. **DROP**：删除数据库对象。
4. **TRUNCATE**：快速清空表中的所有数据，但保留表结构。

#### 示例代码：

- 创建数据库和表：

  ```
  sql复制代码CREATE DATABASE Company;
  
  CREATE TABLE employees (
      id INT PRIMARY KEY,
      name VARCHAR(50),
      department VARCHAR(30),
      salary DECIMAL(10, 2)
  );
  ```

- 修改表结构：

  ```
  sql复制代码ALTER TABLE employees ADD hire_date DATE; -- 添加新列
  ALTER TABLE employees DROP COLUMN hire_date; -- 删除列
  ```

- 删除表：

  ```
  sql
  
  
  复制代码
  DROP TABLE employees;
  ```

- 清空表：

  ```
  sql
  
  
  复制代码
  TRUNCATE TABLE employees; -- 快速清空表数据
  ```

------

### **2. 数据操作语言 (DML: Data Manipulation Language)**

**定义**：
DML 用于对表中的数据进行增、删、改操作。

#### 常用命令：

1. **INSERT**：向表中插入新数据。
2. **UPDATE**：更新表中的现有数据。
3. **DELETE**：删除表中的数据。

#### 示例代码：

- 插入数据：

  ```
  sql复制代码INSERT INTO employees (id, name, department, salary)
  VALUES (1, 'Alice', 'HR', 50000.00);
  ```

- 更新数据：

  ```
  sql复制代码UPDATE employees
  SET salary = 55000.00
  WHERE id = 1; -- 更新 ID 为 1 的员工薪资
  ```

- 删除数据：

  ```
  sql复制代码DELETE FROM employees
  WHERE department = 'HR'; -- 删除 HR 部门的员工
  ```

------

### **3. 数据查询语言 (DQL: Data Query Language)**

**定义**：
DQL 用于查询数据库中的数据，是 SQL 的核心功能。

#### 常用命令：

1. **SELECT**：从一个或多个表中检索数据。
2. **子句支持**：包括 WHERE、ORDER BY、GROUP BY、HAVING、JOIN 等。

#### 示例代码：

- 基本查询：

  ```
  sql
  
  
  复制代码
  SELECT name, salary FROM employees;
  ```

- 条件查询：

  ```
  sql复制代码SELECT name, salary 
  FROM employees 
  WHERE salary > 40000;
  ```

- 排序和分组：

  ```
  sql复制代码SELECT department, AVG(salary) AS avg_salary
  FROM employees
  GROUP BY department
  HAVING AVG(salary) > 50000
  ORDER BY avg_salary DESC;
  ```

- 多表连接：

  ```
  sql复制代码SELECT e.name, d.department_name
  FROM employees e
  INNER JOIN departments d ON e.department = d.id;
  ```

------

### **4. 数据控制语言 (DCL: Data Control Language)**

**定义**：
DCL 用于控制用户对数据库的访问权限。

#### 常用命令：

1. **GRANT**：授予用户权限。
2. **REVOKE**：撤销用户权限。

#### 示例代码：

- 授予权限：

  ```
  sql
  
  
  复制代码
  GRANT SELECT, INSERT ON employees TO user1;
  ```

- 撤销权限：

  ```
  sql
  
  
  复制代码
  REVOKE INSERT ON employees FROM user1;
  ```

------

### **5. 事务控制语言 (TCL: Transaction Control Language)**

**定义**：
TCL 用于管理事务，确保数据的一致性和完整性。事务是指一组 SQL 操作，它们作为一个单元被提交或回滚。

#### 常用命令：

1. **COMMIT**：提交事务，将更改保存到数据库。
2. **ROLLBACK**：回滚事务，撤销未提交的更改。
3. **SAVEPOINT**：设置事务的回滚点。
4. **SET TRANSACTION**：配置事务属性，如隔离级别。

#### 示例代码：

- 开始事务、提交和回滚：

  ```
  sql复制代码BEGIN TRANSACTION;
  UPDATE employees SET salary = 60000.00 WHERE id = 1;
  COMMIT; -- 提交事务
  ```

  ```
  sql复制代码BEGIN TRANSACTION;
  UPDATE employees SET salary = 50000.00 WHERE id = 1;
  ROLLBACK; -- 撤销更改
  ```

- 使用保存点：

  ```
  sql复制代码BEGIN TRANSACTION;
  UPDATE employees SET salary = 70000.00 WHERE id = 1;
  SAVEPOINT savepoint1; -- 设置保存点
  UPDATE employees SET salary = 80000.00 WHERE id = 2;
  ROLLBACK TO savepoint1; -- 回滚到保存点
  COMMIT;
  ```

------

## **3. SQL 的主要特性**

#### **1. 声明式**

SQL 使用声明式语法，描述“做什么”而不是“如何做”。

#### **2. 与数据库无关**

SQL 是一个标准语言，支持大多数关系数据库管理系统（如 MySQL、PostgreSQL、Oracle、Microsoft SQL Server）。

#### **3. 强大的查询能力**

支持复杂的多表查询、子查询、聚合函数等。

#### **4. 支持事务**

保证操作的原子性、一致性、隔离性和持久性（ACID）。

------

## **4. SQL 的常用子句和功能**

#### **1. SELECT 子句**

- **WHERE**：条件筛选。
- **ORDER BY**：排序。
- **GROUP BY**：分组。
- **HAVING**：对分组数据进行过滤。
- **JOIN**：多表连接（INNER JOIN、LEFT JOIN、RIGHT JOIN、FULL JOIN）。

示例：

```
sql复制代码SELECT department, AVG(salary) AS avg_salary
FROM employees
GROUP BY department
HAVING AVG(salary) > 50000
ORDER BY avg_salary DESC;
```

#### **2. 函数支持**

- **聚合函数**：SUM、AVG、COUNT、MAX、MIN。
- **字符串函数**：CONCAT、SUBSTRING、UPPER、LOWER。
- **日期函数**：NOW、DATE_ADD、DATEDIFF。
- **数学函数**：ROUND、ABS、RAND。