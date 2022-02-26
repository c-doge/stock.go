#!/bin/bash

CURRENT_PATH="$(cd `dirname $0`; pwd)"

VERSION_GO="`dirname $CURRENT_PATH`/cmd/version.go"

COMMIT_HASH=`git log -1 --pretty=format:%H`

COMMIT_HASH=${COMMIT_HASH:0:16}

BRANCH_NAME=`git symbolic-ref HEAD 2>/dev/null | cut -d'/' -f 3`

BUILD_TIME=`date "+%Y-%m-%d %H:%M:%S"`

#echo "path: $current_path"
#echo "commit: $commit_hash"
#echo "version: $VERSION_GO"
#echo "git branch: $BRANCH_NAME"
#echo "build time: $BUILD_TIME"

echo -e "package main\n" > $VERSION_GO
echo -e "var _buildTime string = \"$BUILD_TIME\";\n" >> $VERSION_GO
echo -e "var _gitBranch string = \"$BRANCH_NAME-$COMMIT_HASH\";\n" >> $VERSION_GO
echo -e "var _version string = \"v0.0.1\";\n" >> $VERSION_GO
