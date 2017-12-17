#!/usr/bin/python

import colorsys, random, sys
from redshift import run

wheel = None
size = 150
speed = 1
brightness = 1.0
saturation = 1.0

try:
	size = int(sys.argv[1])
except: pass

def rainbow(frame):
	global wheel, size, speed
	if wheel is None:
		wheel = generate_wheel(size or len(frame))
	else:
		wheel = wheel[-speed:] + wheel[:-speed]
	adjust_speed()
	if speed == 0:
	    speed = 1
	return wheel

i = 0
def adjust_speed():
	global speed, i
	if i % 150 == 0:
		speed += random.randint(-1, 1)
	elif i % 75 == 0 and random.randint(1, 2) == 1:
		if speed > 3:
			speed -= 1
		elif speed < -3:
			speed += 1
	i += 1

def generate_wheel(size):
	wheel = []
	for i in range(size):
		hue = float(i) / size
		rgb = colorsys.hsv_to_rgb(hue, saturation, brightness)
		wheel.append((int(255*rgb[0]), int(255*rgb[1]), int(255*rgb[2])))
	return wheel

run(rainbow)
