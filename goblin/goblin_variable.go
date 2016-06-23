package goblin

const (
	ADD_RULE                = iota
	DELETE_RULE
	READ_RULE

	RULE_ADD_LOCATION       = "/rule/add"
	RULE_DELETE_LOCATION    = "/rule/delete"
	RULE_READ_LOCATION      = "/rule/read"

	SERVER_ADD_LOCATION     = "/server/add"
	SERVER_DELETE_LOCATION  = "/server/delete"
	SERVER_READ_LOCATION    = "/server/read"

	GOBLIN_ADD_LOCATION     = "/add"
	GOBLIN_DELETE_LOCATION  = "/del"

	DEFAULT_GOBLIN_LOCATION = "/goblin"
	DEFAULT_API_LOCATION    = "/nginx"

	EMPTY_IP                = "0.0.0.0"
	EMPTY_UUID              = "0"
	EMPTY_UID               = "0"
)
