package main

import(
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type AudioMixer struct {
	mainTheme *audio.Player
	lineClear *audio.Player
}

func (mixer *AudioMixer) Play() {
	mixer.mainTheme.Play()
}

func (mixer *AudioMixer) Restart() {
	mixer.mainTheme.Rewind()
	mixer.mainTheme.Play()
}

func (mixer *AudioMixer) IsPlaying() bool {
	return mixer.mainTheme.IsPlaying()
}

func (mixer *AudioMixer) Stop() {
	mixer.mainTheme.Pause()
	mixer.mainTheme.Rewind()
}

func (mixer *AudioMixer) ClearLine() {
	mixer.lineClear.Rewind()
	mixer.lineClear.Play()
}

func CreateAudioPlayer() AudioMixer {
	audioContext, err := audio.NewContext(44100)
	if err != nil {
		log.Fatal(err)
	}

	file, err := ebitenutil.OpenFile("assets/audio/soundtrack.wav")
	if err != nil {
		log.Fatal(err)
	}

	decodedSound, err := wav.Decode(audioContext, file)
	if err != nil {
		log.Fatal(err)
	}

	player, err := audio.NewPlayer(audioContext, decodedSound)

	file2, err := ebitenutil.OpenFile("assets/audio/clear.wav")
	if err != nil {
		log.Fatal(err)
	}

	decodedSound2, err := wav.Decode(audioContext, file2)
	if err != nil {
		log.Fatal(err)
	}

	effectplayer, err := audio.NewPlayer(audioContext, decodedSound2)
	
	return AudioMixer{
		mainTheme: player,
		lineClear: effectplayer,
	}
}
