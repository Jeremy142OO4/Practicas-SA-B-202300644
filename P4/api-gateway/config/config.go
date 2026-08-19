package config

import "os"

type Configuracion struct {
	Puerto      string
	AuthURL     string
	CatalogURL  string
	OrdersURL   string
	PaymentsURL string
}

func CargarConfiguracion() Configuracion {

	puerto := os.Getenv("PORT")

	if puerto == "" {
		puerto = "4000"
	}

	authURL := os.Getenv("AUTH_URL")

	if authURL == "" {
		authURL = "http://auth-ms:4001"
	}

	catalogURL := os.Getenv("CATALOG_URL")

	if catalogURL == "" {
		catalogURL = "http://catalog-ms:4002"
	}

	ordersURL := os.Getenv("ORDERS_URL")

	if ordersURL == "" {
		ordersURL = "http://orders-ms:4003"
	}

	paymentsURL := os.Getenv("PAYMENTS_URL")

	if paymentsURL == "" {
		paymentsURL = "http://payments-ms:4004"
	}

	return Configuracion{
		Puerto:      puerto,
		AuthURL:     authURL,
		CatalogURL:  catalogURL,
		OrdersURL:   ordersURL,
		PaymentsURL: paymentsURL,
	}
}
