#!/usr/bin/bash
#
# This script:
# - updates the 'Next version' in breaking-changes.md to match the specified version
# - add back the 'Next version' section, with template

set -eou pipefail

# used - universal sed
# sed that works on macos and linux.
used() {
  case $(uname -s) in
    *[Dd]arwin* | *BSD* ) sed "$@";;
    *) sed "$@";;
  esac
}

# Ugly but needed for MacOS. Sed on MacOS fails when calling sed --help with an illegal argument exception.
# Exploit this to check if the user has a bad sed version and exit early.
# if ! sed --help >/dev/null 2>&1; then
#   echo "sed does not look like GNU sed. Are you running on MacOS perhaps? Install GNU sed, make it so you can call it with sed and retry"
#   exit 8
# fi

updateFile() {
  local file
  file=$1
  # offset of placeholder content after Next version
  local offset
  offset=$2

  newfile="newfile.md"

  # get the line where we will introduce back the next version section
  LN=$(grep -n "## Next version" "$file" | head -1 | grep -Eo '^[^:]+')
  # create anchor from version string
  anchor=$(echo "$VERSION" | tr . -)
  # create the new file content
  # take lines up to next version section
  head -n "$(($LN + $offset))" "$file" > "$newfile"
  # introduce the new version header
  echo "## $VERSION [$anchor]" >> "$newfile"
  echo "% **Release date:** Month day, year" >> "$newfile"
  # Add the rest of the file, skip placeholder content, replace next section with version.
  tail -n+$(($LN + $offset)) "$file" | \
    used "s|## Next version.*|## $VERSION [$anchor]|g" /dev/stdin >> "$newfile"
  # replace breaking change file with new content
  mv "$newfile" "$file"
}

updateBreakingChanges() {
  # offset is number of lines below ## Next version + 1
  updateFile "./docs/release-notes/breaking-changes.md" 11
}

updateDeprecations() {
  # offset is number of lines below ## Next version + 1
  updateFile "./docs/release-notes/deprecations.md" 9
}

updateIndex() {
  file=./docs/release-notes/index.md
  newfile="newfile.md"

  # get the line where we will introduce back the next version section
  LN=$(grep -n "## version.next" "$file" | head -1 | grep -Eo '^[^:]+')
  # create anchor from version string
  anchor=$(echo "$VERSION" | tr . -)
  # create the new file content
  # take lines up to next version section
  head -n "$(($LN - 1))" "$file" > "$newfile"
  # add template
  echo "## Next version [elastic-apm-next-release-notes]" >> "$newfile"
  echo "% **Release date:** Month day, year" >> "$newfile"
  echo "" >> "$newfile"
  echo "### Features and enhancements [elastic-apm-next-features-enhancements]" >> "$newfile"
  echo "% * 1 sentence describing the change. ([#PR number](https://github.com/elastic/apm-server/pull/PR number))" >> "$newfile"
  echo "" >> "$newfile"
  echo "### Fixes [elastic-apm-next-fixes]" >> "$newfile"
  echo "% * 1 sentence describing the change. ([#PR number](https://github.com/elastic/apm-server/pull/PR number))" >> "$newfile"
  echo "" >> "$newfile"
  # introduce the new version header
  echo "## $VERSION [$anchor]" >> "$newfile"
  echo "% **Release date:** Month day, year" >> "$newfile"
  # Add the rest of the file, skip placeholder content, replace next section with version.
  tail -n+$(($LN + 2)) "$file" | \
    used "s|elastic-apm-next-features-enhancements|$anchor-features-enhancements|g" /dev/stdin | \
    used "s|elastic-apm-next-fixes|$anchor-fixes|g" /dev/stdin >> "$newfile"
  # replace breaking change file with new content
  mv "$newfile" "$file"
}


VERSION=${1:-}

test -z VERSION && {
  echo "VERSION is required"
  exit 2
}

echo ">> updating changelog for $VERSION"

updateBreakingChanges
updateDeprecations
updateIndex
