# ⏣ Titan Surface
**The Native Capability Layer for TitanPL**

Surface is a high-performance native extension layer for the Titan framework, built entirely in **Go**. It handles data-heavy, IO-bound, and system-level tasks outside the JavaScript runtime to ensure sub-millisecond response times and rock-solid stability.

---

## ⚡ Performance Philosophy
Surface exists to eliminate the "JavaScript Tax" on heavy operations:
1. **Parallelism**: While Node is single-threaded, Surface uses Go's scheduler to handle I/O, Network, and Image tasks across all available CPU cores.
2. **Zero-Reflect JSON**: Bypass standard reflection for manual byte-buffer JSON assembly.
3. **Background Pre-fetching**: Reads CSV data in parallel goroutines so JS never waits for disk I/O.
4. **Connection Pooling**: Reuses persistent TCP/TLS tunnels for bulk SMTP operations.

---

## 📖 API Reference

### 🖼️ Image Module (`image`)
High-performance native image processing using Lanczos3 interpolation. Supports **Local Paths**, **Remote URLs**, and **Zero-Disk Base64 Output**.

#### `image.resize(opts)`
Resizes an image natively. If `width` or `height` is 0, aspect ratio is preserved.
```javascript
// Local Resize
image.resize({
  src: "./input.jpg",
  dist: "./output.jpg",
  width: 800,
  quality: 90
});

// Remote Zero-Disk Resize (Returns Base64)
const { base64 } = image.resize({
  src: "https://example.com/photo.png",
  width: 300
});
```

#### `image.crop(opts)`
Fills the dimensions and crops from the center (Perfect for uniform square thumbnails).
```javascript
image.crop({
  src: "user_upload.jpg",
  width: 500,
  height: 500
});
```

---

### 📊 CSV Module (`csv`)
Ultra-fast, stateful CSV engine with pre-fetching.

#### `csv.open(path, options)`
Opens a CSV and starts background pre-fetching.
```javascript
const h = csv.open("./data.csv", {
  header: true,
  mode: "object",
  select: ["id", "email"]
});
```

#### `csv.readAll(handle)`
Moves the entire loop into Go for zero JS overhead.
```javascript
const records = csv.readAll(h);
```

---

### 📧 SMTP Module (`smtp`)
Enterprise-grade delivery with native rendering and connection pooling.

#### `smtp.bulk(options)`
Parallel delivery using native goroutines.
```javascript
smtp.bulk({
  ...creds,
  emails: [{ to: "user1@abc.com", body: "Msg 1" }],
  concurrency: 10
});
```

#### `smtp.renderFile(path, data)`
Renders Go templates directly from disk.
```javascript
const body = smtp.renderFile("./welcome.tmpl", { name: "Soham" });
```

---

## 🌍 Real-World Use Cases

### 1. Zero-Disk User Profile On-boarding
Process user uploads without filling your server's disk with temporary garbage files.
```javascript
import { image } from "@titanpl/surface";

export function onProfileUpload(req) {
  // Use a Pinterest or S3 URL as source
  const result = image.crop({
    src: req.body.imageUrl,
    width: 200,
    height: 200
  });

  // Save base64 string directly to User profile in DB
  return db.users.update(req.userId, { 
    avatar: result.base64 
  });
}
```

### 2. High-Performance Marketing Pipeline
Combine CSV, Templating, and SMTP for mass personalization.
```javascript
import { csv, smtp, image } from "@titanpl/surface";

export function sendCampaign() {
  const h = csv.open("leads.csv", { mode: "object" });
  const leads = csv.readAll(h);
  csv.close(h);

  const jobs = leads.map(lead => ({
    to: lead.email,
    body: smtp.renderFile("promo.tmpl", { name: lead.name })
  }));

  return smtp.bulk({ ...creds, emails: jobs, concurrency: 20 });
}
```

### 3. Automated Social Media Card Generator
Scrape remote images and crop them natively for social previews.
```javascript
export function generateSocialPreview(req) {
  return image.crop({
    src: "https://images.unsplash.com/photo-12345",
    dist: "./public/previews/card_1.jpg",
    width: 1200,
    height: 630
  });
}
```

---
Built with ❤️ by the TitanPL Team.
