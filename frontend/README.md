# Frontend Application (React + TypeScript + Vite)

## 1. Overview & Technology Stack

The frontend is an isolated Single-Page Application (SPA) designed to communicate with the Go HTTP backend strictly through standard REST APIs.

* **Node Version**: `20.x+` (verified with `v24.x`).
* **Package Manager**: `npm` (v10+).
* **Framework**: React 18+ with TypeScript (strict mode enabled).
* **Build Tooling**: Vite with Hot Module Replacement (HMR).
* **Routing**: React Router v6+.
* **HTTP Client**: Centralized Axios client with typed response envelopes.
* **Testing**: Vitest, React Testing Library, and jsdom.
* **Linting & Formatting**: ESLint (Flat Config) and Prettier.

---

## 2. Directory Structure

```text
frontend/
├── public/                      # Static public files (favicons, manifest)
├── src/
│   ├── app/                     # Shell application and root styling
│   ├── assets/                  # Bundled images, icons, and SVG assets
│   ├── components/              # Shared generic UI components
│   │   ├── common/              # Cross-cutting error boundaries, fallbacks
│   │   ├── layout/              # Application layout (Header, Footer, Layout)
│   │   └── ui/                  # Atomic design system (Button, Card, Badge, Spinner)
│   ├── config/                  # Frontend environment configuration & validation
│   ├── features/                # Vertical domain slices
│   │   └── health/              # System health & backend connectivity feature
│   │       ├── components/      # Feature-specific UI components
│   │       ├── hooks/           # Custom feature data-fetching hooks
│   │       ├── services/        # Feature HTTP API calls
│   │       └── types/           # Domain TypeScript definitions
│   ├── hooks/                   # Generic reusable React hooks (useAsync, etc.)
│   ├── lib/                     # Centralized Axios API client configuration
│   ├── pages/                   # Top-level route pages (HomePage, NotFoundPage)
│   ├── routes/                  # Route definitions
│   ├── services/                # Shared networking layer
│   ├── stores/                  # Global state management
│   ├── types/                   # Shared TypeScript definitions (API envelopes)
│   ├── utils/                   # Pure helper functions (formatting, date tools)
│   └── main.tsx                 # React DOM mount point
├── tests/                       # Unit and component test suites
├── nginx.conf                   # Nginx reverse proxy and SPA fallback routing
├── Dockerfile                   # Multi-stage production container build
├── package.json                 # Dependency manifest and scripts
├── tsconfig.json                # TypeScript compiler configuration
├── vite.config.ts               # Bundler and Vitest configuration
└── README.md                    # This file
```

---

## 3. Local Development

### 3.1 Setup & Installation
```bash
cd frontend
cp .env.example .env
npm install
```

### 3.2 Running the Development Server
```bash
npm run dev
```
The application will be accessible at `http://localhost:5173`.

---

## 4. Environment Variables

Environment variables must be prefixed with `VITE_` to be exposed to the browser bundle:

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `VITE_API_BASE_URL` | `http://localhost:8080/api/v1` | Root URL for the backend Go REST API |
| `VITE_APP_TITLE` | `Production Scaffold` | Application browser title |

> **SECURITY WARNING**: Vite environment variables are embedded into public client-side JavaScript. Never store database credentials, private keys, or server secrets in `frontend/.env`.

---

## 5. Testing, Linting & Typechecking

```bash
# Run Vitest unit & component test suite
npm run test

# Run tests in watch mode
npm run test:watch

# Run TypeScript compiler typecheck (no emit)
npm run typecheck

# Run ESLint validation
npm run lint

# Format codebase with Prettier
npm run format
```

---

## 6. Production Build & Docker

```bash
# Generate optimized production assets in dist/
npm run build

# Preview production build locally
npm run preview
```

### Building Docker Image
```bash
docker build -t frontend:latest .
docker run -p 3000:80 frontend:latest
```
Visit `http://localhost:3000` to view the production Nginx container.
