# This is code for Cal Flow and Frame

#Plese install OpenCv3.2 first, details:http://docs.opencv.org/2.4/doc/tutorials/introduction/linux_install/linux_install.html;
#and install gcc, g++, yasm if need
#run build.sh to build the bin for Cal Flow and Frame
#use the export_frames as following:

./export_frames -i Media_URL -o OutputDir -s WidthxHeight -c frame_count -interval check_time_interval -postfix picfmt -log logoutputflag

Ex:

./export_frames -i test.mp4 -o ./pics -s 256x256 -c 21 -interval 10 -postfix jpg -log 1


===

* 2017.08.13 改动

DoFrameExport函数增加接收任意位置解码功能

输入为改定视频，起始位置，需要的视频帧数，解码得到对应帧数图像
起始位置的输入为0-1之间的float，代表起始点在视频中的位置

decoder_worker.cpp 用于接收解码工作序列然后顺序执行解码请求
工作流程：
读取解码序列文件
按照序列依次调用DoFrameExport函数解码对应N帧图像
将解码的图像写到对应内存位置/Redis/落盘
