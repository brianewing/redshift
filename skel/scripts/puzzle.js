#!/usr/bin/env node

const redshift = require('./redshift');

let map = null;

redshift((frame) => {
    if(map == null || Math.random() > 0.995) {
        map = shuffle(upTo(frame.length));
        
        // randomly point some pixels back at themselves
        for(let i=0; i<map.length; i++)
            if(Math.random() > 0.28)
                map[i] = i;
    }
    
    const copy = frame.slice();
    for(let i=0; i<copy.length; i++)
        frame[i] = copy[map[i]];
});

function upTo(n) {
    const a = Array(n);
    for(let i=0; i<a.length; i++)
        a[i] = i;
    return a;
}

function shuffle(array) {
    let counter = array.length;
    while (counter > 0) {
        let index = Math.floor(Math.random() * counter);
        let temp = array[--counter];
        array[counter] = array[index];
        array[index] = temp;
    }
    return array;
}
