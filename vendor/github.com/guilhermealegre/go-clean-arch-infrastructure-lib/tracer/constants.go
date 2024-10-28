package tracer

// Tracer tags
const (
	TracerTagRequestBody       = "request.body"
	TracerTagResponseBody      = "response.body"
	TracerTagError             = "error"
	TracerTagRedisCmd          = "redis.cmd"
	TracerTagQuery             = "query"
	TracerTagEventName         = "event.name"
	TracerTagHttpMethod        = "http.method"
	TracerTagParams            = "params"
	TracerTagTracer            = "tracer"
	TracerTagPath              = "path"
	TracerTagContentType       = "content-type"
	TracerTagKey               = "key"
	TracerTagBucket            = "bucket"
	TracerTagStatusCode        = "otel.status_code"
	TracerTagStatusDescription = "otel.status_description"
)
