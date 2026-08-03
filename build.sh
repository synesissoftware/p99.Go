#! /bin/bash

ScriptPath=$0
Dir=$(cd $(dirname "$ScriptPath"); pwd)


mkdir -p "$Dir/scratch/build-artefacts/"


# library
echo "building library ..."
go build -v "$Dir"


# examples
if [ -d "$Dir/examples" ];then

    echo "building examples ..."
    for example_dir in "$Dir"/examples/*/; do
        if [ -f "${example_dir}main.go" ]; then
            example_name=$(basename "$example_dir")
            echo "building example ${example_name} ..."
            go build -v -o "$Dir/scratch/build-artefacts/${example_name}" "$example_dir"
        fi
    done
fi


# tests
if [ -d "$Dir/tests" ];then

    echo "building tests ..."
    for test_dir in "$Dir"/tests/*/; do
        if [ -f "${test_dir}main.go" ]; then
            test_name=$(basename "$test_dir")
            echo "building test ${test_name} ..."
            go build -v -o "$Dir/scratch/build-artefacts/${test_name}" "$test_dir"
        fi
    done
fi


# ############################## end of file ############################# #
