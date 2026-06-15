## ADDED Requirements

### Requirement: The Vue frontend SHALL be relocated into a `frontend/` service directory
The entire Vue/Vite application (source, tests, and configuration) SHALL be moved from the repo root into a new `frontend/` directory. The relocation SHALL use `git mv` for all paths to preserve commit history. After relocation, all existing dashboard routes SHALL render identically to their pre-move behavior.

#### Scenario: Vite build alias resolves correctly after relocation
- **WHEN** the frontend is built from `frontend/`
- **THEN** `vite.config.ts` resolves `@` to `./src` relative to `frontend/vite.config.ts` without manual path changes

#### Scenario: TypeScript path mapping resolves after relocation
- **WHEN** `vue-tsc --noEmit` runs from `frontend/`
- **THEN** `tsconfig.json` resolves `@/*` to `./src/*` and includes `src/**/*.ts`, `src/**/*.vue`, and `tests/**/*.ts` relative to `frontend/tsconfig.json`

#### Scenario: shadcn-vue theme CSS resolves after relocation
- **WHEN** shadcn-vue components are rendered
- **THEN** `components.json` points `css` to `src/app/styles.css` relative to `frontend/components.json`, and Tailwind v4 theming continues to work

#### Scenario: Vite entry point resolves after relocation
- **WHEN** the Vite dev server starts from `frontend/`
- **THEN** `index.html` serves as the entry point and its `<script type="module" src="/src/main.ts">` resolves correctly within the Vite project root at `frontend/`

#### Scenario: Test imports resolve after relocation
- **WHEN** vitest runs from `frontend/`
- **THEN** all test files in `frontend/tests/` import from `../src/...` without path errors

#### Scenario: Design-system documentation reflects new path
- **WHEN** a future executor reads `docs/design-system-v2.md`
- **THEN** the global theme surface reference states `frontend/src/app/styles.css` instead of `src/app/styles.css`

#### Scenario: Root package.json no longer exists after relocation
- **WHEN** the relocation is complete
- **THEN** no `package.json` exists at the repo root; all Node dependencies and scripts live in `frontend/package.json`