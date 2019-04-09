CRYPTO_PATH = normalCrypto/image/payload
DISTSERVER_PATH = distServer/image/payload
RAND_PATH = Rand_Number_Assess-master

all: normalcrypto_docker distserver_docker

normalcrypto_docker: $(CRYPTO_PATH)/normalcrypto
	@cp $(RAND_PATH)/libtest.so $(CRYPTO_PATH)
	@docker build -t normalcrypto normalCrypto/image

$(CRYPTO_PATH)/normalcrypto:
	@echo $(CRYPTO_PATH)
	@mkdir -p $(CRYPTO_PATH)
	@go build -o $(CRYPTO_PATH)/normalCrypto normalCrypto

distserver_docker: $(DISTSERVER_PATH)/distserver
	@docker build -t distserver distServer/image

$(DISTSERVER_PATH)/distserver:
	@echo $(DISTSERVER_PATH)
	@mkdir -p $(DISTSERVER_PATH)
	@go build -o $(DISTSERVER_PATH)/distServer distServer

.PHONY: clean
clean :
	@rm -rf $(CRYPTO_PATH)
	@rm -rf $(DISTSERVER_PATH)