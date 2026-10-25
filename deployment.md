---
layout: center
transition: slide-left
class: "text-center"
---

# Deployment

นำ Go REST API จากเครื่องพัฒนาขึ้นใช้งานบน Render

---
layout: two-cols-header
---

# เส้นทางจาก Local สู่ Production

1. เตรียม config และ secrets ให้เหมาะกับ production
2. สร้าง migration ให้ครบ แล้วเตรียม PostgreSQL
3. Build และ deploy API เป็น Render Web Service
4. ตรวจสอบ API, database และไฟล์ที่อัปโหลด

<br />

**หลักคิด:** โค้ดชุดเดียวกัน แต่ค่าตั้งต้นและบริการภายนอกต่างกันตาม environment

---
layout: two-cols-header
---

# Development กับ Production Configuration

:::left::

ในเครื่องพัฒนา `.env` ช่วยโหลดค่าจากไฟล์ผ่าน `godotenv` ส่วน production ควรรับค่าจาก System Environment ที่ตั้งใน Render

โค้ดปัจจุบันเรียก `godotenv.Load()` ทุก environment และคืน error เมื่อไม่มี `.env` จึงต้องปรับให้โหลดไฟล์เฉพาะตอน dev ก่อน deploy

:::right::

<Console fontSize="10px" lineHeight="1.2" filename="config.go">
```go
if os.Getenv("APP_ENV") != "production" {
    _ = godotenv.Load()
}
```
</Console>

production ไม่ต้องมีไฟล์ `.env` บนเครื่องเซิร์ฟเวอร์

---
layout: two-cols-header
---

# Environment Variables และ Secrets

:::left::

ตั้งค่า Environment Variables ในหน้า Render ของ Web Service

- `DB_URL` — Internal Database URL ของ Render PostgreSQL
- `JWT_SECRET` — secret สำหรับเซ็น JWT ที่สุ่มและเก็บเป็นความลับ
- `JWT_ACCESS_EXPIRY` — อายุ access token เช่น `5m`
- `JWT_REFRESH_EXPIRY` — อายุ refresh token เช่น `2160h`
- `PORT` — Render กำหนดให้ service ใช้

ห้าม commit `.env` หรือใช้ `JWT_SECRET=secret` ซึ่งเป็น fallback ในโค้ดปัจจุบันบน production

:::right::

<Console fontSize="10px" lineHeight="1.2" filename="config.go">
```go
DB_URL=postgres://...
JWT_SECRET=<ค่าลับที่สุ่มใหม่>
JWT_ACCESS_EXPIRY=5m
JWT_REFRESH_EXPIRY=2160h
```
</Console>

---
layout: two-cols-header
---

# CORS: อนุญาตให้ Frontend เรียก API

:::left::

CORS เป็นกติกาที่ Browser ใช้ควบคุมการเรียก API ข้าม origin เช่น frontend กับ API อยู่คนละ domain

กำหนด origin ของ frontend ให้ชัดเจน รวมถึง methods และ headers ที่ต้องใช้ หากส่ง cookie หรือ credentials ต้องตั้งค่าฝั่ง server และ client ให้สอดคล้องกัน

**โค้ดปัจจุบันยังไม่มี CORS middleware** จึงต้องเพิ่มและทดสอบกับ origin จริงก่อน deploy frontend

:::right::

```text
Frontend
https://shop.example.com
        │  API request
        ▼
Go API
https://api.example.com
        │
        └─ อนุญาต origin ที่ระบุไว้
```

หลีกเลี่ยงการอนุญาตทุก origin ใน production โดยไม่จำเป็น

---
layout: two-cols-header
---

# Build Go Project

:::left::

ในโปรเจกต์นี้ entry point ของ API อยู่ที่ `cmd/api`

```bash
go build -o app ./cmd/api
```

คำสั่งนี้ compile package ที่จำเป็นและสร้าง executable ชื่อ `app`

บน Render กำหนด **Build Command** เป็น:

```bash
go build -o app ./cmd/api
```

และ **Start Command** เป็น:

```bash
./app
```

:::right::

<Console fontSize="10px" lineHeight="1.2" filename="cmd/api/main.go">
```go
if err := r.Run(":" + cfg.PORT); err != nil {
    log.Fatalf("failed to run server: %v", err)
}
```
</Console>

ให้ server bind กับ `PORT` ที่ platform กำหนด ไม่กำหนด port production แบบตายตัว

---
layout: two-cols-header
---

# เตรียม Database Migrations ให้ครบ

:::left::

ระหว่าง dev เราอาจใช้ `atlas schema apply` เปลี่ยน schema บน database โดยตรง การเปลี่ยนแปลงนั้น **ไม่ได้สร้างไฟล์ migration ให้เอง**

ก่อน deploy ให้เทียบ models กับ migrations แล้วสร้าง migration ที่ขาดให้ครบ จากนั้นตรวจ SQL และ commit ไฟล์ migration เข้า Git

```bash
task migrate:diff -- add_missing_schema
task migrate:status
```

ตัวอย่างนี้ใช้ Atlas ตาม `atlas.hcl` และ Taskfile ของโปรเจกต์

**ข้อควรระวัง:** `atlas.hcl` ปัจจุบันอ่าน `DB_URL` จากไฟล์ `.env` โดยตรง การ export `DB_URL` ของ Render อย่างเดียวอาจไม่เปลี่ยน URL ที่ Atlas ใช้ ต้องปรับ config ให้รับ System Env หรือส่ง config/file ชั่วคราวที่ไม่ commit ก่อนรัน migration

:::right::

```text
แก้ GORM models
      ↓
สร้าง migration ด้วย Atlas
      ↓
ตรวจ SQL + commit ไฟล์
      ↓
apply migration ไปยัง production DB
```

อย่าพึ่งพา schema ที่เคยแก้ด้วยมือหรือ `schema apply` บน dev database เพราะ production ต้องสร้าง schema ซ้ำได้จาก migration ใน repository

---
layout: two-cols-header
---

# สร้าง Render PostgreSQL

1. สร้าง PostgreSQL database ใน Render
2. เลือก region เดียวกับ Web Service
3. คัดลอก **Internal Database URL** สำหรับ service ที่อยู่ใน Render เดียวกัน
4. นำ URL ไปตั้งเป็น `DB_URL` ใน Web Service

Internal URL ช่วยให้ service คุยกันผ่าน private network; อย่านำค่า URL ที่มีรหัสผ่านไปใส่ใน source code หรือสไลด์

---
layout: two-cols-header
---

# ตั้งค่า Render Web Service

:::left::

สร้าง **Web Service** จาก Git repository แล้วกำหนด:

- **Runtime:** Go
- **Build Command:** `go build -o app ./cmd/api`
- **Start Command:** `./app`
- **Region:** ตรงกับ PostgreSQL
- **Environment Variables:** `DB_URL`, `JWT_SECRET` และค่า JWT expiry

ตั้งค่า environment ก่อน deploy ครั้งแรก เพื่อให้แอปเชื่อมต่อ database และสร้าง token ได้

:::right::

```text
Git repository
      │ build
      ▼
Render Web Service ─── Internal Network ─── Render PostgreSQL
      │
      └── PORT จาก Render
```

---
layout: two-cols-header
---

# Apply Migrations และ Deploy

:::left::

ก่อนปล่อยเวอร์ชันใหม่:

1. ตรวจว่า migration ที่ขาดถูกสร้างและ commit แล้ว
2. ตั้ง Atlas ให้ชี้ไปยัง Render PostgreSQL อย่างปลอดภัย แล้วรัน `task migrate:apply`
3. ตรวจผล migration และ log
4. Deploy Web Service และตรวจว่าเริ่มทำงานได้

แอปปัจจุบันเชื่อมต่อ DB ตอนเริ่ม แต่ไม่ได้เรียก Atlas หรือ migrate อัตโนมัติ

:::right::

```bash
task migrate:status
task migrate:apply
```

ตรวจให้แน่ใจว่า Atlas CLI และ config ที่ migration ต้องใช้พร้อมใน environment ที่ใช้รันคำสั่ง หากยังไม่พร้อม ให้รันจากเครื่องหรือ job ที่เข้าถึง Render DB ได้อย่างปลอดภัย และอย่า commit ไฟล์ที่มี production credentials

---
layout: two-cols-header
---

# ตรวจสอบหลัง Deploy

- ดู Build และ Runtime Logs ใน Render
- ตรวจว่า service เริ่มฟังที่ `PORT` ได้
- เรียก API endpoint จริง เช่น `GET /api/v1/products`
- ลองเส้นทางที่ต้องใช้ฐานข้อมูลและ JWT
- ทดสอบการเรียกจาก frontend เพื่อยืนยัน CORS
- ตรวจว่าภาพที่อัปโหลดเปิดผ่าน `/assets/...` ได้

ปัจจุบันยังไม่มี health-check endpoint แยก หากต้องการ health check อัตโนมัติควรเพิ่ม route เช่น `/health` และตั้งค่าใน service

---
layout: two-cols-header
---

# Smoke Test: หน้าเว็บอ่านรายการสินค้า

:::left::

สร้างหน้า HTML ทดสอบจาก origin ที่อนุญาตใน CORS แล้วใช้ `fetch()` เรียก API ที่ deploy ไว้

API ปัจจุบันเปิดรายการสินค้าที่ `GET /api/v1/products` และตอบกลับเป็น JSON array

:::right::

<Console fontSize="9px" lineHeight="1.15" filename="index.html">
```html
<ul id="products"></ul>
<p id="status"></p>
<script>
const api = "https://your-api.onrender.com";
const list = document.querySelector("#products");

fetch(`${api}/api/v1/products`)
  .then((response) => {
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    return response.json();
  })
  .then((products) => {
    list.replaceChildren(...products.map((product) => {
      const item = document.createElement("li");
      item.textContent = `${product.name} — ฿${product.price}`;
      return item;
    }));
  })
  .catch((error) => {
    document.querySelector("#status").textContent = error.message;
  });
</script>
```
</Console>

เมื่อหน้าเว็บและ API อยู่คนละ origin การทดสอบนี้ช่วยยืนยันทั้ง API response และ CORS

---
layout: two-cols-header
---

# ข้อจำกัดของ Render Free Tier

- Web Service หยุดทำงานหลังไม่มี traffic 15 นาที และอาจใช้เวลาประมาณหนึ่งนาทีในการตื่น
- Filesystem เป็นแบบ ephemeral: ไฟล์ในเครื่องหายเมื่อ restart หรือ deploy ใหม่
- Free PostgreSQL มีพื้นที่ 1 GB และหมดอายุหลังสร้าง 30 วัน
- Free PostgreSQL ไม่มี backup; หลังหมดอายุมีช่วงเวลา 14 วันให้อัปเกรดก่อนข้อมูลถูกลบ
- มีข้อจำกัดด้านชั่วโมงใช้งานและทรัพยากร

เหมาะกับการเรียน ทดลอง และ hobby project มากกว่าระบบ production ที่ต้องพร้อมใช้งานและเก็บข้อมูลระยะยาว

---
layout: two-cols-header
---

# Upload Files: ใช้ Object Storage

:::left::

โค้ดปัจจุบันบันทึกไฟล์ลง `./uploads` และให้บริการผ่าน `/assets` ไฟล์ที่อยู่บน Render Free จะหายเมื่อ service restart หรือ deploy ใหม่

อย่าใช้ filesystem ของ Web Service เป็นที่เก็บไฟล์ถาวร

แนวทางคืออัปโหลดไฟล์ไปยัง object storage เช่น **AWS S3** หรือ **Cloudflare R2** แล้วเก็บ object key หรือ URL ไว้ใน PostgreSQL

:::right::

```text
Client ── upload ──▶ Go API ──▶ S3 / Cloudflare R2
                                      │
                                      └── เก็บ URL/key ใน PostgreSQL
```

แยกข้อมูล metadata ใน database ออกจากไฟล์ binary เพื่อให้ไฟล์อยู่ได้แม้ API restart หรือ deploy ใหม่

---
layout: center
class: "text-center"
---

# Deployment Checklist

**Config:** production ใช้ System Env และตั้ง secrets จริง

**Database:** สร้าง migration ที่ขาดให้ครบ แล้ว apply กับ PostgreSQL

**Service:** build, start ด้วย `PORT` และตรวจ logs/API/CORS

**Files:** ย้าย uploads ไป S3, Cloudflare R2 หรือ object storage

**Free tier:** เข้าใจ cold start, ephemeral filesystem และอายุฐานข้อมูล 30 วัน
