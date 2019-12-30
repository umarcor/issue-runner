#!/bin/sh

set -e

gstart () {
  printf "$@\n"
}
gend () {
  :
}

[ -n "$CI" ] && {
  gstart () {
    printf "::group::$@\n"
  }
  gend () {
    echo '::endgroup::'
  }
} || echo "INFO: not in CI"

#---

gstart "Check for uncommitted changes"
git diff --exit-code --stat -- . ':!node_modules' \
|| (echo "::error:: found changed files after build" && exit 1)
gend

gstart "Update files in branch gha-tip"
cp action.yml dist/
head -$(grep -n "## Development" README.md | cut -f1 -d:) README.md | sed '$d' > dist/README.md
cd dist
git init
git checkout --orphan gha-tip
git add .
gend

gstart "Commit changes"
git config --local user.email "tip@gha"
git config --local user.name "GHA"
git commit -a -m "update $GITHUB_SHA"
gend

gstart "Push to origin"
git remote add origin "https://x-access-token:$GITHUB_TOKEN@github.com/$GITHUB_REPOSITORY"
git push origin +gha-tip
gend
