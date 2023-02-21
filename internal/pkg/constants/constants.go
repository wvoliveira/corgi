package constants

const VERSION = "0.0.1"

// TODO: put theses keywords in database
// So, we can update in real time.
// Ref: https://www.mediavine.com/keyword-anti-targeting/
var BLOCKED_KEYWORDS = []string{"crash", "attack", "terrorist", "suicide", "nazi", "killed", "porn", "explosion", "rape", "death", "isis", "shooting", "bomb", "dead", "murder", "terror", "kill", "sex", "massacre", "gun"}

// Cache pattern keys
var (
	TOKEN_AUTH     = "token_auth_%s"
	TOKEN_PERSONAL = "token_personal_%s"
)
