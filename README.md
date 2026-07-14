# Full-Stack Application Suite

This repo contains four independently runnable applications. Each project owns its runtime, tests, documentation, and local development workflow, with credential-free defaults for core functionality.

| Project | Product | Primary stack | Credential-free local path |
| --- | --- | --- | --- |
| [Agent AI](./agent-ai/README.md) | PDF-grounded question answering with optional web comparison | Node.js 22, Express 5, React 19, Vite, MCP | Extractive local retrieval and answers |
| [OnlineOrder](./onlineorder/README.md) | Secure restaurant browsing, cart, and checkout | Java 21, Spring Boot 3, PostgreSQL, React 18, Ant Design | In-memory H2 demo profile; PostgreSQL through Compose |
| [SocialAI](./socialai/README.md) | Authenticated media collection with AI image generation | Go, React 18, Vite, pluggable storage/search/AI adapters | JSON persistence, local media, generated SVG previews |
| [Spotify Local](./spotify/README.md) | Native album, favorites, and music-playback experience | Kotlin, Ktor, Android Compose, Hilt, Retrofit, Room, Media3 | Generated local covers and royalty-free WAV fixtures |

## Application walkthroughs

The first three walkthroughs were recorded from live local builds and follow the real application paths. Spotify's walkthrough is rendered from the checked-in fixtures and the app's Compose state because it requires an Android runtime; the animation labels this explicitly and does not present itself as emulator footage.

<table>
  <tr>
    <td width="50%">
      <a href="./agent-ai/README.md"><img src="./docs/assets/demos/agent-ai-demo.gif" alt="Agent AI live local workflow: upload a PDF, ask a grounded question, and inspect page sources" width="100%"></a><br>
      <strong>Agent AI</strong><br>
      <sub>Upload a real PDF → ask a grounded question → expand the matching page source.</sub>
    </td>
    <td width="50%">
      <a href="./onlineorder/README.md"><img src="./docs/assets/demos/onlineorder-demo.gif" alt="OnlineOrder live local workflow: sign in, browse a menu, add items, and check out" width="100%"></a><br>
      <strong>OnlineOrder / Lai Food</strong><br>
      <sub>Sign in → browse the seeded menu → aggregate quantity and total → check out.</sub>
    </td>
  </tr>
  <tr>
    <td width="50%">
      <a href="./socialai/README.md"><img src="./docs/assets/demos/socialai-demo.gif" alt="SocialAI live local workflow: generate an image, publish it, and open it from the collection" width="100%"></a><br>
      <strong>SocialAI</strong><br>
      <sub>Generate with the local AI adapter → publish → browse and open the collection item.</sub>
    </td>
    <td width="50%">
      <a href="./spotify/README.md"><img src="./docs/assets/demos/spotify-demo.gif" alt="Source-driven Spotify Compose walkthrough: browse Home, open a playlist, play and favorite a song, then view Favorites" width="100%"></a><br>
      <strong>Spotify Local</strong><br>
      <sub>Source-driven Compose walkthrough: Home → playlist → playback and favorite → Favorites.</sub>
    </td>
  </tr>
</table>

Static poster frames are stored beside the animations in [`docs/assets/demos`](./docs/assets/demos/).

## System flow charts

These charts follow the concrete request, storage, and playback paths implemented in each project.

### Agent AI

```mermaid
flowchart LR
    user["User"] --> ui["React document and chat UI"]
    ui -->|"PDF multipart upload"| api["Express security, CORS, and rate limits"]
    api --> pdf["Multer memory upload and PDF parser"]
    pdf --> store["Expiring in-memory DocumentStore"]
    store --> chunks["Page-aware chunks"]
    ui -->|"documentId + question"| api
    chunks --> rank["Lexical ranker: top 4 chunks"]
    api --> rank
    rank --> local["Local extractive answer"]
    rank -.-> openai["Optional OpenAI Responses synthesis"]
    api -.-> mcp["Optional includeWeb → MCP client over stdio"]
    mcp --> search["MCP search server → SerpAPI"]
    local --> response["Answer + page sources"]
    openai --> response
    search --> response
    response --> ui
```

### OnlineOrder / Lai Food

```mermaid
flowchart LR
    customer["Customer"] --> react["React + Ant Design client"]
    react -->|"Vite proxy or Nginx /api"| security["Spring Security: session + CSRF"]
    security --> controllers["Auth, menu, customer, and cart controllers"]
    controllers --> services["Transactional services"]
    services <--> cache["Spring Cache: Caffeine in PostgreSQL mode"]
    services --> jdbc["Spring Data JDBC repositories"]
    jdbc --> profile{"Active data profile"}
    profile -->|"demo"| h2["In-memory H2 + seeded schema"]
    profile -->|"default / Compose"| postgres["PostgreSQL + idempotent seed data"]
    services -->|"checkout"| clear["Delete order items and reset cart total"]
    clear --> react
```

### SocialAI

```mermaid
flowchart LR
    user["Signed-in user"] --> web["React routes + token storage"]
    web -->|"JSON or multipart + Bearer JWT"| middleware["Go HTTP security, CORS, and auth middleware"]
    middleware --> handlers["Auth, post, search, and AI handlers"]
    handlers --> repo{"Repository adapter"}
    repo --> memory["Memory"]
    repo --> file["JSON file"]
    repo --> elastic["Elasticsearch"]
    handlers --> media{"Media adapter"}
    media --> localMedia["Local filesystem"]
    media --> gcs["Google Cloud Storage"]
    handlers --> generator{"Image generator"}
    generator --> svg["Deterministic local SVG"]
    generator --> images["OpenAI Images"]
    generator -->|"preview object"| media
    handlers -->|"publish metadata"| repo
    repo -->|"collection results"| web
```

### Spotify Local

```mermaid
flowchart LR
    user["Android user"] --> compose["Compose screens + Navigation"]
    compose --> vm["Home, Playlist, Favorites, and Player ViewModels"]
    vm --> feedRepo["Feed / playlist repositories"]
    feedRepo --> retrofit["Retrofit client"]
    retrofit --> ktor["Ktor fixture API"]
    ktor --> catalog["Validated feed.json + playlists.json"]
    catalog --> cover["Generated SVG covers"]
    catalog --> audio["Generated WAV + byte-range responses"]
    vm --> favorites["Favorites repository"]
    favorites --> room["Room DAO Flow"]
    vm --> playback["PlaybackController"]
    playback --> media3["Media3 ExoPlayer"]
    room --> compose
    media3 --> compose
```

## Repository layout

```text
full-stack-projects/
├── README.md                  Monorepo overview and verification guide
├── docs/assets/demos/         README animations and static poster frames
├── agent-ai/
│   ├── client/                React/Vite browser app
│   ├── server/                Express API, RAG, PDF, and MCP layers
│   ├── pnpm-lock.yaml         Shared workspace dependency lock
│   └── docker-compose.yml     Nginx web + API stack
├── onlineorder/
│   ├── backend/               Spring Boot/Gradle API
│   ├── frontend/              React/Vite/Ant Design client
│   └── docker-compose.yml     PostgreSQL + API + Nginx web stack
├── socialai/
│   ├── backend/               Go API and adapter implementations
│   ├── web/                   React/Vite client
│   ├── data/ and media/       Ignored local persistence locations
│   ├── .github/workflows/     Go and web CI
│   └── docker-compose.yml     API + Nginx web stack
└── spotify/
    ├── backend/               Independent Ktor/Gradle API
    ├── android/               Independent Android/Gradle application
    ├── scripts/               Fixture smoke and offline structure checks
    └── compose.yaml           Containerized Ktor API
```

There is intentionally no shared application runtime or root dependency graph. Each project owns its lockfile, Gradle wrapper, environment template, tests, and deployment files. That keeps builds reproducible and lets developers work on one project without installing the other stacks.

The root README is the repository landing page, while each application keeps its own dependency graph, lockfiles, environment template, tests, deployment files, and subtree-specific ignore rules.

Several services default to port `8080`—OnlineOrder, SocialAI, and Spotify's backend—so run them one at a time or change their documented port settings. Agent AI uses API port `5001` by default.

## Project highlights

### Agent AI

Agent AI is a local-first retrieval-augmented generation application. A user uploads a real PDF, receives an opaque expiring document session, asks questions, and sees page-aware source excerpts with each grounded answer. The responsive React interface includes drag-and-drop upload, conversation history, dictation, speech playback, keyboard operation, and user-facing error states.

The Express API performs PDF validation and extraction, deterministic passage ranking, bounded uploads and questions, rate limiting, CORS allowlisting, and in-memory session expiry. Local mode needs no credentials. Optional provider boundaries add OpenAI Responses synthesis and an MCP stdio child server backed by SerpAPI without exposing keys to the browser.

Key capabilities:

- real document ingestion and source provenance rather than canned chat responses;
- a testable local fallback around optional AI and search providers;
- per-document isolation instead of a process-global upload path;
- security headers, origin controls, limits, expiry, and early deletion;
- containerized Nginx frontend and private API routing.

### OnlineOrder / Lai Food

OnlineOrder is a transactional full-stack food-ordering application. Customers can register, sign in, browse seeded restaurants and menus, add quantities to a personal cart, see precise decimal totals, and complete a demo checkout. Checkout atomically clears the active cart; payment processing and a persisted order ledger are outside this project's scope. PostgreSQL schema constraints and indexes enforce the domain model, while idempotent seed scripts make repeated local starts safe.

The Java 21 backend uses Spring Boot Web, Security, Data JDBC, Validation, Actuator, Caffeine, BCrypt, server-side sessions, and CSRF protection. The React 18/Ant Design client handles authentication, menu browsing, cart interactions, loading, and failures. Docker Compose connects PostgreSQL, a non-root API image, and an Nginx frontend.

Key capabilities:

- session fixation protection and CSRF on state-changing browser requests;
- passwords in request bodies and BCrypt hashes at rest;
- authenticated, per-customer carts with mutation-driven Spring Cache eviction—Caffeine in the default/PostgreSQL profile and the simple in-memory provider in the demo profile;
- H2-backed service and HTTP integration tests independent of Docker;
- production JAR and frontend bundle verification.

### SocialAI

SocialAI is an authenticated media-sharing application. Users can register, sign in, upload images or videos, search by creator and all caption terms, delete only their own posts, generate an AI image preview, discard it, or publish it into the collection.

The Go backend deliberately uses the standard library and interfaces around three replaceable boundaries: repository (`memory`, JSON file, or Elasticsearch), media (`local` or Google Cloud Storage), and image generation (local SVG or OpenAI). The React app provides register, login, creation, generation, publishing, and collection flows. Local defaults need no cloud account.

Key capabilities:

- PBKDF2 password hashing and strictly validated, expiring HS256 tokens;
- upload size/type validation, storage traversal protection, CORS, and ownership checks;
- race-tested repository and HTTP behavior;
- optional Elasticsearch, GCS, and OpenAI adapters behind local implementations;
- a verified end-to-end browser journey from registration through AI publishing.

### Spotify Local

Spotify Local pairs a Ktor fixture API with a native Android application. The app presents feed sections, navigates to playlist details, persists favorite albums in Room, and controls Media3/ExoPlayer through an activity-scoped floating player with play, pause, progress, and seek behavior.

The Android client uses Compose, MVVM, `StateFlow`, Hilt, Retrofit, Navigation Compose, Room, Coil, and a playback interface that can be replaced in unit tests. The Ktor server preserves the documented feed/playlist/song contracts while generating deterministic SVG covers and five-second WAV tracks at request time, avoiding copyrighted binaries and external media hosting.

Key capabilities:

- API, repository, database, ViewModel, navigation, and playback separation;
- local favorites that survive process restarts;
- shared playback UI across Home, Favorites, and Playlist destinations;
- complete debug APK, Android-test APK, Hilt/Room code generation, and lint verification;
- credential-free, referentially validated sample media.

## Testing and verification

The following commands and results were verified on 2026-07-11. Test counts refer to executed test cases/functions, not just test files.

### Agent AI — 21 tests

```bash
cd agent-ai
pnpm test
pnpm check
pnpm build
pnpm audit --audit-level high
```

- Server: 18 passing Node tests.
- Client: 3 passing Vitest interaction tests.
- Syntax/type-adjacent checks, production client build, and high-severity dependency audit passed.
- The localhost UI was rendered in a browser with no console errors.

Integration coverage: the credential-free PDF/RAG path is covered by the test suite. Compose requires Docker, while live OpenAI and SerpAPI calls require user-owned keys.

### OnlineOrder — 24 tests

```bash
cd onlineorder/backend
./gradlew clean test bootJar

cd ../frontend
pnpm install --frozen-lockfile
pnpm audit --audit-level=high
pnpm test
pnpm build
```

- Backend: 15 passing tests—7 authenticated API integration cases and 8 service cases.
- Frontend: 9 passing Vitest cases.
- Spring Boot JAR, production web bundle, and high-severity dependency audit passed.

Integration coverage: the backend suite uses its PostgreSQL-compatible H2 profile. Run the Compose stack for a PostgreSQL integration smoke test.

### SocialAI — 26 tests

```bash
cd socialai/backend
go test -race ./... -count=1
go vet ./...
go build -o /tmp/socialai-server ./cmd/server

cd ../web
pnpm test
pnpm build
pnpm audit --audit-level high
```

- Backend: 20 passing Go test functions under the race detector.
- Frontend: 2 passing files / 6 passing Vitest cases.
- `go vet`, server binary build, production web bundle, and high-severity dependency audit passed.
- A live local API and browser flow completed register → login → local image generation → publish → collection with no console errors.

Integration coverage: local adapters are covered by the test suite. Elasticsearch, GCS, and OpenAI integration checks require their external services or credentials; Compose requires Docker.

### Spotify Local — 13 executed tests, plus one packaged instrumentation test

```bash
cd spotify/backend
./gradlew test
PUBLIC_BASE_URL=http://127.0.0.1:8080 ./gradlew run

# In a second terminal while the backend is running:
cd ..
BASE_URL=http://127.0.0.1:8080 ./scripts/check-fixtures.sh
python3 scripts/validate-project.py

cd android
./gradlew testDebugUnitTest
./gradlew assembleDebug
./gradlew lintDebug
./gradlew assembleDebugAndroidTest
```

- Backend: 7 passing Ktor route, error, serialization, audio, byte-range, cover, health, and fixture tests.
- Android JVM: 6 passing Retrofit-contract, Home ViewModel, and Player ViewModel tests.
- Live Ktor smoke check passed `/health`, `/feed`, `/playlists`, `/playlist/1`, RIFF audio validation, and a partial-content byte-range request.
- Debug APK assembled and its v2 signature verified with `apksigner`.
- Android lint completed with 0 errors; its 15 warnings were newer-dependency notices rather than code defects.
- The Room DAO instrumentation test compiled and packaged into the Android-test APK.

Integration coverage: backend audio and Android JVM behavior are covered locally. Run `connectedDebugAndroidTest` on an emulator or physical device for Room instrumentation and hands-on Media3 playback; Compose requires Docker.

### Test total

The verified suites executed **84 tests**: 21 Agent AI, 24 OnlineOrder, 26 SocialAI, and 13 Spotify backend/JVM tests. Browser smoke checks, static checks, builds, lint, audits, race detection, fixture validation, APK verification, and the compiled Room instrumentation test are additional evidence and are not included in that number.

## Repository hygiene

The source-level audit found no credentials or private keys in the project trees:

- There are no real `.env`, keystore, PEM, service-account, or local Android property files under the four project trees.
- Every `.env.example` contains blank credential fields or safe localhost values only.
- Agent AI's upload directory, and SocialAI's data/media directories, contain only `.gitkeep`; no user uploads or persisted records were found.
- OpenAI, SerpAPI, PostgreSQL, JWT, Elasticsearch, and cloud credential values remain external configuration.
- Gradle Wrapper JARs, pnpm lockfiles, SocialAI's CI workflow, and Spotify's exported Room schema are intentional source artifacts.
- Java, Go, Node.js, and Android toolchains are external prerequisites and are not committed to the repository.

Generated dependencies, build products, Gradle caches, local SDKs, and verification binaries are excluded from version control:

| Project | Tracked generated material | Repository policy |
| --- | --- | --- |
| Agent AI | None | Ignore rules cover future `node_modules` and `dist` output |
| OnlineOrder | None | Ignore rules cover future frontend and Gradle output |
| SocialAI | None | Local `data`/`media` payloads and web output remain ignored |
| Spotify | None | Ignore rules cover Android SDK metadata, build output, and Gradle/Kotlin caches |

The applications share one repository for discovery while remaining independently buildable and runnable.

## Running the projects

Start with the individual project README linked in the table above. Each contains its exact prerequisites, configuration, API contract, direct-development commands, Docker option, security notes, and troubleshooting guidance.

For a clean review workflow:

1. Choose one project and copy its `.env.example` to `.env` only when its local instructions require it.
2. Install dependencies from its checked-in lockfile or use its Gradle wrapper.
3. Run that project's verified test commands before starting services.
4. Use its local, credential-free adapters first.
5. Enable paid/cloud adapters only with scoped credentials stored outside source control.

The applications are independent projects rather than microservices that depend on one another.
