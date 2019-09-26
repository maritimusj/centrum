package ep6v2

import "fmt"

type Model struct {
	ID      string
	Version string
}

func (model *Model) fetchData(conn modbusClient) error {
	data, err := conn.ReadHoldingRegisters(0, 4)
	if err != nil {
		return err
	}

	model.ID = string([]byte{data[1], data[0], data[3], data[2], data[5], data[4]})
	model.Version = fmt.Sprintf("v%.2f", (float32(data[6])*100+float32(data[7]))/100)

	return nil
}
