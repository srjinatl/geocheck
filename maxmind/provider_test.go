package maxmind

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

const dbName = "../data/GeoLite2-Country.mmdb"

func TestCountryByIP(t *testing.T) {

    type testCase struct {
        name string
        ip string
        expectedCountry string
        expectedError error
    }
    tests := []testCase{
        {
            name: "TestUnroutableIP",
            ip: "192.168.1.2",
            expectedCountry: "",
            expectedError: nil,
        },
        {
            name: "TestValidUSIP",
            ip: "67.166.247.50",
            expectedCountry: "US",
            expectedError: nil,
        },
        {
            name: "TestInvalidIP",
            ip: "a.b.c.d",
            expectedCountry: "",
            expectedError: errors.New("invalid ip address - unable to obtain country"),
        },
    }

    // get db
    db, err := NewGeoLocationProvider(dbName)
    require.Nil(t, err, "Unexpected error creating db")
    require.NotNil(t, db, "Unexpected nil db returned from constructor")

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            country, err := db.GetCountryForIp(test.ip)
            assert.Equal(t, test.expectedError, err)
            assert.Equal(t, test.expectedCountry, country)
        })
    }
}
