package kafkagosaur

import (
	"errors"
	"syscall/js"

	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

func SASLMechanism(config js.Value) (sasl.Mechanism, error) {
	var mechanism sasl.Mechanism = nil
	var err error = nil

	if saslConfig := config.Get("sasl"); !saslConfig.IsUndefined() {
		username := saslConfig.Get("username").String()
		password := saslConfig.Get("password").String()

		switch saslConfig.Get("mechanism").String() {
		case "PLAIN":
			mechanism = plain.Mechanism{
				Username: username,
				Password: password,
			}
		case "SCRAM-SHA-512":
			mechanism, err = scram.Mechanism(
				scram.SHA512,
				username,
				password,
			)
		case "SCRAM-SHA-256":
			mechanism, err = scram.Mechanism(
				scram.SHA256,
				username,
				password,
			)
		default:
			err = errors.New("unknown SASL mechanism")
		}
	}

	return mechanism, err
}
