package infra

import "github.com/rs/zerolog/log"

type SubscriberLogger struct{}

func (sl SubscriberLogger) NewField(fieldname string) {
	log.Info().Str("name", fieldname).Msg("new field")
}

func (sl SubscriberLogger) FirstNonMaskedValue(fieldname string, value any) {
	log.Info().Str("name", fieldname).Msg("unmasked value detected")
}
