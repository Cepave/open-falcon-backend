#!/bin/bash
#set -eu
export TERM=xterm-color

##############################################################################
## Bash Colors
red=`tput setaf 1`
green=`tput setaf 2`
yellow=`tput setaf 3`
white=`tput setaf 7`
bold=`tput bold`
reset=`tput sgr0`
separator=$(echo && printf '=%.0s' {1..100} && echo)

## Logging Finctions
log() {
  if [[ "$@" ]]; then echo "${bold}${green}[LOG `date +'%T'`]${reset} ${green}$@${reset}";
  else echo; fi
}

warning() {
  echo "${bold}${yellow}[WARNING `date +'%T'`]${reset} ${yellow}$@${reset}";
}

error() {
  echo "${bold}${red}[ERROR `date +'%T'`]${reset} ${red}$@${reset}";
}
##############################################################################
