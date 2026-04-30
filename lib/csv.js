import { ext } from "../ext.js";

/**
 * Open CSV for reading
 * @returns {string} handler
 */
export function csvOpen(path, opts = {}) {
  const res = ext.call("csv_open", { path, ...opts });
  return res.handler;
}

/**
 * Read next chunk
 * @returns {{rows: any[], done: boolean}}
 */
export function csvNext(handler, opts = {}) {
  return ext.call("csv_next", { handler, ...opts });
}

/**
 * Read all data from CSV
 */
export function csvReadAll(handler) {
  return ext.call("csv_read_all", { handler });
}

/**
 * Create CSV for writing
 * @returns {string} handler
 */
export function csvCreate(path, opts = {}) {
  const res = ext.call("csv_create", { path, ...opts });
  return res.handler;
}

/**
 * Write chunk
 */
export function csvWrite(handler, rows) {
  const data = Array.isArray(rows) ? rows : [rows];
  return ext.call("csv_write", { handler, rows: data });
}

/**
 * Close handler
 */
export function csvClose(handler) {
  return ext.call("csv_close", { handler });
}