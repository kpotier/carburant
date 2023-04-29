# carburant

[![Go Reference](https://pkg.go.dev/badge/github.com/kpotier/carburant.svg)](https://pkg.go.dev/github.com/kpotier/carburant)

The carburant repository fetches gas price data from the French government API and enables users to view the evolution of gas prices for each service station. It manages the data for each station and allows users to compare gas prices and find the cheapest gas stations near their location.

# Installation

1. `go get github.com/kpotier/carburant`

2. Get an API key from [https://cloud.maptiler.com/](https://cloud.maptiler.com/)

3. Put this API key to `public/src/global.ts`

4. `npm install && npm run build`

5. `go run cmd/main.go`

# Screenshots

![preview](screen1.png)
![preview](screen2.png)

# License

Ministère de l’Économie, des Finances et de la Souveraineté industrielle et numérique - Données originales téléchargées sur [https://data.economie.gouv.fr/explore/dataset/prix-des-carburants-en-france-flux-instantane-v2/](https://data.economie.gouv.fr/explore/dataset/prix-des-carburants-en-france-flux-instantane-v2/), mise à jour le 21 avril 2023.
