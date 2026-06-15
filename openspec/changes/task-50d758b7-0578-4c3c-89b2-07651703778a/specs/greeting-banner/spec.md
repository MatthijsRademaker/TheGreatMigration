# greeting-banner Specification

## REMOVED Requirements

### Requirement: REQ-GREET-01 Banner Visibility
**Reason**: Pipeline validation is complete — the banner was test scaffolding, not a permanent feature.
**Migration**: Remove the greeting banner `<div>` from `src/home/HomeView.vue`. The home page section will no longer have a banner as its first child.

### Requirement: REQ-GREET-02 Consistent Styling
**Reason**: Pipeline validation is complete — the banner styling classes (`bg-muted/40 rounded-lg border p-4 border-l-4 border-primary text-sm`) are no longer needed.
**Migration**: Remove the banner `<div>` and its associated Tailwind utility classes from `src/home/HomeView.vue`. No other components use this specific combination.

### Requirement: REQ-GREET-03 No Side Effects
**Reason**: The banner itself is being removed, so the constraint on not introducing side effects is moot.
**Migration**: No migration needed — the requirement was about what the original change must NOT do. With the banner removed, the constraint is satisfied by absence.

### Requirement: REQ-GREET-04 Build Verification
**Reason**: Pipeline validation is complete. The banner was the only feature this spec governed.
**Migration**: No migration needed — the build verification requirement applied to the original test task's acceptance. The removal itself must still pass `scripts/precommit-run` and `npm run build`.
