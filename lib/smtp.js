import { ext } from "../index.js";


/**
 * Send email
 */
export function smtpSend(opts) {
  return ext.call("smtp_send", opts);
}

/**
 * Send multiple emails concurrently using native worker pool
 */
export function smtpBulkSend(opts) {
  return ext.call("smtp_bulk_send", opts);
}

/**
 * Render HTML template natively using Go tpl engine
 */
export function smtpRender(template, data = {}) {
  return ext.call("smtp_render", { template, data });
}

/**
 * Render HTML template file natively from disk
 */
export function smtpRenderFile(path, data = {}) {
  return ext.call("smtp_render_file", { path, data });
}