#!/usr/bin/env node

const redshift = require('./redshift')

function bgAnimation(numLeds=60, fps=60) {
	return new Animation(numLeds, [
		new Clear,
		new Mood,
		new MirrorLarson([0, 0, 255]),
	])
}

class Animation {
	constructor(numberOfLeds, effects=[]) {
		this.buffer = new Array(numberOfLeds).fill(null).map(() => [0, 0, 0])
		this.effects = effects
	}

	render() {
		for(let i=0; i<this.effects.length; i++) {
			this.effects[i].render(this.buffer)
		}
		return this.buffer
	}
}

class Clear {
	render(buffer) {
		for(let i=0; i<buffer.length; i++) {
			buffer[i] = [0, 0, 0]
		}
	}
}

class Mood {
	constructor() {
	    this.brightness = 0
	    this.delta = 2
	}

	render(buffer) {
		if(this.brightness <= 2) {
			this.color = this.newColor()
			this.delta = 2
		} else if(this.brightness >= 255) {
			this.delta = -2
		}

		let adjustedColor = this.color.map((x, _) => Math.round(x*this.brightness/255))
		this.brightness += this.delta

		for(let i=0; i<buffer.length; i++)
			buffer[i] = adjustedColor
	}

	newColor() {
		return [0, 0, 0].map(() => Math.floor(Math.random() * 255))
	}
}

class MirrorLarson {
	constructor(color) {
		this.color = color
		this.position = 0
		this.direction = 1 // -->
	}

	render(buffer) {
		if(this.position >= buffer.length / 2)
			this.direction = -1
		else if(this.position <= 0)
			this.direction = 1

		this.position += this.direction

		buffer[this.position] = this.color
		buffer[buffer.length-this.position-1] = this.color
	}
}

let animation // e.g. = bgAnimation(20)
redshift(frame => {
    if(!animation)
        animation = bgAnimation(frame.length)

    return animation.render()
})
