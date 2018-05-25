#! /bin/sh

cd $(dirname $0)

# Create temporal dir
idir="$(echo "$1" | sed -e 's/\//--/g' | sed -e 's/#/--/g')"
mkdir -pv "$idir" && cd "$idir"

# Get extract_sources.py script
[ -f ../extract_sources.py ] &&  cp ../extract_sources.py ./ || \
  curl -o extract_sources.py -L https://github.com/1138-4EB/issue-runner/raw/master/extract_sources.py

# Detect python, fall back to docker python:slim-stretch
export py="$(command -v python)"
if [ "$py+x" = "+x" ]; then
  export py="$(command -v python3)"
  if [ "$py+x" = "+x" ]; then
    export py="docker run --rm -tv /$(pwd):/src -w //src python:slim-stretch python"
  fi
fi

# Execution command
cmd="$py extract_sources.py $1"
echo "cmd: $cmd"

# Extract sources
$cmd

# Prepare sim.sh
if [ -f "sim.sh" ]; then
  sed -i.bak 's/\r$//g' sim.sh
  chmod +x sim.sh
fi

# Prepare and execute run.sh
if [ -f "run.sh" ]; then
  sed -i.bak 's/\r$//g' run.sh
  chmod +x run.sh
  ./run.sh
  exit 0
fi

echo "'run.sh' not found!"
exit 1
