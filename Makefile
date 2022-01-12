lambda-build:
	make -C hsn-code-service lambda-build

deploy:
	make -C infra deploy
	