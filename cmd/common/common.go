package common

type (
	CTX_CLIENT    struct{}
	CTX_API_BASE  struct{}
	CTX_AUTH_BASE struct{}
	CTX_FORMAT    struct{}
	CTX_VERSION   struct{}
)

const (
	FORMAT_JSON  = "json"
	FORMAT_HUMAN = "human"
)
