.PHONY: clean

TARGET=tensority

$(TARGET): cSimdTs.o
	cd .. && go test 

%.o: %.cpp
	g++ -o $@ -c $^ -std=c++11 -pthread -mavx2 -O3 -fopenmp -D_USE_OPENMP

clean:
	rm -f *.o *.so *.a