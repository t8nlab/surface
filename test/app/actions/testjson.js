import { json } from "@titanpl/surface";
import { path } from "@titanpl/native";

export default function testjson(req) {
  const steps = [];
  try {
    const stdJsonPath = path.resolve("../data.json");
    const bridgeCsvPath = path.resolve("../bridged.csv");
    steps.push(`Paths resolved: ${stdJsonPath}, ${bridgeCsvPath}`);
    
    // ✅ 1. CREATE STANDARD NESTED JSON
    steps.push("Creating standard JSON...");
    const wh = json.create(stdJsonPath);
    for (let i = 1; i <= 50; i++) {
      json.write(wh, { id: i, name: `User ${i}`, email: `user${i}@titan.pl` });
    }
    json.close(wh);
    steps.push("Standard JSON created and closed");

    // ✅ 2. PATH-BASED EXTRACTION
    steps.push("Opening standard JSON for streaming...");
    const rh = json.open(stdJsonPath, { format: "json" });
    const chunk = json.next(rh, { size: 5 });
    steps.push(`Fetched ${chunk.rows.length} rows successfully`);
    json.close(rh);

    // ✅ 3. NATIVE BRIDGE (JSON -> CSV)
    steps.push("Executing Native Bridge: JSON to CSV...");
    const bridgeRes = json.toCSV(stdJsonPath, bridgeCsvPath);
    steps.push(`Bridge success: ${bridgeRes.success}`);

    return {
      success: true,
      steps: steps,
      sample_row: chunk.rows[0],
      bridge_path: bridgeCsvPath
    };
  } catch (err) {
    return { 
      success: false, 
      error: err.message, 
      steps: steps 
    };
  }
}
