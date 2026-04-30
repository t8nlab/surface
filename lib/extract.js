import { ext } from "../ext.js";

/**
 * Extract raw HTML from a URL
 * @param {string} url 
 */
export function extractHtml(url) {
  return ext.call("extract_html", { url });
}

/**
 * Extract all links from a URL
 * @param {string} url 
 */
export function extractLinks(url) {
  return ext.call("extract_links", { url });
}

/**
 * Extract SEO/OG Metadata from a URL
 * @param {string} url 
 */
export function extractMeta(url) {
  return ext.call("extract_meta", { url });
}
