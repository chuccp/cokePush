package user

type Page struct {
	Num  int
	List []*PageUser
}

func (p *Page)Size() int {
	return len(p.List)
}

func NewPage() *Page {
	return &Page{0, make([]*PageUser, 0)}
}

type PageUser struct {
	UserName       string
	MachineAddress string
	CreateTime     string
}

func NewPageUser(UserName string, MachineAddress string, CreateTime string) *PageUser {
	return &PageUser{UserName: UserName, MachineAddress: MachineAddress, CreateTime: CreateTime}
}
