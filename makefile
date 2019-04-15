# This makefile defines the following targets
#
#   - all (default) - builds all targets
#   - clean -remove all file produced by this file

#PROJECT_PATH = github.com/HanDaXia/BlockChainSafeTesting
CRYPTO_PATH = normalCrypto/image/payload
DISTSERVER_PATH = distServer/image/payload
RAND_PATH = Rand_Number_Assess-master
MESSAGEHUB_PATH = messagehub/image/payload


all: normalcrypto_docker distserver_docker messagehub_docker

normalcrypto_docker: $(CRYPTO_PATH)/normalcrypto
	@cd $(RAND_PATH) && make && cd ..
	@cp $(RAND_PATH)/libtest.so $(CRYPTO_PATH)
	@cp -rf $(RAND_PATH)/templates $(CRYPTO_PATH)
	@docker build -t normalcrypto normalCrypto/image

$(CRYPTO_PATH)/normalcrypto:
	@echo $(CRYPTO_PATH)
	@mkdir -p $(CRYPTO_PATH)
	@go build -o $(CRYPTO_PATH)/normalCrypto ./normalCrypto

distserver_docker: $(DISTSERVER_PATH)/distserver
	@docker build -t distserver distServer/image

$(DISTSERVER_PATH)/distserver:
	@echo $(DISTSERVER_PATH)
	@mkdir -p $(DISTSERVER_PATH)
	@go build -o $(DISTSERVER_PATH)/distServer ./distServer

messagehub_docker:$(MESSAGEHUB_PATH)/messagehub
	@docker build -t messagehub messagehub/image

$(MESSAGEHUB_PATH)/messagehub:
	@mkdir -p $(MESSAGEHUB_PATH)
	@go build -o $(MESSAGEHUB_PATH)/messagehub ./messagehub


.PHONY: clean
clean :
	@rm -rf $(CRYPTO_PATH) $(DISTSERVER_PATH) $(RAND_PATH)/obj $(RAND_PATH)/libtest.so $(MESSAGEHUB_PATH)
