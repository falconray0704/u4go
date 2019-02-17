#!/bin/bash

set -o nounset
#set -o errexit

#set -x

run_native_func()
{
    ./outBin -word=opt -numb=7 -fork -svar=flag
    echo "ret:$?"

    ./outBin -word=opt
    echo "ret:$?"

    ./outBin -word=opt a1 a2 a3
    echo "ret:$?"

    ./outBin -word=opt a1 a2 a3 -numb=7
    echo "ret:$?"

    ./outBin -h
    echo "ret:$?"

    ./outBin -wat
    echo "ret:$?"

    ./outBin -loop=true -word=opt a1 a2 a3 looping
    echo "ret:$?"
}

run_docker_func()
{
    docker run --rm -v $(pwd)/tmp:/myApp/tmp myapp:falcon
}

run_clean_docker_datas_func()
{
    rm -rf tmp/*
}

usage()
{
    echo "Run native:"
    echo "./run.sh lc"
    echo ""
    echo "Run docker:"
    echo "./run.sh dk"
    echo ""
    echo "Run clean datas:"
    echo "./run.sh clean"
}

[ $# -lt 1 ] && usage && exit

mkdir -f ./tmp

case $1 in
    lc) echo "Run native..."
        run_native_func
        ;;
    dk) echo "Run in docker..."
        run_docker_func
        ;;
    clean) echo "Clean datas..."
        run_clean_docker_datas_func
        ;;
    *) echo "Unknown command!"
        usage
        ;;
esac



