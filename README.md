# AgentSpec

**The contract layer for AI-powered software development.**

AgentSpec lets your team define feature requirements as a structured YAML spec — a single source of truth that every AI agent (PM, frontend, backend, QA) can consume.

## The Problem

When multiple AI agents work on the same feature:

- **PRD drift** — AI-generated code diverges from the PRD over time
- **Cross-system blind spots** — Agent A changes model X without knowing it breaks Agent B's page Y
- **Context mismatch** — Each agent sees different, often stale, context
- **Late-stage discovery** — QA finds spec-code misalignment that should have been caught at design time

## The Solution

```
┌─────────────┐
│  spec.yaml  │  ← Single source of truth
└──────┬──────┘
       │
  ┌────┴────┐
  │agentspec│  ← CLI tool
  └────┬────┘
       │
  ┌────┼────────┼────────┼────┐
  ▼    ▼        ▼        ▼    ▼
 init validate generate context check/diff
```

One YAML file defines **models**, **behaviors**, **APIs**, **UI pages**, and **acceptance criteria**. Each role gets a filtered, Markdown view of the same spec.

## Quick Start

```bash
# Install
go install github.com/chenzhaohao/agentspec/cmd/agentspec@latest

# Initialize a new spec
agentspec init

# Validate your spec
agentspec validate

# Generate code scaffolding
agentspec generate --target backend
agentspec generate --target frontend
agentspec generate --target test

# Generate role-specific agent context
agentspec context --role backend
agentspec context --role frontend
agentspec context --role qa
agentspec context --role pm

# Check code conformance
agentspec check --source ./src

# Analyze spec changes
agentspec diff --old spec-v1.yaml --new spec-v2.yaml
agentspec diff --old git:main:spec.yaml --new spec.yaml
```

## Spec Format

```yaml
kind: FeatureSpec
version: "1.0"
name: task-management
display_name: Task Management
owner: product-team
description: A collaborative task management system

models:
  Task:
    description: A unit of work
    fields:
      id:
        type: uuid
        required: true
      title:
        type: string
        required: true
        constraints:
          min_length: 1
          max_length: 300
      status:
        type: "enum:task-status"
        required: true
    states:
      field: status
      values: [todo, in_progress, in_review, done]
      transitions:
        todo: [in_progress]
        in_progress: [in_review, todo]
        in_review: [done, in_progress]
        done: []

enums:
  task-status:
    values:
      - name: todo
        value: todo
      - name: in_progress
        value: in_progress
      - name: in_review
        value: in_review
      - name: done
        value: done

behaviors:
  create-task:
    description: Create a new task
    actor: User
    target_model: Task
    preconditions:
      - description: User must be a project member
    effects:
      - description: Task created with status 'todo'
    errors:
      - code: PERMISSION_DENIED
        message: Not a project member
        http_status: 403
    cross_impact:
      - target: Project
        description: Project's updated_at refreshes

apis:
  create-task:
    method: POST
    path: /api/projects/{project_id}/tasks
    behavior: create-task
    auth: true

ui:
  pages:
    task-board:
      type: detail
      path: /projects/{project_id}/tasks
      data_source: list-tasks

acceptance:
  - id: AC-001
    title: User can create a task
    behaviors: [create-task]
    apis: [create-task]
    priority: high
```

## Commands

| Command | Description |
|---------|-------------|
| `init` | Create a new spec file interactively |
| `validate` | Check spec for structural, semantic, and cross-reference errors |
| `generate` | Generate backend (Go), frontend (TypeScript/React), or test scaffolding |
| `context` | Generate role-filtered Markdown for AI agent system prompts |
| `check` | Scan source code to verify conformance with spec |
| `diff` | Compare two spec versions with downstream impact analysis |

## Role-Based Context

The `context` command generates tailored Markdown for each role:

| Role | Models | Behaviors | APIs | UI | Acceptance |
|------|--------|-----------|------|----|------------|
| **backend** | Full (with constraints) | Full (with pre/post) | Full | — | — |
| **frontend** | Types only | Summary only | Full | Full | — |
| **qa** | Full | Full | Full | Full | Full |
| **pm** | Types only | Full | Summary | Full | Full |

## Validation

Three-pass validation engine:

1. **Structural** — Required fields, valid types
2. **Semantic** — ref/enum resolution, state machine integrity, behavior-API links
3. **Cross-reference** — Orphaned models/behaviors, UI→API consistency, acceptance criteria coverage

## Spec Type System

| Spec Type | Go Type | TypeScript Type |
|-----------|---------|-----------------|
| `string` | `string` | `string` |
| `int` | `int64` | `number` |
| `float` | `float64` | `number` |
| `bool` | `bool` | `boolean` |
| `uuid` | `string` | `string` |
| `datetime` | `time.Time` | `string` |
| `ref:Model` | `*Model` | `Model` |
| `enum:Name` | `Name` | `Name` |
| `[]type` | `[]type` | `type[]` |

## Example

See [`examples/task-management/spec.yaml`](examples/task-management/spec.yaml) for a complete example with 4 models, 6 behaviors, 10 APIs, 5 UI pages, and 6 acceptance criteria.

## License

[Apache License 2.0](LICENSE)
