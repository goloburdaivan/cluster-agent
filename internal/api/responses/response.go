package responses

type Response[T any] struct {
	Data  T      `json:"data"`
	Error string `json:"error,omitempty"`
}

func Success[T any](data T) Response[T] {
	return Response[T]{
		Data: data,
	}
}

func Error(msg string) Response[any] {
	return Response[any]{
		Error: msg,
	}
}
