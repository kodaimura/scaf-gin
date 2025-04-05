package common

func GetPayload(c *gin.Context) Payload {
	pl := c.Keys[CONTEXT_KEY_JWT_PAYLOAD]
	if pl == nil {
		return Payload{}
	}
	return pl.(Payload)
}

func GetCustomClaims(c *gin.Context) map[string]interface{} {
	pl := GetPayload(c)
	return pl.CustomClaims
}

func GetAccountId(c *gin.Context) string {
	value, _ := GetPayload(c)["account_id"]
	return value
}

func GetAccountName(c *gin.Context) string {
	value, _ := GetPayload(c)["account_name"]
	return value
}