import { csv, smtp } from "@titanpl/surface";

export default function bulkmail(req) {
  const start = Date.now();
  const csvPath = "../app/bulk_test.csv";

  // // 1. GENERATE CSV with 200 rows natively
  // const hCreate = csv.create(csvPath, { headers: ["to", "subject", "body"] });
  // const mockRows = [];
  // for (let i = 1; i <= 200; i++) {
  //   mockRows.push({
  //     to: "ezetapp@gmail.com",
  //     subject: `Native Bulk Test #${i}`,
  //     body: `<h1>Message ${i}</h1><p>Sent via Surface Parallel Engine.</p>`
  //   });
  // }
  // csv.write(hCreate, mockRows);
  // csv.close(hCreate);

  // 2. READ CSV back (High speed)
  const hRead = csv.open(csvPath, { mode: "object" });
  const emails = csv.readAll(hRead);
  csv.close(hRead);

  // 3. BLAST EMAILS in parallel (using Go worker pool)
  const settings = {
    host: "smtp.gmail.com",
    port: 587,
    username: "clashersoham07@gmail.com",
    password: "jjke wzkr tyfs aeod", 
    from: "clashersoham07@gmail.com"
  };

  try {
    const results = smtp.bulk({
      ...settings,
      emails,
      concurrency: 10 // Use 10 parallel connections for speed
    });

    const end = Date.now();
    return {
      success: true,
      totalSent: results.length,
      timeTaken: `${end - start}ms`,
      results: results.slice(0, 5) // Return first 5 results as proof
    };
  } catch (err) {
    return { success: false, error: err.message };
  }
}
