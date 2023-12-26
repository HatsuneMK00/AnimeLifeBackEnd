#!/bin/bash

# Check the number of parameters
if [ "$#" -ne 2 ]; then
  echo "Illegal number of parameters. Usage: ./script.sh [version] [minor/patch]"
  exit 1
fi

# Read the version and choice from the command line parameters
version=$1
choice=$2

# Split the version into major, minor, and patch
major=$(echo $version | awk -F. '{print $1}')
minor=$(echo $version | awk -F. '{print $2}')
patch=$(echo $version | awk -F. '{print $3}' | cut -d '-' -f1)

# Increment the minor or patch version based on the user's choice
if [ "$choice" == "minor" ]; then
  minor=$((minor+1))
  patch=0
elif [ "$choice" == "patch" ]; then
  patch=$((patch+1))
else
  echo "Invalid choice. Please choose either 'minor' or 'patch'."
  exit 1
fi

# Reassemble and print the new version
new_version="$major.$minor.$patch"
echo $new_version