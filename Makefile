.PHONY: clean all

all: clean nola

nola:
	go build -o nola .

clean:
	rm -f nola