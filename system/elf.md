### Linux ELF 文件

* 分类

  * 可执行文件
  * 共享目标文件(.so)
  * 可重定位文件(.o)
  * 系统core dump 文件

* 文件结构

  * 分段(.text(代码段),.bss(全局未定义变量或未定义静态变量),.data(全局定义变量和局部变量))

    ![image-20211222215943751](/home/geekwu/.config/Typora/typora-user-images/image-20211222215943751.png)

  * 分段的目的

    * 共享地址减少内存使用
    * 9/1原则方便增大程序命中率
    * 划分不同的区域权限(.text可执行可读，.data可读可写)

  * 段构成

  * 各种表