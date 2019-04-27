# Redshift LED Light Server & Effects Library

Addressable LEDs can be used to create [art installations](https://www.youtube.com/watch?v=7_MhvCsibZg), [interactive clothing](https://learn.adafruit.com/category/flora), [real-time music visualisations](https://www.youtube.com/watch?v=tnYHr8YYkiM) and more.

This project is an attempt to make programming LEDs more convenient and fun, with a real-time remote control interface / web-based IDE. It's still in design stages but useful enough to experiment with.

In the end, I hope this will be enable more people to experiment with programming and controlling lighting / art installations in their own homes and other spaces.

[github.com/brianewing/redshift-web](https://github.com/brianewing/redshift-web)

## Features

* Compose, layer and blend an effect chain with real time feedback, save + reload
* Growing list of built-in effects like Rainbow, Larson Scanner, Game of Life, Mood, Stripe, Strobe, Gamma, Brightness, Sepia
* Building block effects such as Layer, Mirror, Layout, Switch, Toggle, Trigger
* Effect parameters controllable in real time via UI, web socket, OPC, OSC, MIDI, time, tween functions and JavaScript expressions
* External effect support with simple stdio interface (pixel buffer bytes in/out)
	* Web IDE has syntax support for Python, JavaScript and CoffeeScript. Extensions welcome
	* Any language can be used, so long as it can handle stdio
	* Hot reload whenever the executable file / script changes
* Flexible [OpenPixelControl](https://github.com/zestyping/openpixelcontrol) input and forwarding - server and effects are configurable using system exclusive commands
* Render to ws2811/ws2812(b) LEDs (Raspberry Pi supported, pull requests welcome for other hardware)
* Web UI in development, rendering at 60 FPS with a real time script editor and minimalist, responsive controls

## Getting Started

Working with real LEDs and electricity takes a lot of patience, responsibility & care.

You are responsible for your own safety and hardware, and this software comes without any warranty.

If you just want to experiment with the web UI and start making some effects, skip to Prerequisites!

For information about setting up ws2812 LEDs with a Raspberry Pi, see [jgarff/rpi_ws281x](https://github.com/jgarff/rpi_ws281x) and [the Adafruit website on NeoPixel best practices](https://learn.adafruit.com/adafruit-neopixel-uberguide/best-practices).

### Prerequisites

The software requires a working Go environment to build (see [https://golang.org/doc/install](https://golang.org/doc/install))

ws2811/2(b) is supported on Raspberry Pis running Linux (`go build +ws2811`) with the [jgarff/rpi_ws281x](https://github.com/jgarff/rpi_281x) library installed.

Make sure rpi_ws281x headers (.h) and objects (.so / .a) have been copied to /usr/local/include and /usr/local/lib respectively if you intend to use the `+ws2811` build tag.

### Installing

To set up Redshift, assuming `$GOPATH` is `~/go` (default):

0. `$ go get github.com/brianewing/redshift`
0. `$ cd ~/go/src/github.com/brianewing/redshift`
0. `$ go build`
0. `$ ./redshift`

It's recommended to install and run the [web UI](https://github.com/brianewing/redshift-web) too.

With the server and UI running, you should see something like this:

![Screenshot](https://i.imgur.com/HuYXYPA.png)

## Running tests

The server includes a few benchmarks to compare the performance of different hardware and catch regressions.

To run them:

0. `cd ~/go/src/github.com/brianewing/redshift`
0. `go test -bench=.`

## Built with

* [lucasb-eyer/go-colorful](https://github.com/lucasb-eyer/go-colorful) - for colour blending and working with colour systems
* [gorilla/websocket](https://github.com/gorilla/websocket) - web socket implementation (see `server/web_socket_server.go`)
* [jgarff/rpi_ws281x](https://github.com/jgarff/rpi_ws281x) - driving ws2811/2(b)'s with the RPi PWM

## Contributing

Contributions are welcome! Hack away!

## Versioning

The project uses [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/brianewing/redshift/tags). 

## Authors

* Brian Ewing ([github.com/brianewing](https://github.com/brianewing))

See also the list of [contributors](https://github.com/brianewing/redshift/contributors) who participated in this project, if they are not already listed here.

## License

This project is licensed under the Affero GPL v3 License - see the LICENSE file for details

## Acknowledgements

* [Sara van der Valk](http://www.bananenmelk.nl), thanks for your support :)
* [scanlime/fadecandy](https://github.com/scanlime/fadecandy) inspiration
* [hyperion-project/hyperion](https://github.com/hyperion-project/hyperion) ^^
