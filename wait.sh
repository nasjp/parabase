#!/bin/sh

until mysqladmin ping -h mysql -u root --password=password --silent;
do
  sleep 1
done
