#!/bin/bash
mkdir -p $INPUT_DATA_PATH # create input data directory
wget --directory-prefix=$INPUT_DATA_PATH $INPUT_DATA_URL # download input data into the input data directory
(cd $INPUT_DATA_PATH && eval $INPUT_DATA_TRANSFORM_COMMAND) # execute the transform command in the input data directory
mkdir -p $OUTPUT_DATA_PATH # create the output data directory
(cd $REPOSITORY_PATH && eval $REPOSITORY_RUN_COMMAND) # execute the run command in the repo root directory