import { ext } from "../index.js";

/**
 * Execute multiple image operations in a single pass (Pipeline)
 */
export function imageProcess(opts) {
  return ext.call("image_process", opts);
}

/**
 * Process multiple images in parallel using native worker pool
 */
export function imageBatch(opts) {
  return ext.call("image_batch", opts);
}

/**
 * Native Image Resizing (Legacy)
 */
export function imageResize(opts) {
  return ext.call("image_resize", opts);
}

/**
 * Native Smart Cropping (Legacy)
 */
export function imageCrop(opts) {
  return ext.call("image_crop", opts);
}
