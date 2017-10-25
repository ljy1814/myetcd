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


##### 接口分析
- registerDevice
    - 开发者调用,注册设备
    - 普通用户调用则绑定设备
- registerDevices
    - 产品维度的公钥私钥
    - 产品信息
    - device_warehouse  // 注册
    - device_record // 激活
- getPublicKey
    - key NullString
- getDeviceLogicalId
    - device_atom　//　存逻辑id,物理id
- deviceHeartbeat
    - device_record
- activateDevice
    - event 打印到日志,而后再由日志处理器处理,发给消息队列
    - 短链激活,内部激活(不分配token)


内部激活
    设备不存在
        新增一条激活数据
    设备已存在
        更新设备激活数据

外部激活
    短链,校验签名,设备授权
    设备尚未被激活
        激活不存在的设备,分配token
    设备曾经激活过
        设备是长连接激活则token不存在
        token存在则直接返回token
        否则执行激活操作,分配token

反激活
    修改激活状态

设备信息
    设备记录与属性

设备授权
    根据产品的授权模式分别验证
    白名单
        检查设备gid是否在白名单里面
        不是则检查设备的状态是否为导入,device_record,激活更改状态
        第一次认证但未导入白名单
更改设备属性
    设备自身通过inner service更新
    开发者更新

设备下线
    更新心跳


#### 重点模块
- 入库
    - registerDevice
        - 直接卖出没有激活
    - registerDevices
    - listStock
    - stockStatistics
    - applyRegistDevice
    - listDevices
- 设备激活
    - activateDevice
    - antiactivateDevice
- 密钥管理
    - generateRsaKey
    - generateAesKey
    - registerUnifiedKey
    - getUnifiedKey
- token
    - getToken 
    - getAccessToken
    - getRefreshToken
    - updateAccessToken
- 配额
    - createLicenseQuota
    - deleteLicenseQuota
    - getLicenseQuota
    - modifyQuota
    - getLicenseQuotaHistory
    - getLicenseQuotaHistoryCount
- 设备授权
    - authDeviceLicense
- 设备属性
    - getDeviceProperty
    - updateDeviceProperty

#### 数据库表
- 主域相关
    - domain_device_record
    - domain_device_warehouse
    - domain_device_atom
    - domain_device_licenses
    - domain_device_property
- 配额相关
    - warehouse_licenses_quota
    - warehouse_license_req_id
    - warehouse_license_quota_history
- 其它
    - warehouse_device_meta
    - warehouse_product_status
