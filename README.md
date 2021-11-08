# Google Maps Platform

## Find Place
```go
client := gapi.New("API_KEY")

p, err := client.FindPlaceByString(context.Background(), "Champ de Mars, 5 Av. Anatole France, 75007 Paris", gapi.FindPlaceDefaultFields)
if err != nil {
	log.Fatal(err.Error())
}
```