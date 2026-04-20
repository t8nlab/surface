import { createExt } from "./utils/native";
import { csvOpen, csvNext, csvReadAll, csvCreate, csvWrite, csvClose } from "./lib/csv.js";
import { smtpSend, smtpBulkSend, smtpRender, smtpRenderFile } from "./lib/smtp.js";
import { imageResize, imageCrop, imageProcess, imageBatch } from "./lib/image.js";
import { jsonOpen, jsonNext, jsonClose, jsonCreate, jsonWrite, jsonStringify, jsonToCsv } from "./lib/json.js";

export const ext = createExt("@titanpl/surface");

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
  close: jsonClose,
  create: jsonCreate,
  write: jsonWrite,
  stringify: jsonStringify,
  toCSV: jsonToCsv,
};

export default { csv, smtp, image, json };