# Mock Server

Simulates the NYC Cares API for local development and CI. Runs as a plain HTTP server locally (`:3001`) or as a Lambda Function URL in the `-ci` environment.

## Endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| `POST` | `/user/login` | none | Issues a session cookie |
| `GET` | `/api/registrations/dashboard/upcoming/{userId}/{page}` | cookie | Returns future projects |
| `GET` | `/api/registrations/dashboard/today/{userId}/{page}` | cookie | Returns today's projects (filters by current date) |
| `GET` | `/api/messenger/channel/{channelId}/messages` | cookie | Returns channel messages |
| `POST` | `/api/messenger/channel/{channelId}/messages/post` | cookie | Sends a message |
| `POST` | `/api/messenger/create-pin-message/{campaignId}` | cookie | Pins a message |
| `GET` | `/` | none | Health check |

## Stateless Project Injection

There is no shared state. Project data flows through the session cookie:

1. **Login** — caller passes `mock_projects` form field: a base64-encoded JSON array of `{name, date, id}` objects.
2. **Cookie** — login sets `session: mock-session:<base64>` encoding that payload.
3. **Upcoming/today projects** — reads the cookie, decodes the projects, and builds the response. `/upcoming` returns all projects; `/today` filters to only those whose date matches the current date.

If `mock_projects` is omitted, upcoming projects returns a single default project dated 6 days from now.

**Date override**: set `NYCARES_CURRENT_DATE=YYYY-MM-DD` to shift what "now" means inside `MockUpcomingResponse`.

## Package Layout

```
main.go                        — entry point; gorilla/mux router; Lambda vs HTTP mode
middleware/requirecookie.go    — 401 if no session cookie present
routes/
  login.go                     — /user/login handler + cookie encoding
  admin.go                     — GetProjectsFromCookie helper (cookie → []ProjectConfig)
  upcomingprojects.go          — /upcoming and /today project route handlers
  messages.go                  — send, pin, channel message handlers
  mockresponses/
    schedulemock.go            — builds UpcomingResponse from []ProjectConfig
    messagesmock.go            — builds message-related responses
utils/utils.go                 — ValidateInternalID, ValidateUUID, NewUUID
```

## Adding a New Endpoint

1. Add a handler file under `routes/`.
2. Register it in `main.go` via `routes.RegisterXxxRoute(r)`.
3. If it needs auth, wrap the handler with `middleware.RequireCookieMiddleware`.
4. Add a mock response builder under `routes/mockresponses/` if the response is non-trivial.
