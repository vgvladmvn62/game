# Bullseye Hardware

## Overview

This part contains implementation of hardware client.

## Prerequisites

Bullseye shelves with cables to connect to Raspberry Pi. Additional hub might be
useful as RPi does not support many USB ports.

## Installation

One way of installing this software on each shelf is so connect it to the
computer and with Arduino IDE installed run the following script:

`ls /dev/cu.usb* && arduino --port $(ls /dev/cu.usb*) --upload slab_hw.ino`

## Configuration

After plugging shelves into Raspberry Pi ports and turning RPi on
configuration is as it follows - placing hand on a shelf assigns it's ID,
starting from 0. After performing this step color of the shelf should change
to green. If configuration was performed in the past then it is stored in
RPi database and no more actions are required. In any other case changing assigned IDs
is possible and approachable.

## Usage

Software running on shelves performs all actions on it's own so no additional
steps are required from a user or administrator perspective.