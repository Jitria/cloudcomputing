// server.js
const express = require('express');
const { createServer } = require('http');
const { Server } = require('socket.io');
const pty = require('node-pty');

const app = express();
const server = createServer(app);
const io = new Server(server);
const port = 3000;

app.use(express.static('public'));

io.on('connection', (socket) => {
    console.log('Client connected');

    const shell = pty.spawn('sh', [], {
        name: 'xterm-color',
        cols: 80,
        rows: 24,
        cwd: process.env.HOME,
        env: process.env
    });

    shell.on('data', (data) => {
        socket.emit('output', data);
    });

    socket.on('input', (input) => {
        shell.write(input + '\r');
    });

    socket.on('disconnect', () => {
        shell.kill();
        console.log('Client disconnected');
    });
});

server.listen(port, () => {
    console.log(`Server running at http://localhost:${port}`);
});
