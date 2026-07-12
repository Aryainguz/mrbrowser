import server from './dist/server/server.js';
const req = new Request('http://localhost/');
const res = await server.fetch(req);
console.log("Status:", res.status);
console.log("Headers:", res.headers);
console.log("Body preview:", (await res.text()).slice(0, 100));
