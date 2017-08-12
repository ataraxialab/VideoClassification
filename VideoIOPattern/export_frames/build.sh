tar -xvf ffmpeg-3.3.2.tar.gz
cd ffmpeg-3.3.2
cd ffmpeg-3.3.2
./configure --enable-gpl --prefix=../ --enable-static --disable-shared --disable-ffprobe --disable-ffmpeg --disable-ffplay --disable-ffserver --disable-lzma
make -j4
make install
cd ..
cd ..
make
