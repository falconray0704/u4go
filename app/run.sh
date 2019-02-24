#!/bin/bash

set -o nounset
#set -o errexit

#set -x

run_native_func()
{

    ./outBin
    #./outBin -c="./configs/appCfgs.yaml"
    echo "ret:$?"

}

run_native_cmd_test_func()
{
    set -x

    ./outBin -c="./configs/appCfgs.yaml" -word=opt -numb=7 -fork -svar=flag
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -word=opt
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -word=opt a1 a2 a3
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -word=opt a1 a2 a3 -numb=7
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -h
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -wat
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -loop=true -word=opt a1 a2 a3 looping
    echo "ret:$?"

    set +x
}

run_docker_func()
{
    docker run --rm -v $(pwd)/logsData:/myApp/logsData myapp:falcon
}

run_clean_docker_datas_func()
{
    rm -rf logsData/*
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

mkdir -p ./logsData

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



