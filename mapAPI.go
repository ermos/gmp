package gmp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	StaticMapsTypeRoadMap = "roadmap"
	StaticMapsTypeSatellite = "satellite"
	StaticMapsTypeTerrain = "terrain"
	StaticMapsTypeHybrid = "hybrid"
	StaticMapsRegionFrance = "fr"
	StaticMapsFormatPNG = "png"
	StaticMapsFormatPNG8 = "png8"
	StaticMapsFormatPNG32 = "png32"
	StaticMapsFormatGIF = "gif"
	StaticMapsFormatJPG = "jpg"
	StaticMapsFormatJPGBaseline = "jpg-baseline"
	StaticMapsMarkerSizeDefault = ""
	StaticMapsMarkerSizeTiny = "tiny"
	StaticMapsMarkerSizeMid = "mid"
	StaticMapsMarkerSizeSmall = "small"
)

type StaticMaps struct {
	// Default MapType is RoadMap
	MapType string
	// Default Format is PNG
	Format 	string
	Region string
	Center string
	Zoom int
	Size StaticMapsSize
	Markers []StaticMapsMarker
	//TODO API Path
}

type StaticMapsSize struct{
	Width int
	Height int
}

type StaticMapsMarker struct {
	Color string
	Label string
	Size  string
	IconURL string
}

func (c *Client) GetStaticMapsURL(sm StaticMaps) string {
	queries := make(map[string]string)

	if sm.Region != "" {
		queries["region"] = sm.Region
	}

	if sm.Format != "" {
		queries["format"] = sm.Format
	}

	if sm.Center != "" {
		queries["center"] = sm.Center
	}

	if sm.Zoom != 0 {
		queries["zoom"] = strconv.Itoa(sm.Zoom)
	}

	if sm.Size.Width != 0 && sm.Size.Height != 0 {
		queries["size"] = fmt.Sprintf("%dx%d", sm.Size.Width, sm.Size.Height)
	}

	if sm.MapType != "" {
		queries["maptype"] = sm.MapType
	}

	if len(sm.Markers) != 0 {
		var markers string
		for i, v := range sm.Markers {
			if i != 0 {
				markers += "&markers="
			}

			var val []string

			if v.IconURL != "" {
				val = append(val, fmt.Sprintf("icon:%s", v.Size))
			}

			if v.Color != "" {
				val = append(val, fmt.Sprintf("color:%s", v.Color))
			}

			if v.Label != "" {
				val = append(val, fmt.Sprintf("label:%s", v.Label))
			}

			if v.Size != "" {
				val = append(val, fmt.Sprintf("size:%s", v.Size))
			}

			markers += strings.Join(val, "|")
		}

		queries["markers"] = markers
	}

	queries["key"] = c.apiKey

	return c.buildUrl("https://maps.googleapis.com/maps/api/staticmap", queries)
}

func (c *Client) GetStaticMapsImage(ctx context.Context, sm StaticMaps) (img []byte, err error) {
	url := c.GetStaticMapsURL(sm)
	client := &http.Client {}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("receive %d response from %s", res.StatusCode, url)
		return
	}

	return ioutil.ReadAll(res.Body)
}