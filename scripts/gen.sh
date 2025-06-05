#!/bin/bash

# 获取脚本名称（可能包含路径）
script_name="$0"

# 获取脚本所在的目录的绝对路径
root_dir="$(dirname $(dirname "$(readlink -f "$0")"))"

echo "root_dir: ${root_dir} 参数：$1 $2"
# 判断参数
if [ $# -ne 2 ]; then
  echo "Usage: gen.sh <template> <name>"
  return
fi
template=$1
name=$2
cd $root_dir

if [ -f ./app ]; then
  ./app gen $1 $2
else
 go run cmd/main.go gen $1 $2
fi
