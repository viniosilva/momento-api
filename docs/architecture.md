# Architecture — Decisions (V1)

## 1) High-level architecture

**Backend**
- Single Go API (modular monolith)
- Exposes REST only

**Frontend**
- React app (TanStack Router)
- Talks directly to the API (no BFF)
- 100% client-side (no SSR)

---

## 2) Authentication & Accounts

**Login**
- Email + password
- No magic link (v1)
- No OAuth (v1)

**Session strategy**
- JWT access token
- Refresh token flow
- Email verification required in V1

---

## 3) Users, Roles & Permissions

**Roles (V1)**
- Event Owner
- Participant

**Permissions**

| Action | Owner | Participant |
|---|---|---|
| Invite people | ✅ | ❌ |
| Delete photos | ✅ | ❌ |
| Edit event | ✅ | ❌ |

---

## 4) Events (Core Domain)

- Unlimited events per user
- No participant limit (V1)
- Invitations via **public invite link**

---

## 5) Photos (Critical Feature)

**Storage**
- S3-compatible object storage

**Upload flow**
- Signed URL upload (client → storage)

**Processing**
- Image compression: **Yes**
- Thumbnails: **No (V1)**
- EXIF/metadata: **No**

**Limits**
- Max photo size defined in V1

---

## 6) Database

- MongoDB
- No complex search required
- No soft delete
- Users can permanently delete data
- Events can be archived

---

## 7) Initial Scalability

- Designed for **hobby scale**

---

## 8) Infra & Deploy

- Cloud provider: **Hetzner**
- Docker: **Yes**
- CI/CD from day one: **Yes**
- File storage in same provider: **Yes**

---

## 9) Observability (Initial)

- Structured logs: **Yes**
- Error tracking: **Later**
- Metrics: **Later**