#!/usr/bin/env node
let redshift = require('./redshift')

/*
 * Adapted from github.com/scanlime/fadecandy
 * (examples/node/strip_redblue.js)
 */

redshift((frame) => {
    for (let pixel=0; pixel<frame.length; pixel++) {
        let t = pixel * 0.2 + Date.now() * 0.002;
        frame[pixel][0] = 128 + 96 * Math.sin(t);
        frame[pixel][1] = 128 + 96 * Math.sin(t + 0.1);
        frame[pixel][2] = 128 + 96 * Math.sin(t + 0.3);
    }
});