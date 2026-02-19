.PHONY: build test testacc install clean

BINARY_NAME=terraform-provider-htpasswd
PROVIDER_NAMESPACE=loafoe
PROVIDER_HOST=registry.opentofu.org

build:
	go build -o $(BINARY_NAME)

test:
	go test ./... -v

testacc:
	TF_ACC=1 \
	TF_ACC_TERRAFORM_PATH=$$(which tofu) \
	TF_ACC_PROVIDER_NAMESPACE=$(PROVIDER_NAMESPACE) \
	TF_ACC_PROVIDER_HOST=$(PROVIDER_HOST) \
	go test ./htpasswd -v -timeout 30m

install: build
	mkdir -p ~/.terraform.d/plugins/$(PROVIDER_HOST)/$(PROVIDER_NAMESPACE)/htpasswd/0.0.1/$$(go env GOOS)_$$(go env GOARCH)
	cp $(BINARY_NAME) ~/.terraform.d/plugins/$(PROVIDER_HOST)/$(PROVIDER_NAMESPACE)/htpasswd/0.0.1/$$(go env GOOS)_$$(go env GOARCH)/

clean:
	rm -f $(BINARY_NAME)
