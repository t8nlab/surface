import { ext } from "../index.js";

/**
 * Native Image Resizing
 */
export function imageResize(opts) {
  return ext.call("image_resize", opts);
}

/**
 * Native Smart Cropping (Fill and Crop from center)
 */
export function imageCrop(opts) {
  return ext.call("image_crop", opts);
}
