package types

type GrantMsg struct {
	Grantee     string
	Granter     string
	SpendLimit  string
	Expiration  string
	Period      string
	PeriodLimit string
	AllowedMsg  []string
}

type RevokeGrantMsg struct {
	Grantee string
	Granter string
}

type QueryGrantMsg struct {
	Grantee string
	Granter string
}
