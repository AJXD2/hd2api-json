package main

type Faction struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

type Warbond struct {
	Index           int    `json:"index"`
	Name            string `json:"name"`
	Id              string `json:"id"`
	CreditsToUnlock int    `json:"credits_to_unlock"`
}
