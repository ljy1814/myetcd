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
