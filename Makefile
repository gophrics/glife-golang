.PHONY : build


build:
	(cd profile; make build;)

run:
	make build
	go run main.go