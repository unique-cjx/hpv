# hpv
hpv疫苗

_用于发布约苗平台订阅人数较少的社区到qq群，qq机器人发布消息用的是[go-cqhttp包](https://raw.githubusercontent.com/Mrs4s/go-cqhttp/master/README.md)_

> 目前只支持抢苗平台后期会扩展其他平台

- 初始启动项目根目录下会生成个`depart_data.json`记录已经发送过的社区
- 如果要修改qq群号或者抢苗的城市，在`config`目录下的配置文件自行修改
- 程序运行是会启动三个常驻goroutine，功能分别是
    1. 检查到低于订阅人数的社区发布消息到通道
    2. 获取通道的最新消息发布到qq群
    3. 定时更新`wxtoken`避免获取平台的数据失败
