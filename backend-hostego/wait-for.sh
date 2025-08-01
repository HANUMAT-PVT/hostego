#!/bin/sh

HOST=$1
shift
CMD="$@"
exec $CMD
echo "⏳ Waiting for $HOST..."

# while ! nc -z $(echo $HOST | cut -d: -f1) $(echo $HOST | cut -d: -f2); do
#   sleep 1
# done

echo "✅ $HOST is available, running: $CMD"

