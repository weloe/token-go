package constant

// ctx.Response constant variable
const (
	AccessControlExposeHeaders = "Access-Control-Expose-Headers"

	SetCookie = "Set-Cookie"
)

// persist timeout constant variable
const (
	// NeverExpire does not expire
	NeverExpire int64 = -1
	// NotValueExpire does not exist
	NotValueExpire int64 = -2
)

const (
	TokenName = "Tokengo"
)

const (
	BeReplaced int = -4
	BeKicked   int = -5
	BeBanned   int = -6
)
