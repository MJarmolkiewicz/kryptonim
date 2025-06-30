## 💱 Kryptonim – Currency & Crypto Exchange API

Prosta aplikacja w Go z wykorzystaniem frameworka Gin, służąca do pobierania kursów walut oraz przeliczania kryptowalut.

## 📋 Wymagania

- Go 1.24.2+
- Docker (opcjonalnie)
- Konto i klucz API do [OpenExchangeRates.org](https://openexchangerates.org)

## 🚀 Uruchomienie
Lokalnie

    Ustaw zmienną środowiskową OXR_API_KEY z Twoim kluczem:

export OXR_API_KEY=your_api_key

Uruchom aplikację:

    go run main.go

Z wykorzystaniem Dockera

    Zbuduj obraz:

docker build -t go-currency-app .

Uruchom kontener:

    docker run -p 8080:8080 --env OXR_API_KEY=your_api_key go-currency-app

🧪 Przykłady użycia API
Request:

GET /rates?currencies=<your-currency-list>

Response (JSON):

[
  { "from": "USD", "to": "GBP", "rate": 1.0 },
  { "from": "GBP", "to": "USD", "rate": 1.0 },
  { "from": "USD", "to": "EUR", "rate": 1.0 },
  { "from": "EUR", "to": "USD", "rate": 1.0 },
  { "from": "EUR", "to": "GBP", "rate": 1.0 },
  { "from": "GBP", "to": "EUR", "rate": 1.0 }
]

Request:

GET /rates?exchange=?from=<from-crypto>&to=<to-crypto>&amount=<ammount>

Response (JSON):

[
  { "from": "EUR", "to": "GBP", "rate": 1.0 },
  { "from": "GBP", "to": "EUR", "rate": 1.0 }
]

⚙️ Dostępne endpointy

    GET /rates?currencies=... – zwraca kursy między podanymi walutami

    GET /exchange?from=...&to=...&amount=... – przeliczanie wartości kryptowalut 
