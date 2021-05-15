package coreLib

//CryptoInfo - Struct for handling GET requests from Coinmarketcap
type CryptoInfo struct {
	Status map[string]interface{} `json:"status"`
	//Data   []interface{}          `json:"data"`
	Data []struct {
		Id         int     `json:"id"`
		Name       string  `json:"name"`
		Symbol     string  `json:"symbol"`
		MaxSupply  float64 `json:"max_supply"`
		CircSupply float64 `json:"circulating_supply"`
		TotSupply  float64 `json:"total_supply"`
		Rank       int32   `json:"cmc_rank"`
		Quote      struct {
			Usd struct {
				Price        float64 `json:"price"`
				Vol24        float64 `json:"volume_24h"`
				PercentChg24 float64 `json:"percent_change_24h"`
				PercentChg7d float64 `json:"percent_change_7d"`
				MarketCap    float64 `json:"market_cap"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

//CryptoInternal - Internal representation for cryptocurrencies
type CryptoInternal struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	MaxSupply    float64 `json:"max_supply"`
	CircSupply   float64 `json:"circulating_supply"`
	TotSupply    float64 `json:"total_supply"`
	Rank         int32   `json:"cmc_rank"`
	Price        float64 `json:"USD_price"`
	Vol24        float64 `json:"volume_24h"`
	PercentChg24 float64 `json:"percent_change_24h"`
	PercentChg7d float64 `json:"percent_change_7d"`
	MarketCap    float64 `json:"market_cap"`
}

//VolumeWebhook - Webhook struct for volume changes/trends
type VolumeWebhook struct {
	Url               string  `json:"url"`
	Number            string  `json:"phone_number"` //For sending webhook notifications to phone number
	Id                int     `json:"id"`
	Name              string  `json:"name"`
	Symbol            string  `json:"symbol"`
	StartVol          float64 `json:"starting_volume"`
	CurrentVol        float64 `json:"current_volume"`
	PercentThreshold  float64 `json:"percentage_threshold"`
	CurrentPercentage float64 `json:"current_percenatge"`
	HasTriggered      bool    `json:"webhook_has_triggered"` //Volume webhook will be deleted upon invocation
	WebhookID         string  `json:"webhook_id"`
}

//PriceWebhook - Webhook for checking if a price point has been reached
type PriceWebhook struct {
	Url             string  `json:"url"`
	Number          string  `json:"phone_number"` //For sending webhook notifications to phone number
	Name            string  `json:"name"`
	Symbol          string  `json:"symbol"`
	StartPrice      float64 `json:"start_price"` //For figuring out if price point is up or down
	CurrentPrice    float64 `json:"current_price"`
	TargetPrice     float64 `json:"target_price"`
	IsPriceIncrease bool    `json:"is_price_increase"`
	HasTriggered    bool    `json:"webhook_has_triggered"`
}
