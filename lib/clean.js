import { ext } from "../ext.js";

/**
 * Validate Emails in a file natively
 * @param {string} path - Path to file
 */
export function cleanValidateEmails(path) {
  return ext.call("clean_validate_emails", { path });
}

/**
 * Normalize phone numbers
 * @param {string[]|object} input - Array of phones OR object with {path, out}
 */
export function cleanNormalizePhones(input) {
  if (Array.isArray(input)) {
    return ext.call("clean_normalize_phones", { phones: input });
  }
  return ext.call("clean_normalize_phones", input);
}

/**
 * Remove duplicate rows natively
 * @param {string} src - Source file path
 * @param {string} out - Destination file path
 */
export function cleanRemoveDuplicates(src, out) {
  return ext.call("clean_remove_duplicates", { src, out });
}

/**
 * Perform multiple cleaning operations in one native pass
 * @param {object} opts - { src, out, normalize, dedup }
 */
export function cleanProcess(opts) {
  return ext.call("clean_process", opts);
}
