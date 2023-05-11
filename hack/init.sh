#!/bin/bash

# variables instead of env vars for now
repo_path=example
data_path=data
output_path=output
url=https://huggingface.co/datasets/aliabd/hello-world/raw/main/data.csv
transfrom_cmd="echo 'hello'"
run_cmd="python3 example.py"

# git pull alternative for now
mkdir $repo_path
cd $repo_path
echo "print('hello world')
with open('data/data.csv', 'r') as f:
    with open('output/output.csv', 'w') as f2:
        f2.write('hello world\n')
        for line in f:
            print(line)
            f2.write(line)
     " > example.py
# data is our data path, for now, will be replaced
mkdir -p $data_path
absolute_repo_path=$(pwd)

cd $data_path
# replace with url env var
wget $url
# eval should be filled with the command to run the code (transfrom command)
eval $transfrom_cmd
cd $absolute_repo_path
mkdir -p $output_path
eval $run_cmd
