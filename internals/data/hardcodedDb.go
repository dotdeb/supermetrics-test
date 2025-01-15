package data

type HardCodedDb struct {
}

func (h *HardCodedDb) GetUsers() []User {
	return []User{
		{
			Id:       "1",
			Username: "john.doe",
		},
		{
			Id:       "2",
			Username: "jane.smith",
		},
	}
}
