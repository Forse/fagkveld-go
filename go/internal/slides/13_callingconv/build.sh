#!/bin/sh

set -ex

go build -o internal/slides/13_callingconv/bin/goadd internal/slides/13_callingconv/main.go && objdump -M intel -d internal/slides/13_callingconv/bin/goadd -S > internal/slides/13_callingconv/bin/goadd.asm
gcc -O3 -g3 -o internal/slides/13_callingconv/bin/cadd internal/slides/13_callingconv/main.c && objdump -M intel -d internal/slides/13_callingconv/bin/cadd -S > internal/slides/13_callingconv/bin/cadd.asm

./internal/slides/13_callingconv/bin/goadd
./internal/slides/13_callingconv/bin/cadd
