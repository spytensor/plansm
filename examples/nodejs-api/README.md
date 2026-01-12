# Node.js API Example

This example demonstrates using plansm to build a Node.js REST API with test-driven development.

## What This Example Shows

- Test-driven step verification
- File existence checks
- Pattern matching in code
- HTTP endpoint verification
- Step dependencies
- Path restrictions

## Setup

```bash
# Initialize the project with plansm
plansm init --claude

# The plan.json defines 5 steps to build the API
cat plan.json
```

## The Plan

### STEP_001: Project Setup
- Creates package.json
- Installs Express
- **Proof**: package.json exists and contains express

### STEP_002: Basic Server
- Creates src/server.js with Express app
- **Proof**:
  - File exists
  - Contains express import
  - Server starts on port 3000

### STEP_003: Health Endpoint
- Adds /health endpoint
- **Proof**:
  - GET /health returns 200
  - Response contains {"status": "ok"}

### STEP_004: Users API
- Implements /api/users (GET, POST)
- **Proof**:
  - Unit tests pass (npm test)
  - Pattern matching for CRUD functions

### STEP_005: Integration Tests
- End-to-end API tests
- **Proof**:
  - Integration test suite passes
  - All endpoints return correct responses

## Running the Example

```bash
# Step 1: Setup
npm init -y
npm install express jest supertest

plansm verify --current  # Should pass
plansm advance

# Step 2: Basic Server
# Create src/server.js with Express
cat > src/server.js << 'EOF'
const express = require('express');
const app = express();
const PORT = process.env.PORT || 3000;

app.use(express.json());

if (require.main === module) {
  app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
  });
}

module.exports = app;
EOF

plansm verify --current  # Should pass
plansm advance

# Step 3: Health Endpoint
# Add health endpoint to server.js
# ... (continue implementing each step)

# After each step:
plansm verify --current
plansm advance
```

## With Claude Code

1. Open this directory in Claude Code
2. Run `/pstatus` to see the plan
3. Run `/pwork` to see current step
4. Implement the current step
5. Run `/pverify` to check
6. Run `/pnext` to advance

Claude will:
- Focus on one step at a time
- Can't mark complete without passing tests
- Automatically advances when verified

## Key Lessons

1. **Proof-Based**: Each step has machine-checkable verification
2. **Dependencies**: Later steps locked until earlier ones verified
3. **Focused Scope**: `allow_paths` restricts where Claude can edit
4. **Test-Driven**: Tests define "done"
5. **Low Token**: Claude only sees current step, not full plan

## Verification Examples

### Command (Tests)
```json
{
  "type": "command",
  "cmd": "npm test",
  "expect": {"exit_code": 0}
}
```

### File Exists
```json
{
  "type": "file_exists",
  "file": "src/server.js"
}
```

### File Contains Pattern
```json
{
  "type": "file_contains",
  "file": "src/server.js",
  "pattern": "const express = require"
}
```

### HTTP Check
```json
{
  "type": "http",
  "url": "http://localhost:3000/health",
  "expect": {"http_status": 200}
}
```

## What Happens If You Cheat?

Try manually editing plan.json to mark STEP_002 as VERIFIED without implementing it:

```bash
# Edit plan.json: "status": "VERIFIED"
git add plan.json
git commit -m "fake completion"
```

**Result**: Pre-commit hook warns you!

```
⚠️  WARNING: plan.json status or current_step has changed

Status and current_step should only be modified by:
  - plansm verify (sets status to VERIFIED/FAILED)
  - plansm advance (updates current_step)
```

## CI Integration

In CI, run:

```bash
plansm verify --all
```

This ensures:
- All VERIFIED steps can still be verified
- No steps were fake-completed
- Plan and code are synchronized

## Next Steps

- Try [Python FastAPI example](../python-fastapi/)
- Read [Concept documentation](../../docs/concept.md)
- See [Claude Code integration guide](../../docs/claude-code.md)
