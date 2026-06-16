## Context

The people module currently requires a client-supplied ID (`p1`, `p2`, …) when creating a person. Rooms and tasks already use backend-generated sequential IDs via PostgreSQL sequences:

```
rooms: 'room-' || nextval('rooms_areas_id_seq')   → room-1, room-2, …
tasks: 'task-' || nextval('backlog_tasks_id_seq') → task-100, task-101, …
```

The people table schema (`id TEXT PRIMARY KEY`) naturally supports the same pattern. The existing seed data uses IDs `p1`..`p8`, so the sequence should be seeded from the max existing number and use `p` prefix for consistency with existing records.

## Goals / Non-Goals

**Goals:**

- `POST /api/people` accepts only `name` and `initials`; `id` is generated server-side
- New PostgreSQL sequence `people_id_seq` with a 2-phase migration: create sequence, then seed from max existing `p{N}`
- The `CreatePerson` sqlc query uses `'p' || nextval('people_id_seq')` matching the room/task pattern
- The frontend "Add a person" form removes the ID input field
- Existing demo data (`p1`..`p8`) continues to work without migration or ID changes
- The OpenAPI snapshot and generated frontend client are refreshed

**Non-Goals:**

- No changes to how existing `p1`..`p8` IDs are stored or referenced
- No changes to other person endpoints (`PUT`, `DELETE`, availability)
- No changes to the Person type response shape (still returns `id`)
- No changes to the `people-availability-integration` capability or the form layout beyond removing the ID field
- No retroactive re-generation of existing IDs

## Decisions

### Decision 1: PostgreSQL sequence with `p` prefix (matching room/task pattern)

| Option | Assessment |
|---|---|
| **Sequence + `p` prefix** (`'p' \|\| nextval(...)`) | ✅ Matches existing `p1`..`p8` seed data idiom. ✅ No prefix collision — users wouldn't normally type this. ✅ Same SQL pattern as rooms and tasks. |
| UUIDv4 | ❌ Not human-friendly for debugging/API use. ❌ Doesn't match any existing entity ID pattern. |
| Keep manual ID entry | ❌ The original problem — forces users to invent slugs. |
| Name-derived slug | ❌ Duplicate risks, special-character handling, no existing pattern. |

**Decision**: Use `'p' || nextval('people_id_seq')` — consistent with the established `prefix || nextval()` convention.

### Decision 2: Two-phase migration matching the room sequence precedent

Follow the exact pattern set by `007_add_rooms_areas_sequence.sql`:

- **Phase 1 (up)**: `CREATE SEQUENCE IF NOT EXISTS people_id_seq`, then `SELECT setval(...)` from `MAX(CAST(SUBSTRING(id FROM 2) AS INTEGER))` from the `people` table.
- **Phase 2 (down)**: `DROP SEQUENCE IF EXISTS people_id_seq`.

### Decision 3: Remove `id` from the API body (not make it optional)

| Option | Assessment |
|---|---|
| **Remove `id` from body** — `POST /api/people` accepts only `name` + `initials` | ✅ Clean contract, no ambiguity. ✅ Matches the room API (which also doesn't accept an ID). ❌ Breaking change (requires frontend and OAPI update). |
| Make `id` optional — fall back to auto-gen if absent | ❌ More complex validation logic. ❌ Two code paths to test. ❌ "Optional override" doesn't have a use case. |

**Decision**: Remove `id` from `CreatePersonInput.Body` entirely. The handler no longer checks for `id` or performs a duplicate-ID check (sequence guarantees uniqueness).

### Decision 4: Frontend removes only the ID field, keeps the rest of the form

The form keeps Name and Initials as-is. The `newId` ref, its `<label>` + `<Input>`, and the conditional `createError` for empty ID are removed. The `handleCreate` function drops the `id` field from the mutation body.

## Risks / Trade-offs

| Risk | Severity | Mitigation |
|---|---|---|
| **Concurrent creates** could race on `setval` seed if migration runs while writes are in-flight | Low | Migration runs in a goose transaction. The sequence is created first, then seeded atomically. After migration completes, `nextval` is guaranteed unique. |
| **Breaking API change** may cause frontend build failures if types aren't regenerated | Medium | The task list explicitly includes regenerating the OpenAPI snapshot and client types before the frontend code change. CI will catch mismatches. |
| **`p` prefix collision** if a user had manually typed an ID like `p999` that overlaps with the sequence | Low | The sequence is seeded from `MAX(...)`, so it starts after the highest existing number. No collision. |
| **Downgrade/rollback** would need to restore the `id` field in the API | Low | The migration has a `DROP SEQUENCE` down. The old handler code returns if the migration is rolled back and the old frontend is deployed. |
| **No duplicate-ID check** after removing the manual check | None (eliminated) | Sequence guarantees uniqueness — the primary key constraint enforces it. A `409` from a sequence collision is effectively impossible. |
