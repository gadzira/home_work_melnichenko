// +build !bench

package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("empty data", func(t *testing.T) {
		emptydata := `{}`
		result, err := GetDomainStat(bytes.NewBufferString(emptydata), "")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("invalid json", func(t *testing.T) {
		errdata := `{"Id":6,"Name":"First Name","Username":"ugnetator1","Email":"sssddd@browsecat.com","Phone":"8-974-394-34-11","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
				{"Id":7,"Name":"Second Name","Username":"ugnetator2","Email":"rrrttt@linktype.com","Phone":8-495-224-22-22,"Password":"acSBF5","Address":"Russell Trail 61"}
				{"Id":8,"Name":"Third Name","Username":"ugnetator3","Email":"uuuiii@linktype.com","Phone":"8-495-224-22-22","Password":"tyVBN9","Address":"Monterey Park 39"}`
		result, err := GetDomainStat(bytes.NewBufferString(errdata), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 1,
			"linktype.com":  1,
		}, result)
	})
}
