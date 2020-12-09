#!/bin/bash
status="$(curl -Is ${GATEWAY_RELOAD_ENDPOINT} | head -1)"
validate=( $status )
if [ -z "$validate" ]
  then
  echo "empty server status, likely blank or invalid URL"
  exit -1
fi 

if [ ${validate[-2]} == "200" ]; then
  echo "Gateway reload succeeded"
  exit 0
else
  echo "Gateway reload failed"
  exit 255
fi
