package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var (
	device string
)

// playbackCmd represents the playback command
var playbackCmd = &cobra.Command{
	Use:   "playback",
	Short: "Control spotify plaback",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	log.Println("Playback settings")
	// },
}

func init() {
	playbackCmd.PersistentFlags().StringVarP(&device, "device", "d", "", "Specify a specific device to control")

	playbackCmd.AddCommand(playCmd)
	playbackCmd.AddCommand(pauseCmd)
	playbackCmd.AddCommand(nextCmd)
	playbackCmd.AddCommand(prevCmd)
	rootCmd.AddCommand(playbackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playbackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playbackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Resume playback",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.PlayOpt(nil)
		if err != nil {
			log.Println(err)
		}
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause playback",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.PauseOpt(nil)
		if err != nil {
			log.Println(err)
		}
	},
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Play next song",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.NextOpt(&spotify.PlayOptions{})
		if err != nil {
			log.Println(err.Error())
		}
		getCurrentSong()
	},
}

var prevCmd = &cobra.Command{
	Use:   "previous",
	Short: "Play previous song",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.PreviousOpt(nil)
		if err != nil {
			log.Println(err.Error())
		}
		getCurrentSong()
	},
}

func getCurrentSong() {
	currPlaying, _ := client.PlayerCurrentlyPlayingOpt(nil)
	fmt.Printf("Now playing %s by %s", currPlaying.Item.Name, currPlaying.Item.Artists[0].Name)
}

func getCurrentDevices() {
	devices, err := client.PlayerDevices()
	for _, device := range devices {
		log.Println(device.Name)
		log.Println(device.ID)
	}
	log.Printf("%v", err)
}
