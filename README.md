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
|transcoding_input_url  |string |转码输入URL                      |
|transcoding_output_url  |string |转码输出URL                      |
|transcoding_params  |string |转码参数                      |
|transcoding_start_time |string |转码开始时间                         |
|transcoding_last_active_time |string |转码最近状态检查时间                         |
|transcoding_state_info |string |转码状态信息                         |

###### 接口示例
http://host:port/v1/transcodings?output_url=cnRtcDovLzU4LjIwMC4xMzEuMjoxOTM1L2xpdmV0di9odW5hbnR2

"cnRtcDovLzU4LjIwMC4xMzEuMjoxOTM1L2xpdmV0di9odW5hbnR2" 为"rtmp://58.200.131.2:1935/livetv/hunantv"的 urlbase64编码


###### 返回
- Body
```
{
  "code": 200,
  "data": "730781",
  "message": "OK"
}
```
