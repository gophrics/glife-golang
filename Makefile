.PHONY : build


build:
	(cd profile; make build;)
	(cd location; make build;)
	go build -o build/main main.go
run:
	./build/profile & ./build/location & ./build/main
