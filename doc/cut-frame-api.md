# 视频截帧API

提供视频截帧API，对截帧进行管理。主要功能：提供不同的方式进行视频截帧，并提供访问截帧结果的能力。

## 架构

```

                                              +-------------------+
                                              | process           |
                   request frame&flow ----->  | frame&flow server |
                                              +-------------------+
                                                        |  store result
                                                        v
                                              +-------------------+
                                              | frame&flow result |
                   consume frame&flow <-----  | server            |
                                              +-------------------+
```

## API

### 请求截帧

根据截帧模式对视频进行截帧，并保存结果。

#### 请求

```
POST /frame
Content-Type: application/json

{
  "pattern": "random" // TODO 有哪些模式
}
```

#### 返回

```
// 成功
200

// 请求参数非法
400

// 服务器错误
500
```

#### 获取截帧结果

根据截帧模式，获取截帧数据。

#### 请求

```
GET /frame
Content-Type: application/json

{
  "pattern": "random", // TODO 有哪些模式
  "count": 1 // 1 - 100
}

```

#### 返回

```
200
Content-Type: application/json

{
	// TODO 返回什么数据？
}
```

### 其他（查询历史，删除数据==，但是非核心可以放后面做）

