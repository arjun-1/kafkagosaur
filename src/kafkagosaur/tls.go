package kafkagosaur

import (
	"crypto/tls"
	"crypto/x509"

	"errors"
	"syscall/js"
)

func TLSConfig(config js.Value) (*tls.Config, error) {
	tlsJs := config.Get("tls")
	tlsConfig := &tls.Config{}

	handleBoolean := func() (*tls.Config, error) {
		if tlsJs.Bool() {
			return tlsConfig, nil
		} else {
			return nil, nil
		}
	}

	handleObject := func() (*tls.Config, error) {
		if keyPairJs := config.Get("keyPair"); !keyPairJs.IsUndefined() {
			key := keyPairJs.Get("key").String()
			cert := keyPairJs.Get("cert").String()

			keyPair, err := tls.X509KeyPair([]byte(key), []byte(cert))
			if err != nil {
				return nil, err
			}
			tlsConfig.Certificates = []tls.Certificate{keyPair}
		}

		if caJs := tlsJs.Get("ca"); !caJs.IsUndefined() {
			caCertPool := x509.NewCertPool()

			if !caCertPool.AppendCertsFromPEM([]byte(caJs.String())) {
				return nil, errors.New("CA certificate could not be parsed")
			} else {
				tlsConfig.RootCAs = caCertPool
				return tlsConfig, nil
			}
		}

		if insecureSkipVerifyJs := tlsJs.Get("insecureSkipVerify"); !insecureSkipVerifyJs.IsUndefined() {
			tlsConfig.InsecureSkipVerify = insecureSkipVerifyJs.Bool()
		}

		return tlsConfig, nil
	}

	if !tlsJs.IsUndefined() {
		switch tlsJs.Type() {
		case js.TypeObject:
			return handleObject()
		default:
			return handleBoolean()
		}
	}

	return nil, nil
}
