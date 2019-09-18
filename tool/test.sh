#!/usr/bin/env sh

set -e

enable_color() {
  ENABLECOLOR='-c '
  ANSI_RED="\033[31m"
  ANSI_GREEN="\033[32m"
  ANSI_YELLOW="\033[33m"
  ANSI_BLUE="\033[34m"
  ANSI_MAGENTA="\033[35m"
  ANSI_GRAY="\033[90m"
  ANSI_CYAN="\033[36;1m"
  ANSI_DARKCYAN="\033[36m"
  ANSI_NOCOLOR="\033[0m"
}

disable_color() { unset ENABLECOLOR ANSI_RED ANSI_GREEN ANSI_YELLOW ANSI_BLUE ANSI_MAGENTA ANSI_CYAN ANSI_DARKCYAN ANSI_NOCOLOR; }

enable_color

print_start() {
  if [ "x$2" != "x" ]; then
    COL="$2"
  elif [ "x$BASE_COL" != "x" ]; then
    COL="$BASE_COL"
  else
    COL="$ANSI_MAGENTA"
  fi
  printf "${COL}${1}$ANSI_NOCOLOR\n"
}

gstart () {
  print_start "$@"
}
gend () {
  :
}

if [ -n "$GITHUB_EVENT_PATH" ]; then
  export CI=true
fi

[ -n "$CI" ] && {
  gstart () {
    printf '::[group]'
    print_start "$@"
    SECONDS=0
  }

  gend () {
    duration=$SECONDS
    echo '::[endgroup]'
    printf "${ANSI_GRAY}took $(($duration / 60)) min $(($duration % 60)) sec.${ANSI_NOCOLOR}\n"
  }
} || echo "INFO: not in CI"

#--

cd $(dirname $0)/..

if [ -z "$CI" ]; then
  if [ -f issue-runner ]; then
    rm -rf issue-runner
  fi
  cd tool
  go build -o ../issue-runner
  cd ..
fi

set +e

for t in \
  `ls ./__tests__/md/*.md` \
  'https://raw.githubusercontent.com/1138-4EB/issue-runner/master/test/vunit_py.md' \
  'VUnit/vunit#337' \
  'ghdl/ghdl#579' \
  'ghdl/ghdl#584' \
; do
  gstart "[test] $t"
  ./issue-runner -y -c "$t"
  gend
done

gstart "[test] mixed"
  ./issue-runner -y -c \
  'https://raw.githubusercontent.com/1138-4EB/issue-runner/master/test/vunit_py.md' \
  vunit_sh.md \
  'VUnit/vunit#337'
gend

gstart "[test] stdin cat hello001"
  cat ./__tests__/md/hello001.md | ./issue-runner -
gend

gstart "[test] stdin cat hello003"
  cat ./__tests__/md/hello003.md | ./issue-runner -y -
gend

gstart "[test] stdin < hello003"
  ./issue-runner -y - < ./__tests__/md/hello003.md
gend