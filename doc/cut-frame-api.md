<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [视频截帧API](#%E8%A7%86%E9%A2%91%E6%88%AA%E5%B8%A7api)
  - [架构](#%E6%9E%B6%E6%9E%84)
  - [截帧方式](#%E6%88%AA%E5%B8%A7%E6%96%B9%E5%BC%8F)
  - [光流计算方式](#%E5%85%89%E6%B5%81%E8%AE%A1%E7%AE%97%E6%96%B9%E5%BC%8F)
  - [API](#api)
    - [截帧](#%E6%88%AA%E5%B8%A7)
      - [请求](#%E8%AF%B7%E6%B1%82)
      - [返回](#%E8%BF%94%E5%9B%9E)
    - [计算光流](#%E8%AE%A1%E7%AE%97%E5%85%89%E6%B5%81)
      - [请求](#%E8%AF%B7%E6%B1%82-1)
      - [返回](#%E8%BF%94%E5%9B%9E-1)
    - [获取截帧结果](#%E8%8E%B7%E5%8F%96%E6%88%AA%E5%B8%A7%E7%BB%93%E6%9E%9C)
      - [请求](#%E8%AF%B7%E6%B1%82-2)
      - [返回](#%E8%BF%94%E5%9B%9E-2)
    - [获取光流结果](#%E8%8E%B7%E5%8F%96%E5%85%89%E6%B5%81%E7%BB%93%E6%9E%9C)
      - [请求](#%E8%AF%B7%E6%B1%82-3)
      - [返回](#%E8%BF%94%E5%9B%9E-3)
    - [截帧数据使用结束](#%E6%88%AA%E5%B8%A7%E6%95%B0%E6%8D%AE%E4%BD%BF%E7%94%A8%E7%BB%93%E6%9D%9F)
      - [请求](#%E8%AF%B7%E6%B1%82-4)
      - [返回](#%E8%BF%94%E5%9B%9E-4)
    - [光流数据使用结束](#%E5%85%89%E6%B5%81%E6%95%B0%E6%8D%AE%E4%BD%BF%E7%94%A8%E7%BB%93%E6%9D%9F)
      - [请求](#%E8%AF%B7%E6%B1%82-5)
      - [返回](#%E8%BF%94%E5%9B%9E-5)
    - [其他（查询历史，删除数据==，但是非核心可以放后面做）](#%E5%85%B6%E4%BB%96%E6%9F%A5%E8%AF%A2%E5%8E%86%E5%8F%B2%E5%88%A0%E9%99%A4%E6%95%B0%E6%8D%AE%E4%BD%86%E6%98%AF%E9%9D%9E%E6%A0%B8%E5%BF%83%E5%8F%AF%E4%BB%A5%E6%94%BE%E5%90%8E%E9%9D%A2%E5%81%9A)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



# 视频截帧API

提供视频截帧API，对截帧进行管理。主要功能：提供不同的方式进行视频截帧，并提供访问截帧结果的能力。



## 架构

```

                                              +-------------------+
                   request frame&flow ----->  | frame&flow server | <------ consume frame&flow
                                              +-------------------+
                                                        |  store result
                                                        v
                                              +-------------------+
                                              | storage server    |
                                              +-------------------+
```



## 截帧方式

四种基本的截帧方式:

1. 选择一个视频随机的截取一段连续的原始图像
2. 选择一个视频随机截取一段连续的光流图像

3. 选择一个视频同时的返回视频中某一段的光流和原始图像


4. 在一个视频中, 按照视频的时间顺序从前往后均匀选取时间点,用以上3种方式的一种去截取图像.


## 光流计算方式

参考 **截帧方式**



## API

### 截帧

根据截帧模式对视频进行截帧，并保存结果。

#### 请求

```
POST /frame
Content-Type: application/json

{
  "pattern": "random|sample",    // random: 随机截取一段，sample: 时间范围，然后在时间范围内随机
  "op": "start|stop",            // start: 开始截帧生产数据, stop: 结束截帧
  "params": {
    "duration": 100,             // 截帧的时长
    "interval": 100              // 播放时间间隔，sample时需要
  }
}
```

#### 返回

```
// 成功
200

// 请求参数非法
400

// 服务器错误
5xx

Content-Type: application/json
{
  "message": "error message"
}
```



### 计算光流

根据截帧模式对视频进行计算光流，并保存结果。

#### 请求

```
POST /flow
Content-Type: application/json

{
  "pattern": "random|sample",    // random: 随机截取一段，sample: 时间范围，然后在时间范围内随机
  "op": "start|stop",            // start: 开始计算光流生产数据, stop: 结束计算
  "params": {
    "duration": 100,             // 计算光流的连续视频时长
    "interval": 100              // 播放时间间隔，sample时需要
  }
}
```

#### 返回

```
// 成功
200

// 请求参数非法
400

// 服务器错误
5xx

Content-Type: application/json
{
  "message": "error message"
}
```



### 获取截帧结果

根据截帧模式，获取截帧数据。

#### 请求

```
GET /frame
Content-Type: application/json

{
  "pattern": "random|sample",  // 同上
  "from": 0,                   // 起始位置
  "count": 1                   // 1 - 100
}
```

#### 返回

```
// 成功
200
Content-Type: application/json

{
 "frames":[{
    "idx": 1000,
    "label": 999,0,
    "image_path": "path/of/image"
 }]
}

// 请求参数非法
400

// 服务器错误
5xx

Content-Type: application/json
{
  "message": "error message"
}
```



### 获取光流结果

根据计算光流算法，获取光流数据。

#### 请求

```
GET /flow
Content-Type: application/json

{
  "pattern": "random|sample",  // 同上
  "from": 0,                   // 起始位置
  "count": 1                   // 1 - 100
}
```

#### 返回

```
// 成功
200
Content-Type: application/json

{
	// TODO 返回什么数据？
}

// 请求参数非法
400

// 服务器错误
5xx

Content-Type: application/json
{
  "message": "error message"
}
```



### 截帧数据使用结束

为了删除生成的截帧数据

#### 请求

```
DELETE /frame
Content-Type: application/json

{
  "id": ["frame1", "frame2"] // 截帧id，在获取接口中返回
}
```

#### 返回

```
// 成功
200

// 服务器错误
5xx

Content-Type: application/json
{
  "message": "error message"
}
```



### 光流数据使用结束

为了删除生成的光流数据

#### 请求

```
DELETE /flow
Content-Type: application/json

{
  "id": ["flow1", "flow2"] // 光流id，在获取接口中返回
}
```

#### 返回

```
// 成功
200

// 服务器错误
5xx

Content-Type: application/json
{
  "message": "error message"
}
```



### 其他（查询历史，删除数据==，但是非核心可以放后面做）

