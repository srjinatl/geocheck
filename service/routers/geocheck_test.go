package routers

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/srjinatl/geocheck/log"
    "github.com/srjinatl/geocheck/maxmind"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGeoCheck_NoUrlParms(t *testing.T) {
    router := gin.Default()
    // create handler
    provider, err := maxmind.NewGeoLocationProvider("../../data/GeoLite2-Country.mmdb")
    require.Nil(t, err, "Unexpected error creating new geo location provider")
    require.NotNil(t, provider, "Unexpected nil provider returned")
    handler := NewGeoCheckHandler(log.NewLogger("test", true), provider)
    require.NotNil(t, handler, "Unexpected nil handler returned")

    router.GET("/geocheck", handler.GeoCheck)

    // perform get without any url parms
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/geocheck", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    expected := "{\"input\":{\"ip_address\":\"\",\"white_list\":\"\"},\"result\":{\"status\":\"Error\",\"country\":\"\"},\"error_msg\":\"ipaddr url parameter required\"}"
    assert.Equal(t, expected, w.Body.String(), "Unexpected response content")
}

func TestGeoCheck_BadIPAddr(t *testing.T) {
    router := gin.Default()
    // create handler
    provider, err := maxmind.NewGeoLocationProvider("../../data/GeoLite2-Country.mmdb")
    require.Nil(t, err, "Unexpected error creating new geo location provider")
    require.NotNil(t, provider, "Unexpected nil provider returned")
    handler := NewGeoCheckHandler(log.NewLogger("test", true), provider)
    require.NotNil(t, handler, "Unexpected nil handler returned")

    router.GET("/geocheck", handler.GeoCheck)

    // perform get without any url parms
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/geocheck?countries=US,GB&ipaddr=a.b.c.d", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    expected := "{\"input\":{\"ip_address\":\"a.b.c.d\",\"white_list\":\"US,GB\"},\"result\":{\"status\":\"Error\",\"country\":\"\"},\"error_msg\":\"invalid ip address - unable to obtain country\"}"
    assert.Equal(t, expected, w.Body.String(), "Unexpected response content")
}

func TestGeoCheck_IPInWhiteList(t *testing.T) {
    router := gin.Default()
    // create handler
    provider, err := maxmind.NewGeoLocationProvider("../../data/GeoLite2-Country.mmdb")
    require.Nil(t, err, "Unexpected error creating new geo location provider")
    require.NotNil(t, provider, "Unexpected nil provider returned")
    handler := NewGeoCheckHandler(log.NewLogger("test", true), provider)
    require.NotNil(t, handler, "Unexpected nil handler returned")

    router.GET("/geocheck", handler.GeoCheck)

    // perform get without any url parms
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/geocheck?countries=US,GB&ipaddr=67.166.247.56", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    expected := "{\"input\":{\"ip_address\":\"67.166.247.56\",\"white_list\":\"US,GB\"},\"result\":{\"status\":\"Valid\",\"country\":\"US\"}}"
    assert.Equal(t, expected, w.Body.String(), "Unexpected response content")
}

func TestGeoCheck_IPNotInWhiteList(t *testing.T) {
    router := gin.Default()
    // create handler
    provider, err := maxmind.NewGeoLocationProvider("../../data/GeoLite2-Country.mmdb")
    require.Nil(t, err, "Unexpected error creating new geo location provider")
    require.NotNil(t, provider, "Unexpected nil provider returned")
    handler := NewGeoCheckHandler(log.NewLogger("test", true), provider)
    require.NotNil(t, handler, "Unexpected nil handler returned")

    router.GET("/geocheck", handler.GeoCheck)

    // perform get without any url parms
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/geocheck?countries=SE,GB&ipaddr=67.166.247.56", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    expected := "{\"input\":{\"ip_address\":\"67.166.247.56\",\"white_list\":\"SE,GB\"},\"result\":{\"status\":\"Invalid\",\"country\":\"US\"}}"
    assert.Equal(t, expected, w.Body.String(), "Unexpected response content")
}
