.PHONY : build


build:
	(cd profile; make build;)
	(cd location; make build;)

run:
	make build
	(cd profile; $(MAKE) run) &
	(cd location; $(MAKE) run) &
	(cd common; $(MAKE) run) &
	go run main.go