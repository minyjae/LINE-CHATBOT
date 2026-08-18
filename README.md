# Go Starter

Starter project สำหรับสร้าง REST API ด้วย Go โดยวางโครงสร้างแบบ **Hexagonal Architecture (Ports & Adapters)** เพื่อแยก business logic ออกจาก framework/infrastructure ให้ชัดเจน พร้อม config, database, JWT auth, validator และ response helper ที่ตั้งค่าไว้ให้แล้ว เหลือแค่เพิ่ม domain ของตัวเอง

## Tech Stack

- [Go](https://go.dev/) 1.25
- [Fiber v2](https://gofiber.io/) — HTTP framework
- [GORM](https://gorm.io/) + PostgreSQL — ORM / database
- [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) — JWT authentication
- [go-playground/validator v10](https://github.com/go-playground/validator) — request validation
- [godotenv](https://github.com/joho/godotenv) — โหลดค่า config จาก `.env`
- Docker & Docker Compose (API + PostgreSQL + pgAdmin)

## โครงสร้างโปรเจกต์

```
cmd/
  api/               # entry point (main.go): bootstrap config, db, fiber app
internal/
  adapters/
    http/
      handlers/      # controller ของแต่ละ domain (รับ request → เรียก service)
      middleware/     # เช่น JWT auth middleware
      routes/         # ลงทะเบียน route ของแต่ละ domain
    presistance/
      models/         # GORM model (โครงสร้างตาราง)
      repositories/   # implementation ของ repository ports (คุยกับ database)
  config/            # โหลด .env, เชื่อมต่อ database, migration
  core/
    domain/
      entities/       # business entity (plain Go struct, ห้าม import framework)
      ports/
        repositories/ # interface ที่ repository ต้อง implement
        services/     # interface ที่ service ต้อง implement
    services/         # business logic (implement service ports)
types/                 # shared type/DTO
utils/                 # response helper, JWT, password hash, pagination, generic controller
validators/            # request body validation ผ่าน go-playground/validator
```

แต่ละ layer มีไฟล์ `example.go` (และ `init.go` ในบาง folder) เป็นตัวอย่าง/แม่แบบไว้ดูตอนสร้าง domain ใหม่

## Requirements

- Go 1.25+
- PostgreSQL (หรือใช้ผ่าน Docker Compose ก็ได้)
- Docker + Docker Compose (ถ้าต้องการรันแบบ container)

## Getting Started

### 1. Clone และติดตั้ง dependencies

```bash
git clone <repo-url>
cd go-starter
go mod download
```

### 2. ตั้งค่า Environment Variables

สร้างไฟล์ `.env` ที่ root ของโปรเจกต์ (ดู `internal/config/config.go` สำหรับ field ทั้งหมด):

```env
# App
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Database (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_SSL=disable

# JWT
JWT_SECRET=change-me-to-a-random-32-char-secret
JWT_EXPIRES_IN=24h

# Migration
AUTO_MIGRATE=true

# Admin (จำเป็นเมื่อ APP_ENV=production)
ADMIN_EMAIL=admin@example.com

# pgAdmin (ใช้เฉพาะตอนรันผ่าน docker-compose)
PGADMIN_DEFAULT_EMAIL=admin@example.com
PGADMIN_DEFAULT_PASSWORD=admin
PGADMIN_PORT=5051
```

> ในโหมด `production` ระบบจะ validate ว่า `DB_PASSWORD`, `JWT_SECRET` (ต้องยาว ≥ 32 ตัวอักษร) และ `ADMIN_EMAIL` ต้องถูกตั้งค่าไว้ ไม่งั้น server จะไม่ start

### 3. รันแบบ local

```bash
go run ./cmd/api
```

Server จะรันที่ `http://localhost:${APP_PORT}` (ค่า default คือ `8080`)

### 4. รันด้วย Docker Compose

```bash
docker compose up --build
```

จะรัน 3 services: `api` (port ตาม `APP_PORT`), `postgres` (port ตาม `DB_PORT`) และ `pgadmin` (port ตาม `PGADMIN_PORT`)

## Database Migration

ค่า default จะรัน `AutoMigrate` อัตโนมัติเมื่อ `APP_ENV=development` หรือเมื่อตั้ง `AUTO_MIGRATE=true` — เพิ่ม model ที่ต้อง migrate ได้ที่ `runMigration()` ใน [internal/config/database.go](internal/config/database.go)

```go
func runMigration(db *gorm.DB) {
    err := db.AutoMigrate(&models.User{}, &models.List{})
    ...
}
```

## วิธีเพิ่ม Domain ใหม่

1. **Entity** — เพิ่ม business struct ที่ `internal/core/domain/entities/<name>.go` (ห้าม import framework ใดๆ)
2. **Ports** — ประกาศ interface ที่ `internal/core/domain/ports/repositories/` และ `internal/core/domain/ports/services/`
3. **Repository** — implement repository port ที่ `internal/adapters/presistance/repositories/<name>.go` และเพิ่ม GORM model ที่ `internal/adapters/presistance/models/`
4. **Service** — implement service port (business logic) ที่ `internal/core/services/<name>.go`
5. **Handler** — สร้าง controller ที่ `internal/adapters/http/handlers/<name>/` (รับ service ผ่าน constructor)
6. **Route** — ลงทะเบียน endpoint ที่ `internal/adapters/http/routes/<name>.go` แล้ว wire เข้ากับ `SetupRoute` ใน [routes.go](internal/adapters/http/routes/routes.go)
7. **Wire ใน main.go** — สร้าง repository → service → ส่งเข้า `routes.SetupRoute` ที่ [cmd/api/main.go](cmd/api/main.go)

ดูตัวอย่างเต็มรูปแบบได้จากไฟล์ `example.go` ในแต่ละ layer

## Utilities ที่มีให้แล้ว

| Package | หน้าที่ |
|---|---|
| `utils.Response` | Response helper มาตรฐาน (`Item`, `Created`, `BadRequest`, `NotFound`, `Unauthorized`, `ValidateFailed`, `ErrorHandler` ฯลฯ) |
| `utils.GenerateJWT` / `ValidateJWT` | ออกและตรวจสอบ JWT token |
| `utils` (password) | hash / compare password |
| `utils.Pagination` | struct ช่วยจัดการ pagination |
| `validators.ValidateStruct[T]` | middleware validate request body ด้วย struct tag (`binding`) |
| `middleware.JWTProtected` | middleware ตรวจสอบ `Authorization: Bearer <token>` แล้วเก็บ `user_id`, `email`, `display_name` ไว้ใน `c.Locals` |

## License

ยังไม่ได้กำหนด license — เพิ่มไฟล์ `LICENSE` เองตามความเหมาะสมของโปรเจกต์
