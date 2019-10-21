# mediaConv
直播流转码服务，提供RTSP，RTMP 流媒体转码服务

---
接口API规范
---
**1\. 查询转码实例**
###### 接口功能
> 查询转码实例信息
###### URL
> /v1/transcoding
###### 支持格式
> JSON
###### HTTP请求方式
> GET
> |参数|必选|类型|说明|
|:-----  |:-------|:-----|-----                               |
|name    |ture    |string|请求的项目名                          |
|type    |true    |int   |请求项目的类型。1：类型一；2：类型二 。|
————————————————
