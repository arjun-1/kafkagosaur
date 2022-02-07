package kafkagosaur

import (
	"errors"
	"syscall/js"

	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

func SASLMechanism(config js.Value) (sasl.Mechanism, error) {
	if saslConfig := config.Get("sasl"); !saslConfig.IsUndefined() {
		username := saslConfig.Get("username").String()
		password := saslConfig.Get("password").String()

		switch saslConfig.Get("mechanism").String() {
		case "PLAIN":
			return plain.Mechanism{
				Username: username,
				Password: password,
			}, nil
		case "SCRAM-SHA-512":
			return scram.Mechanism(
				scram.SHA512,
				username,
				password,
			)
		case "SCRAM-SHA-256":
			return scram.Mechanism(
				scram.SHA256,
				username,
				password,
			)
		default:
			return nil, errors.New("unknown SASL mechanism")
		}
	}

	return nil, nil
}
