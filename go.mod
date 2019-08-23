module github.com:tencentyun/tfdoc

go 1.12

require (
	github.com/hashicorp/terraform v0.12.7
	github.com/terraform-providers/terraform-provider-tencentcloud v1.16.0
)

replace (
	github.com/Unknwon/com v0.0.0-20190804042917-757f69c95f3e => github.com/unknwon/com v0.0.0-20190804042917-757f69c95f3e
	github.com/golang/lint v0.0.0-20190409202823-959b441ac422 => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	github.com/hashicorp/vault-plugin-auth-pcf v0.0.0-20190821162840-1c2205826fee => github.com/hashicorp/vault-plugin-auth-cf v0.0.0-20190821162840-1c2205826fee
	github.com/nats-io/go-nats v1.8.1 => github.com/nats-io/nats.go v1.8.1
	github.com/openzipkin/zipkin-go-opentracing v0.4.2 => github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.2
	sourcegraph.com/sourcegraph/go-diff v0.5.1 => github.com/sourcegraph/go-diff v0.5.1
)
