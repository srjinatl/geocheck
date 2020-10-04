// # Database URL
// https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=XXeKEmWLYrOFUVdA&suffix=tar.gz

package maxmind

import (
    "errors"
    "fmt"
    "net"

    maxminddb "github.com/oschwald/geoip2-golang"
)


type GeoLocationProvider interface {
    GetCountryForIp(ip string) (string, error)
    Close()
}

type GeoLocationProviderImpl struct {
    maxMindReader *maxminddb.Reader
}

func NewGeoLocationProvider(dbName string) (db GeoLocationProvider, err error) {
    reader, err := maxminddb.Open(dbName)
    if err != nil {
        return
    }
    db = &GeoLocationProviderImpl{maxMindReader: reader}
    return
}

func (db *GeoLocationProviderImpl) Close() {
    if db.maxMindReader != nil {
        db.maxMindReader.Close()
    }
}

func (db *GeoLocationProviderImpl) GetCountryForIp(ip string) (country string, err error) {
    // parse ip - check for nil result if ip is badly formed
    ipAddr := net.ParseIP(ip)
    if ipAddr == nil {
        return "", errors.New("invalid ip address - unable to obtain country")
    }

    dbCountry, err := db.maxMindReader.Country(ipAddr)
    if err != nil {
        return
    }

    if dbCountry == nil {
        err = fmt.Errorf("country not available for %s", ip)
        return
    }
    fmt.Println(dbCountry)

    return dbCountry.Country.IsoCode, nil
}
