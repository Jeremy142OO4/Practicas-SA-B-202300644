import http from "k6/http";
import { check } from "k6";

export const options = {
  stages: [
    { duration: "20s", target: 50 },
    { duration: "20s", target: 100 },
    { duration: "20s", target: 200 },
    { duration: "60s", target: 200 },
    { duration: "20s", target: 0 },
  ],
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<500"],
  },
};

const baseUrl = __ENV.BASE_URL || "http://127.0.0.1:4000";
const body = JSON.stringify({
  query: "query { videojuegos { idVideojuego nombre activo } }",
});

export default function () {
  const response = http.post(`${baseUrl}/catalog/graphql`, body, {
    headers: { "Content-Type": "application/json" },
  });

  check(response, {
    "respuesta HTTP 200": (result) => result.status === 200,
  });
}
