#!/usr/bin/env sh

set -e

if ! git status -s | grep -v "vendor/" | wc -l | grep -qE "^\s*0$"; then
  echo "There are changed files after running '${1}'"
	exit 1
fi
