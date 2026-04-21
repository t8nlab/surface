import { csv, json } from "@titanpl/surface";

/**
 * Action: Cloud Stream Test (Full Data)
 * Demonstrates how to fetch ALL rows from a cloud source using both
 * the one-shot 'readAll' method and the chunked 'next' loop.
 */
export default function cloud_stream() {
  const steps = [];
  try {
    // --- 1. Cloud CSV: Fetching EVERYTHING at once ---
    const csvUrl = "https://raw.githubusercontent.com/datasciencedojo/datasets/master/titanic.csv";
    steps.push(`Attempting to open Cloud CSV: ${csvUrl}`);
    
    const csvHandler = csv.open(csvUrl, {
      header: true,
      mode: "object"
    });
    
    steps.push("Using 'csv.readAll()' to fetch the entire dataset natively...");
    const allCsvRows = csv.readAll(csvHandler);
    
    steps.push(`Success: Fetched ALL ${allCsvRows.length} rows from Titanic dataset`);
    csv.close(csvHandler);

    // --- 2. Cloud JSON: Fetching via Chunked Loop ---
    const jsonUrl = "https://raw.githubusercontent.com/vega/vega/master/docs/data/movies.json";
    steps.push(`\nAttempting to open Cloud JSON: ${jsonUrl}`);
    
    const jsonHandler = json.open(jsonUrl, { format: "json" });
    
    let totalJsonRecords = 0;
    let chunkCount = 0;
    
    steps.push("Streaming JSON in chunks of 1000 records to show progress...");
    while (true) {
      const chunk = json.next(jsonHandler, { size: 1000 });
      totalJsonRecords += chunk.rows.length;
      chunkCount++;
      
      if (chunk.done) {
        steps.push(`Reached end of stream at chunk #${chunkCount}`);
        break;
      }
    }
    
    steps.push(`Success: Fetched TOTAL ${totalJsonRecords} movie records`);
    json.close(jsonHandler);

    return {
      success: true,
      message: "Full cloud data streaming verified",
      results: {
        csv: {
          total_rows: allCsvRows.length,
          method: "readAll"
        },
        json: {
          total_records: totalJsonRecords,
          method: "next_loop",
          chunks_processed: chunkCount
        }
      },
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
