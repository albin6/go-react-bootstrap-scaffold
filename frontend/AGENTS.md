# Frontend AI Agent Guide

## 1. Absolute Directives

* **Frontend must never import backend source code.**
* **Frontend must never access PostgreSQL.**
* **Frontend must never depend on pgx or SQLC.**
* **Frontend communicates with backend exclusively through HTTP APIs.**
* **Frontend environment variables are public and must not contain secrets.**

---

## 2. Architecture & Design Rules

* **Strict Separation**: The frontend knows nothing about Go packages, struct definitions, SQL files, or database schemas. All interaction occurs over the HTTP REST contract (`/api/v1/...`).
* **Feature-Oriented Slices**:
  Organize domain logic in `src/features/<feature-name>/`:
  * `components/`: UI components rendered for this feature.
  * `hooks/`: Custom React hooks encapsulating feature state and queries.
  * `services/`: API calls invoking the centralized HTTP client.
  * `types/`: Strongly typed TypeScript interfaces.
* **Component Primitives**: Generic UI widgets (buttons, cards, badges, modal dialogs) belong in `src/components/ui/`.
* **Layouts**: Page scaffolding belongs in `src/components/layout/`.
* **State Management**: Prefer local component state (`useState`) and custom feature hooks for data fetching. For application-wide state, use dedicated stores in `src/stores/`.
* **API Client**: Never call `fetch()` or `axios` directly inside React components. All requests flow through `src/lib/api-client.ts` or feature services.
* **Strict TypeScript**: Never use `any`. Always specify types for props, state, and API envelopes.

---

## 3. Testing Requirements

* **Unit & Component Tests**: Write tests with Vitest and `@testing-library/react` under `tests/` or alongside components in `__tests__/`.
* **User-Centric Testing**: Test user interactions and visual states (loading, error, empty, success) rather than internal state details.
* **Mocking**: Mock API services at the service boundary rather than mocking Axios internals.

---

## 4. Verification Checklist

Before declaring a frontend task complete:
* [ ] `npm run lint` passes with 0 errors and 0 warnings.
* [ ] `npm run typecheck` passes cleanly.
* [ ] `npm run test` passes all tests.
* [ ] `npm run build` generates `dist/` without errors.
