package responseStructs

type GuildRoster struct {
	Links   Links    `json:"_links"`
	Guild   Guild    `json:"guild"`
	Members []Member `json:"members"`
}

type Guild struct {
	Key     Link    `json:"key"`
	Name    string  `json:"name"`
	ID      int     `json:"id"`
	Realm   Realm   `json:"realm"`
	Faction Faction `json:"faction"`
}

type Member struct {
	Character Character `json:"character"`
	Rank      int       `json:"rank"`
}

type Character struct {
	Key           Link          `json:"key"`
	Name          string        `json:"name"`
	ID            int           `json:"id"`
	Realm         Realm         `json:"realm"`
	Level         int           `json:"level"`
	PlayableClass PlayableClass `json:"playable_class"`
	PlayableRace  PlayableRace  `json:"playable_race"`
	Faction       Faction       `json:"faction"`
}

type PlayableClass struct {
	Key Link `json:"key"`
	ID  int  `json:"id"`
}

type PlayableRace struct {
	Key Link `json:"key"`
	ID  int  `json:"id"`
}
