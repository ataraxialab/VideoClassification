# 1. 随机序列生成 - rander.py

* 1 遍历UCF-101文件夹，
获取文件列表，类别列表
得到每个类别有多少个文件
每个类别每个文件有多少帧，长宽，fps
得到每个文件对应的类别Idx和类别name
[
	{
		clsidx: 0,
		clsname: "jump",
		root: "/disk2/data/UCF-101/UCF-101",
		videos:[
		{idx:0,name:"vid-001.avi", "nbFrames":100, "fps":30, "iWidth":256, "iHeight":256},
		{idx:1,name:"vid-002.avi"},
		...
		]
	},
	...

]

* 2 生成3个随机序列，分别对应classIdx，videoIdx，frameIdx

struct {
	float clsIdx;
	float vidIdx;
	float frmIdx;
}


* 3 保存文件
	视频信息文件
	随机序列文件
