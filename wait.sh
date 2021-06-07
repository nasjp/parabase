#!/bin/sh

set -e

cmd="$@"

echo "Waiting for MySQL"
until mysql -h mysql -u root --password=password &> /dev/null
do
        >&2 echo -n "."
        sleep 1
done

>&2 echo "MySQL is Up - executing command"
exec $cmd
