package main

import (
	"CryptoNotify/api"
	lib "CryptoNotify/coreLib"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/api/option"
)

var Version string = "v1"                    //Version of service
var Root string = "/cryptonotify/" + Version //URL root path
var VolHook = Root + "/trends"               //Registration of volume webhooks
var PointHook = Root + "/pricetarget"        //Registration of price/volume point webhooks
var PortFolio = Root + ""                    //Registration of portfolio webhooks

var mock bool = true //If mocking the api or not

func port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// Function for polling, caching the response and then going through the webhooks
func cryptoPolling() {
	for {
		time.Sleep(15 * time.Second)

		if mock {
			lib.GetMock()
		} else {
			//REAL API
		}
		lib.UpdateInternalMap()
		api.CheckVolumeWebhooks()
		api.CheckPriceWebhooks()
	}
}

func setupRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,          //Logging API requests
		middleware.RedirectSlashes, //For redirecting slashed URLs
	)

	router.Get("/", api.RootPage)
	router.Get(Root+"/currency/{currency_code}", api.CurrencyHandler)
	//router.Get(Root+"/country/{country_name}", api.CasesCountry)
	//router.Get(Root+"/policy/{country_name}", api.PolicyEnd)
	//router.Get(Root+"/diag", api.Diag)

	router.Route(VolHook, func(r chi.Router) {
		r.Post("/", api.VolumeWebhookReg) //Handling of webhooks to .../notifications
		r.Delete("/{id}", api.WebhookVolumeDel)
		r.Delete("/", api.WebhookVolumeDel)
		r.Get("/{id}", api.GetVolumeWebhook)
		r.Get("/", api.AllVolumeWebhooks)

	})

	router.Route(PointHook, func(r chi.Router) {
		r.Post("/", api.PriceWebhookReg) //Handling of webhooks to .../notifications
		r.Delete("/{id}", api.WebhookPriceDel)
		r.Delete("/", api.WebhookPriceDel)
		r.Get("/{id}", api.GetPriceWebhook)
		r.Get("/", api.AllPriceWebhooks)

	})

	return router
}

func main() {
	fmt.Println("running:")

	//Firebase initialization start
	api.Ctx = context.Background()
	sa := option.WithCredentialsFile("./cloud-project-dd1b4-firebase-adminsdk-py7gd-a693745d88.json")
	app, err := firebase.NewApp(api.Ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	api.Client, err = app.Firestore(api.Ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer api.Client.Close()
	//Firebase initialization end
	///TEST
	/*client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "100")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "fa238227-46eb-4bc2-8e66-37c50f341fdb")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))*/
	///
	lib.GetMock()
	lib.UpdateInternalMap()
	go cryptoPolling()

	port := port()
	router := setupRoutes()

	log.Fatal(http.ListenAndServe(":"+port, router))
}
