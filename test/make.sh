#!/bin/bash

./build.sh

#
# test 1
#

echo "{\"username\": \"24D8101308H\", \"password\": \"fykcFSVQpNH6@y2\"}" > settings.json
echo "y" | ./kadai make 01 a b x > /dev/null

# check directory
out=$(ls -R 01)
expected="01:
inputFiles
kadai01a.c
kadai01b.c
kadai01x.c

01/inputFiles:
inputa1.txt
inputb1.txt
inputx1.txt"
if [ "$out" = "$expected" ]; then
    echo OK
else 
    echo NO
    echo "expected: '$expected'"
    echo "out: '$out'"
fi

#
# test 2
#

