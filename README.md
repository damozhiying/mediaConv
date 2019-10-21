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

###### 返回字段
> 
|返回字段|字段类型|说明                              |
|:-----   |:------|:-----------------------------   |
|status   |int    |返回结果状态。0：正常；1：错误。   |
|transcoding_req  |string | 转码请求                      |
|transcoding_start_time |string |转码开始时间                         |
|transcoding_last_active_time |string |转码最近坚持时间                         |
|transcoding_state_info |string |转码状态信息                         |

————————————————
