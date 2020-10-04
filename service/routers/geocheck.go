package routers

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/srjinatl/geocheck/log"
    "github.com/srjinatl/geocheck/maxmind"
    "go.uber.org/zap"
)

// GeoCheckHandler .
type GeoCheckHandler struct {
    log *log.Logger
    provider maxmind.GeoLocationProvider
}

type geoCheckInput struct {
    InputAddress string `json:"ip_address"`
    InputCountryWhiteList string `json:"white_list"`
}

type geoCheckResult struct {
    Status string `json:"status"`
    Country string `json:"country"`
}

type geoResponse struct {
    Input geoCheckInput `json:"input"`
    Result geoCheckResult `json:"result"`
    ErrorMsg string `json:"error_msg,omitempty"`
}

// NewGeoCheckHandler returns handler for geo check endpoints
func NewGeoCheckHandler(log *log.Logger, provider maxmind.GeoLocationProvider) *GeoCheckHandler{
    return &GeoCheckHandler{log: log, provider: provider}
}

// GeoCheck validates whether specified ip address is within list of specified countries
func (h *GeoCheckHandler) GeoCheck(c *gin.Context) {
    ipAddr := c.Query("ipaddr")
    countries := c.Query("countries")

    // make sure ip address was specified
    if ipAddr == "" {
        c.JSON(http.StatusBadRequest,
            createResponse(ipAddr, "", "Error", "", "ipaddr url parameter required"))
        return
    }

    // make sure countries has been specified - split into slice
    if countries == "" {
        c.JSON(http.StatusBadRequest,
            createResponse(ipAddr, countries, "Error", "", "countries url parameter required"))
        return
    }

    // split country list for check
    countryList := strings.Split(countries, ",")
    h.log.Zap.Debug("Country list element count", zap.Int("countryCount", len(countryList)))

    // call provider to get country for ip
    ipCountry, err := h.provider.GetCountryForIp(ipAddr)
    if err != nil {
        c.JSON(http.StatusOK, createResponse(ipAddr, countries, "Error", ipCountry, err.Error()))
        return
    }

    // check country for membership in white list
    resultStr := "Invalid"
    if isCountryInWhitelist(ipCountry, countryList) {
        resultStr = "Valid"
    }
    c.JSON(http.StatusOK, createResponse(ipAddr, countries, resultStr, ipCountry, ""))
}

func createResponse(ipAddr, countries, resultStr, ipAddrCountry, errMsg string) geoResponse {
    resp := geoResponse{
        Input:    geoCheckInput{
            InputAddress:          ipAddr,
            InputCountryWhiteList: countries,
        },
        Result:   geoCheckResult{
            Status: resultStr,
            Country: ipAddrCountry,
        },
    }
    if errMsg != "" {
        resp.ErrorMsg = errMsg
    }
    return resp
}

func isCountryInWhitelist(inputCountry string, whiteList []string) (inWhiteList bool) {
    for _, country := range whiteList {
        if country == inputCountry {
            inWhiteList = true
            break
        }
    }
    return
}