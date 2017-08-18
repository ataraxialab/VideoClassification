<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [视频截帧API](#%E8%A7%86%E9%A2%91%E6%88%AA%E5%B8%A7api)
  - [架构](#%E6%9E%B6%E6%9E%84)
  - [截帧方式](#%E6%88%AA%E5%B8%A7%E6%96%B9%E5%BC%8F)
  - [光流计算方式](#%E5%85%89%E6%B5%81%E8%AE%A1%E7%AE%97%E6%96%B9%E5%BC%8F)
  - [API](#api)
    - [截帧/计算光流](#%E6%88%AA%E5%B8%A7%E8%AE%A1%E7%AE%97%E5%85%89%E6%B5%81)
      - [请求](#%E8%AF%B7%E6%B1%82)
      - [返回](#%E8%BF%94%E5%9B%9E)
    - [获取截帧/光流结果](#%E8%8E%B7%E5%8F%96%E6%88%AA%E5%B8%A7%E5%85%89%E6%B5%81%E7%BB%93%E6%9E%9C)
      - [请求](#%E8%AF%B7%E6%B1%82-1)
      - [返回](#%E8%BF%94%E5%9B%9E-1)

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

### 截帧/计算光流

根据截帧模式对视频进行截帧&计算光流，并保存结果。

#### 请求

```
POST /:target                    // target = frame(截帧) | flow(光流)
Content-Type: application/json

{
  "pattern": "random",           // random: 从视频的某个时间开始随机截取多帧
  "op": "start|stop",            // start: 开始截帧生产数据, stop: 结束截帧
  "params": {                    // op = start 有效
    "count": 100,                // 帧数量
    "offset": 0.1,               // 开始截帧的偏移量，范围：0-1
  }
}
```

#### 返回

```
// 成功
200

// 非法请求
4xx

Content-Type: application/json
{
  "message": "error message"
}

// 服务器错误
5xx

Content-Type: application/json
{
  "message": "error message"
}
```



### 获取截帧/光流结果

根据模式，获取截帧/光流数据。

#### 请求

```
GET /:target/:pattern/:from/:count
// target = frame(截帧) | flow(光流)
// pattern = random
// from >= 0
// count >= 0
```

#### 返回

```
// 成功
200
Content-Type: application/json

[{
    "idx": 1000,
    "label": 999.0,
    "image_path": "path/of/image0"
  }, {
    "idx": 1001,
    "label": 999.0,
    "image_path": "path/of/image1"
  }
]

// 非法请求
4xx

Content-Type: application/json
{
  "message": "error message"
}

// 服务器错误
5xx

Content-Type: application/json
{
  "message": "error message"
}
```
