#!/bin/bash
# wait-for.sh

set -e

host="$1"
shift
cmd="$@"

until mongo --host "$host" --eval "print(\"waited for connection\")"; do
  >&2 echo "MongoDB is unavailable - sleeping"
  sleep 1
done

>&2 echo "MongoDB is up - executing command"
exec $cmd
