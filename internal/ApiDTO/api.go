package ApiDTO

type SignInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
}

type SignInfoJson struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type PickDTO struct {
	PickNum   int        `json:"pickNum"`
	AwardType string     `json:"awardType"`
	Picked    []SignInfo `json:"picked"`
}

type AvatarsDTO struct {
	AvatarNum  int      `json:"avatarNum"`
	AvatarUrls []string `json:"avatarUrls"`
}

type AssignDTO struct {
	Name  string `json:"name"`
	Award string `json:"award"`
}

type AllDTO struct {
	UserNum int        `json:"userNum"`
	Users   []SignInfo `json:"users"`
}

var AwardTypes = []string{"特等奖", "一等奖", "二等奖", "三等奖", "补抽", "未中奖"}
