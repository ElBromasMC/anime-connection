const WebSocket = require('ws');
const http = require('http');

// Creating a WebSocket server on port 8010
const wss = new WebSocket.Server({ port: 8010 });

wss.on('connection', _ => {
    console.info('LiveReload: A new client connected.');
});

// Function to broadcast a reload event to all connected clients
function broadcastReload() {
    wss.clients.forEach(client => {
        if (client.readyState === WebSocket.OPEN) {
            client.send('reload');
        }
    });
}

// Creating an HTTP server on port 8020
const server = http.createServer((req, res) => {
    if (req.method === 'GET') {
        console.info('LiveReload: Received reload trigger via GET request');
        // Triggering reload event
        broadcastReload();
        res.statusCode = 200;
        res.setHeader('Content-Type', 'text/plain');
        res.end('Reload event triggered.\n');
    }
});

server.listen(8020, () => {
    console.info('LiveReload: HTTP server running on port 8020');
});