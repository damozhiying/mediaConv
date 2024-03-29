# mediaConv
直播流转码服务，提供RTSP，RTMP 流媒体转码服务

---
接口API规范
---
**1\. 查询转码实例**
#### 接口功能
> 查询转码实例信息
#### 请求参数
| 参数 | 必选 | 类型 | 说明 |
|:-----  |:-------|:-----|-----|
|output_url  |false|string|转码输出URL的urlbase64值,不带该参数表示查询所有转码实例|
#### 返回字段
> 
|返回字段|字段类型|说明                              |
|:-----   |:------|:-----------------------------   |
|code   |int    |返回结果状态码   |
|message  |string    |结果状态描述信息   |
|transcoding_count   |int    |符合条件的转码实例个数   |
|transcoding_array   |    |符合条件的转码实例数组   |
|transcoding_input_url  |string |转码输入URL                      |
|transcoding_output_url  |string |转码输出URL                      |
|transcoding_params  |string |转码参数                      |
|transcoding_start_time |string |转码开始时间                         |
|transcoding_last_active_time |string |转码最近状态检查时间                         |
|transcoding_state_info |string |转码状态信息                         |


#### Request

- Method: **GET**
- URL: ```/v1/transcodings?output_url={urlBase64(output_url)}```
  - 查询特定的转码实例:```/v1/transcodings?output_url=cnRtcDovLzU4LjIwMC4xMzEuMjoxOTM1L2xpdmV0di9odW5hbnR2XzUwMGs```
  - 查询特定的所有转码实例:```/v1/transcodings```  
- Header:
- Body:
```
```
#### Response


- Body
```
{
  "code": 200,
  "message": "Success",
  "transcoding_count": 1,
  "transcodings_array": [
    {
      "transcoding_input_url": "rtmp://58.200.131.2:1935/livetv/hunantv",
      "transcoding_output_url": "rtmp://58.200.131.2:1935/livetv/hunantv_500k",
      "transcoding_params": "/vb/500k/s/640x360",
      "transcoding_start_time": "2019-10-22 09:13:18.502538 +0800 CST m=+135.682582454",
      "transcoding_last_active_time": "2019-10-22 09:15:43.577268 +0800 CST m=+280.762973568"
      "transcoding_state_info": "frame= 2045 fps= 15 q=20.0 size=   16203kB time=00:02:14.20 bitrate= 989.1kbits/s         
                                speed=0.984x    \r"
    }
  ]
}
```

**2\. 创建转码实例**
#### 接口功能
> 创建转码实例
#### 请求参数
#### 返回字段
> 
|返回字段|字段类型|说明                              |
|:-----   |:------|:-----------------------------   |
|code   |int    |返回结果状态码   |
|message  |string    |结果状态描述信息   |

#### Request

- Method: **POST**
- URL: ```/v1/transcodings```
- Header: Content-Type:application/json
- Body:
```
{
  input_url: "rtmp:/58.200.131.2:1935/livetv/hunantv",
  output_url: "rtmp://58.200.131.2:1935/livetv/hunantv_500k",
  conv_params: /vb/500k
}
```
#### Response


- Body
```
{
  "code": 200,
  "message": "rtmp:/58.200.131.2:1935/livetv/hunantv Transcoding in Preparing"
}
```

**3\. 删除转码实例**
#### 接口功能
> 删除转码实例
#### 请求参数
| 参数 | 必选 | 类型 | 说明 |
|:-----  |:-------|:-----|-----|
|output_url  |true|string|转码输出URL的urlbase64值|
#### 返回字段
> 
|返回字段|字段类型|说明                              |
|:-----   |:------|:-----------------------------   |
|code   |int    |返回结果状态码   |
|message  |string    |结果状态描述信息   |

#### Request

- Method: **DEL**
- URL: ```/v1/transcodings```
  - 查询特定的转码实例:```/v1/transcodings?output_url=cnRtcDovLzU4LjIwMC4xMzEuMjoxOTM1L2xpdmV0di9odW5hbnR2XzUwMGs```
- Header: 
- Body:

#### Response

- Body
```
{
  "code": 200,
  "message": "rtmp:/58.200.131.2:1935/livetv/hunantv_500k Transcoding in Closing"
}
```
```
```
