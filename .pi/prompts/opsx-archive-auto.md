---
description: Automatically archive a completed OpenSpec change. Non-interactive — suitable for swarm worker automation.
swarm: true
agent_types:
  - swarm-worker
---

Automatically archive a completed OpenSpec change.

$ARGUMENTS

**Input**: `$ARGUMENTS` is the exact OpenSpec change name (typically the task ID). This prompt is non-interactive: it never asks questions, never waits for confirmation, and always syncs delta specs before archiving.

**Steps**

1. **Extract change name**

   Use the first token of `$ARGUMENTS` as the exact change name. Do NOT prompt, infer, or ask for selection. If `$ARGUMENTS` is empty, fail immediately.

2. **Check if change directory exists**

   Verify `openspec/changes/<name>/` exists on disk.

   - If the directory does not exist: the change has likely already been archived. Treat this as successful idempotent completion — emit `{"outcome":"finished"}` and stop.
   - If the directory exists: continue.

3. **Check if archive target already exists**

   Compute the archive target: `openspec/changes/archive/YYYY-MM-DD-<name>/` using the current date.

   - If the archive target already exists: treat as idempotent success. The change was already archived under today's date. Emit `{"outcome":"finished"}` and stop.
   - If the archive target does not exist: continue.

4. **Archive using openspec CLI**

   Run the archive command:

   ```bash
   openspec archive --yes <name>
   ```

   This command:
   - Always syncs delta specs into canonical specs (default behavior).
   - Skips interactive confirmation prompts (`--yes`).
   - Moves the change directory to `openspec/changes/archive/YYYY-MM-DD-<name>/`.

   **If the command succeeds:** emit `{"outcome":"finished"}` and stop.

   **If the command fails:**
   - If the error indicates the change directory is missing (already archived): emit `{"outcome":"finished"}`.
   - If the error indicates the archive target already exists: emit `{"outcome":"finished"}`.
   - For any other error (filesystem error, validation failure, permission error): this is an unexpected failure. Report the error verbatim and do NOT emit a finished outcome. The workflow engine will treat this as a failed state.

**Output**

On successful or idempotent completion, emit exactly:

```json
{"outcome":"finished"}
```

On unexpected failure, report the error and stop without emitting the finished outcome.

**Guardrails**

- Never ask the user any questions. Never use interactive tools (AskUserQuestion, confirmation prompts, etc.).
- Never skip sync. Always use the default archive behavior which syncs delta specs.
- Treat already-archived and missing-source-directory as idempotent success.
- Do not check artifact or task completion status — proceed regardless.
- Do not offer sync choices or any interactive path.
- If the change directory does not exist when we first check it, the change is already archived — treat this as success.
