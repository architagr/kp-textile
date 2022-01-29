build:
	make -C ../client-service lambda-build
	make -C ../hsn-code-service lambda-build
	make -C ../item-service lambda-build
	make -C ../quality-service lambda-build
	make -C ../transportor-service lambda-build
	make -C ../vendor-service lambda-build
	make -C ../document-service lambda-build

deploy:
	rm -rf cdk.out
	clear
	cdk deploy --all --require-approval never