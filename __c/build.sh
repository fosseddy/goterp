#!/bin/bash

set -xe

files="main.c scanner.c parser.c lib/mem.c lib/os.c"
outname=main

flags="-g -Werror=declaration-after-statement -Wall -Wextra -pedantic -std=c99"
incl=
libs=

if [[ $1 = "prod" ]]; then
    flags=${flags/-g/-O2}
fi

gcc $flags -o $outname $files $incl $libs
