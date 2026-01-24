#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');

const binName = process.platform === 'win32' ? 'git-swift.exe' : 'git-swift';
const binPath = path.join(__dirname, 'bin', binName);

// Spawn the go binary with all arguments passed through
const child = spawn(binPath, process.argv.slice(2), {
    stdio: 'inherit',
    shell: false
});

child.on('error', (err) => {
    console.error(`Error running git-swift: ${err.message}`);
    process.exit(1);
});

child.on('close', (code) => {
    process.exit(code);
});
