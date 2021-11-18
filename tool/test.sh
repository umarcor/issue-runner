#!/usr/bin/env sh

set -e

enable_color() {
  export ANSI_RED="\033[31m"
  export ANSI_GREEN="\033[32m"
  export ANSI_YELLOW="\033[33m"
  export ANSI_BLUE="\033[34m"
  export ANSI_MAGENTA="\033[35m"
  export ANSI_GRAY="\033[90m"
  export ANSI_CYAN="\033[36;1m"
  export ANSI_DARKCYAN="\033[36m"
  export ANSI_NOCOLOR="\033[0m"
}

disable_color() { unset ANSI_RED ANSI_GREEN ANSI_YELLOW ANSI_BLUE ANSI_MAGENTA ANSI_CYAN ANSI_DARKCYAN ANSI_NOCOLOR; }

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
    printf '::group::'
    print_start "$@"
    SECONDS=0
  }

  gend () {
    duration=$SECONDS
    echo '::endgroup::'
    printf "${ANSI_GRAY}took $((duration / 60)) min $((duration % 60)) sec.${ANSI_NOCOLOR}\n"
  }
} || echo "INFO: not in CI"

#--

cd $(dirname "$0")/..

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
  $(ls ./tests/md/empty*.md) \
; do
  gstart "[test] $t"
  ./issue-runner -y -c "$t"
  exitcode=$?
  gend
  if [ $exitcode -eq 1 ]; then
    printf "${ANSI_GREEN}success${ANSI_NOCOLOR}\n"
  else
    printf "${ANSI_RED}failure [${exitcode}]${ANSI_NOCOLOR}\n"
  fi
done

for t in \
  $(ls ./tests/md/hello*.md) \
; do
  gstart "[test] $t"
  ./issue-runner -y -n -c "$t"
  exitcode=$?
  gend
  if [ $exitcode -eq 0 ]; then
    printf "${ANSI_GREEN}success${ANSI_NOCOLOR}\n"
  else
    printf "${ANSI_RED}failure [${exitcode}]${ANSI_NOCOLOR}\n"
  fi
done

#for t in \
#  `ls ./tests/md/attach*.md` \
#  'https://raw.githubusercontent.com/eine/issue-runner/master/test/vunit-py.md' \
#  'VUnit/vunit#337' \
#  'ghdl/ghdl#579' \
#  'ghdl/ghdl#584' \
#; do
#  gstart "[test] $t"
#  ./issue-runner -y -c "$t"
#  exitcode=$?
#  gend
#  if [ $exitcode -eq 0 ]; then
#    printf "${ANSI_GREEN}success${ANSI_NOCOLOR}\n"
#  else
#    printf "${ANSI_RED}failure [${exitcode}]${ANSI_NOCOLOR}\n"
#  fi
#done
#
#gstart "[test] mixed"
#  ./issue-runner -y -c \
#  'https://raw.githubusercontent.com/eine/issue-runner/master/test/vunit-py.md' \
#  vunit-sh.md \
#  'VUnit/vunit#337'
#gend
#
#gstart "[test] stdin cat hello001"
#  cat ./tests/md/hello001.md | ./issue-runner -
#gend
#
#gstart "[test] stdin cat hello003"
#  cat ./tests/md/hello003.md | ./issue-runner -y -
#gend
#
#gstart "[test] stdin < hello003"
#  ./issue-runner -y - < ./tests/md/hello003.md
#gend
#
