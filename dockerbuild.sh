#!/bin/bash

if [ "$1" == "test" ]; then
  docker-compose run --rm test
elif [ "$1" == "codegen" ]; then
  docker-compose run --rm codegen
else
  docker-compose run --rm build
  docker-compose build --force-rm server
  docker-compose build --force-rm cron
fi

docker-compose down
