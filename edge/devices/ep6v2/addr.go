package ep6v2

import (
	"encoding/binary"
	"fmt"

	"github.com/maritimusj/centrum/edge/devices/modbus"
)

type IPAddr [4]uint8

func (ip *IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

type MAC [6]uint8

func (mac *MAC) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

type Addr struct {
	Ip      IPAddr
	Mask    IPAddr
	Gateway IPAddr
	Mac     MAC
}

func (addr *Addr) fetchData(conn modbus.Client) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			retErr = fmt.Errorf("unexpect error: %#v", err)
			return
		}
	}()

	data, err := conn.ReadHoldingRegisters(0x0020, 18)
	if err != nil {
		return err
	}

	addr.Ip[0] = uint8(binary.BigEndian.Uint16(data[0:]))
	addr.Ip[1] = uint8(binary.BigEndian.Uint16(data[2:]))
	addr.Ip[2] = uint8(binary.BigEndian.Uint16(data[4:]))
	addr.Ip[3] = uint8(binary.BigEndian.Uint16(data[6:]))

	addr.Mask[0] = uint8(binary.BigEndian.Uint16(data[8:]))
	addr.Mask[1] = uint8(binary.BigEndian.Uint16(data[10:]))
	addr.Mask[2] = uint8(binary.BigEndian.Uint16(data[12:]))
	addr.Mask[3] = uint8(binary.BigEndian.Uint16(data[14:]))

	addr.Gateway[0] = uint8(binary.BigEndian.Uint16(data[16:]))
	addr.Gateway[1] = uint8(binary.BigEndian.Uint16(data[18:]))
	addr.Gateway[2] = uint8(binary.BigEndian.Uint16(data[20:]))
	addr.Gateway[3] = uint8(binary.BigEndian.Uint16(data[22:]))

	for i := range addr.Mac {
		addr.Mac[i] = byte(binary.BigEndian.Uint16(data[24+i*2:]))
	}

	return nil
}
