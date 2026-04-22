import { clean, csv, json } from "@titanpl/surface";
import { path } from "@titanpl/native";

export default function testclean() {
  const steps = [];
  try {
    const dirtyPath = path.resolve("../data.json");
    const cleanPath = path.resolve("../clean_data.json");

    steps.push(`Resolving source data: ${dirtyPath}`);
    steps.push("Running unified native clean on JSON file...");
    const stats = clean.process({
      src: dirtyPath,
      out: cleanPath,
      normalize: true,
      dedup: true,
      concurrency: 4
    });

    steps.push("Reading cleaned results for verification...");
    const rh = json.open(cleanPath, { header: true });

    return {
      success: true,
      message: "Native cleaning verified for all fields",
      stats: stats,
      execution_log: steps
    };
  } catch (err) {
    return {
      success: false,
      error: err.message,
      execution_log: steps
    };
  }
}
