package client

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

type DiscordClient interface {
}

type discordClient struct {
	s       *discordgo.Session
	appID   string
	guildID string
}

func NewDiscordClient(s *discordgo.Session) DiscordClient {
	appID := viper.GetString("app-id")
	guildID := viper.GetString("guild-id")
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	// Components are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:

			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
	_, err := s.ApplicationCommandCreate(appID, guildID, &discordgo.ApplicationCommand{
		Name:        "create-game",
		Description: "Создать игру",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
	_, err = s.ApplicationCommandCreate(appID, guildID, &discordgo.ApplicationCommand{
		Name: "selects",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "multi",
				Description: "Multi-item select menu",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "single",
				Description: "Single-item select menu",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "captains",
				Description: "Выбор капитана",
			},
		},
		Description: "Lo and behold: dropdowns are coming",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	return &discordClient{
		s:       s,
		appID:   viper.GetString("app-id"),
		guildID: viper.GetString("guild-id"),
	}
}

var (
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"fd_no": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Huh. I see, maybe some of these resources might help you?",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "📜",
									},
									Label: "Documentation",
									Style: discordgo.LinkButton,
									URL:   "https://discord.com/developers/docs/interactions/message-components#buttons",
								},
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"captains_select_ok": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Капитаны выбраны!",
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"captain_select": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			captain1ID := i.MessageComponentData().Values[0]
			captain2ID := i.MessageComponentData().Values[1]
			content := "Первый кэп @" + fmt.Sprintf("%v", i.MessageComponentData().Resolved.Users[captain1ID]) + "\n" + "Второй кэп @" + fmt.Sprintf("%v", i.MessageComponentData().Resolved.Users[captain2ID])
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
				Components: &[]discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								// Label is what the user will see on the button.
								Label: "Подтвердить выбор",
								// Style provides coloring of the button. There are not so many styles tho.
								Style: discordgo.SuccessButton,
								// Disabled allows bot to disable some buttons for users.
								Disabled: false,
								// CustomID is a thing telling Discord which data to send when this button will be pressed.
								CustomID: "captains_select_ok",
							},
							discordgo.Button{
								Label:    "Отмена",
								Style:    discordgo.DangerButton,
								Disabled: false,
								CustomID: "fd_no",
							},
						},
					},
				},
			})
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong" + err.Error(),
				})
				return
			}
			// err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
			// 	Data: &discordgo.InteractionResponseData{
			// 		Content:
			// 		Components: []discordgo.MessageComponent{
			// 			discordgo.ActionsRow{
			// 				Components: []discordgo.MessageComponent{
			// 					discordgo.Button{
			// 						// Label is what the user will see on the button.
			// 						Label: "Подтвердить выбор",
			// 						// Style provides coloring of the button. There are not so many styles tho.
			// 						Style: discordgo.SuccessButton,
			// 						// Disabled allows bot to disable some buttons for users.
			// 						Disabled: false,
			// 						// CustomID is a thing telling Discord which data to send when this button will be pressed.
			// 						CustomID: "captains_select_ok",
			// 					},
			// 					discordgo.Button{
			// 						Label:    "Отмена",
			// 						Style:    discordgo.DangerButton,
			// 						Disabled: false,
			// 						CustomID: "fd_no",
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// })
		},
	}
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"create-game": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Хостим?",
					// Buttons and other components are specified in Components field.
					Components: []discordgo.MessageComponent{
						// ActionRow is a container of all buttons within the same row.
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									// Label is what the user will see on the button.
									Label: "Да",
									// Style provides coloring of the button. There are not so many styles tho.
									Style: discordgo.SuccessButton,
									// Disabled allows bot to disable some buttons for users.
									Disabled: false,
									// CustomID is a thing telling Discord which data to send when this button will be pressed.
									CustomID: "fd_yes",
								},
								discordgo.Button{
									Label:    "Нет",
									Style:    discordgo.DangerButton,
									Disabled: false,
									CustomID: "fd_no",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"selects": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			minValues := 2
			var response *discordgo.InteractionResponse
			switch i.ApplicationCommandData().Options[0].Name {
			case "captains":
				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Выберите капитанов\n",
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.SelectMenu{
										MinValues:   &minValues,
										MaxValues:   2,
										MenuType:    discordgo.UserSelectMenu,
										Placeholder: "Выберите капитанов",
									},
								},
							},
						},
					},
				}
			}
			// _, err := s.InteractionResponseEdit(i.Interaction, response)
			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				fmt.Println("got err ", err)
			}
			// time.AfterFunc(time.Second*1, func() {
			// 	fmt.Println(i.ApplicationCommandData())
			// 	captain1ID := i.MessageComponentData().Values[0]
			// 	captain2ID := i.MessageComponentData().Values[1]
			// 	content := "Первый кэп @" + fmt.Sprintf("%v", i.MessageComponentData().Resolved.Users[captain1ID]) + "\n" + "Второй кэп @" + fmt.Sprintf("%v", i.MessageComponentData().Resolved.Users[captain2ID])
			// 	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			// 		Content: &content,
			// 		Components: &[]discordgo.MessageComponent{
			// 			discordgo.ActionsRow{
			// 				Components: []discordgo.MessageComponent{
			// 					discordgo.Button{
			// 						// Label is what the user will see on the button.
			// 						Label: "Подтвердить выбор",
			// 						// Style provides coloring of the button. There are not so many styles tho.
			// 						Style: discordgo.SuccessButton,
			// 						// Disabled allows bot to disable some buttons for users.
			// 						Disabled: false,
			// 						// CustomID is a thing telling Discord which data to send when this button will be pressed.
			// 						CustomID: "captains_select_ok",
			// 					},
			// 					discordgo.Button{
			// 						Label:    "Отмена",
			// 						Style:    discordgo.DangerButton,
			// 						Disabled: false,
			// 						CustomID: "fd_no",
			// 					},
			// 				},
			// 			},
			// 		},
			// 	})
			// 	if err != nil {
			// 		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			// 			Content: "Something went wrong" + err.Error(),
			// 		})
			// 		return
			// 	}
			// })
		},
	}
)
