# cat-agent

## 简介
cat-agent是点评CAT server的一个agent。

## 原理
传统的PHP应用无法复用socket连接，直接接入CAT的话可能会将cat server连接数打高。另外接入cat在进行rpc调用前需要生成全局唯一的messageId，PHP要生成的话必须借助redis或其他类型的数据库来共享自增id，请求量高的情况下对数据库的压力也非常大。其次PHP也无法在本地进行cat log的聚合。

为了解决上述的问题，cat-agent通过在PHP应用所在机器上启动一个agent进程，PHP进程通过unix domain socket与agent通信来获取messgeId和发送cat log。cat log发送到agent之后并不是直接发出，而在在本地聚合之后再发出，降低了对cat server的请求次数。agent与cat server之间的连接被复用，降低连接数。

![cat-agent架构](https://github.com/Orlion/cat-agent/blob/main/resources/architecture.png)

## 使用
### 1. 启动agent进程
1. 下载源码
```
$ git clone 
```
2. 进入源码目录，编译为可执行文件
```
$ cd cat-agent
$ go build -o cat-agent
```
3. 修改默认配置，指定domain与cat server的地址
```
$ vim cat-agent.conf.yml
```
4. 运行可执行文件，启动agent进程
```
$ ./cat-agent -conf=cat-agent.conf.yml
```
### 2. PHP接入
请通过composer安装PHP客户端：[github.com/Orlion/cat-agent-php](https://github.com/Orlion/cat-agent-php)完成应用接入

## 其他方案
通过扩展的方式让PHP接入cat可能是性能更好的方案，大部分公司采用的应该都是这种方案。