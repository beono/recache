.PHONY: codereview

test:
	dep ensure
	go test

codereview:
	golint $$(go list ./...)
	go fmt $$(go list ./...)
	go vet $$(go list ./...)