#!/bin/bash
sudo cp dhm /usr/local/bin/
sudo cp dhm.1 /usr/local/share/man/man1/
sudo mandb
echo "DHM installed! Try 'dhm --help' or 'man dhm'"
