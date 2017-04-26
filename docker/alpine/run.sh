#!/bin/bash

set -e

# Launch
echo "start $@"
./open-falcon start --console-output $@
errno=$?
if [ $errno -ne 0 ] ; then
  echo "Failed to start"
  exit 1
fi

# Monitor
modules=($@)
pattern=""
for (( idx=0; idx<$#; ++idx )) ; do
  pattern=$pattern"falcon-"${modules[$idx]}"|"
done
pattern=${pattern:0:-1}
pid=$(pgrep -f "$pattern")
while kill -0 $pid 2>/dev/null ; do
  sleep 10
done
for mod in "$@" ; do
  echo "$mod exited"
done
