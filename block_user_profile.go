package sii

func init() {
	RegisterBlock(&UserProfile{})
}

type UserProfile struct {
	Face             int64      `sii:"face"`
	Brand            Ptr        `sii:"brand"`
	MapPath          string     `sii:"map_path"`
	Logo             Ptr        `sii:"logo"`
	CompanyName      string     `sii:"company_name"`
	Male             bool       `sii:"male"`
	CachedExperience int64      `sii:"cached_experience"`
	CachedDistance   int64      `sii:"cached_distance"`
	UserData         []RawValue `sii:"user_data"`
	ActiveMods       []string   `sii:"active_mods"`
	Customization    int64      `sii:"customization"` // ??? Maybe bit-flags?
	CachedStats      []int64    `sii:"cached_stats"`
	CachedDiscovery  []int64    `sii:"cached_discovery"`
	Version          int64      `sii:"version"`
	OnlineUserName   string     `sii:"online_user_name"`
	OnlinePassword   string     `sii:"online_password"`
	ProfileName      string     `sii:"profile_name"`
	CreationTime     int64      `sii:"creation_time"`
	SaveTime         int64      `sii:"save_time"`

	blockName string
}

func (UserProfile) Class() string { return "user_profile" }

func (u *UserProfile) Init(class, name string) {
	u.blockName = name
}

func (u UserProfile) Name() string { return u.blockName }
