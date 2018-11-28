package module

type Peer struct {
	PeerParam
	ContainerId      string `json:containerId`
	CouchId          string `json:"couchId"`
	CouchContainerId string `json:couchContainerId`
}
type Org struct {
	OrgName     string `json:"orgName"`
	OrgId       string `json:"orgId"`
	AnchorIp    string `json:"anchorIp"`
	AnchorPort  int    `json:"anchorPort"`
	PeerNumber  int    `json:"peerNumber"`
	Peers       []Peer `json:"peers"`
	CaIp        string `json:"caIp"`
	CaId        string `json:"caId"`
	CaPort      int    `json:"caPort"`
	CaUser      string `json:"caUser"`
	CaPwd       string `json:"caPwd"`
	ContainerId string `json:containerId`
}
type Order struct {
	OrderParam
	ContainerId string `json:containerId`
}
type Channel struct {
	ChannelId   string   `json:"channelId"`
	CreateTime  int64    `json:"createTime"`
	IncludeOrgs []string `json:"includeOrgs"`
	ChainCodes  []CC     `json:"chaincodes"`
}
type CC struct {
	CCGitUrl  string `json:"ccGitUrl"`
	CCName    string `json:"ccName"`
	CCPath    string `json:"ccPath"`
	CCVersion string `json:"ccVersion"`
	Using     bool   `json:"using"`
}
type ChainCode struct {
	ChannelId string `json:"channelId"`
	CCGitUrl  string `json:"ccGitUrl"`
	CCName    string `json:"ccName"`
	CCPath    string `json:"ccPath"`
	CCVersion string `json:"ccVersion"`
	ID        string `json:"id"`
	PID       string `json:"pid"`
}
type Define struct {
	ProjectName      string    `json:"projectName"`
	CreateTime       int64     `json:"createTime"`
	Manager          string    `json:"manager"`
	Domain           string    `json:"domain"`
	ProjectUser      string    `json:"projectUser"`
	ProjectPassword  string    `json:"projectPassword"`
	ExploderUser     string    `json:"exploderUser"`
	ExploderPassword string    `json:"exploderPassword"`
	Status           string    `json:"status"`
	NetWork          string    `json:"network"`
	Consensus        string    `json:"consensus"`
	ID               string    `json:"id"`
	OrderName        string    `json:"orderName"`
	OrderId          string    `json:"orderId"`
	KafkaIp          string    `json:"kafkaIp"`
	HasApp           bool      `json:"hasApp"`
	AddChannels      []Channel `json:"addChannels"`
	Orgs             []Org     `json:"orgs"`
	Orders           []Order   `json:"orders"`
}
type DefineList []Define

func (d DefineList) Len() int           { return len(d) }
func (d DefineList) Less(i, j int) bool { return d[i].CreateTime < d[j].CreateTime }
func (d DefineList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

type Config struct {
	Host            string                   `json:"host"`
	Port            string                   `json:"port"`
	CurOrgId        string                   `json:"curOrgId"`
	OrderId         string                   `json:"orderId"`
	Consensus       string                   `json:"consensus"`
	Manager         bool                     `json:"manager"`
	Orderers        []string                 `json:"orderers"`
	Orgs            []string                 `json:"orgs"`
	Channels        []map[string]interface{} `json:"channels"`
	ExpireTime      string                   `json:"jwt_expiretime"`
	ProjectName     string                   `json:"projectName"`
	CcSrcPath       string                   `json:"CC_SRC_PATH"`
	KeyValueStore   string                   `json:"keyValueStore"`
	EventWaitTime   string                   `json:"eventWaitTime"`
	CaUser          string                   `json:"caUser"`
	CaSecret        string                   `json:"caSecret"`
	ProjectPassword string                   `json:"projectPassword"`
	NetworkConfig   map[string]interface{}   `json:"network-config"`
}

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	TelNo    string `json:"telno"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Active   bool   `json:"active"`
	Avator   string `json:"avator"`
}

type RtUser struct {
	User
	ProjectNum int `json:"projectNum"`
}

type VerifyCode struct {
	Random     string `json:"random"`
	VerifyCode string `json:"verifycode"`
}

type TelCode struct {
	ID         string `json:"id"`
	TelNo      string `json:"telno"`
	VerifyCode string `json:"verifycode"`
}

type Announce struct {
	ID    string `json:"id"`
	Body  string `json:"body"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
}

type Feedback struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	TelNo string `json:"telno"`
	Body  string `json:"body"`
}

type Log struct {
	ID          string `json:"id"`
	UserName    string `json:"username"`
	Time        int64  `json:"time"`
	Interaction string `json:"interaction"`
	Data        string `json:"data"`
}
