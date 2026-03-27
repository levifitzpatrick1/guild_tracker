package responseStructs

type GuildProfile struct {
	Links             Links   `json:"_links"`
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Faction           Faction `json:"faction"`
	AchievementPoints int     `json:"achievement_points"`
	MemberCount       int     `json:"member_count"`
	Realm             Realm   `json:"realm"`
	Crest             Crest   `json:"crest"`
	Roster            Link    `json:"roster"`
	Achievements      Link    `json:"achievements"`
	CreatedTimestamp  int64   `json:"created_timestamp"`
	Activity          Link    `json:"activity"`
	NameSearch        string  `json:"name_search"`
}

type Links struct {
	Self Link `json:"self"`
}

type Link struct {
	Href string `json:"href"`
}

type Realm struct {
	Key  Link   `json:"key"`
	Name string `json:"name,omitempty"`
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

type Faction struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
}

type Crest struct {
	Emblem     Emblem     `json:"emblem"`
	Border     Border     `json:"border"`
	Background Background `json:"background"`
}

type Emblem struct {
	ID    int   `json:"id"`
	Media Media `json:"media"`
	Color Color `json:"color"`
}

type Border struct {
	ID    int   `json:"id"`
	Media Media `json:"media"`
	Color Color `json:"color"`
}

type Background struct {
	Color Color `json:"color"`
}

type Media struct {
	Key Link `json:"key"`
	ID  int  `json:"id"`
}

type Color struct {
	ID   int  `json:"id"`
	RGBA RGBA `json:"rgba"`
}

type RGBA struct {
	R int     `json:"r"`
	G int     `json:"g"`
	B int     `json:"b"`
	A float64 `json:"a"`
}
