const fs = require('fs');

const MAX_BUFFER_SIZE = 65536;

/*
 * Redshift JS helper library for external scripts
 * (c) 2017 Brian J Ewing
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

module.exports = function run(fn) {
	const fd = fs.openSync('/dev/stdin', 'r');
	const strip = []; // [[R,G,B],[R,G,B]..]
	
	console.log = console.error;
	
	for(;;) {
		if(read(fd, strip) != 0) {
			fn(strip); // mutates the buffer
			write(strip); // write it back to stdout
		} else {
			process.exit(0); // eof, let's quit
		}
	}
}

/* Pipe io */

function read(fd, strip) {
	const buffer = new Buffer(MAX_BUFFER_SIZE);
	const bytesRead = fs.readSync(fd, buffer, 0, buffer.length, null);
	unpack(Buffer.from(buffer.buffer, 0, bytesRead), strip);
	return bytesRead;
}

function write(buffer) {
	let byteArray = pack(buffer);
	process.stdout.write(byteArray);
}

/* Buffer serialization */

function pack(buffer) {
	let bytes = [];
	for(let i=0; i<buffer.length; i++) {
		bytes.push(buffer[i][0]);
		bytes.push(buffer[i][1]);
		bytes.push(buffer[i][2]);
	}
	let byteArray = new ArrayBuffer(bytes.length);
	new Uint8Array(byteArray).set(bytes);
	return Buffer.from(byteArray);
}

function unpack(byteArray, buffer) {
	let len = byteArray.length;
	for(let i=0; i<len; i=i+3) {
		if(!buffer[i/3])
			buffer[i/3] = [];
		buffer[i/3][0] = byteArray[i];
		buffer[i/3][1] = byteArray[i+1];
		buffer[i/3][2] = byteArray[i+2];
	};
	if(buffer.length > len/3)
		buffer.length = len/3; // truncate extra pixels from last frame
	return buffer;
}
