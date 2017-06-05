package json

import (
	"database/sql/driver"
	"fmt"
	"net"
)

// This type of time would be serialized to UNIX time for MarshalJSON()
type IP net.IP

// Implement MarshalJSON() here if needed
// func (ip IP) MarshalJSON() ([]byte, error) {
// }

// UnmarshalJSON parses the JSON-encoded IP string
func (ip *IP) UnmarshalJSON(data []byte) error {
	str := UnmarshalToJson(data).MustString()
	parsed := net.ParseIP(str)
	if parsed == nil {
		return fmt.Errorf("Cannot parse IP string: %s\n", str)
	}
	*ip = IP(parsed)
	return nil
}

func (ip IP) Value() (driver.Value, error) {
	// If ip is not an IPv4 address, To4 returns nil.
	if v4 := net.IP(ip).To4(); v4 != nil {
		return []byte(v4), nil
	}
	// ip must be IPv6 address
	return []byte(net.IP(ip).To16()), nil
}
