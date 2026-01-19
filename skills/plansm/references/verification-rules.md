# Verification Rules Reference

plansm supports four types of machine-verifiable proof rules: command, file_exists, file_contains, and http. Each rule must pass for a step to be marked VERIFIED.

## Rule Types

### 1. command - Execute Shell Command

Runs a bash command and checks the exit code and/or output.

**Required fields:**
- `type`: `"command"`
- `cmd`: Shell command to execute

**Optional fields:**
- `expect.exit_code`: Expected exit code (default: 0)
- `expect.stdout_contains`: Regex pattern that stdout must match

**Examples:**

```json
{
  "type": "command",
  "cmd": "npm test",
  "expect": {
    "exit_code": 0
  }
}
```

```json
{
  "type": "command",
  "cmd": "git status",
  "expect": {
    "exit_code": 0,
    "stdout_contains": "nothing to commit"
  }
}
```

### 2. file_exists - Check File Existence

Verifies that a file or directory exists at the specified path.

**Required fields:**
- `type`: `"file_exists"`
- `file`: Path to file/directory (relative to project root)

**Example:**

```json
{
  "type": "file_exists",
  "file": "src/components/DarkModeToggle.tsx"
}
```

### 3. file_contains - Check File Content

Verifies that a file contains content matching a regex pattern.

**Required fields:**
- `type`: `"file_contains"`
- `file`: Path to file (relative to project root)
- `pattern`: Regex pattern to search for

**Example:**

```json
{
  "type": "file_contains",
  "file": "src/index.ts",
  "pattern": "export.*DarkModeToggle"
}
```

**Pattern notes:**
- Uses grep extended regex syntax
- Case-sensitive by default
- `.` matches any character except newline
- `.*` matches any characters
- Use `\\` to escape special regex chars

### 4. http - Check HTTP Response

Makes an HTTP GET request and verifies the response.

**Required fields:**
- `type`: `"http"`
- `url`: Full URL to request

**Optional fields:**
- `expect_status`: Expected HTTP status code (default: 200)
- `expect_body`: Regex pattern that response body must match

**Examples:**

```json
{
  "type": "http",
  "url": "http://localhost:3000/api/health",
  "expect_status": 200,
  "expect_body": "ok"
}
```

```json
{
  "type": "http",
  "url": "http://localhost:3000/api/users",
  "expect_status": 200,
  "expect_body": "\\[.*\\]"
}
```

## Combining Multiple Rules

A step can have multiple verification rules. ALL must pass:

```json
{
  "id": "STEP_003",
  "objective": "Implement and test login feature",
  "status": "PENDING",
  "verify": [
    {
      "type": "file_exists",
      "file": "src/auth/login.ts"
    },
    {
      "type": "file_contains",
      "file": "src/auth/login.ts",
      "pattern": "export.*login"
    },
    {
      "type": "command",
      "cmd": "npm test -- login",
      "expect": {
        "exit_code": 0
      }
    },
    {
      "type": "http",
      "url": "http://localhost:3000/api/login",
      "expect_status": 401
    }
  ]
}
```

## Best Practices

1. **Start Simple**: Begin with `file_exists`, then add `file_contains`, finally `command`
2. **Test First**: Write verification rules before implementation (like TDD)
3. **Be Specific**: Use precise patterns to avoid false positives
4. **Layer Verification**:
   - File exists (basic)
   - File contains expected exports (structure)
   - Tests pass (functionality)
5. **Avoid Brittleness**: Don't check exact line numbers or whitespace
6. **Use Dependencies**: Order steps so simpler verifications come first

## Anti-Patterns to Avoid

❌ **Too Loose**: `{"type": "file_exists", "file": "."}`
✅ **Specific**: `{"type": "file_exists", "file": "src/auth/login.ts"}`

❌ **No Verification**: Relying on LLM self-reporting
✅ **Always Verify**: Every step needs at least one verification rule

❌ **Manual Status**: Editing plan.json status fields by hand
✅ **Machine Only**: Let verify.sh update status based on test results

## See Also

- [SKILL.md](../SKILL.md) - Main skill documentation
- [templates/plan-template.json](../templates/plan-template.json) - Example plan structure
