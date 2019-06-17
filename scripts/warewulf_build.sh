#!/bin/sh

source /usr/local/bin/common.sh
make_data_dirs
load_configs

rsync -av /buildtime-statedir/ /data/
wwinit database

warewulf-sync build
