#!/usr/bin/env js

/*
 * Redshift JS helper library for external scripts
 * (c) 2018 Brian J Ewing
 * @license AGPL 3.0
 *
 * Usage:
 *   const redshift = require('./redshift');
 *   redshift((buffer) => {
 *     // This effect function will be called once every frame
 *     // Implement your effect by modifying the buffer!
 *     // ( buffer is an array of leds like [[r,g,b], ...] )
 *   });
 */

const fs = require('fs');

const MAX_BUFFER_SIZE = 65536;

function run(fn) {
	const fd = fs.openSync('/dev/stdin', 'r');
	
	console.log = console.error;
	
	for(;;) {
		const input = read(fd);
		if(input) {
			let result = eval(input);
			write(result);
		} else {
			process.exit(0); // eof, let's quit
		}
	}
}

run()

/* Pipe io */

function read(fd, strip) {
	let buffer = new Buffer(MAX_BUFFER_SIZE);
	let bytesRead = fs.readSync(fd, buffer, 0, buffer.length, null);
	return Buffer.from(buffer.buffer, 0, bytesRead).toString('utf8');
}

function write(result) {
	process.stdout.write(JSON.stringify(result));
}
