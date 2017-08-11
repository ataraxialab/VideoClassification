OBJS=export_frames.o decode_frames.o trans_frames.o
INCLUDES=-I ./ffmpeg-3.3.2/include 
CXXFLAGS=-D__STDC_CONSTANT_MACROS `pkg-config --cflags opencv`
LDFLAGS=-L ./ffmpeg-3.3.2/lib -lavformat -lavcodec -lavutil -lrt -lswscale -lswresample -lpthread -lm -lz -ldl
LIBS=`pkg-config --libs opencv`

export_frames:$(OBJS)
	g++ $(CXXCFLAGS) -o export_frames $(OBJS) $(LIBS) $(LDFLAGS) 
export_frames.o:export_frames.cpp
	g++ $(CXXFLAGS) -c export_frames.cpp
decode_frames.o:decode_frames.cpp
	g++ $(CXXFLAGS) $(INCLUDES) -c decode_frames.cpp
trans_frames.o:trans_frames.cpp
	g++ $(CXXFLAGS) -c trans_frames.cpp

clean:
	-rm export_frames ${OBJS} -rf
