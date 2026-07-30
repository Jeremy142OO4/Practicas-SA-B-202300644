package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

func DecodificarJSONEstricto(
	data []byte,
	destino any,
) error {
	decoder := json.NewDecoder(
		bytes.NewReader(data),
	)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destino); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(
		err,
		io.EOF,
	) {
		if err == nil {
			return errors.New(
				"el cuerpo debe contener un único objeto JSON",
			)
		}

		return err
	}

	return nil
}
