package gmp

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Place struct {
	Candidates []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			Viewport struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		Icon                string `json:"icon"`
		IconBackgroundColor string `json:"icon_background_color"`
		IconMaskBaseUri     string `json:"icon_mask_base_uri"`
		Name                string `json:"name"`
		PlaceId             string `json:"place_id"`
		PlusCode            struct {
			CompoundCode string `json:"compound_code"`
			GlobalCode   string `json:"global_code"`
		} `json:"plus_code"`
		Types []string `json:"types"`
	} `json:"candidates"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

// FindPlaceDefaultFields return default fields for FindPlace search
// fields possibility :
//    formatted_phone_number, international_phone_number, opening_hours,
//    website, address_component, adr_address, business_status, formatted_address,
//    geometry, icon, icon_mask_base_uri, icon_background_color, name, permanently_closed (deprecated),
//    photo, place_id, plus_code, type, url, utc_offset, vicinity, price_level, rating, review, user_ratings_total
var FindPlaceDefaultFields = []string{"name", "formatted_address", "geometry"}

// FindPlaceByPhoneNumber takes a phone number and returns a place.
func (c *Client) FindPlaceByPhoneNumber(ctx context.Context, phoneNumber string, fields []string) (Place, error) {
	return c.findplace(ctx, "phonenumber", phoneNumber, fields)
}

// FindPlaceByString takes a text input and returns a place.
func (c *Client) FindPlaceByString(ctx context.Context, text string, fields []string) (Place, error) {
	return c.findplace(ctx, "textquery", text, fields)
}

func (c *Client) findplace(ctx context.Context, inputType, input string, fields []string) (p Place, err error) {
	if len(fields) == 0 {
		fields = FindPlaceDefaultFields
	}

	url := c.buildUrl("https://maps.googleapis.com/maps/api/place/findplacefromtext/json", map[string]string{
		"fields":    strings.Join(fields, ","),
		"input":     input,
		"inputtype": inputType,
	})

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		return
	}

	if p.Status == "INVALID_REQUEST" {
		return Place{}, errors.New(p.ErrorMessage)
	}

	return
}
