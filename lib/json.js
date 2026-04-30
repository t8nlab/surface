import { ext } from "../ext.js";

/**
 * Open JSON/JSONL for Native Streaming
 */
export function jsonOpen(path, opts = {}) {
  const res = ext.call("json_open", { path, ...opts });
  return res.handler;
}

/**
 * Fetch next chunk of JSON records
 */
export function jsonNext(handler, opts = {}) {
  return ext.call("json_next", { handler, ...opts });
}

/**
 * Read all data from JSON stream
 */
export function jsonReadAll(handler) {
  return ext.call("json_read_all", { handler });
}

/**
 * Close JSON stream
 */
export function jsonClose(handler) {
  const res = ext.call("json_close", { handler });
  return res.success;
}

/**
 * Create JSON/JSONL for writing
 */
export function jsonCreate(path, opts = {}) {
  const res = ext.call("json_create", { path, ...opts });
  return res.handler;
}

/**
 * Write record to stream
 */
export function jsonWrite(handler, data, opts = {}) {
  const res = ext.call("json_write", { handler, data, ...opts });
  return res.success;
}

/**
 * Fast JSON Stringify (Native)
 */
export function jsonStringify(data) {
  const res = ext.call("json_stringify", { data });
  return res.json;
}

/**
 * Native Cross-Engine Bridge: JSON -> CSV
 * Streams JSON records into a CSV file with zero JS overhead.
 */
export function jsonToCsv(jsonPath, csvPath, opts = {}) {
  return ext.call("json_to_csv", { path: jsonPath, out: csvPath, ...opts });
}
