<div align="center">
  <img src="./apps/frontend/public/upsync-logo.svg" height="120" alt="UpSync Logo" />
  <h1>UpSync</h1>
  <p><strong>Fast, Secure, and Temporary File Sharing.</strong></p>
  <p>
    <a href="#features">Features</a> •
    <a href="#architecture">Architecture</a> •
    <a href="#setup">Local Setup</a> •
    <a href="#api-reference">API Overview</a>
  </p>
</div>

---

**UpSync** is a full-stack, production-ready temporary file sharing application. Upload any file up to 50 MB to receive an instant, secure, and short-lived download link. All files and their metadata automatically self-destruct after your chosen expiry time.

## 🚀 Features

- **Blazing Fast Uploads:** Drag and drop support with real-time UI progress bars.
- **Custom Expirations:** Select expiry intervals from 1 hour up to 7 days.
- **Auto-Deletion:** Go cron-job continuously sweeps and destroys expired database records and storage objects.
- **High Security:** 60-second signed download links prevent hotlinking. The backend relies solely on the Supabase Service Role (bypassing public API access).
- **Strong Typing & Error Recovery:** Complete TypeScript frontend and Go backend, integrated with robust React Error Boundaries and application-wide Toast notifications.

---

## 🏗 Architecture

The project is structured as a native monorepo containing an efficient RESTful **Go backend** and a modern **React + Vite frontend**.

### Tech Stack

| Backend              | Frontend                 | Cloud               |
| -------------------- | ------------------------ | ------------------- |
| Golang 1.26          | React 18                 | Supabase PostgreSQL |
| Gin HTTP framework   | TypeScript               | Supabase Storage    |
| robfig/cron/v3       | Vite                     |                     |
| Custom Service Layer | Context API + CSS Tokens |                     |

### Folder Structure

```text
/UpSync
├── SUPABASE_SETUP.md          # Guide: Creates Postgres tables and buckets
└── apps/
    ├── backend/               # 🐹 Golang API Server
    │   ├── cmd/server/        # Entry point: main.go
    │   ├── internal/
    │   │   ├── apierr/        # Structured HTTP errors handling
    │   │   ├── config/        # Environment variable parsing
    │   │   ├── database/      # Supabase REST client (No 3rd-party SDKs)
    │   │   ├── handlers/      # Gin HTTP controllers
    │   │   ├── middleware/    # CORS, Logger, MaxBodySize, Panic Recovery
    │   │   ├── models/        # Go structs for DB and JSON mapping
    │   │   ├── scheduler/     # Background cleanup cron system
    │   │   └── services/      # Business logic & Database coordination
    │   ├── go.mod
    │   └── .env.example
    │
    └── frontend/              # ⚛️ React Vite App
        ├── src/
        │   ├── api/           # Typed Axios client & interceptors
        │   ├── components/    # Reusable UI fragments (DropZone, Toast, etc.)
        │   ├── context/       # Global State (ToastProvider)
        │   ├── pages/         # Page/Route components (Upload, Share)
        │   ├── types/         # Typescript interface definitions
        │   ├── App.tsx        # React Router topology
        │   └── index.css      # Core CSS tokens & styles
        ├── vite.config.ts     # Proxy & dev configuration
        └── .env.example
```

---

## 🛠 Local Setup

Follow these exact instructions to get the backend, frontend, and database running locally.

### 1. Configure Supabase Server

Follow the instructions provided in [`SUPABASE_SETUP.md`](./SUPABASE_SETUP.md) located in the project root to set up:

- The `files` Postgres table with correct column structures.
- The `upsync-files` Storage bucket.

### 2. Run the Backend (Golang)

Open your first terminal and run:

```bash
cd apps/backend
cp .env.example .env

# Edit .env with your Supabase URL and SERVICE_ROLE_KEY
# Ensure PORT=8080 and FRONTEND_URL=http://localhost:5173

go mod tidy
go run ./cmd/server/main.go
```

The Go API will be running on `http://localhost:8080`.

### 3. Run the Frontend (React / Vite)

Open a **second** terminal and run:

```bash
cd apps/frontend
cp .env.example .env

# Edit .env so VITE_API_URL points to the local backend port
# Example: VITE_API_URL=http://localhost:8080

npm install
npm run dev
```

The User Interface will be running on `http://localhost:5173`.

---

## 🛣 API Reference

All backend functionality is exposed purely as a structured JSON API on `/api/files`.

| Method | Endpoint                  | Description                                                             |
| ------ | ------------------------- | ----------------------------------------------------------------------- |
| `POST` | `/api/files/upload`       | Accepts a Multipart-Form with max `50MB`. Generates a UUID and uploads. |
| `GET`  | `/api/files/:id`          | Returns parsed file metadata (used by Share Page).                      |
| `GET`  | `/api/files/:id/download` | Returns a strict 60-second temporary `downloadUrl`.                     |

---

## 🔒 Security Posture

1. **No Public API Access:** The frontend never connects to Supabase directly. Supabase Row Level Security (RLS) is bypassed securely via the server-only `SERVICE_ROLE_KEY`.
2. **Short-Lived Signatures:** Download URLs generated by Go are temporary and expire precisely within 60 seconds, completely preventing hotlinking.
3. **Hard Memory Caps:** A `MaxBodySize` middleware layer enforces explicit rejection on inputs exceeding 50 MB limits directly preventing buffer overloads.
4. **Automated Cleanups:** The Go Cron Scheduler (`robfig/cron/v3`) silently manages and purges orphaned Postgres rows and blob storage objects every 15 minutes.

<br />
<div align="center">
  <p>Built with ❤️ utilizing <strong>Go</strong> and <strong>React</strong>.</p>
</div>
