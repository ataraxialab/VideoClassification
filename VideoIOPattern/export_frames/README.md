# This is code for Cal Flow and Frame

#Plese install OpenCv3.2 first, details:http://docs.opencv.org/2.4/doc/tutorials/introduction/linux_install/linux_install.html;
#and install gcc, g++, yasm if need
#run build.sh to build the bin for Cal Flow and Frame
#use the export_frames as following:

./export_frames -i Media_URL -o OutputDir -c frame_count -interval check_time_interval -postfix picfmt -log logoutputflag

Ex:

./export_frames -i test.mp4 -o ./pics -c 21 -interval 10 -postfix jpg -log 1
