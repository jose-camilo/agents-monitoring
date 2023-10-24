#! /bin/bash
docker run --name agents -e MYSQL_ROOT_PASSWORD=verysecuresecret -e MYSQL_DATABASE=agents -p 3306:3306 -d mysql:8.0