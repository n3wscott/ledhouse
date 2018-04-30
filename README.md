# ledhouse

LED House runs on an Ardunio with a strip of 10 NeoPixels. It is intended to be
a representation of home light automation for my talk at Kubecon EU 2018
entitled [Kubernetes as an Abstraction Layer for a Connected Home](http://sched.co/DqwC)

I built this because it is much more portable to demo than bringing Hue Lights
and a hub. 

This currently speaks at 115200 Baud by default.

### Protocol

This connects on the serial port, `LED` is used to key off of, `\n` is used to
signal that the command is finished.

Control a light: `LED<Index><Red><Green><Blue>\n` where,

Index is a 2 digit hex value 0-10
Red, Green, and Blue are 2 diget hex values 0-255.

Examples:

LED 5, Red: `LED05ff0000`

LED 0, Blue: `LED000000ff`

LED 3, White: `LED03ffffff`

To clear the ledhouse: `CLEAR\n`

### Exercise

To see if everything is working, use the exercise tool, example:

`go run exercise.go --serial=/dev/cu.usbmodem1441`
