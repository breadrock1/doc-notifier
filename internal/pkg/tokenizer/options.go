package tokenizer

type TokenizerMode int

const (
	Assistant TokenizerMode = iota
	LangChain
	None
)

type Options struct {
	Mode         TokenizerMode
	Address      string
	ChunkedFlag  bool
	ChunkSize    int
	ChunkOverlap int
}

func GetModeFromString(mode string) TokenizerMode {
	switch mode {
	case "assistant":
		return Assistant
	case "langchain":
		return LangChain
	case "none":
		return None
	default:
		return None
	}
}
