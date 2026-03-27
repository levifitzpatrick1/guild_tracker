package responseStructs

type CharacterProfessions struct {
	Links       Links             `json:"_links"`
	Character   Character         `json:"character"`
	Primaries   []ProfessionGroup `json:"primaries"`
	Secondaries []ProfessionGroup `json:"secondaries"`
}

type ProfessionGroup struct {
	Profession ProfessionDetails `json:"profession"`
	Tiers      []ProfessionTier  `json:"tiers"`
}

type ProfessionDetails struct {
	Key  Link   `json:"key"`
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type ProfessionTier struct {
	SkillPoints    int      `json:"skill_points"`
	MaxSkillPoints int      `json:"max_skill_points"`
	Tier           TierInfo `json:"tier"`
	KnownRecipes   []Recipe `json:"known_recipes,omitempty"`
}

type TierInfo struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Recipe struct {
	Key  Link   `json:"key"`
	Name string `json:"name"`
	ID   int    `json:"id"`
}
