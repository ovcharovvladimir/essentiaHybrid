#!/bin/bash

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
if [[ $PWD = *"travis"* ]];
then
    solution="$GOPATH/go/src"
    #"${PWD%/go/src/*}"
else
    solution="${PWD%/go/src/*}/go/src"
fi

#"$GOPATH"
root="$PWD"
echo "GOPATH: $GOPATH"
echo "Project Working Dir: $PWD"
echo "Workspace: $workspace"
echo "Solution: $solution"
dir="$workspace/src/github.com/ovcharovvladimir"


#"$GOPATH"
root="$PWD"
echo "GOPATH: $GOPATH"
echo "Project Working Dir: $PWD"
echo "Workspace: $workspace"
echo "Solution: $solution"
dir="$workspace/src/github.com/ovcharovvladimir"
 



if [ ! -L "$dir/essentiaHybrid" ]; then
    mkdir -p "$dir"
    cd "$dir"
    ln -s ../../../../../. essentiaHybrid
    cd "$root"
fi
pth="$workspace/src/github.com/"
#github.com/dancannon/gorethink
if [ ! -L "$pth/dancannon" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/dancannon/. dancannon
    cd "$root"
fi
#github.com/sirupsen/logrus
pth="$workspace/src/github.com/"

if [ ! -L "$pth/sirupsen" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/sirupsen/. sirupsen
    cd "$root"
fi
#github.com/x-cray/logrus-prefixed-formatter
pth="$workspace/src/github.com"

if [ ! -L "$pth/x-cray" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/x-cray/. x-cray
    cd "$root"
fi
#
#github.com/cenkalti/backoff
pth="$workspace/src/github.com"

if [ ! -L "$pth/cenkalti" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/cenkalti/. cenkalti
    cd "$root"
fi
#github.com/opentracing/opentracing-go
pth="$workspace/src/github.com"

if [ ! -L "$pth/opentracing" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/opentracing/. opentracing
    cd "$root"
fi
#gopkg.in/rethinkdb
pth="$workspace/src/gopkg.in"

if [ ! -L "$pth/rethinkdb" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/gopkg.in/rethinkdb/. rethinkdb
    cd "$root"
fi
#gopkg.in/fatih/pool.v2
pth="$workspace/src/gopkg.in"

if [ ! -L "$pth/fatih" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/gopkg.in/fatih/. fatih
    cd "$root"
fi
#gopkg.in/rethinkdb/rethinkdb-go.v5/encoding
pth="$workspace/src/gopkg.in"

if [ ! -L "$pth/rethinkdb" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/gopkg.in/rethinkdb/rethinkdb-go.v5/. rethinkdb-go.v5
    cd "$root"
fi
#golang.org/x/sys
pth="$workspace/src/golang.org"

if [ ! -L "$pth/x" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/golang.org/x/. x
    cd "$root"
fi
#github.com/mgutz/ansi
pth="$workspace/src/github.com"

if [ ! -L "$pth/mgutz" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/mgutz/. mgutz
    cd "$root"
fi
#github.com/mattn/go-colorable
pth="$workspace/src/github.com"

if [ ! -L "$pth/mattn" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/mattn/. mattn
    cd "$root"
fi
#github.com/golang/protobuf
pth="$workspace/src/github.com"

if [ ! -L "$pth/golang" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/golang/. golang
    cd "$root"
fi
#github.com/hailocab/go-hostpool
pth="$workspace/src/github.com"

if [ ! -L "$pth/hailocab" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/hailocab/. hailocab
    cd "$root"
fi
#github.com/pkg/profile
pth="$workspace/src/github.com"

if [ ! -L "$pth/pkg" ]; then
    mkdir -p "$pth"
    cd "$pth"
    ln -s $solution/github.com/pkg/. pkg
    cd "$root"
fi



GOPATH="$workspace"
export GOPATH 


# Run the command inside the workspace.
cd "$dir/essentiaHybrid"
PWD="$dir/essentiaHybrid"
echo "----- $@"
# Launch the arguments with the configured environment.
exec "$@"
