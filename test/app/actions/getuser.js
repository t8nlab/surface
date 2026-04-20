import { csv } from "@titanpl/surface";

export function getuser(req) {
  const start = Date.now();

  // 1. Create a new CSV file
  const h2 = csv.create("../app/hui.csv", {
    headers: ["count", "name", "email"],
  });

  csv.write(h2, [
    {
      count: 1,
      name: "John",
      email: "john@example.com",
    },
    {
      count: 2,
      name: "Jane",
      email: "jane@example.com",
    }
  ]);

  // 3. Always close the writer
  csv.close(h2);

  // 4. Open and read back
  const h = csv.open("../app/hui.csv", {
    header: true,
    inferTypes: false,
    mode: "object",
  });

  try {
    const users = csv.readAll(h);
    const end = Date.now();
    return {
      users,
      time: end - start
    };
  } finally {
    csv.close(h);
  }
}