#!/bin/bash

previous_version=""

replace() {
  latest_version=$(./latest_tag_fetcher)
  echo "Latest tag at the repository is $latest_version"
  previous_version=$(cat "$1"/values.yaml | grep -i "version:"| cut -d ":" -f2 | xargs echo -n)
  echo "Previous Latest tag at the repository is $previous_version"
  echo "Replacing $previous_version with $latest_version in $1/values.yaml"
  sed -i'' "s/$previous_version/$latest_version/g" "$1"/values.yaml
  export NEXT_VERSION=${latest_version}
}
