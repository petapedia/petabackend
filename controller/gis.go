package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetRegion(respw http.ResponseWriter, req *http.Request) {
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid "
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(req)
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	var longlat model.LongLat
	err = json.NewDecoder(req.Body).Decode(&longlat)
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Body tidak valid"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	filter := bson.M{
		"border": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longlat.Longitude, longlat.Latitude},
				},
			},
		},
	}
	region, err := atdb.GetOneDoc[model.Region](config.Mongoconn, "region", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, region)
		return
	}
	at.WriteJSON(respw, http.StatusOK, region)
}

func GetGEOJSONRoadsIntersectRegion(respw http.ResponseWriter, req *http.Request) {
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid "
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(req)
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	var longlat model.LongLat
	err = json.NewDecoder(req.Body).Decode(&longlat)
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Body tidak valid"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	filter := bson.M{
		"border": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longlat.Longitude, longlat.Latitude},
				},
			},
		},
	}
	region, err := atdb.GetOneDoc[model.Region](config.Mongoconn, "region", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, region)
		return
	}
	//mendapatkan intersect dari roadsa
	filter = bson.M{
		"geometry": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Polygon",
					"coordinates": region.Border.Coordinates,
				},
			},
		},
	}
	roadsintersectregion, err := atdb.GetAllDoc[[]model.Roads](config.Mongoconn, "roads", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, roadsintersectregion)
		return
	}
	geojsonroads := model.GeoJSONRoads{
		Type:     "FeatureCollection",
		Features: roadsintersectregion,
	}
	at.WriteJSON(respw, http.StatusOK, geojsonroads)
}

func GetRoads(respw http.ResponseWriter, req *http.Request) {
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid "
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(req)
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	var longlat model.LongLat
	err = json.NewDecoder(req.Body).Decode(&longlat)
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Body tidak valid"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	filter := bson.M{
		"geometry": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longlat.Longitude, longlat.Latitude},
				},
				"$maxDistance": 600,
			},
		},
	}
	roads, err := atdb.GetAllDoc[[]model.Roads](config.Mongoconn, "roads", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, roads)
		return
	}
	at.WriteJSON(respw, http.StatusOK, roads)
}

func GetGeoJSONRoads(respw http.ResponseWriter, req *http.Request) {
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid "
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(req)
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	var longlat model.LongLat
	err = json.NewDecoder(req.Body).Decode(&longlat)
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Body tidak valid"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	filter := bson.M{
		"geometry": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longlat.Longitude, longlat.Latitude},
				},
				"$maxDistance": 600,
			},
		},
	}
	roads, err := atdb.GetAllDoc[[]model.Roads](config.Mongoconn, "roads", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, roads)
		return
	}
	geojsonroads := model.GeoJSONRoads{
		Type:     "FeatureCollection",
		Features: roads,
	}
	at.WriteJSON(respw, http.StatusOK, geojsonroads)
}
