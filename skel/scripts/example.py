#!/usr/bin/env python3

from redshift import run, log
import random, sys

class Example:
	index = 0
	color = [0, 0, 0]

	def __call__(self, frame):
		self.jump_around()
		self.change_color()
		frame[self.index % len(frame)] = self.color
		frame[(self.index + 1) % len(frame)] = self.darken(self.color)
		frame[(self.index - 1) % len(frame)] = self.darken(self.color)
	
	def jump_around(self):
		if random.randint(0, 30) == 0:
			self.index += random.randint(-5, 5)
	
	def change_color(self):
		# random orangy colour
		self.color[0] = random.randint(0, 255)
		self.color[1] = random.randint(0, 100)
		self.color[2] = random.randint(0, 25)

	def darken(self, color, factor=2.0):
		return map(lambda n: int(n/factor), color)

def series(effects):
	def run(frame):
		for effect in effects:
			effect(frame)
	return run

n = int(sys.argv[1]) if len(sys.argv) > 1 else 1

effects = [Example() for _ in range(n)]
run(series(effects))
