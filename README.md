# ⏣ Titan Surface

Surface is a high-performance and top-level utils provider native extension for the TitanPl framework, built entirely in **Go**. It handles data-heavy, IO-bound, and system-level tasks outside the JavaScript runtime to ensure sub-millisecond response times and rock-solid stability.

---

## 📖 Complete API Reference

### 🖼️ Image Module (`image`)
Professional-grade image processing using Lanczos3 interpolation.

| Function | Description |
| :--- | :--- |
| `image.resize(opts)` | Resizes an image locally or from a URL. Returns Base64 if `dist` is omitted. |
| `image.crop(opts)` | Fills a square/rectangle and center-crops excess. Perfect for thumbnails. |

**Example: Local to Local**
```javascript
image.resize({ src: "big.jpg", dist: "small.jpg", width: 800 });
```

**Example: URL to Base64 (Zero-Disk)**
```javascript
const { base64 } = image.crop({ 
  src: "https://site.com/user.jpg", 
  width: 200, 
  height: 200 
});
```

---

### 📊 CSV Module (`csv`)
Sub-millisecond CSV engine with background pre-fetching.

| Function | Description |
| :--- | :--- |
| `csv.open(path, opts)` | Opens file and starts native background reader. |
| `csv.next(h, opts)` | Fetches next chunk from native buffer (v. fast). |
| `csv.readAll(h)` | Returns ALL data in one native call (Zero JS overhead). |
| `csv.create(path, opts)`| Creates a new CSV file natively. |
| `csv.write(h, rows)` | Writes record(s) to the file. Supports single object or array. |
| `csv.close(h)` | Closes handle and kills background threads. |

**Example: Ultimate Bulk Read**
```javascript
const h = csv.open("data.csv", { header: true, mode: "object" });
try {
  const data = csv.readAll(h);
} finally {
  csv.close(h);
}
```

---

### 📧 SMTP Module (`smtp`)
Enterprise email system with connection pooling and native rendering.

| Function | Description |
| :--- | :--- |
| `smtp.send(opts)` | Sends a single email. Supports STARTTLS and Direct SSL. |
| `smtp.bulk(opts)` | Sends multiple emails concurrently via parallel worker pool. |
| `smtp.render(tpl, data)` | Renders a Go HTML template string natively. |
| `smtp.renderFile(path, d)` | Renders a `.tmpl` or `.html` file directly from disk. |

**Example: Native Rendering**
```javascript
const body = smtp.renderFile("./welcome.tmpl", { name: "Soham" });
smtp.send({ ...creds, body, raw: true });
```

---

## 🌍 Real-World Industrial Workflows

### ⚡ Case 1: Ultra-Fast OTP System
Send transactional OTPs in microseconds by pre-rendering templates natively.
```javascript
import { smtp } from "@titanpl/surface";

export function sendOTP(req) {
  // Render OTP template natively in Go
  const body = smtp.render("<h1>Your OTP is: {{.code}}</h1>", { 
    code: Math.floor(1000 + Math.random() * 9000) 
  });

  // Fast Native Delivery
  return smtp.send({
    host: "smtp.gmail.com",
    username: "...",
    password: "...",
    to: req.userEmail,
    subject: "Your Login Code",
    body
  });
}
```

### 🖼️ Case 2: Production Profile Picture Pipeline
Process user uploads from URLs and save the optimized Base64 string directly to the DB.
```javascript
export function updateAvatar(req) {
  const result = image.crop({
    src: req.imageUrl,
    width: 250,
    height: 250,
    quality: 90
  });

  return db.users.update(req.userId, { avatar: result.base64 });
}
```

### 📈 Case 3: Massive Bulk Personalization
Combining all modules into a high-speed data pipeline.
```javascript
import { csv, smtp, path } from "@titanpl/surface";

export function marketingCampaign() {
  // 1. Read 100k leads natively
  const h = csv.open("leads.csv", { mode: "object" });
  const leads = csv.readAll(h);
  csv.close(h);

  // 2. Prep Rendered Data
  const tpl = "../app/emails/promo.tmpl";
  const jobs = leads.map(l => ({
    to: l.email,
    body: smtp.renderFile(tpl, { name: l.name })
  }));

  // 3. Parallel Blast via 10 workers
  return smtp.bulk({ ...creds, emails: jobs, concurrency: 10, raw: true });
}
```

---
Built with ❤️ by the TitanPL Team.
