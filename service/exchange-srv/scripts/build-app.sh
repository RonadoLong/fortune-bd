#!/bin/bash
gocmd=$1
$gocmd version
CGO_ENABLED=0 GOOS=linux $gocmd build  -o  app
