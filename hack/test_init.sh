#!/bin/bash

# variables instead of env vars for now

input_data_path=example/data # location the input data gets downloaded into, defined by user
input_data_url=https://huggingface.co/datasets/aliabd/hello-world/raw/main/data.csv # URL from which to download the input data, defined by user
transfrom_cmd="echo 'hello'" # command to run after downloading the input data (e.g. for unpacking), defined by user

repo_path=example # location of the source code, hardcoded
run_cmd="python3 example.py" # command to run in order to execute the program, defined by user

output_data_path=example/output # location that holds the output data after the execution, defined by user

# setup code for initial testing, remove later
mkdir $repo_path
echo "print('hello world')
with open('./data/data.csv', 'r') as f:
    with open('./output/output.csv', 'w') as f2:
        f2.write('hello world')
        for line in f:
            print(line)
            f2.write(line)
     " > $repo_path/example.py

mkdir -p $input_data_path # create input data directory
wget --directory-prefix=$input_data_path $input_data_url # download input data into the input data directory
(cd $input_data_path && eval $transfrom_cmd) # execute the transform command in the input data directory
mkdir -p $output_data_path # create the output data directory
(cd $repo_path && eval $run_cmd) # execute the run command in the repo root directory
