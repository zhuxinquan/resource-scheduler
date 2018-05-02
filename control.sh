#!/bin/bash
workspace=$(cd $(dirname $0) && pwd)
cd $workspace

app=resource-scheduler-agent
conf=config/conf.cfg
pidfile=var/app.pid
logfile=log/app.log
gitversion=.gitversion

mkdir -p var &>/dev/null
mkdir -p log &>/dev/null

## opt
function build() {
    # 设置golang环境变量
    echo -e "`go version`"
    export GOPATH=$GOPATH:$workspace

    go build -o $app main.go
    local lcode=$?
    if [[ "$lcode" == "0" ]]; then
        echo "build ok"
    else
        echo "build error"
    fi
}

function start() {
    nohup ./$app -c $conf >>$logfile 2>&1 &
    echo $! > $pidfile
    echo "start ok, pid=$!"
}

function stop() {
    kill `get_pid`
    echo "stoped"
}

function dev() {
    stop
    export GOPATH=$GOPATH:$workspace/Godeps/_workspace:$workspace
    go fmt ./src/... #格式化代码
    git log -1 --pretty=%h > gitversion
    go run main.go -c  config/dev.cfg $1
}

#function _godep() {
#    # 将当前依赖的代码文件全部拷贝至Godeps的workspace中, 并进行独立提交
#    export GOPATH=$GOPATH:$workspace
#    godep save .
#}

function shutdown() {
    pid=`get_pid`
    kill -9 $pid
    echo "stoped"
}

function restart() {
    stop
    start
}

## other
function status() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo -n "running, pid="
        cat $pidfile
    else
        echo "stoped"
    fi
}

function version() {
    ./$app -v
}

function tailf() {
    tail -f $logfile
}

##### 对于首次部署或新迁移的程序，需要注意此函数
function get_pid() {
    if [ -f $pidfile ];then
        cat $pidfile
    else
        # 如果不存在pid文件，有可能是在其他地方部署过, 方便其他程序迁移至此地
        # 第一次迁移，需要关注如何填写
        pid=`ps aux | grep $app | grep $conf | awk '{print $2}'`
        if [ "x_$pid" != "x_" ]; then
            echo $pid > $pidfile
            cat $pidfile
        fi
    fi
}

## internal
function check_pid() {
    pid=`get_pid`
    if [ "x_" != "x_$pid" ]; then
        running=`ps -p $pid|grep -v "PID TTY" |wc -l`
        return $running
    fi
    return 0
}

## usage
function usage() {
    echo "$0 build|pack|start|stop|restart|status|tail|version|dev"
}

## main
action=$1
case $action in
    ## build
    "build" )
        build
        ;;
    ## opt
    "start" )
        start
        ;;
    "stop" )
        stop
        ;;
    "kill" )
        shutdown
        ;;
    "restart" )
        restart
        ;;
    ## other
    "status" )
        status
        ;;
    "version" )
        version
        ;;
    "tail" )
        tailf
        ;;
    "dev" )
        dev $2
        ;;
    "pid" )
        get_pid
        ;;
    * )
        usage
        ;;
esac