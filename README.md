# 💍 Digital Wedding Invitation

A personalized digital wedding invitation built as a fully static Vue frontend backed by Go serverless functions on Vercel, with Google Sheets as the data store. Every invited guest receives a unique URL; all visitors share one public guestbook.

> Built to the specification in `DIGITAL WEDDING INVITATION WEBSITE.txt` (PRD).

---

## Table of Contents

- [Features](#features)
- [Screenshots](#screenshots)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Google Sheets Setup](#google-sheets-setup)
- [Environment Variables](#environment-variables)
- [Local Development](#local-development)
- [Deployment](#deployment)
- [Templates & Design Customization](#templates--design-customization)
- [API Reference](#api-reference)
- [Guest Management](#guest-management)
- [Troubleshooting](#troubleshooting)
- [Roadmap](#roadmap)
- [License](#license)

---

## Features

- **Personalized invitations** — `/andi-pratama` greets the guest by name; `/` greets them as *Anonymous*.
- **Group invitations** — `/keluarga-besar-thanos` displays a shared greeting for a family.
- **Shared public guestbook** — every visitor reads and writes to one guestbook, visible from every URL.
- **Bride-friendly spreadsheet workflow** — guests are managed by adding names to column A of a private Google Sheet. Technical identifiers are generated and maintained by the backend.
- **Three independent templates** — switch the active design by changing one line in `App.vue`.
- **Centralized configuration** — wedding content, images, and optional sections live in one file.
- **Backend-enforced identity** — sender names are resolved from URL slugs, never trusted from the client.
- **Rate limiting, validation, pagination** — the guestbook is protected against spam and abuse.
- **Mobile-first responsive design** — every template tested across desktop, tablet, and phone.

---

## Screenshots

> Add your own images to `docs/screenshots/` and reference them here once you have a real wedding design.

| Template A | Template B | Template C |
|---|---|---|
| _(minimalist editorial)_ | _(floral romantic)_ | _(traditional card)_ |

---

## Architecture

```
┌─────────────────────┐
│  Browser (guest)    │
│  GET /andi-pratama  │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐       ┌──────────────────┐
│  Vue 3 + Vite       │──────▶│  Go (Vercel fns) │
│  frontend/dist      │ /api  │  backend/api/*.go│
└─────────────────────┘       └────────┬─────────┘
                                       │ Google Sheets API
                                       ▼
                              ┌──────────────────┐
                              │  Google Sheets   │
                              │  Guests | Messages│
                              └──────────────────┘
```

**Key boundaries:**

- The frontend never touches Google Sheets directly.
- The backend resolves the guest's identity from the URL slug — the client cannot impersonate another guest.
- Templates are visual shells that receive data via a fixed prop contract; all data logic lives in `InvitationController.vue`.
- Shared CSS tokens live in `styles/variables.css`; template-specific layout lives inside each template folder.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend framework | Vue 3 (`<script setup>`) |
| Build tool | Vite 5 |
| Backend | Go 1.22 |
| Serverless runtime | Vercel Go functions |
| Data store | Google Sheets (via `google.golang.org/api/sheets/v4`) |
| Auth | Google service account (JWT) |
| Deployment | Vercel |

---

## Project Structure

```
wedding-invitation/
├── README.md
├── vercel.json                   # Multi-service config + rewrites
├── .env.example
├── .gitignore
│
├── backend/
│   ├── go.mod                    # module backend
│   ├── go.sum
│   ├── cmd/
│   │   └── dev/
│   │       └── main.go           # Local dev API server (not deployed)
│   ├── api/                      # Vercel-facing function folders
│   │   ├── health/index.go       # GET    /api/health
│   │   ├── messages/index.go     # GET    /api/messages  (list)
│   │   │                         # POST   /api/messages  (anonymous submit)
│   │   ├── guest/index.go        # GET    /api/guests/:slug
│   │   └── guest_messages/index.go
│   │                             # POST   /api/guests/:slug/messages
│   └── internal/
│       ├── config/config.go      # Env var loading
│       ├── httputil/httputil.go  # JSON + CORS + slug extraction
│       ├── repo/repo.go          # Wires config + sheets client
│       └── sheets/
│           ├── client.go         # Google Sheets API client
│           ├── guests.go         # Guest sync, slug generation, cache
│           └── messages.go       # Message read/append, IP hashing
│
└── frontend/
    ├── package.json
    ├── vite.config.js            # /api proxy → :3000
    ├── index.html
    ├── public/images/            # Placeholder SVGs; replace with photos
    └── src/
        ├── main.js
        ├── App.vue               # ← active template selector
        ├── app/
        │   └── InvitationController.vue
        ├── templates/
        │   ├── TemplateA/{index.vue, style.css}
        │   ├── TemplateB/{index.vue, style.css}
        │   └── TemplateC/{index.vue, style.css}
        ├── components/
        │   ├── Guestbook.vue
        │   ├── Countdown.vue
        │   ├── Gallery.vue
        │   ├── LoadingState.vue
        │   └── ErrorState.vue
        ├── composables/
        │   └── useInvitation.js
        ├── services/
        │   └── api.js
        ├── config/
        │   └── wedding.config.js  # ← all wedding content
        └── styles/
            ├── variables.css      # Design tokens
            ├── global.css         # Reset, base typography
            └── responsive.css     # Shared responsive rules
```

---

## Prerequisites

- **Node.js** ≥ 18
- **Go** ≥ 1.22
- **Git**
- **Vercel CLI** (`npm i -g vercel`) — only for deployment
- A **Google account** with access to Google Cloud Console and Google Sheets

---

## Google Sheets Setup

The backend expects one private spreadsheet with two tabs.

### 1. Create the spreadsheet

Go to <https://sheets.google.com> → **Blank spreadsheet** → rename to `Wedding Invitation`.

Copy the **Spreadsheet ID** from the URL:

```
https://docs.google.com/spreadsheets/d/1aBcD.../edit
                                       └─────┘
                                       this is GOOGLE_SHEET_ID
```

### 2. Create the two tabs

Rename `Sheet1` → **`Guests`**. Add a second tab → **`Messages`**. Names must match exactly (case-sensitive).

### 3. Set up the `Guests` tab

Headers in row 1, in this order:

| A1 | B1 | C1 | D1 |
|---|---|---|---|
| `Guest Name` | `guest_id` | `slug` | `created_at` |

Fill column A with guest names only. Leave B–D blank — the backend populates them on first read.

| Guest Name |
|---|
| Andi Pratama |
| Budi Santoso |
| Keluarga Besar Thanos |

> **Tip:** hide columns B–D once populated. The bride only ever needs to see `Guest Name`.

### 4. Set up the `Messages` tab

Headers in row 1, in this order:

| A1 | B1 | C1 | D1 | E1 | F1 |
|---|---|---|---|---|---|
| `message_id` | `guest_id` | `display_name` | `message` | `created_at` | `ip_hash` |

No data rows needed. Messages are appended automatically.

### 5. Create a Google Cloud service account

1. <https://console.cloud.google.com> → **New Project** → name it `wedding-invitation`.
2. Enable the **Google Sheets API**.
3. **APIs & Services → Credentials → Create credentials → Service account**.
4. Name it `wedding-backend`, skip optional steps, click **Done**.
5. Open the account → **Keys → Add key → Create new key → JSON**.
6. Download the JSON file. **Treat it as a secret.**

### 6. Share the spreadsheet with the service account

Open the spreadsheet → **Share** → paste the service account's `client_email` (from the JSON) → set role to **Editor** → uncheck "Notify people" → **Share**.

> Without this step, every API call returns `403 PERMISSION_DENIED`.

---

## Environment Variables

The backend reads three required variables. Copy `.env.example` → `.env` at the repo root:

```env
GOOGLE_SERVICE_ACCOUNT_EMAIL=wedding-backend@your-project.iam.gserviceaccount.com
GOOGLE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEv...\n-----END PRIVATE KEY-----\n"
GOOGLE_SHEET_ID=1aBcD...xyz123...eFg
```

| Variable | Where to find it |
|---|---|
| `GOOGLE_SERVICE_ACCOUNT_EMAIL` | The `client_email` field in the downloaded JSON |
| `GOOGLE_PRIVATE_KEY` | The `private_key` field — **paste as-is** including `\n` escapes |
| `GOOGLE_SHEET_ID` | The ID from the spreadsheet URL |

**Optional overrides** (defaults shown):

```env
GUEST_LIST_TAB=Guests
MESSAGES_TAB=Messages
PAGE_SIZE=20
MAX_MESSAGE_LENGTH=500
RATE_LIMIT_SECONDS=30
```

> **Never commit `.env`.** It's in `.gitignore` by default. If the private key leaks, delete it in the Cloud Console → create a new key → update the env var → redeploy.

---

## Local Development

The project requires **two processes**: the Go API and the Vite frontend.

### Terminal 1 — Go API

```bash
cd backend
go mod tidy          # first time only
go run ./cmd/dev
```

Expected output:

```
Go API listening on http://localhost:3000
```

The dev server loads `.env` from the repo root automatically via `godotenv`.

### Terminal 2 — Vite frontend

```bash
cd frontend
npm install          # first time only
npm run dev
```

Open <http://localhost:5173>. The Vite proxy forwards `/api/*` to `http://localhost:3000`.

### Smoke test

```bash
# Health
curl -s http://localhost:3000/api/health
# → {"status":"ok","time":"..."}

# Guest resolution (also populates the Guests tab on first call)
curl -s http://localhost:3000/api/guests/andi-pratama
# → {"guest_id":"g_...","slug":"andi-pratama","guest_name":"Andi Pratama"}

# Guestbook read
curl -s "http://localhost:3000/api/messages?limit=5"
# → {"messages":[],"next_cursor":null}

# Guestbook write (guest)
curl -s -X POST http://localhost:3000/api/guests/andi-pratama/messages \
  -H 'Content-Type: application/json' \
  -d '{"message":"Selamat menempuh hidup baru!"}'
# → 201 Created

# Guestbook write (anonymous)
curl -s -X POST http://localhost:3000/api/messages \
  -H 'Content-Type: application/json' \
  -d '{"message":"Congratulations!"}'
# → 201 Created
```

In a browser:

- <http://localhost:5173/andi-pratama> → personalized greeting.
- <http://localhost:5173/> → "Anonymous" greeting.
- <http://localhost:5173/does-not-exist> → error page.

### Why not `vercel dev`?

Vercel's local dev server does not provide a Go runtime — `vercel dev` will fail with `No dev server available for service "backend"`. Vercel compiles and runs Go functions **only in the cloud**. Use the two-terminal workflow above for local iteration; use `vercel --prod` (or `vercel` for a preview URL) for cloud testing.

---

## Deployment

### 1. Push to GitHub

```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/<you>/wedding-invitation.git
git push -u origin main
```

### 2. Import to Vercel

1. <https://vercel.com/new> → **Import Git Repository** → select the repo.
2. **Framework Preset:** Vite.
3. **Root Directory:** `frontend` (or leave default — `vercel.json` handles the build).
4. **Environment Variables:** add all three `GOOGLE_*` variables for **Production, Preview, and Development**.
5. **Deploy**.

### 3. Custom domain (optional)

Project → **Settings → Domains** → add `wedding-domain.com`. Follow the DNS instructions.

### 4. Subsequent deploys

Push to `main` → Vercel deploys automatically. Or run:

```bash
vercel --prod
```

### 5. Verify production

```bash
curl -s https://your-domain.com/api/health
curl -s https://your-domain.com/api/guests/andi-pratama
```

Open `https://your-domain.com/andi-pratama` in a browser — the invitation should render with the correct greeting.

---

## Templates & Design Customization

### Change the active template

Edit `frontend/src/App.vue`, one line:

```js
// DESIGN EDIT: Change this value to select the active invitation template.
const activeTemplate = 'TemplateA' // 'TemplateA' | 'TemplateB' | 'TemplateC'
```

No backend or config changes required.

### Add a new template

1. Create `frontend/src/templates/TemplateD/` with `index.vue` and `style.css`.
2. Implement the same prop/event contract as the existing templates (see `TemplateA/index.vue`).
3. Import it in `App.vue` and add it to the `templates` object.
4. Change `activeTemplate` to `'TemplateD'`.

### Change wedding content

All wedding text, dates, and image paths live in:

```
frontend/src/config/wedding.config.js
```

```js
export const weddingConfig = {
  couple: { partnerOne: 'Andi', partnerTwo: 'Sari' },
  event: {
    date: '2027-01-01',
    time: '10:00',
    timezone: 'Asia/Makassar',
    venue: 'The Grand Ballroom',
    address: 'Jl. Contoh Raya No. 123, Makassar',
    mapsUrl: 'https://maps.google.com/?q=Grand+Ballroom+Makassar'
  },
  images: {
    hero: '/images/example.svg',
    gallery: ['/images/example-gallery-1.svg', /* ... */]
  },
  sections: {
    introduction: true,
    gallery: true,
    countdown: true,
    guestbook: true
  },
  introduction: 'With joy in our hearts...',
  footer: { closing: 'With love,', signature: 'Andi & Sari', credit: 'Made with love' }
}
```

### Replace images

Drop your files in `frontend/public/images/` and update the paths in `wedding.config.js`. No template edits needed.

### Shared design tokens

Colors, fonts, and spacing live in:

```
frontend/src/styles/variables.css
```

Change a token once → all templates that reference it update.

### Locate customization points

Search the repo for these comment markers:

| Marker | Purpose |
|---|---|
| `DESIGN EDIT:` | Visual/token changes |
| `CONTENT EDIT:` | Text and image path changes |
| `BEHAVIOR CHANGE:` | Animation or interaction logic |

### Handover reference

| Requested change | File |
|---|---|
| Change active template | `frontend/src/App.vue` |
| Add a new template | `frontend/src/templates/TemplateD/` |
| Change wedding details | `frontend/src/config/wedding.config.js` |
| Replace hero image | `frontend/public/images/` + config |
| Change shared palette | `frontend/src/styles/variables.css` |
| Change Template A layout | `frontend/src/templates/TemplateA/style.css` |
| Change guestbook behavior | `frontend/src/components/Guestbook.vue` + `services/api.js` |
| Change guest identity logic | `backend/internal/sheets/guests.go` |
| Change Google Sheets integration | `backend/internal/sheets/` |
| Change API communication | `frontend/src/services/api.js` |

---

## API Reference

Base URL: `https://your-domain.com/api`

### `GET /api/health`

Liveness check.

```json
{ "status": "ok", "time": "2026-10-30T12:00:00Z" }
```

### `GET /api/guests/:slug`

Resolve a personalized invitation.

**200 OK**

```json
{
  "guest_id": "g_a1b2c3d4",
  "slug": "andi-pratama",
  "guest_name": "Andi Pratama"
}
```

**404 Not Found** — unknown slug.

### `GET /api/messages?cursor=...&limit=20`

Public guestbook, newest first.

**200 OK**

```json
{
  "messages": [
    {
      "message_id": "m_xyz123",
      "guest_id": "g_a1b2c3d4",
      "display_name": "Andi Pratama",
      "message": "Selamat menempuh hidup baru!",
      "created_at": "2027-01-01T10:00:00+08:00"
    }
  ],
  "next_cursor": "m_prevmsgid"
}
```

`next_cursor` is `null` when no more pages exist. `limit` is capped at 100.

### `POST /api/guests/:slug/messages`

Submit a message as the resolved guest.

**Request**

```json
{ "message": "Selamat menempuh hidup baru!" }
```

**201 Created** → returns the created message.

**400** invalid / empty / too long · **404** unknown slug · **429** rate-limited.

### `POST /api/messages`

Submit a message as **Anonymous** (from the root URL).

Same body and responses as the personalized endpoint, with `guest_id` omitted and `display_name` set to `"Anonymous"`.

### Status codes

| Code | Meaning |
|---|---|
| 200 | Success |
| 201 | Message created |
| 400 | Invalid request (empty, too long, malformed JSON) |
| 404 | Guest not found |
| 429 | Rate limit exceeded |
| 500 | Internal error |
| 503 | Dependency unavailable |

---

## Guest Management

The bride (or wedding organizer) manages the guest list entirely inside Google Sheets.

| Task | Action |
|---|---|
| **Add a guest** | Type the name in **column A** of the `Guests` tab. Do not touch B, C, or D. |
| **Rename a guest** | Edit **column A**. The slug and guest ID stay stable, so old URLs still work. |
| **Delete a guest** | Delete the entire row. Old guestbook messages retain their `display_name`. |
| **Duplicate names** | Handled automatically — the backend generates unique slugs (`andi-pratama`, `andi-pratama-2`). |
| **Sort / filter rows** | Safe. IDs are stored as cell values, not row numbers, so sorting never breaks identity. |
| **View messages** | The `Messages` tab is append-only; the backend owns it. Do not edit rows manually. |
| **See invitations** | A guest's personalized URL is `https://your-domain.com/<slug>`. The slug is visible in column C. |

> The `guest_id`, `slug`, and `created_at` columns are populated the first time a guest's invitation is requested. A newly added name may show blanks until someone opens its URL — or you can trigger a sync manually with `curl https://your-domain.com/api/guests/<slug>`.

---

## Troubleshooting

<details>
<summary><b><code>package backend/internal/... is not in std</code></b></summary>

The `go.mod` module line does not match your import prefix. Verify:

```bash
cd backend
go env GOMOD          # must print a real path
head -1 go.mod        # must match the prefix of every internal import
```

If `go.mod` says `module backend`, imports must be `backend/internal/...`. If it says `module wedding-invitation`, imports must be `wedding-invitation/internal/...`.
</details>

<details>
<summary><b><code>malformed import path "backend/api/guests/[slug]"</code></b></summary>

Bracketed folder names are invalid in Go. Delete any `[slug]` directories under `backend/api/`. Routes with `:slug` are declared in `vercel.json` rewrites and dispatched via query string.
</details>

<details>
<summary><b><code>found packages X and Y in .../api</code></b></summary>

Two `.go` files with different `package` names live in the same directory. Move each handler into its own subfolder (`api/guest/index.go`, `api/messages/index.go`, etc.).
</details>

<details>
<summary><b><code>[DEV_SERVICE_START_FAILED] No dev server available for service "backend"</code></b></summary>

Vercel's local dev server does not support Go. Run the backend with `go run ./cmd/dev` and the frontend with `npm run dev` in separate terminals. See [Local Development](#local-development).
</details>

<details>
<summary><b><code>403 PERMISSION_DENIED</code> from Google Sheets</b></summary>

The spreadsheet is not shared with the service account. Open **Share** and add the `client_email` from the JSON key with **Editor** role.
</details>

<details>
<summary><b><code>invalid PEM block</code> / <code>crypto/rsa: verification error</code></b></summary>

`GOOGLE_PRIVATE_KEY` was pasted with real newlines instead of `\n` escapes, or quotes were stripped. Re-paste from the JSON file, wrap in double quotes, keep `\n` as literal two characters. The backend converts them at runtime via `strings.ReplaceAll`.
</details>

<details>
<summary><b>Every API call returns <code>404 Not Found</code> in production</b></summary>

Check `vercel.json` rewrites:

```json
{ "source": "/api/guests/:slug", "destination": "/backend/api/guest/index.go?slug=:slug" }
```

The path after `destination` must match the actual file location in your repo. Verify by looking at the deploy's function list in the Vercel dashboard.
</details>

<details>
<summary><b>Personalized URL returns 404 after refresh</b></summary>

SPA fallback rewrite is missing. Confirm `vercel.json` contains:

```json
{ "source": "/((?!api/).*)", "destination": "/index.html" }
```

The negative lookahead keeps API routes from being swallowed.
</details>

---

## Roadmap

The initial release intentionally omits the following (per PRD § 1.2):

- ❌ RSVP submission & attendance tracking
- ❌ Guest authentication
- ❌ Private guestbook messages
- ❌ Admin dashboard
- ❌ Online gifts / payments
- ❌ CMS

Ideas for a future version:

- Optional RSVP via a third spreadsheet tab
- WhatsApp share button on the invitation page
- Background music with a mute toggle
- QR code on printed cards linking to each guest's URL

---

## License

MIT — see `LICENSE` for details. Wedding photographs and personal content are not covered by the MIT license; they remain the property of the couple.

---

<p align="center">
  Made with love · <a href="https://vercel.com">Deployed on Vercel</a>
</p>