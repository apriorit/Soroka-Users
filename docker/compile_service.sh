#!/usr/bin/env bash

cd ../
./build_service.sh
cd docker
./run_composer.sh up --build