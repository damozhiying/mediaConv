# mediaConv
直播流转码服务，提供RTSP，RTMP 流媒体转码服务

---
接口API规范
---
**1\. 查询转码实例**
###### 接口功能
> 查询转码实例信息
###### URL
> /v1/transcodings

###### 支持格式
> JSON

###### HTTP请求方式
> GET

###### 请求参数
>
| 参数 | 必选 | 类型 | 说明 |
|:-----  |:-------|:-----|-----|
|output_url  |false|string|转码输出URL|
