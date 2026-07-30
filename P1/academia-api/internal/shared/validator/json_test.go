package validator

import "testing"

func TestDecodificarJSONEstricto(
	t *testing.T,
) {
	t.Parallel()

	type cuerpo struct {
		Estado string `json:"estado"`
	}

	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "acepta campos conocidos",
			json:    `{"estado":"finalizada"}`,
			wantErr: false,
		},
		{
			name: "rechaza campos adicionales",
			json: `{
				"estado":"finalizada",
				"titulo":"no permitido"
			}`,
			wantErr: true,
		},
		{
			name: "rechaza múltiples objetos",
			json: `{"estado":"finalizada"}
				{"estado":"registrada"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var destino cuerpo

			err := DecodificarJSONEstricto(
				[]byte(tt.json),
				&destino,
			)

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"error = %v, se esperaba error = %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
