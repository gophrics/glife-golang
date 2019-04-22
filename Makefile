.PHONY : build


build:
	(cd profile; make build;)
	(cd location; make build;)

run:
	make build
	go run main.go