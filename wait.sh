#!/bin/sh

until mysqladmin ping -h$MYSQL_HOST -P$MYSQL_PORT -u$MYSQL_USER -p$MYSQL_PASSWORD;
do
  sleep 1
done
