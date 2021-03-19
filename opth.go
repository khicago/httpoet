package httpoet


func OHContentPB() Option {
	return OSetHeader("Content-Type", "application/x-protobuf")
}

func OHContentJson() Option {
	return OSetHeader("Content-Type", "application/json")
}

func OHContentString() Option {
	return OSetHeader("Content-Type", "application/json")
}

func OHContentForm() Option {
	return OSetHeader("Content-Type", "application/x-www-form-urlencoded")
}

func OHUserAgent(us string) Option {
	return OSetHeader("User-Agent", us)
}
