package main

import (
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
)

func main() {

	_ = viper.BindEnv("SLACK_VERIFICATION_TOKEN")
	_ = viper.BindEnv("SLACK_TOKEN")
	_ = viper.BindEnv("PORT")

	http.HandleFunc("/stamp", func(w http.ResponseWriter, r *http.Request) {
		vToken := viper.Get("SLACK_VERIFICATION_TOKEN")
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !s.ValidateToken(vToken.(string)) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch s.Command {
		case "/stamp":
			token := viper.Get("SLACK_TOKEN")
			eName := s.Text
			eName = strings.Replace(eName, ":", "", -1)
			api := slack.New(token.(string))
			emojis, err := api.GetEmoji()
			if err != nil {
				log.Printf("[ERROR] GetEmoji: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			eUrl := emojis[eName]

			u, err := api.GetUserInfo(s.UserID)
			if err != nil {
				log.Printf("[ERROR] GetUserInfo: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			// memo: textに何か入れないと怒られるからとりあえず半角スペースを入れといた
			a := slack.Attachment{
				Color:    "FFF",
				ImageURL: eUrl,
				Text:     " ",
			}

			_, _, err = api.PostMessage(s.ChannelID, slack.MsgOptionUsername(s.UserName), slack.MsgOptionIconURL(u.Profile.Image192), slack.MsgOptionAttachments(a))
			if err != nil {
				log.Printf("[ERROR] Postmessage: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	p := viper.Get("PORT")
	log.Printf("[INFO] Server listening %s\n", p.(string))
	if err := http.ListenAndServe(":"+p.(string), nil); err != nil {
		log.Printf("[ERROR] %v", err)
	}
}
