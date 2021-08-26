# hpv
hpv疫苗发布订阅

_发布平台订阅人数较少的社区到指定的qq群，qq机器人发布消息用的是[go-cqhttp包](https://raw.githubusercontent.com/Mrs4s/go-cqhttp/master/README.md)_

> 目前只支持抢苗平台后期会扩展其他平台
- 初始启动程序，根目录下会生成个`depart_data.json`用于记录已发送过的社区
- 如果要修改发布消息的qq群或抢苗的城市，请自行修改配置文件
- 程序运行时会启动三个tasks，功能分别是
    1. *dispatch_mess_task*: 检查到低于订阅人数的社区发布消息到通道 
    2. *send_mess_task*: 获取通道的最新消息发布到qq群
    3. *refresh_mess_task*: 定时更新`wxtoken`避免获取平台的数据失败
