package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

type discordInteractionRequest struct {
	Type int             `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

const (
	discordInteractionRequestPing                           = 1
	discordInteractionRequestApplicationCommand             = 2
	discordInteractionRequestMessageComponent               = 3
	discordInteractionRequestApplicationCommandAutocomplete = 4
	discordInteractionRequestModalSubmit                    = 5
)

type discordInteractionRequestApplicationCommandData struct {
	Name    string                                              `json:"name"`
	Options []discordInteractionRequestApplicationCommandOption `json:"options,omitempty"`
}

func (command *discordInteractionRequestApplicationCommandData) getStringOption(name string) string {
	for _, option := range command.Options {
		if option.Name != name {
			continue
		}

		return fmt.Sprintf("%v", option.Value)
	}

	return ""
}

type discordInteractionRequestApplicationCommandOption struct {
	Name  string `json:"name"`
	Value any    `json:"value,omitempty"`
}

type discordInteractionResponse struct {
	Type int `json:"type"`
	Data any `json:"data,omitempty"`
}

const (
	discordInteractionResponsePong                                 = 1
	discordInteractionResponseChannelMessageWithSource             = 4
	discordInteractionResponseDeferredChannelMessageWithSource     = 5
	discordInteractionResponseDeferredUpdateMessage                = 6
	discordInteractionResponseUpdateMessage                        = 7
	discordInteractionResponseApplicationCommandAutocompleteResult = 8
	discordInteractionResponseModal                                = 9
	discordInteractionResponsePremiumRequired                      = 10
)

type discordInteractionResponseMessageData struct {
	Content string `json:"content,omitempty"`
}

type discordInteractionAuth struct {
	applicationPublicKey ed25519.PublicKey
}

func newDiscordInteractionAuth() discordInteractionAuth {
	applicationPublicKey := make([]byte, ed25519.PublicKeySize)
	if _, err := hex.Decode(applicationPublicKey, []byte(os.Getenv("DISCORD_APPLICATION_PUBLIC_KEY"))); err != nil {
		log.Fatalf("failed to parse environment variable DISCORD_APPLICATION_PUBLIC_KEY: %v", err)
	}

	return discordInteractionAuth{
		applicationPublicKey: applicationPublicKey,
	}
}

func (auth *discordInteractionAuth) verify(ctx *gin.Context) (bool, error) {
	signature := make([]byte, ed25519.SignatureSize)
	if _, err := hex.Decode(signature, []byte(ctx.GetHeader("X-Signature-Ed25519"))); err != nil {
		return false, fmt.Errorf("failed to parse header X-Signature-Ed25519: %w", err)
	}

	timestamp := ctx.GetHeader("X-Signature-Timestamp")

	body, err := peekBody(ctx)
	if err != nil {
		return false, err
	}

	return ed25519.Verify(auth.applicationPublicKey, []byte(timestamp+string(body)), signature), nil
}

func (auth *discordInteractionAuth) middleware(ctx *gin.Context) {
	if verified, err := auth.verify(ctx); err != nil {
		ctx.String(http.StatusBadRequest, "failed to parse credentials: %v", err)
		ctx.Abort()
	} else if !verified {
		ctx.String(http.StatusUnauthorized, "access denied")
		ctx.Abort()
	}
}

func discordMount(engine *gin.Engine) {
	discordInteractionAuth := newDiscordInteractionAuth()

	engine.POST("/webhooks/discord/interactions", discordInteractionAuth.middleware, discordHandleInteraction)
}

func discordHandleInteraction(ctx *gin.Context) {
	var request discordInteractionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.String(http.StatusBadRequest, "failed to parse request body: %v", err)

		return
	}

	switch request.Type {
	case discordInteractionRequestPing:
		discordHandleInteractionRequestPing(ctx)
	case discordInteractionRequestApplicationCommand:
		discordHandleInteractionRequestApplicationCommand(ctx, &request)
	default:
		ctx.String(http.StatusNotImplemented, "interaction type %v not implemented", request.Type)
	}
}

func discordHandleInteractionRequestPing(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &discordInteractionResponse{
		Type: discordInteractionResponsePong,
		Data: nil,
	})
}

func discordHandleInteractionRequestApplicationCommand(ctx *gin.Context, request *discordInteractionRequest) {
	var command discordInteractionRequestApplicationCommandData
	if err := json.Unmarshal(request.Data, &command); err != nil {
		ctx.String(http.StatusBadRequest, "failed to parse interaction data: %v", err)

		return
	}

	switch command.Name {
	case "roll":
		ctx.JSON(http.StatusOK, discordHandleCommandRoll(&command))
	default:
		ctx.String(http.StatusBadRequest, "unrecognized command: %v", command.Name)
	}
}

func discordHandleCommandRoll(command *discordInteractionRequestApplicationCommandData) *discordInteractionResponse {
	return &discordInteractionResponse{
		Type: discordInteractionResponseChannelMessageWithSource,
		Data: discordInteractionResponseMessageData{
			Content: roll(command.getStringOption("formula")),
		},
	}
}

func discordEscapeMarkdown(input string) string {
	var output strings.Builder

	for _, currentRune := range input {
		// Skip anything non-printable. Avoids things like ANSI escape
		// codes which Discord does actually support.
		if !unicode.IsGraphic(currentRune) {
			continue
		}

		// "\" + any ascii punctuation character works, so we do it to
		// all in case Discord ever desides to extend their Markdown
		// even more.
		//
		// "\" + any other character will include the "\" literally so
		// we cannot just insert "\" before every character in the input
		// string and call it a day.
		if strings.ContainsRune("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", currentRune) {
			output.WriteByte('\\')
		}

		output.WriteRune(currentRune)
	}

	return output.String()
}
