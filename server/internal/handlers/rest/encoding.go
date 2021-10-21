package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/holmes89/thoth/internal"
	"github.com/rs/zerolog/log"
)

// EncodeResponse writes json to response writer
func EncodeResponse(w http.ResponseWriter, response interface{}, err error, msg ...string) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		var t string
		if response != nil {
			t = reflect.TypeOf(response).String()
		}
		log.Error().Err(err).Msgf("unable to process entity: %s", t)
		errmsg := processError(err)
		w.WriteHeader(errmsg.Id)
		response = errmsg
	}

	if err := enc.Encode(response); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Error().AnErr("error", err).Str("body", fmt.Sprintf("%+v", response)).Msg("unable to encode response")
	}
}

func processError(err error, message ...string) internal.ErrorObj {
	desc := err.Error()
	msg := internal.ErrorObj{
		Description: &desc,
	}
	if len(message) > 0 {
		msg.Description = &message[0]
	}
	switch err {
	default:
		msg.Name = "server error"
		msg.Id = http.StatusInternalServerError
	}
	return msg
}
