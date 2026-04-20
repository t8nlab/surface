# ⏣ Titan Surface
**The Native Capability Layer for TitanPL**

Surface is an ultra-optimized native extension provider for the TitanPl framework, built entirely in **Go**. It handles data-heavy, IO-bound, and system-level tasks outside the JavaScript runtime to ensure sub-millisecond response times and rock-solid stability.

---

## 🏎️ Core Architecture (The Go Advantage)

Surface eliminates the "JavaScript Tax" by moving heavy workloads to the native layer.

- **Non-Blocking I/O**: While Node.js waits for disk/network interrupts, Surface uses Go routines (lightweight threads) to pre-process data in parallel.
- **Zero-Copy Buffering**: We use `bytes.Buffer` and manual JSON builders in Go to avoid the memory fragmentation common in garbage-collected environments.
- **Atomic Operations**: We chain multiple logical steps (like Image Resize + Crop + Format) in one execution cycle, keeping data in the CPU's L2 cache.

---

## 📊 Module: CSV Streaming Engine
The fastest way to process massive tabular datasets in the Titan ecosystem.

### Under the Hood: The Pre-fetcher
When you call `csv.open`, Surface spawns a background "Look-ahead" thread. While your JavaScript is processing row #1, the Go thread is already parsing and buffering rows #1,001 to #2,000. When you call `csv.next`, the data is delivered instantly from RAM.

### Functions & Examples

#### `csv.open(path, options)`
Initializes the streamer and background pre-fetcher.
```javascript
const h = csv.open("./data.csv", {
  header: true,      // Auto-extract column names
  mode: "object",    // Output as { key: value }
  inferTypes: true,  // Convert "true" -> true, "123" -> 123
  select: ["id", "name"] // High Performance: only parse what you need
});
```

#### `csv.next(handle, options)`
Streaming access to the native buffer.
```javascript
const chunk = csv.next(h, { size: 500 });
console.log(chunk.rows); // Array of 500 pre-fetched records
```

#### `csv.readAll(handle)`
Flushes the entire file through Go into JS in one pass.
```javascript
const allUsers = csv.readAll(h); // The absolute fastest way to read
```

#### `csv.create(path, options)`
Creates a new CSV with native buffering.
```javascript
const wh = csv.create("./export.csv", { headers: ["name", "score"] });
```

#### `csv.write(handle, data)`
Writes records with near-zero latency.
```javascript
csv.write(wh, { name: "Soham", score: 100 });
csv.write(wh, [{ name: "John", score: 80 }, { name: "Jane", score: 95 }]);
```

#### `csv.close(handle)`
Kills background threads and releases OS file locks.
```javascript
try { ... } finally { csv.close(h); }
```

### 💡 Industrial Use Case: Data Warehouse Ingestion
Syncing 10 million records from an old CSV into a modern DB.
> **Problem**: Standard Node `fs` and `fast-csv` bloat the memory and lock the event loop for minutes.
> **Solution**: Use `csv.open` with a 10,000-row `next()` loop. Surface handles the 1GB file reading in the background while your JS keeps the server responsive.

---

## 📧 Module: SMTP Communication Engine
Enterprise-grade delivery with persistent connection pooling.

### Under the Hood: The Pooling Logic
Unlike standard libraries that "Login -> Send -> Logout" for every email, `smtp.bulk` maintains a **Persistent TLS Tunnel**. Workers stay logged in and stream emails through an open socket, bypassing Gmail/Exchange security rate-limits.

### Functions & Examples

#### `smtp.send(options)`
Delivers a single email. Supports implicit SSL (465) and STARTTLS (587).
```javascript
smtp.send({
  host: "smtp.gmail.com",
  username: "...", password: "...",
  to: "client@test.com",
  subject: "Order #123",
  body: "<h1>Thanks!</h1>"
});
```

#### `smtp.bulk(options)`
Uses a native worker pool to send hundreds of emails concurrently.
```javascript
smtp.bulk({
  ...creds,
  emails: jobs,
  concurrency: 10 // Maximize throughput with 10 parallel connections
});
```

#### `smtp.render(template, data)` / `renderFile(path, data)`
Native Go `html/template` engine. Fast and XSS-safe.
```javascript
const body = smtp.renderFile("./otp.tmpl", { code: "9182" });
```

### 💡 Industrial Use Case: High-Speed OTP & SaaS Notifications
Delivering millions of per-user notifications.
> **Problem**: Node-based templating is slow for 100k renders.
> **Solution**: Move rendering to Go. `smtp.renderFile` executes in microseconds, and `smtp.bulk` blasts them via parallel tunnels. No lag, no connection errors.

---

## 🖼️ Module: Atomic Image Processing
The ultimate media engine for modern web apps.

### Under the Hood: The Atomic Pipeline
Every time you resize/crop an image, it must be decoded into pixels and then re-encoded into a format (JPG/PNG). `image.process` allows you to chain 5 actions while the image is **already decoded**, skipping 4 redundant encode/decode cycles.

### Functions & Examples

#### `image.process(options)`
Multi-step atomic processing.
```javascript
const result = image.process({
  src: "input.png",
  out: "thumb.webp", // Optional: returns base64 if omitted
  format: "webp",
  steps: [
    { action: "resize", width: 800 },
    { action: "grayscale" },
    { action: "crop", width: 300, height: 300 }
  ]
});
```

#### `image.batch(options)`
Parallel media generator via a shared native worker pool.
```javascript
image.batch({
  concurrency: 8,
  items: [
    { src: "a.jpg", out: "a_sm.jpg", width: 100 },
    { src: "b.jpg", out: "b_sm.jpg", width: 100 }
  ]
});
```

### 💡 Industrial Use Case: Cloud-Scale Image Optimization
Generating 5 different social media sizes for every user upload.
> **Problem**: Uploading a 5MB image and generating 5 thumbnails in Node.js can take 5+ seconds and 200MB of RAM.
> **Solution**: Direct URL streaming in Surface. Go fetches the 5MB file, renders all variations in parallel via goroutines, and returns optimized WebP strings in <800ms with constant RAM usage.

---
Built with ❤️ by the TitanPL Team.
