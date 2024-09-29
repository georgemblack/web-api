go install go.uber.org/mock/mockgen@latest
mockgen -source=pkg/api/api.go -destination=pkg/testutil/mocks.go -package=testutil
