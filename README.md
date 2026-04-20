# ⏣ TitanPL Surface
### The Definitive Native Function Catalog

Surface is an ultra-optimized high level modules provider native extension for the TitanPl framework, built entirely in **Go**. It handles data-heavy, IO-bound, and system-level tasks outside the JavaScript runtime to ensure sub-millisecond response times and rock-solid stability.


---

## 📊 CSV Module (`csv`)

### `csv.open(path, opts)`
Opens a CSV file for streaming.
```javascript
// Opens file with automatic type inference
const handler = csv.open("./data.csv", { header: true, inferTypes: true });
```

### `csv.next(handler, opts)`
Fetches the next chunk.
```javascript
const chunk = csv.next(handler, { size: 500 });
// Returns: { rows: [...], done: false }
```

### `csv.readAll(handler)`
Memory-intensive but ultra-fast dump.
```javascript
const allData = csv.readAll(handler);
```

### `csv.create(path, opts)`
Creates a new CSV with fixed headers.
```javascript
const wh = csv.create("./export.csv", { headers: ["id", "status"] });
```

### `csv.write(handler, rows)`
Buffered native writing.
```javascript
csv.write(wh, [{ id: 1, status: "active" }, { id: 2, status: "pending" }]);
```

### `csv.close(handler)`
Releases file descriptors and flushes memory.
```javascript
csv.close(wh);
```

---

## 💎 JSON Module (`json`)

### `json.open(path, opts)`
Native token-walking reader.
```javascript
// Path walking: skips bytes to reach 'logs.errors' natively
const rh = json.open("./app.json", { fpath: "logs.errors[*]" });
```

### `json.next(handler, opts)`
Native streaming for massive JSON files.
```javascript
const logs = json.next(rh, { size: 10 });
```

### `json.create(path, opts)`
Stream writer for JSON or JSONL.
```javascript
const wh = json.create("./logs.jsonl", { format: "jsonl" });
```

### `json.write(handler, data)`
Serialized write (Standard `[...]` or `JSONL`).
```javascript
json.write(wh, { event: "click", time: Date.now() });
```

### `json.stringify(data)`
Native Go-based serialization (faster than JS).
```javascript
const nativeStr = json.stringify({ complex: "object" });
```

### `json.toCSV(jsonPath, csvPath)`
The Native Bridge. Streams records directly from JSON to CSV.
```javascript
json.toCSV("./data.json", "./static/data.csv");
```

### `json.close(handler)`
Safety cleanup.
```javascript
json.close(rh);
```

---

## 📧 SMTP Module (`smtp`)

### `smtp.send(config, email)`
Immediate TLS delivery.
```javascript
await smtp.send(config, {
  to: "user@titan.pl",
  subject: "OTP",
  body: "Code: 1234"
});
```

### `smtp.bulk(config, emails)`
Parallel worker-pool delivery.
```javascript
// Send 10,000 emails concurrently
smtp.bulk(config, [{ to: "a@a.com", ... }, { to: "b@b.com", ... }]);
```

### `smtp.render(html, data)`
XSS-Safe Go Template rendering.
```javascript
const body = smtp.render("<h1>Hello {{.name}}</h1>", { name: "Titan" });
```

### `smtp.renderFile(path, data)`
File-based templating.
```javascript
const body = smtp.renderFile("./templates/welcome.html", { name: "User" });
```

---

## 🖼️ Image Module (`image`)

### `image.resize(inputPath, opts)`
High-performance single resize.
```javascript
image.resize("./in.jpg", { out: "./out.webp", width: 300 });
```

### `image.crop(inputPath, opts)`
Coordinate or Center cropping.
```javascript
image.crop("./in.jpg", { out: "./cropped.jpg", x: 0, y: 0, w: 100, h: 100 });
```

### `image.process(opts)`
Atomic Pipeline (One pass).
```javascript
image.process({
  input: "./in.png",
  out: "./optimized.webp",
  steps: [
    { resize: { width: 500 } },
    { crop: { w: 200, h: 200 } }
  ]
});
```

### `image.batch(list)`
Massive Parallel processing.
```javascript
image.batch([
  { input: "1.png", out: "1.webp", resize: { width: 50 } },
  { input: "2.png", out: "2.webp", resize: { width: 50 } }
]);
```

---

## 🌍 Real-World Architecture

### Industrial Use Case: AI Data Pipeline
1.  **Extract**: `json.open` (using `fpath`) to pull key data from nested AI response.
2.  **Transform**: Process records via JS logic.
3.  **Load**: `json.toCSV` to bridge the final cleaned data into a spreadsheet for the user.

### Industrial Use Case: Marketing Engine
1.  **Segmentation**: Use `csv.open` to read subscriber list.
2.  **Rendering**: Use `smtp.render` to create personalized emails.
3.  **Dispatch**: Use `smtp.bulk` for concurrent SMTP tunneling.

---

(c) 2026 Titan Planet Team. Speed is the only feature.
