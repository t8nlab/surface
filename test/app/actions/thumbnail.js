import { image } from "@titanpl/surface";
import { path } from "@titanpl/native";

export default function thumbnail(req) {
  try {
    const src = "https://i.pinimg.com/736x/2c/c2/fe/2cc2fe16eed28daf889d3fe5eff629c3.jpg";
    
    // ✅ 1. PIPELINE MODE: Parallel operations in ONE native pass
    const pipelineResult = image.process({
      src: src,
      out: path.resolve("app/complex_thumb.webp"),
      format: "webp",
      quality: 80,
      steps: [
        { action: "resize", width: 800 },
        { action: "grayscale" },
        { action: "blur", sigma: 0.5 },
        { action: "crop", width: 400, height: 400 }
      ]
    });

    // ✅ 2. BATCH MODE: Parallel worker processing for multiple files
    const batchResult = image.batch({
      concurrency: 4,
      items: [
        { src: src, out: path.resolve("app/batch_1.jpg"), width: 100 },
        { src: src, out: path.resolve("app/batch_2.png"), width: 300, format: "png" },
        { src: src, width: 200 }
      ]
    });

    return {
      success: true,
      pipeline: pipelineResult,
      batch: batchResult
    };
  } catch (err) {
    return { success: false, error: err.message };
  }
}
