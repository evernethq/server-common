package errors

import (
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	krhttp "github.com/go-kratos/kratos/v2/transport/http"
)

const (
	baseContentType = "application"
)

func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := errors.FromError(err)
	if strings.Contains(se.Message, "context deadline exceeded") {
		se = errors.FromError(TimeOut)
	}
	if se.Reason == "" {
		se = errors.FromError(Unknown(se.Message))
	}

	codec, _ := krhttp.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	w.WriteHeader(int(se.Code))
	_, _ = w.Write(body)
}

func ContentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}
