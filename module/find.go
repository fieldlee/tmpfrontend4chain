package module

type UserFindResponse struct {
	Docs []User `json:"docs"`
}

type SetupFindResponse struct {
	Docs []Define `json:"docs"`
}

type CCFindResponse struct {
	Docs []ChainCode `json:"docs"`
}

type RandomFindResponse struct {
	Docs []VerifyCode `json:"docs"`
}

type TelFindResponse struct {
	Docs []TelCode `json:"docs"`
}

type AnnFindResponse struct {
	Docs []Announce `json:"docs"`
}

type FeedFindResponse struct {
	Docs []Feedback `json:"docs"`
}

type LogFindResponse struct {
	Docs []Log `json:"docs"`
}
