#!/bin/bash

if [[ "$1" == "" || "$2" == "" ]]; then
    echo -e "Usage:\t$0 X Y"
    exit 1
fi

MONITOR_SIZE="$(xrandr | grep mm | head -n 1 | awk '{ print $(NF-2) " " $(NF) }' | sed 's/mm//g')"
CVT="$(cvt ${1} ${2} | grep Modeline)"
echo "{"
echo "  services.xserver.monitorSection = ''"
echo "    " ${CVT}
echo "    " Option '"PreferredMode"' $(echo ${CVT} | awk '{ print $2 }')
echo "    " DisplaySize ${MONITOR_SIZE} # In millimeters
echo "  '';"
echo "}"
