import http from 'node:http';
import fs from 'node:fs';
import readline from 'node:readline';
import path from 'node:path';

const port = 3000;
const csvPath = path.resolve('../app/users.csv');

http.createServer(async (req, res) => {
  if (req.url === '/users') {
    const start = Date.now();
    const users = [];
    
    try {
      const fileStream = fs.createReadStream(csvPath);
      const rl = readline.createInterface({
        input: fileStream,
        crlfDelay: Infinity
      });

      let isHeader = true;
      let headers = [];

      for await (const line of rl) {
        const parts = line.split(',');
        if (isHeader) {
          headers = parts;
          isHeader = false;
          continue;
        }
        
        const user = {};
        headers.forEach((h, i) => {
          user[h.trim()] = parts[i]?.trim();
        });
        users.push(user);
        
        if (users.length >= 10000) break;
      }

      const end = Date.now();
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({
        time: `${(end - start).toFixed(2)}ms`,
        count: users.length,
        users
      }));
    } catch (err) {
      res.writeHead(500);
      res.end(err.message);
    }
  } else {
    res.writeHead(404);
    res.end('Not Found');
  }
}).listen(port, () => {
  console.log(`Node.js Native server listening at http://localhost:${port}`);
  console.log(`CSV Path: ${csvPath}`);
});
