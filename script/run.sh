#!/usr/bin/env bash

set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "run project"
    echo
    echo "Usage:"
    echo "  $Command [flags]"
    echo
    echo "Flags:"
    echo "  -b, --build         build all before running"
    echo "  -c, --code          build go code before running"
    echo "  -g, --grpc          build grpc before running"
    echo "  -s, --static        build static before running"
    echo "  -h, --help          help for $Command"
}

ARGS=`getopt -o hbcgs --long help,build,code,grpc,static -n "$Command" -- "$@"`
eval set -- "${ARGS}"
build=0
grpc=0
static=0
code=0
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
        ;;
        -b|--build)
            build=1
            shift 1
        ;;
        -c|--code)
            code=1
            shift 1
        ;;
        -s|--static)
            static=1
            shift 1
        ;;
        -g|--grpc)
            grpc=1
            shift 1
        ;;
        --)
            shift
            break
        ;;
        *)
            echo Error: unknown flag "$1" for "$Command"
            echo "Run '$Command --help' for usage."
            exit 1
        ;;
    esac
done

if [[ $build != 0 ]];then
    "$BashDir/grpc.sh"
    "$BashDir/document.sh"
    "$BashDir/go.sh"
else
    if [[ $grpc != 0 ]];then
        "$BashDir/grpc.sh"
    fi
    if [[ $static != 0 ]];then
        "$BashDir/document.sh"
    fi
    if [[ $code != 0 ]];then
        "$BashDir/go.sh"
    fi
fi
cd "$Dir/bin"
args=(
    ./"$Target" web -d --no-upgrade
)
exec="${args[@]}"
echo $exec
eval "$exec"
