import { ext } from "../ext.js";

/**
 * Perform an HTTP request with full control (Axios-like)
 * @param {string|Object} urlOrConfig 
 * @param {Object} [config] 
 * @returns {Promise<{status: number, statusText: string, headers: Record<string, string>, data: any, url: string, ok: boolean}>}
 */
export function request(urlOrConfig, config = {}) {
  let finalConfig = {};
  if (typeof urlOrConfig === "string") {
    finalConfig = { url: urlOrConfig, ...config };
  } else {
    finalConfig = urlOrConfig;
  }
  
  return ext.call("http_request", finalConfig);
}

/**
 * Shorthand for GET request
 * @param {string} url 
 * @param {Object} [config] 
 */
export function get(url, config = {}) {
  return request({ ...config, url, method: "GET" });
}

/**
 * Shorthand for POST request
 * @param {string} url 
 * @param {any} [data] 
 * @param {Object} [config] 
 */
export function post(url, data, config = {}) {
  return request({ ...config, url, data, method: "POST" });
}

/**
 * Shorthand for PUT request
 * @param {string} url 
 * @param {any} [data] 
 * @param {Object} [config] 
 */
export function put(url, data, config = {}) {
  return request({ ...config, url, data, method: "PUT" });
}

/**
 * Shorthand for DELETE request
 * @param {string} url 
 * @param {Object} [config] 
 */
export function del(url, config = {}) {
  return request({ ...config, url, method: "DELETE" });
}

/**
 * Shorthand for PATCH request
 * @param {string} url 
 * @param {any} [data] 
 * @param {Object} [config] 
 */
export function patch(url, data, config = {}) {
  return request({ ...config, url, data, method: "PATCH" });
}

const http = {
  request,
  get,
  post,
  put,
  delete: del,
  patch,
};

export default http;
