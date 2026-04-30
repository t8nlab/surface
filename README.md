# ⏣ TitanPL Surface
### The Definitive Native Function Catalog

Surface is an ultra-optimized high level modules provider native extension for the TitanPl framework, built entirely in **Go**. It handles data-heavy, IO-bound, and system-level tasks outside the JavaScript runtime to ensure sub-millisecond response times and rock-solid stability.


---

## 📊 CSV Module (`csv`)

### `csv.open(path, opts)`
Opens a CSV file for streaming. Supports **local file paths** and **public Cloud URLs**.
```javascript
// Local file with automatic type inference
const handler = csv.open("./data.csv", { header: true, inferTypes: true });

// Cloud URL (Zero-copy stream directly from the internet)
const cloud = csv.open("https://example.com/data.csv", { header: true });
```

### `csv.next(handler, opts)`
Fetches the next chunk of records from the native pre-fetch buffer.
```javascript
const chunk = csv.next(handler, { size: 500 });
// Returns: { rows: [...], done: false, mode: "object" }
```

### `csv.readAll(handler)`
Ultra-fast native dump. Loads the entire remaining contents of the CSV into memory in one go.
```javascript
const allData = csv.readAll(handler);
```

### `csv.create(path, opts)`
Creates or overwrites a CSV file for writing with fixed headers.
```javascript
const wh = csv.create("./export.csv", { headers: ["id", "name", "status"], delimiter: "," });
```

### `csv.write(handler, rows)`
Buffered native writing. Accepts an array of objects or an array of arrays.
```javascript
csv.write(wh, [
  { id: 1, name: "Titan", status: "active" }, 
  { id: 2, name: "Planet", status: "pending" }
]);
```

### `csv.close(handler)`
Releases file descriptors and flushes any pending native memory.
```javascript
csv.close(wh);
```

---

## 💎 JSON Module (`json`)

### `json.open(path, opts)`
Native token-walking reader. Supports **local files** and **public URLs**.
```javascript
// Path walking: skips bytes to reach 'logs.errors' natively before streaming
const rh = json.open("./app.json", { fpath: "logs.errors[*]" });
```

### `json.next(handler, opts)`
Native streaming for massive JSON or JSONL files.
```javascript
const logs = json.next(rh, { size: 10 });
```

### `json.readAll(handler)`
Loads the entire remaining JSON stream into a single JavaScript array.
```javascript
const allRecords = json.readAll(rh);
```

### `json.create(path, opts)`
Creates a stream writer for JSON or JSONL output.
```javascript
const wh = json.create("./logs.jsonl", { format: "jsonl" });
```

### `json.write(handler, data, opts)`
Serialized native write. Handles standard JSON array structure or line-delimited (JSONL).
```javascript
json.write(wh, { event: "click", time: Date.now() }, { format: "jsonl" });
```

### `json.stringify(data)`
Ultra-fast native Go-based serialization (significantly faster than `JSON.stringify` for large objects).
```javascript
const nativeStr = json.stringify({ complex: "object", data: [1, 2, 3] });
```

### `json.toCSV(jsonPath, csvPath, opts)`
The Native Bridge. Streams records directly from JSON to CSV natively.
```javascript
json.toCSV("./data.json", "./static/data.csv", { fpath: "items[*]" });
```

### `json.close(handler)`
Closes the stream and releases resources.
```javascript
json.close(rh);
```

---

## 📧 SMTP Module (`smtp`)

### `smtp.send(opts)`
Immediate native TLS/SSL delivery via Go SMTP engine.
```javascript
smtp.send({
  host: "smtp.gmail.com",
  port: 587,
  username: "user@gmail.com",
  password: "app-password",
  to: "client@titan.pl",
  from: "system@titan.pl",
  subject: "OTP Verification",
  body: "Your code is: 4829"
});
```

### `smtp.bulk(opts)`
Massive parallel worker-pool delivery. Uses native Go concurrency to send thousands of emails simultaneously.
```javascript
smtp.bulk({
  host: "smtp.titan.pl",
  port: 587,
  username: "system",
  password: "...",
  emails: [{ to: "a@a.com", body: "Msg 1" }, { to: "b@b.com", body: "Msg 2" }],
  concurrency: 10 // Opens 10 parallel SMTP tunnels
});
```

### `smtp.render(template, data)`
Renders a Go HTML/Text template natively. XSS-Safe and ultra-fast.
```javascript
const html = smtp.render("<h1>Hello {{.Name}}</h1>", { Name: "Titan Planet" });
```

### `smtp.renderFile(path, data)`
Reads and renders a Go template file directly from disk.
```javascript
const body = smtp.renderFile("./templates/welcome.html", { User: "Ezet" });
```

---

## 🖼️ Image Module (`image`)

### `image.resize(opts)`
High-performance native resizing. Supports JPG, PNG, and WebP.
```javascript
image.resize({ 
  src: "./photo.jpg", 
  out: "./thumb.webp", 
  width: 300, 
  quality: 80,
  format: "webp" 
});
```

### `image.crop(opts)`
Smart center cropping or coordinate-based cropping.
```javascript
image.crop({ 
  src: "./in.jpg", 
  out: "./square.jpg", 
  width: 400, 
  height: 400 
});
```

### `image.process(opts)`
Atomic Pipeline. Perform multiple operations (Resize, Crop, Blur, Grayscale) in **one single native pass**.
```javascript
image.process({
  src: "https://site.com/large.jpg",
  out: "./optimized.webp",
  steps: [
    { action: "resize", width: 800 },
    { action: "grayscale" },
    { action: "blur", sigma: 0.5 },
    { action: "crop", width: 400, height: 400 }
  ]
});
```

### `image.batch(opts)`
Massive Parallel processing using a native worker pool.
```javascript
image.batch({
  concurrency: 4,
  items: [
    { src: "1.png", out: "1.webp", width: 50 },
    { src: "2.png", out: "2.webp", width: 50 }
  ]
});
```

## 🧹 Data Cleaning Module (`clean`)

### `clean.validateEmails(path)`
Natively validates emails in a massive file using high-performance Go regex engines.
```javascript
const stats = clean.validateEmails("./users_export.csv");
// Returns: { valid: 920, invalid: 80 }
```

### `clean.normalizePhones(phones)`
Normalizes an array of phone numbers to E.164 format natively.
```javascript
const cleanPhones = clean.normalizePhones(["(555) 123-4567", "1-555-987-6543"]);
// Returns: ["+5551234567", "+15559876543"]
```

### `clean.removeDuplicates(src, out)`
Natively removes duplicate rows from a file. Extremely fast for datasets with millions of rows.
```javascript
const res = clean.removeDuplicates("./dirty.csv", "./cleaned.csv");
// Returns: { processed: 1000, duplicates: 34, saved: 966 }
```

### `clean.process(opts)`
The flagship cleaning engine. Processes millions of rows in a single native pass using a **Parallel Worker Pool** (powered by Go routines).
- Performs **Deep Normalization**: Trims all fields and formats phones to E.164.
- Performs **Native Deduplication**: High-speed row comparison using thread-safe maps.
```javascript
const stats = clean.process({
  src: "./dirty_data.csv",
  out: "./clean_data.csv",
  normalize: true, 
  dedup: true  ,
  concurrency: 4,
});
// Returns: { processed: 1000000, duplicates: 342, workers: 16, success: true }
```

### `clean.validateEmails(path)`
Natively validates email syntax across an entire file.
```javascript
const stats = clean.validateEmails("./leads.csv");
```
---

---

## 🔗 Web Extraction Module (`extract`)

### `extract.html(url)`
Fetches raw HTML from a public URL natively using Go's optimized HTTP stack.
```javascript
const rawHtml = extract.html("https://google.com");
```

### `extract.links(url)`
Extracts all unique unique links (`href` attributes) from a URL natively.
```javascript
const links = extract.links("https://titanpl.vercel.app");
```

### `extract.meta(url)`
Extracts SEO and OpenGraph metadata natively from any public URL.
```javascript
const seo = extract.meta("https://github.com");
// Returns: { "og:title": "GitHub", "description": "...", ... }
```

---

## 🌐 HTTP Module (`http`)

### `http.get(url, opts)`
Native high-performance GET request with **Axios-like API**. Automatically handles query parameters and headers.
```javascript
const res = http.get("https://api.example.com/data", {
  params: { limit: 10 },
  headers: { "Authorization": "Bearer ..." }
});

console.log(res.data); // Automatically parsed if JSON
```

### `http.post(url, data, opts)`
Sends a POST request. **Auto-serializes** JavaScript objects to JSON and sets appropriate headers.
```javascript
const res = http.post("https://api.example.com/users", { 
  name: "Titan", 
  role: "Admin" 
});
```

### `http.request(config)`
Generic request method for full control over methods (PUT, DELETE, PATCH) and advanced settings.
```javascript
const res = http.request({
  method: "PUT",
  url: "https://api.example.com/update/1",
  data: { status: "active" },
  timeout: 5000 // 5 second timeout
});
```

---

## 🚀 Pro Examples (Industrial Workflows)

### 1. Cloud Streaming (Zero Disk Overhead)
Stream datasets directly from the internet without saving to server disk.
```javascript
import { csv } from "@titanpl/surface";

export default function cloud_action() {
  const url = "https://raw.githubusercontent.com/.../data.csv";
  const handler = csv.open(url, { header: true });
  
  // Stream via chunked loop (Best for GB+ files)
  while (true) {
    const chunk = csv.next(handler, { size: 1000 });
    if (chunk.done) break;
  }
  
  csv.close(handler);
}
```

### 2. High-Speed Bulk OTP Delivery
Combine CSV streaming and SMTP bulk delivery for extreme performance.
```javascript
import { csv, smtp } from "@titanpl/surface";

export default function bulk_otp() {
  const handler = csv.open("./subscribers.csv", { mode: "object" });
  const emails = csv.readAll(handler);
  csv.close(handler);

  // Surface blasts 10 emails simultaneously via Go worker pool
  return smtp.bulk({
    host: "smtp.titan.pl",
    port: 587,
    username: "system",
    password: "...",
    emails: emails.map(row => ({
      to: row.email,
      subject: "Login Code",
      body: `Code: ${Math.random().toString().slice(2, 6)}`
    })),
    concurrency: 10
  });
}
```

## 🌍 Architecture

### Industrial Use Case: Lead Generation & Data Cleaning
1.  **Extract**: `extract.links(url)` to natively crawl a target site for possible leads.
2.  **Verify**: `clean.validateEmails(file)` to ensure all extracted contact data is deliverable.
3.  **Deliver**: `smtp.bulk(opts)` to send ultra-fast verified notifications.

---

(c) 2026 Titan Planet Team.
