// public/script.js
document.addEventListener('DOMContentLoaded', () => {
    const terminal = document.getElementById('terminal');
    const input = document.getElementById('input');
    const socket = io();

    socket.on('output', (data) => {
        terminal.innerHTML += `<div>${data.replace(/\n/g, '<br>')}</div>`;
        terminal.scrollTop = terminal.scrollHeight;
    });

    input.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            const command = input.value;
            if (command.trim() !== '') {
                terminal.innerHTML += `<div>> ${command}</div>`;
                input.value = '';
                socket.emit('input', command);
            }
        }
    });
});
