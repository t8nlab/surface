import { csvOpen, csvNext, csvReadAll, csvCreate, csvWrite, csvClose } from "./lib/csv.js";
import { smtpSend, smtpBulkSend, smtpRender, smtpRenderFile } from "./lib/smtp.js";
import { imageResize, imageCrop, imageProcess, imageBatch } from "./lib/image.js";
import { jsonOpen, jsonNext, jsonClose, jsonCreate, jsonWrite, jsonStringify, jsonToCsv, jsonReadAll } from "./lib/json.js";
import { cleanValidateEmails, cleanNormalizePhones, cleanRemoveDuplicates, cleanProcess } from "./lib/clean.js";
import { extractHtml, extractLinks, extractMeta } from "./lib/extract.js";
import httpSrf from "./lib/http.js";


// Compatibility layer
export const csv = {
  open: csvOpen,
  next: csvNext,
  readAll: csvReadAll,
  create: csvCreate,
  write: csvWrite,
  close: csvClose,
};

export const smtp = {
  send: smtpSend,
  bulk: smtpBulkSend,
  render: smtpRender,
  renderFile: smtpRenderFile,
};

export const image = {
  resize: imageResize,
  crop: imageCrop,
  process: imageProcess,
  batch: imageBatch,
};

export const json = {
  open: jsonOpen,
  next: jsonNext,
  readAll: jsonReadAll,
  close: jsonClose,
  create: jsonCreate,
  write: jsonWrite,
  stringify: jsonStringify,
  toCSV: jsonToCsv,
};

export const clean = {
  validateEmails: cleanValidateEmails,
  normalizePhones: cleanNormalizePhones,
  removeDuplicates: cleanRemoveDuplicates,
  process: cleanProcess,
};

export const extract = {
  html: extractHtml,
  links: extractLinks,
  meta: extractMeta,
};

export const http = httpSrf;

export default { csv, smtp, image, json, clean, extract, http };