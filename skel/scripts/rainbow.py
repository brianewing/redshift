#!/usr/bin/python

import colorsys, random, sys
from redshift import run, log

class Rainbow:
	wheel = None
	size = 150
	speed = 1
	brightness = 1.0
	saturation = 1.0
	
	def __call__(self, frame):
		if self.wheel is None:
			self.wheel = self.generate_wheel(self.size or len(frame))
		else:
			self.wheel = self.wheel[-self.speed:] + self.wheel[:-self.speed]
			
		return self.wheel
	
	def generate_wheel(self, size):
		wheel = []
		
		for i in range(size):
			hue = float(i) / size
			rgb = colorsys.hsv_to_rgb(hue, self.saturation, self.brightness)
			wheel.append((int(255*rgb[0]), int(255*rgb[1]), int(255*rgb[2])))

		return wheel

effect = Rainbow()

try:
	effect.size = int(sys.argv[1])
	effect.speed = int(sys.argv[2])
except:
	pass

log("Starting rainbow effect, size=%r, speed=%r" % (effect.size, effect.speed))

run(effect)
