# go-rocket
## Overview
Model rocket telemetry and control system.  Go-rocket includes software for
a base station and a package loaded into the rocket.

## Rocket package
The rocket package consists of a Raspberry Pi Zero W with a GPS receiver
connected via the serial port, altimeter, accelerometer and compass
connected via I2C bus, a camera and ground link via a wireless serial
link attached with USB (currently an XBee Pro 900 Mhz).

## Ground station
The ground station consists of a Raspberry Pi Zero W with a GPS receiver
connected via USB (serial port), and a wireless serial link to the rocket
attached with USB (currently an XBee Pro 900 Mhz).  The ground software
has a basic UI for monitoring rocket and ground position and starting the camera
via the wireless link.

