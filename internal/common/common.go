package common

func GetPayload(c *gin.Context) Payload {
	pl := c.Keys[CONTEXT_KEY_JWT_PAYLOAD]
	if pl == nil {
		return Payload{}
	}
	return pl.(Payload)
}