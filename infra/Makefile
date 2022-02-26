build:
	make -C ../vendor-service lambda-build
	make -C ../document-service lambda-build
	make -C ../item-service lambda-build
	make -C ../quality-service lambda-build
	make -C ../transportor-service lambda-build
	make -C ../organization-service lambda-build

	
deploy:
	make build
	rm -rf cdk.out
	clear
	cdk deploy --all --require-approval never