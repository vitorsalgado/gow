package worker

import "github.com/rs/zerolog/log"

func info(msg string) {
	log.Info().
		Timestamp().
		Str(CtxKey, CtxValue).
		Msg(msg)
}
