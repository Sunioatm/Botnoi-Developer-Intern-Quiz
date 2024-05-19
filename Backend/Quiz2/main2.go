package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sunioatm/main2/models"

	"github.com/gofiber/fiber/v2"
)

func getPokemonData(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		var requestBody struct {
			ID interface{} `json:"id"`
		}

		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}
		id = fmt.Sprintf("%v", requestBody.ID)
	}

	if id == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ID is required")
	}

	pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", id)
	pokemonResp, err := http.Get(pokemonURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch Pokémon data")
	}
	defer pokemonResp.Body.Close()
	pokemonBody, err := io.ReadAll(pokemonResp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read Pokémon response")
	}
	var pokemon models.PokemonResponse
	json.Unmarshal(pokemonBody, &pokemon)

	formURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-form/%s/", id)
	formResp, err := http.Get(formURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch Pokémon form data")
	}
	defer formResp.Body.Close()
	formBody, err := io.ReadAll(formResp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read Pokémon form response")
	}
	var pokemonForm models.PokemonFormResponse
	json.Unmarshal(formBody, &pokemonForm)

	result := models.CombinedResponse{
		Stats:   pokemon.Stats,
		Name:    pokemonForm.Name,
		Sprites: pokemonForm.Sprites,
	}

	return c.JSON(result)
}

func main() {
	app := fiber.New()

	// I'm not sure why POST?
	app.Post("/pokemon", getPokemonData)

	// GET is more appropriate
	app.Get("/pokemon/:id", getPokemonData)

	log.Fatal(app.Listen("127.0.0.1:3000"))
}
