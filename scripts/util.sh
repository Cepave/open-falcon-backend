#!/usr/bin/env bash

# A POSIX variable
OPTIND=1         # Reset in case getopts has been used previously in the shell.

# Initialize variables
container_list="agent aggregator alarm dashboard fe graph hbs judge links nodata portal query sender task transfer"
callback=""
sub_command=""
options=""
args=""
verbose=false

#
#
# Utility function
#
#
function msg() {
  skip=${1}+1
  echo "[$2][${FUNCNAME[$skip]}] ${@:3}" >&2
}

function vmsg() {
  if $verbose ; then
    msg ${1}+1 "${@:2}"
  fi
}

function invalid_option() {
  msg 1 "ERR" "Invalid option: $@"
}

function require_option() {
  msg 1 "ERR" "Option $@ requires an argument."
}

function usage() {
  echo "[MSG] The script provides some useful utilities in testing env."
  echo "      $0 [sub_command] [options...] [args...]"
  echo "[CMD] Docker"
  echo "      -R             [-f][all] docker-compose rm, -f for init.yml; all for both yml files."
  echo "      -S             [-f][all] docker-compose stop, -f for init.yml; all for both yml files."
  echo "      -l             [-f][-t NUM][args...] docker logs a contaner. (default: null)"
  echo "      Control"
  echo "      -c             [args...] control status check. (default: all)"
  echo "      -s             [args...] control start. (default: null)"
  echo "      -r             [args...] control restart. (default: null)"
  echo "      -C             [-f][args...] clean data under /home/openfalcon/ (default: all)"
  echo "[OPT] Flag"
  echo "      -f             [option] add -f flag"
  echo "      -t NUM         [option] add --tail=NUM"
  echo "      -h             [option] help"
  echo "      -v             [option] verbose"
  echo "[ARG] Name of containers"
  exit 0
}

function parse() {
  while getopts ":RSlcsrCft:hv" opt; do
    case $opt in
      #
      # [Sub_command]
      #
      # Docker
      R)
        sub_command="remove"
        ;;
      S)
        sub_command="stop"
        ;;
      l)
        sub_command="logs"
        ;;
      # Control
      c)
        sub_command="check"
        ;;
      s)
        sub_command="start"
        ;;
      r)
        sub_command="restart"
        ;;
      C)
        sub_command="clean"
        ;;
      #
      # [Option]
      #
      f)
        options+="-f "
        ;;
      t)
        options+="--tail=$OPTARG "
        ;;
      h)
        usage
        ;;
      v)
        verbose=true
        vmsg 0 "MSG" "${!verbose@}: $verbose"
        ;;
      \?)
        invalid_option "-$OPTARG"
        exit 1
        ;;
      :)
        require_option "-$OPTARG"
        exit 1
        ;;
    esac
  done
}

function is_substring() {
  sub=$1
  str=${@:2}
  if [[ "$str" == *"$sub"* ]]; then
    return 1
  fi
  return 0
}

# Do callback for each container
function ioc() {
  vmsg 1 "MSG" "${FUNCNAME[1]}() the following containers: [$@]"
  for i in $@; do
    is_substring "$i" "$container_list"
    # Check func's return value
    if [[ $? == 1 ]]; then
      eval $callback $i
    else
      msg 1 "ERR" "No such container [$i]"
    fi
  done
}

function check_one() {
  container=$1
  output=$(docker exec $container /home/$container/control status)

  # Check the exit status of last command
  if [[ $? != 0 ]]; then
    echo "[ERROR] $container"
  # Match the substring in the output
  elif [[ $output == *"stoped"* ]]; then
    echo "[STOP!] $container"
  else
    echo "[PASS.] $container"
  fi
}

function check() {
  callback=check_one
  # All container
  if [[ $1 == "" ]]; then
    ioc $container_list
  else
    # Some containers
    ioc $@
  fi
}

function control() {
  if [[ $# == 0 ]]; then
    msg 0 "ERR" "Invalid sub_command [$sub_command] and args [$args]"
    return
  fi
  # Some containers
  ioc $@
}

function start_one() {
  container=$1
  docker exec $container /home/$container/control start
}

function restart_one() {
  container=$1
  docker exec $container /home/$container/control restart
}

function log_one() {
  container=$1
  docker logs $options $container
}

function clean_one() {
  container=$1
  msg 0 "MSG" "sudo rm -r $options /home/openfalcon/$container"
  sudo rm -r $options /home/openfalcon/$container
}

function clean() {
  callback=clean_one
  # All
  if [[ $1 == "" ]]; then
    msg 0 "MSG" "sudo rm -r $options /home/openfalcon/*"
    sudo rm -r $options /home/openfalcon/*
  else
    # Some containers
    ioc $@
  fi
}

function compose_one() {
  cmd=$@
  # Default
  if [[ $options == "" ]]; then
    docker-compose $cmd
  # -f init.yml
  elif [[ $options == "-f " ]]; then
    docker-compose -f init.yml $cmd
  else
    msg 0 "ERR" "[$sub_command]: args [$args]"
  fi
}

function compose() {
  opt=$1
  cmd=${@:2}

  if [[ $opt == "all" ]]; then
    options=""
    compose_one $cmd
    options="-f "
    compose_one $cmd
  else
    compose_one $cmd
  fi
}

#
#
# Main function
#
#

function main() {
  parse $@

  # Shift params & print msgs
  vmsg 0 "MSG" "${!sub_command@}: [$sub_command]"
  vmsg 0 "MSG" "${!options@}: [$options]"
  shift $((OPTIND-1))
  args=$@
  vmsg 0 "MSG" "${!args@}: [$args]"

  case $sub_command in
    check)
      check $args
      ;;
    start)
      callback=start_one
      control $args
      ;;
    restart)
      callback=restart_one
      control $args
      ;;
    logs)
      callback=log_one
      control $args
      ;;
    clean)
      clean $args
      ;;
    stop)
      compose $1 "stop"
      ;;
    remove)
      compose $1 "rm -f"
      ;;
    *)
      msg 0 "ERR" "Invalid sub_command [$sub_command] and args [$args]"
      ;;
  esac

}

main $@

# End of file