main:
	go generate
	go install 

example:
	cd examples && startapp example . --config=config.yaml
	cd examples/example/clients/ios && make
