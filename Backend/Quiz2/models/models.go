package models

type PokemonStats struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"stat"`
}

type PokemonResponse struct {
	Stats []PokemonStats `json:"stats"`
}

type PokemonSprites struct {
	BackDefault      string `json:"back_default"`
	BackFemale       string `json:"back_female"`
	BackShiny        string `json:"back_shiny"`
	BackShinyFemale  string `json:"back_shiny_female"`
	FrontDefault     string `json:"front_default"`
	FrontFemale      string `json:"front_female"`
	FrontShiny       string `json:"front_shiny"`
	FrontShinyFemale string `json:"front_shiny_female"`
}

type PokemonFormResponse struct {
	Name    string         `json:"name"`
	Sprites PokemonSprites `json:"sprites"`
}

type CombinedResponse struct {
	Stats   []PokemonStats `json:"stats"`
	Name    string         `json:"name"`
	Sprites PokemonSprites `json:"sprites"`
}
