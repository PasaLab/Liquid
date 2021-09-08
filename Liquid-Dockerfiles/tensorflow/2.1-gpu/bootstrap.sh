#!/bin/bash

if ! [[ -z "${repo}" ]]; then
  if [[ -z "${branch}" ]]; then
    git clone $repo /workspace
  else
    git clone -b $branch $repo /workspace
  fi
fi

if [ -d /workspace ]; then
  cd /workspace
fi

#sleep infinity

# use eval because commands likes `key=value command` would cause file not found error when using $@, but this eval will ruin current environment
eval $@

code=$?

# Persist output
python /etc/save.py

exit $code
