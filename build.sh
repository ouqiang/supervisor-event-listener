#!/usr/bin/env bash

# set -x -u
# 构建应用, supervisor-event-listener.tar.gz
# ./build.sh -p linux -a amd64
# 参数含义
# -p 指定平台(linux|darwin)
# -a 指定体系架构(amd64|386), 默认amd64


TEMP_DIR=`date +%s`-temp-`echo $RANDOM`

# 目标平台 linux,darwin
OS=''
# 目标平台架构
ARCH=''
# 应用名称
APP_NAME='supervisor-event-listener'
# 可执行文件名
EXEC_NAME=''
# 压缩包名称
COMPRESS_FILE=''


# -p 平台 -a 架构
while getopts "p:a:" OPT;
do
    case $OPT in
        p) OS=$OPTARG
        ;;
        a) ARCH=$OPTARG
        ;;
    esac
done

if [[ -z  $OS ]];then
    echo "平台不能为空"
    exit 1
fi

if [[ $OS && $OS != 'linux' && $OS != 'darwin' ]];then
    echo '平台错误，支持的平台 linux darmin(osx)'
    exit 1
fi

if [[ -z $ARCH ]];then
    ARCH='amd64'
fi

if [[ $ARCH != '386' && $ARCH != 'amd64' ]];then
    echo 'arch错误，仅支持 386 amd64'
    exit 1
fi

echo '开始编译'

GOOS=$OS GOARCH=$ARCH go build -ldflags '-w'

if [[ $? != 0 ]];then
    exit 1
fi
echo '编译完成'


EXEC_NAME=${APP_NAME}
COMPRESS_FILE=${APP_NAME}-${OS}-${ARCH}.tar.gz

mkdir -p $TEMP_DIR/$APP_NAME
if [[ $? != 0 ]]; then
    exit 1
fi

# 需要打包的文件
PACKAGE_FILENAME=(supervisor-event-listener.ini ${EXEC_NAME})

echo '复制文件到临时目录'
# 复制文件到临时目录
for i in ${PACKAGE_FILENAME[*]}
do
    cp -r $i $TEMP_DIR/$APP_NAME
done


echo '压缩文件'
# 压缩文件
cd $TEMP_DIR

tar czf $COMPRESS_FILE *
mv $COMPRESS_FILE ../
cd ../

rm $EXEC_NAME
rm -rf $TEMP_DIR

echo '打包完成'
echo '生成压缩文件--' $COMPRESS_FILE