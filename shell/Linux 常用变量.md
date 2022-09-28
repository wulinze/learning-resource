### Linux 常用变量

* $0：返回当前执行的命令名，即第一个参数
* $n：对应命令的第几个参数
* $#：命令参数的个数，即执行命令的时候后面跟着的参数的个数
* $*：命令行的所有参数，参数作为一个整体
* $@：命令行的所有参数，参数作为一个一个的个体
* $?：返回上一条命令的退出状态
* $$：返回当前进程的进程号

* $!：返回最后一个后台进程的进程号

~~~shell
printf "The complete list is %s\n" "$$"
printf "The complete list is %s\n" "$!"
printf "The complete list is %s\n" "$?"
printf "The complete list is %s\n" "$*"
printf "The complete list is %s\n" "$@"
printf "The complete list is %s\n" "$#"
printf "The complete list is %s\n" "$0"
printf "The complete list is %s\n" "$1"
printf "The complete list is %s\n" "$2"
~~~

结果：

The complete list is 24601
The complete list is 
The complete list is 0
The complete list is 
The complete list is 
The complete list is 0
The complete list is test.sh
The complete list is 
The complete list is 

