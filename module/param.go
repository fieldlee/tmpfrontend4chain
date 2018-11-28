package module

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Genarate struct {
	ID string `json:"id"` //文档_id
}

type OnlyID struct {
	ID string `json:"id"` //文档_id
}

type ChannelTx struct {
	ID          string   `json:"id"` //文档_id
	ChannelId   string   `json:"channelId"`
	IncludeOrgs []string `json:"includeOrgs"`
}

type ChainParam struct {
	ID        string `json:"id"`
	ChannelId string `json:"channelId"`
}

type ChainCodeParam struct {
	ID               string `json:"id"`
	ChannelId        string `json:"channelId"`
	ChainCodeName    string `json:"chaincodeName"`
	ChainCodeVersion string `json:"chaincodeVersion"`
}

type IPANDID struct {
	ID string `json:"id"` //文档_id
	IP string `json:"ip"` //需要检验的ip
}

type Sdk struct {
	SDKGitUrl string `json:"sdkGitUrl"`
	Path      string `json:"path"`
	ID        string `json:"id"`
	PID       string `json:"pid"`
}

type Token struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}

type PeerParam struct {
	PeerIp        string `json:"peerIp"`
	PeerId        string `json:"peerId"`
	PostPort      int    `json:"postPort"`
	EventPort     int    `json:"eventPort"`
	JoinCouch     bool   `json:"joinCouch"`
	CouchUsername string `json:"couchUsername"`
	CouchPassword string `json:"couchPassword"`
	CouchPort     int    `json:"couchPort"`
}
type OrgParam struct {
	OrgName    string      `json:"orgName"`
	OrgId      string      `json:"orgId"`
	PeerNumber int         `json:"peerNumber"`
	Peers      []PeerParam `json:"peers"`
	CaIp       string      `json:"caIp"`
	CaId       string      `json:"caId"`
	CaPort     int         `json:"caPort"`
}
type OrderParam struct {
	OrderName string `json:"orderName"`
	OrderId   string `json:"orderId"`
	OrderIp   string `json:"orderIp"`
	OrderPort int    `json:"orderPort"`
}
type DefineParam struct {
	ProjectName     string       `json:"projectName"`
	ID              string       `json:"id"`
	ProjectPassword string       `json:"projectPassword"`
	Domain          string       `json:"domain"`
	NetWork         string       `json:"network"`
	Consensus       string       `json:"consensus"`
	OrderName       string       `json:"orderName"`
	OrderId         string       `json:"orderId"`
	KafkaIp         string       `json:"kafkaIp"`
	Orgs            []OrgParam   `json:"orgs"`
	Orders          []OrderParam `json:"orders"`
}

type AddOrgParam struct {
	PID       string     `json:"pid"`
	ChannelId string     `json:"channelId"`
	AddOrgs   []OrgParam `json:"orgs"`
}

type ExplorderParam struct {
	ID string `json:"id"`
}

type TXParam struct {
	ID       string `json:"id"`
	Blocknum string `json:"blocknum"`
}
