#### redirect

负载均衡
监控gateway,mqtt的meta
上一个接入点不为空则选择上一个,新接入的选一个最小的ap,非第一次接入则比较负载最小的节点和上一个节点的负载差距,小于20%或者上一个接入点负载小于50%则仍选上一个
接入点包括 service gateway两个机器
mqtt协议  tcp,tcps(tls加密),ws(websocket)

采用tcp连接,多线程模型,调用gateway的接口处理channel里面的socket
