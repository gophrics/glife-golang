#!/bin/bash

# turn on bash's job control
set -m

# Start the primary process and put it in the background
./bin/profile &

fg %1