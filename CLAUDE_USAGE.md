# Claude 自动化工作流指南

本文档说明如何让 Claude 在 Claude Code 中自动遵循 plansm 计划执行，**无需用户手动干预**。

## 核心理念

**目标**：Claude 读取 plan.json → 自动执行当前步骤 → 自动验证 → 自动推进 → 循环直到完成

**用户角色**：只需编写初始计划，然后让 Claude 自动执行

## 设置步骤（一次性）

### 1. 复制插件到项目

```bash
# 零安装模式（推荐）
cd your-project
cp -r ~/path/to/plansm/.claude-plugin .

# 确保安装了 jq
brew install jq  # macOS
```

### 2. 创建初始计划

编辑 `plan.json`：

```json
{
  "version": 1,
  "current_step": "STEP_001",
  "invariants": [
    "Do not mark VERIFIED without running verification.",
    "Only work on current_step unless explicitly unlocked.",
    "Always verify before advancing."
  ],
  "steps": [
    {
      "id": "STEP_001",
      "objective": "Implement user authentication",
      "status": "PENDING",
      "allow_paths": ["src/auth/**", "tests/**"],
      "verify": [
        {"type": "command", "cmd": "npm test auth", "expect": {"exit_code": 0}}
      ]
    },
    {
      "id": "STEP_002",
      "objective": "Add API rate limiting",
      "status": "LOCKED",
      "depends_on": ["STEP_001"],
      "verify": [
        {"type": "file_contains", "file": "src/middleware/rate-limit.js", "pattern": "rateLimit"},
        {"type": "command", "cmd": "npm test rate-limit"}
      ]
    }
  ]
}
```

### 3. 配置 Claude Code Hook（强烈推荐）

在 Claude Code 中运行 `/hooks` 并添加：

**Stop Hook**:
```
Command: bash .claude-plugin/scripts/verify.sh --current
Blocking: true
```

这会防止 Claude 在验证失败时停止。

## Claude 自动化工作流

### 启动方式

在 Claude Code 中，给 Claude 这个指令：

```
请按照 plan.json 执行任务。流程如下：

1. 运行 /pwork 查看当前步骤
2. 实现当前步骤的目标
3. 运行 /pverify 验证（必须通过）
4. 运行 /pnext 推进到下一步
5. 重复直到所有步骤完成

规则：
- 不允许手动编辑 plan.json 的 status 字段
- 验证失败时必须修复后重新验证，不能跳过
- 只关注 current_step，忽略其他步骤
```

### Claude 的自动执行循环

Claude 会自动执行以下循环：

```
┌─────────────────────────────────────┐
│ 1. /pwork - 查看当前步骤            │
│    → 显示: 目标、允许路径、验证规则│
└─────────────┬───────────────────────┘
              ↓
┌─────────────────────────────────────┐
│ 2. 实现功能                         │
│    → 写代码、修改文件                │
│    → 只在 allow_paths 内工作        │
└─────────────┬───────────────────────┘
              ↓
┌─────────────────────────────────────┐
│ 3. /pverify - 运行验证              │
│    → 执行测试、检查文件              │
│    → 如果失败：返回步骤2修复         │
│    → 如果成功：继续                  │
└─────────────┬───────────────────────┘
              ↓
┌─────────────────────────────────────┐
│ 4. /pnext - 推进到下一步            │
│    → 解锁依赖步骤                    │
│    → 更新 current_step              │
└─────────────┬───────────────────────┘
              ↓
    是否还有待完成步骤？
         ↓ 是 → 返回步骤1
         ↓ 否 → 完成

```

## 实际示例对话

### 示例 1：简单功能实现

**用户**:
```
请按照 plan.json 实现登录功能
```

**Claude 的自动执行**:
```
我来查看当前步骤：
> /pwork

CURRENT_STEP: STEP_001
OBJECTIVE: Implement login API
VERIFY:
  - command: npm test -- login
  - file_contains: src/api/login.js pattern=export.*login

好的，我会实现登录 API。
[实现代码...]

让我验证实现：
> /pverify

STEP STEP_001: OK
  ✓ command: npm test -- login
  ✓ file_contains: src/api/login.js

验证通过！推进到下一步：
> /pnext

advanced to STEP_002

现在开始 STEP_002...
> /pwork
...
```

**用户无需任何干预** - Claude 自动循环直到所有步骤完成。

### 示例 2：验证失败时的自动修复

**Claude 的执行**:
```
> /pverify

STEP STEP_002: FAILED
  ✗ command: npm test -- rate-limit (exit=1)

测试失败了。让我检查错误并修复。
[查看错误，修复代码]

重新验证：
> /pverify

STEP STEP_002: OK
  ✓ command: npm test -- rate-limit

通过！推进：
> /pnext
```

**关键点**：Claude 不能跳过失败的验证，必须修复后才能推进。

## 常见命令

### Claude 会自动使用的命令

- `/pwork` - 显示当前步骤（低 token）
- `/pverify` - 验证当前步骤
- `/pstatus` - 查看所有步骤状态
- `/pnext` - 推进到下一步

### 或直接运行脚本

```bash
bash .claude-plugin/scripts/fsm.sh current
bash .claude-plugin/scripts/verify.sh --current
bash .claude-plugin/scripts/fsm.sh advance
bash .claude-plugin/scripts/fsm.sh status
```

## 防止 Claude 作弊的机制

### 1. **只读 plan.json 的状态字段**

Claude 不能手动编辑：
- `status` - 只能由验证器写入
- `current_step` - 只能由 advance 更新

如果 Claude 试图编辑这些字段，git hook 会警告。

### 2. **强制验证**

Stop hook 确保 Claude 在停止前必须通过验证。

### 3. **机器可验证的 proof**

验证规则都是机器执行的（tests、commands、file checks），不是语言描述。

### 4. **Audit trail**

所有状态变更都在 git 历史中可追踪：

```bash
git log -p plan.json
```

## 最佳实践

### 1. 小步快跑

每个步骤应该是 5-15 分钟的工作量：

```json
{
  "objective": "Add login endpoint",  // ✓ 好：具体、可测试
  "objective": "Implement authentication"  // ✗ 差：太宽泛
}
```

### 2. 明确的验证规则

```json
{
  "verify": [
    {"type": "command", "cmd": "npm test login"}  // ✓ 好：机器可验证
  ]
}

// ✗ 差：自我描述
{
  "verify": [
    {"type": "command", "cmd": "echo done"}
  ]
}
```

### 3. 使用 depends_on 强制顺序

```json
{
  "id": "STEP_003",
  "depends_on": ["STEP_001", "STEP_002"]  // 必须等前两步完成
}
```

### 4. 限制修改范围

```json
{
  "allow_paths": ["src/auth/**", "tests/auth/**"]  // Claude 只能修改这些文件
}
```

## 与传统方式的对比

| 特性 | 传统 Markdown Planning | plansm 自动化 |
|------|----------------------|--------------|
| 计划格式 | Markdown checklist | JSON 状态机 |
| 完成验证 | Claude 自我声明 | 机器验证（tests/commands） |
| 状态更新 | Claude 手动勾选 | 自动（verify 通过后） |
| 推进控制 | Claude 决定 | 只能在 VERIFIED 后推进 |
| 作弊可能 | 高（容易假装完成） | 低（必须通过真实测试） |
| Token 成本 | 高（反复读整个计划） | 低（只读当前步骤） |
| Claude 角色 | 既是执行者又是评判者 | 只是执行者（验证器是评判者） |

## 故障排除

### Claude 不自动执行命令

**解决**：明确告诉 Claude：
```
请使用 /pwork, /pverify, /pnext 命令来执行计划
```

### 验证失败但 Claude 想跳过

**解决**：
1. 设置 Stop hook（blocking）
2. 在 invariants 中明确写：
   ```json
   "Do not advance until verification passes"
   ```

### Claude 修改了 plan.json 的 status

**解决**：安装 git hook：
```bash
cp hooks/pre-commit .git/hooks/
```

## 总结

**用户的工作**：
1. 写一次 plan.json（定义步骤和验证规则）
2. 告诉 Claude "按 plan.json 执行"
3. （可选）Review 最终结果

**Claude 的工作**：
1. 自动查看当前步骤（/pwork）
2. 实现功能
3. 自动验证（/pverify）
4. 自动推进（/pnext）
5. 重复直到完成

**关键优势**：
- ✅ 无需人工干预
- ✅ 不能假装完成（必须通过测试）
- ✅ 低 token（只看当前步骤）
- ✅ 可审计（git 历史）
- ✅ 零安装（只需 jq）

这就是 plansm 的革命性之处：把 LLM 从"自我评判者"变成"被机器评判的执行者"。
