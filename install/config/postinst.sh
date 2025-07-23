#!/bin/bash
set -e # exit on error

# detect active user
USER=${SUDO_USER:-$(logname)}
echo -e "Installing for \e[1m$USER\e[0m"

# run init command at user level (not root)
runuser -l "$USER" -c '/usr/local/bin/pvault init'
