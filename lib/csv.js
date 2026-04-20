import { ext } from "../index.js";


/**
 * Open CSV for reading
 * @returns {string} handle
 */
export function csvOpen(path, opts = {}) {
  const res = ext.call("csv_open", { path, ...opts });
  return res.handle;
}

/**
 * Read next chunk
 * @returns {{rows: any[], done: boolean}}
 */
export function csvNext(handle, opts = {}) {
  return ext.call("csv_next", { handle, ...opts });
}

/**
 * Read all data from CSV
 */
export function csvReadAll(handle) {
  return ext.call("csv_read_all", { handle });
}

/**
 * Create CSV for writing
 * @returns {string} handle
 */
export function csvCreate(path, opts = {}) {
  const res = ext.call("csv_create", { path, ...opts });
  return res.handle;
}

/**
 * Write chunk
 */
export function csvWrite(handle, rows) {
  const data = Array.isArray(rows) ? rows : [rows];
  return ext.call("csv_write", { handle, rows: data });
}

/**
 * Close handle
 */
export function csvClose(handle) {
  return ext.call("csv_close", { handle });
}