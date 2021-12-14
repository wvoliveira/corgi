package location

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"net"
)

// IPv4ToLong convert ipv4 to uint32.
func IPv4ToLong(ip string) (i uint32) {
	_ = binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &i)
	return
}

// IPv6ToLong convert ipv6 to big.Int.
func IPv6ToLong(ip string) (b *big.Int) {
	IPv6Int := big.NewInt(0)
	return IPv6Int.SetBytes(net.ParseIP(ip).To16())
}
