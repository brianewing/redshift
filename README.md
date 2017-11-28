# Redshift LED Light Server & Effects Library

WS2812 LEDs have been [used for art installations](https://www.youtube.com/watch?v=7_MhvCsibZg), [sewn into clothing](https://learn.adafruit.com/category/flora), [programmed for real-time music visualisation](https://www.youtube.com/watch?v=tnYHr8YYkiM), the possibilities are endless.

This project is an attempt to make programming and interacting with these setups more convenient and fun, with a real-time remote control interface and web-based IDE.

In the end, I hope this interface will be simple enough to enable people to control dynamic lighting / art in the home, and maybe serve as an introduction to programming.

[=> Client app coming soon <=](https://github.com/brianewing/redshift-app)

## Features

* Built-in effects (rainbow, larson, brightness, mood, stripe)
* Blending, layering and composition of effects (create your own interesting combinations)
* Script engine & external effect support (simple stdio interface, buffer in/out)
* Effects list & parameters controllable in real time via HTTP
* [OpenPixelControl](https://github.com/zestyping/openpixelcontrol) server
* Renders to WS2811/WS2812(b) LEDs (Raspberry Pi supported, eventually other hardware, pull requests welcome)
* MIDI control coming soon
* Parameter automation coming soon (e.g. tween effect parameters with a sin wave)
* Beautiful web UI coming soon (renders at 60fps, real time script editor, effect controls)

## Getting Started

These instructions will get you a copy of the project up and running for development and testing purposes. Working with real LEDs takes time, patience, a lot of responsibility & care.

Please read and understand the instructions [here](https://github.com/jgarff/rpi_ws281x) and [on the Adafruit website](https://learn.adafruit.com/adafruit-neopixel-uberguide/best-practices) to achieve a safe and reliable setup of ws2812 LEDs on a Raspberry Pi

YOU ARE RESPONSIBLE FOR YOUR OWN HARDWARE, SAFETY, ENVIRONMENT AND BELONGINGS. 
You should seek electrical advice before working with live wires at these currents.

The software comes without a warranty!

### Prerequisites

The software requires a working Go environment to build (see [https://golang.org/doc/install](https://golang.org/doc/install))

ws2811/2(b) is supported on Raspberry Pis running Linux (`go build +ws2811`)

### Installing

To set up Redshift, assuming `$GOPATH` is `~/go` (default):

0. `$ go get github.com/brianewing/redshift`
0. `$ cd ~/go/src/github.com/brianewing/redshift`
0. `$ go build`
0. `$ ./redshift`

It's recommended to install and run the client app too!

With the server and client running, you should see something like this:

![Screenshot](https://i.imgur.com/FVmPan3.png)

## Running tests

The server includes a few benchmarks to compare the performance of different hardware and prevent regressions.

To run these:

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

* **Brian Ewing** ([github.com/brianewing](https://github.com/brianewing))

See also the list of [contributors](https://github.com/brianewing/redshift/contributors) who participated in this project, if they are not already listed here.

## License

This project is licensed under the Affero GPL v3 License - see the LICENSE file for details

## Acknowledgements

* [Sara van der Valk](http://www.bananenmelk.nl), thanks for your help & enthusiasm :)
* [scanlime/fadecandy](https://github.com/scanlime/fadecandy) - effect inspiration and occasional reference
* [hyperion-project/hyperion](https://github.com/hyperion-project/hyperion) - more effect inspiration
