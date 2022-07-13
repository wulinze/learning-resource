### TCP协议

#### TCP三次握手

* 握手过程

  * 第一次握手：Client将SYN置为1,随即生成一个seq发送给Server，进入***SYN_SENT***状态

  * 第二次握手：Server收到Client的SYN=1之后，知道客户端请求建立连接，将自己的SYN置为1,设置ACK=1，设置ack number=seq number+1,随即初始化自己的seq number之后发送给client。进入***SYN_RCVD***状态

  * 第三次握手：Clinet收到Server的ACK=1， seq number为之前发送的seq number+1之后，设置自己的ACK=1,之后设置ack number为server端传回来的seq number+1设为当前的ack number之后进入***ESTABLISHED***状态。完成三次握手。

    ![三次握手](/home/wu/learning-resource/net/picture/tcp握手.png)

* 为什么要三次握手

#### TCP四次挥手

* 为什么要四次挥手

#### TCP应用场景

#### TCP和UDP区别

#### TCP拥塞控制

