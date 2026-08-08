# Spec-kit

constitution
定项目章程/原则（治理约束）
.specify/memory/constitution.md

specify
把想法写成功能规范，开功能分支
specs/N-xxx/spec.md

clarify
找高影响模糊点，问答后写回规范
更新后的 spec.md

plan
技术设计与实施规划
plan.md、research.md、data-model.md、contracts/、quickstart.md

tasks
按依赖拆可执行任务，标并行/MVP
tasks.md

analyze
实施前一致性/覆盖率审查（只读）
分析报告（默认不改文件）

implement
按任务实现并验证
代码、测试；勾选 tasks.md



constitution（项目级，通常一次）
        ↓
specify      → 生成 spec.md（需求）
        ↓
clarify     → 澄清模糊点，写回 spec.md（plan 前建议做）
        ↓
plan        → 生成 plan.md + research.md + data-model.md + contracts/ + quickstart.md
        ↓
tasks       → 生成 tasks.md（可执行任务清单）
        ↓
analyze     → 只读审查 spec/plan/tasks 是否一致（实施前）
        ↓
implement   → 按 tasks.md 写代码、测试、勾任务





# Superpower

 /using-superpowers 默认调用 /brainstorming 生成design.md /using-git-worktrees 提交git（大部分都会用到 这里是首
次） /writing-plans 生成implementation.md /subagent-driven-development 生成代码/requesting-code-review 生成review.txt  如果有问题自动修复  /receiving-code-review 对审查结果逐项修复 /verification-before-completion 勾选完成>任务前进行验证 /test-driven-development TDD思路，生成测试文件 /systematic-debugging 系统化排查方法



其中手动执行的有

/using-superpowers 默认调用

/brainstorming 生成design.md

/writing-plans 生成implementation.md

/subagent-driven-development 生成代码



using-superpowers          ← 总开关（默认开启，先匹配再动手）
        ↓
brainstorming              ←【手动】设计定稿
        ↓  design.md
using-git-worktrees        ←【首次/需隔离时】独立工作区 + git 准备
        ↓
writing-plans              ←【手动】实施计划
        ↓  implementation.md（或 plans/*.md）
subagent-driven-development←【手动】按计划生成代码
        │
        ├─（实现中自动/强制穿插）
        │   test-driven-development      → 测试文件 + 红绿重构
        │   verification-before-completion → 完成宣称前验证
        │   systematic-debugging         → 遇 bug/失败时先根因
        ↓
requesting-code-review     → review.txt（或评审报告）
        ↓ 有问题则自动/逐项修
receiving-code-review      → 对审查项逐条修复
        ↓
verification-before-completion → 最终验证
        ↓
finishing-a-development-branch → merge / PR / 清理（可选收尾）



## bug处理

1. `systematic-debugging` 这是 Bug 的第一入口，必须先用。流程是：

   读取错误信息

   -> 稳定复现

   -> 检查最近改动

   -> 追踪数据流和组件边界

   -> 找到根因

   不要一看到报错就直接改代码。

2. `test-driven-development` 找到根因后，先写一个能复现 Bug 的失败测试：

   写失败测试

   -> 运行并确认确实失败

   -> 修改最少代码

   -> 运行测试确认通过

   -> 运行完整回归测试

3. `verification-before-completion` 修复后、准备声称“已修复”之前使用。需要重新执行实际验证命令，不能只根据代码修改或代理报告判断。

4. `requesting-code-review` Bug 修复涉及多个服务、数据库、认证、Kubernetes 或公共模块时，建议追加代码审查。

5. `finishing-a-development-branch` Bug 修复完成、测试通过后，用于决定：

   - 合并到 `master`
   - 推送并创建 PR
   - 保留分支
   - 丢弃分支