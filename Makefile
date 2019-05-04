.PHONY : build


build:
	(cd profile; make build;)
	(cd location; make build;)
run:
	./build/profile & ./build/location
