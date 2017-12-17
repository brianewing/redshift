#!/usr/bin/env node
let redshift = require('./redshift')

/*
 * Adapted from github.com/scanlime/fadecandy
 * (examples/javascript/strip_redblue.js)
 */

redshift((frame) => {
	for (let i=0; i<frame.length; i++) {
		let t = i * 0.2 + Date.now() * 0.002;
		frame[i][0] = 128 + 96 * Math.sin(t);
		frame[i][1] = 128 + 96 * Math.sin(t + 0.1);
		frame[i][2] = 128 + 96 * Math.sin(t + 0.3);
	}
});
