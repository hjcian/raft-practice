gen-go:
	@buf generate
	@mkdir -p grpcproto
	@mv ./gen/go/proto/* grpcproto/
	@rmdir -p gen/go/proto