#!/bin/bash
#
# SPDX-FileCopyrightText: 2019 SAP SE or an SAP affiliate company and Gardener contributors
# SPDX-License-Identifier: Apache-2.0

set -e

echo "> Check Helm charts"

if [[ -d "$1" ]]; then
  echo "Checking for chart symlink errors"
  BROKEN_SYMLINKS=$(find -L $1 -type l)
  if [[ "$BROKEN_SYMLINKS" ]]; then
    echo "Found broken symlinks:"
    echo "$BROKEN_SYMLINKS"
    exit 1
  fi
  echo "Checking whether all charts can be rendered"
  for chart_dir in $(find charts -type d -exec test -f '{}'/Chart.yaml \;  -print -prune | sort); do
    [ -f "$chart_dir/values-test.yaml" ] && values_files="-f $chart_dir/values-test.yaml" || unset values_files
    helm template $values_files "$chart_dir" 1> /dev/null
  done
fi

echo "All checks successful"
