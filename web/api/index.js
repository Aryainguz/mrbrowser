import server from '../dist/server/server.js';

export default async function handler(req, res) {
  try {
    // 1. Convert Node.js req to Web Request
    const protocol = req.headers['x-forwarded-proto'] || 'https';
    const host = req.headers['x-forwarded-host'] || req.headers['host'];
    const url = new URL(req.url, `${protocol}://${host}`);
    
    // Some headers might be arrays, so we stringify them or pass them as is. 
    // Fetch API Request accepts a plain object for headers.
    const headers = new Headers();
    for (const [key, value] of Object.entries(req.headers)) {
      if (Array.isArray(value)) {
        value.forEach(v => headers.append(key, v));
      } else if (value) {
        headers.set(key, value);
      }
    }

    const init = {
      method: req.method,
      headers: headers,
    };

    if (req.method !== 'GET' && req.method !== 'HEAD') {
      // For Vercel Serverless Functions, req.body is often pre-parsed as an object or string
      if (typeof req.body === 'string') {
        init.body = req.body;
      } else if (Buffer.isBuffer(req.body)) {
        init.body = req.body;
      } else if (typeof req.body === 'object' && req.body !== null) {
        init.body = JSON.stringify(req.body);
      } else {
        init.body = req;
        init.duplex = 'half';
      }
    }

    const request = new Request(url.href, init);

    // 2. Call the TanStack Start fetch handler
    const response = await server.fetch(request);

    // 3. Convert Web Response back to Node.js res
    res.status(response.status);
    
    response.headers.forEach((value, key) => {
      // Avoid setting set-cookie incorrectly if there are multiple
      if (key.toLowerCase() === 'set-cookie') {
        let cookies = res.getHeader('set-cookie') || [];
        if (!Array.isArray(cookies)) cookies = [cookies];
        cookies.push(value);
        res.setHeader('set-cookie', cookies);
      } else {
        res.setHeader(key, value);
      }
    });

    if (response.body) {
      // Vercel Serverless Functions support streaming via Web Streams or async iterables
      const reader = response.body.getReader();
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        res.write(value);
      }
      res.end();
    } else {
      res.end();
    }

  } catch (error) {
    console.error('Vercel API Handler Error:', error);
    res.status(500).setHeader('Content-Type', 'text/plain').send(`SSR Server Error: ${error.message}\n\n${error.stack}`);
  }
}
