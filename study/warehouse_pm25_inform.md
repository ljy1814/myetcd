#### warehouse分析

##### 相关模块
- 调用模块
    - platform
    - product
- 调用warehouse的模块
    - bind
        - registerDevice
        - getDeviceLogicalId
        - antiActivateDevice
    - gateway
        - getDeviceLogicalId
        - getDeviceInfo
        - deviceHeartbeat
        - getPublicKey
        - authDeviceLicense
        - activateDevice
    - redirect
        - getDeviceCount
        - getDeviceInfo
    - router
        - getToken

#### 重点模块
- 入库
    - registerDevice
        - 直接卖出没有激活
        - 开发者调用,注册设备,bind调用用于开发者绑定设备
        - 普通用户调用则绑定设备
    - registerDevices
        - 产品维度的公钥私钥
        - 产品信息
        - device_warehouse  // 密钥授权
        - device_record // 激活
    - listStock
        - 设备列表,包括key,record
    - stockStatistics
        -统计信息,包括设备数量,配额,激活数量
    - applyRegistDevice
        - 获取批次,弃用
    - listDevices
        - 列出设备,包含在线状态
- 设备激活
    - activateDevice
        - event 打印到日志,而后再由日志处理器处理,发给消息队列
        - 短链激活,内部激活(不分配token)
        - 内部激活,gateway
            - 设备不存在
                - 新增一条激活数据
            - 设备已存在
                - 更新设备激活数据

        - 外部激活
            - 短链,校验签名,设备授权
            - 设备尚未被激活
                - 激活不存在的设备,分配token
            - 设备曾经激活过
                - 设备是长连接激活则token不存在
                - token存在则直接返回token
                - 否则执行激活操作,分配token
    - antiactivateDevice
        - 设备不存在或者未激活则不需要反激活
        - 保留认证状态,修改激活状态
- 密钥管理
    - generateRsaKey
        - 生成rsa公钥,烧进固件,用于gateway四次握手
    - generateAesKey
        - 生成aes公钥
    - registerUnifiedKey
        - 注册产品统一密钥
    - getUnifiedKey
        - 获取产品统一密钥
- token
    - getToken 
        - 获取access,refresh token
    - getAccessToken
        - 获取access token
    - getRefreshToken
        - 获取refresh token
    - updateAccessToken
        - 更新access token
- 配额
    - createLicenseQuota
        - 创建配额表并初始化
    - deleteLicenseQuota
        - 删除配额数据
    - getLicenseQuota
        - 获取配额
    - modifyQuota
        - 修改配额,记录历史
    - getLicenseQuotaHistory
        - 配额历史记录
    - getLicenseQuotaHistoryCount
        - 配额历史记录数量
- 设备授权
    - authDeviceLicense
        - 设备授权
            - 根据产品的授权模式分别验证
            - 白名单
                - 检查设备gid是否在白名单里面
                - 不是则检查设备的状态是否为导入,device_record,激活更改状态
                - 第一次认证但未导入白名单
- 设备属性
    - getDeviceProperty
    - updateDeviceProperty

#### 数据库表
- 主域相关
    - domain_device_record
        - 设备记录
    - domain_device_warehouse
        - 设备身份鉴权信息
    - domain_device_atom
        - 设备的物理ID与逻辑Id映射
    - domain_device_property
- 配额相关
    - warehouse_licenses_quota
    - warehouse_license_req_id
    - warehouse_license_quota_history
- 其它
    - warehouse_device_meta
    - warehouse_product_status



#### pm25

##### cron
- 每20分钟启动一次
- 爬虫接口,获取最新数据

##### 接口列表
- getLatestData
- getLatestPM25Data
- getLatestAqi
- getLastDaysAqi 
- getLastHoursAqi 
- getHistoryData 

- 台阶式数据列表

- 获取最近几天的数据,默认７天
- getLastDaysData 
- getLastDaysPM25Data 

- 获取最近几小时的数据
- getLastHoursData
- getLastHoursPM25Data 

- getLatestWeatherData
- getLastDaysWeatherData
- getLastHoursWeatherData


##### 实际接口
- PM25Handler
    - getLatestAqi
        - current_air_condition
    - getLatestData
        - current_air_condition
    - getLastDays
        - history_air_condition
    - getLastDaysAqi
        - history_air_condition
    - getLastHours
        - history_air_condition
    - getLastHoursAqi
        - history_air_condition

    - getLatestWeatherData
        - current_weather_condition
    - getLastDaysWeatherData
        - history_weather_condition
    - getLastHoursWeatherData
        - history_weather_condition


```
{
    "history": [
        {
            "timestamp": "2017-10-28 09",
            "aqi": 0,
            "min": 0,
            "max": 0
        },
        {
            "timestamp": "2017-10-28 10",
            "aqi": 70,
            "min": 70,
            "max": 70
        },
        {
            "timestamp": "2017-10-28 11",
            "aqi": 65,
            "min": 65,
            "max": 65
        },
        {
            "timestamp": "2017-10-28 12",
            "aqi": 65,
            "min": 65,
            "max": 65
        },
        {
            "timestamp": "2017-10-28 13",
            "aqi": 59,
            "min": 59,
            "max": 59
        },
        {
            "timestamp": "2017-10-28 14",
            "aqi": 56,
            "min": 56,
            "max": 56
        },
        {
            "timestamp": "2017-10-28 15",
            "aqi": 56,
            "min": 56,
            "max": 56
        }
    ]
}
```


#### inform

##### sms
- sendSmsMessage,发送短信
    - 选择一个模板,选择一个发送工具,数据使用','分割
- sendSmsMessageByUserList
    - 批量发送短信
- sendSmsMessageByPhoneList
    - 根据手机号列表批量发送

- listSmsTemplate
    - 列出sms模板
- addSmsTemplate
    - 添加sms模板
        - 如果是案例,状态为AUDIT_APPROVED,否则以AUDIT_WAITING状态插入
        - 添加appID等信息
- updateSmsTemplate
    - 更新template信息
- deleteSmsTemplate
    - 删除template

- addSmsInfo
    - if AUDIT_APPROVED, appId, templateId,添加sms信息 
    - else 删除sms信息
    - 更新sms template的status,domain
- getSmsInfo
    - 获取sms和template信息

- listEmailTemplate
    - 列出所以邮件模板
- addEmailTemplate
    - 获取最大的可替换变亮个数
    - 获取最大的email template id
    - 添加email模板
- modifyEmailTemplate
    - 获取email template,检查其状态看是否已通过
    - 更新邮件模板
- deleteEmailTemplate
    - 获取email template,检查其状态看是否已通过,通过则禁止删除
    - 删除邮件模板

- modifyStatus
    - 修改邮件模板状态

- sendEmailByTemplate
    - 获取邮件模板,状态必须已通过模
    - 使用ESP怎挑选一个邮件发送器,否则使用自带的发送器
- sendEmail
    - 发送邮件
- sendEmailByBatch
    - 批量发送
- sendEmailByUserList
    - 批量发送,参数为userId

- reloadEspConfig
    - 重新加载ESP配置

- notifyDeviceUsers
    - 从bind获取用户列表
    - 根据手机类型选择ym或者xg推送
    - 推送消息
- notifyUsers
    - 通知多个用户
- notifyDeviceUsersByBatch
    - 批量通知用户

- setWhiteList
    - 友盟设置白名单

- addNotifyInfo
    - zc_notify_info添加数据
- listNotifyInfo
    - 列出notify accessId, key等信息
- deleteNotifyInfo
    - 删除notify数据
- getNotifyInfo
    - 获取notify数据


##### sql
- zc_inform_sms
```c
     domain: 971
   template: 0
   provider: yzx
     app_id: abcdef123456
template_id: abcdef123456
create_time: 2017-07-17 06:19:16
modify_time: 2017-07-17 06:19:16
```
- zc_inform_sms_template
```c
          domain: 20
              id: 1
   template_name: 示例模版
   template_type: 1
template_content: {1}（验证码），请注意保密，10分钟内有效，若非本人操作，请忽略。
   template_sign: 智云奇点
          status: 4
     create_time: 2016-01-14 09:00:12
```

##### 短信发送器
- ytx
- yzx

##### APP推送
- xg
- ym

##### 邮件发送器
- mailgun
- sendcloud
