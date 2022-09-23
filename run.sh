#!/bin/sh

LEDS=${LEDS:-256}

./redshift -log \
	-davAddr 127.0.0.1:9292 \
	-effectsDavAddr 127.0.0.1:9393 \
	-httpAddr 127.0.0.1:9191 \
	-opcAddr 127.0.0.1:7890 \
	-oscAddr 0.0.0.0:9494 \
	-effectsYaml effects.yaml \
	-leds $LEDS
