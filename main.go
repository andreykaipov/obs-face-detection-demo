package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	sceneitems "github.com/andreykaipov/goobs/api/requests/scene_items"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/typedefs"
)

func main() {
	obs, err := NewOBS(os.Getenv("OBS_HOST"), os.Getenv("OBS_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}

	scene, err := obs.Client.Scenes.GetCurrentScene()
	if err != nil {
		log.Fatal(err)
	}

	obsFaceSource := getenv("OBS_FACE_SOURCE", "Webcam")
	obsCheckInterval := getenv("OBS_CHECK_INTERVAL", "15")
	brbScene := getenv("OBS_BRB_SCENE", "BRB")
	interval, _ := strconv.Atoi(obsCheckInterval)

	activeScene := scene.Name
	swapScene := ""

	if activeScene == swapScene {
		log.Fatalf("Active scene %q is set to brb scene %q already.", activeScene, brbScene)
	}

	log.Printf("Current scene is %s. Using that as our original and active scene.", scene.Name)

	for range time.Tick(time.Duration(interval) * time.Second) {
		face, err := obs.DetectFace(obsFaceSource)
		if err != nil {
			log.Fatal(err)
		}

		if face {
			swapScene = scene.Name
			activeScene = handleFace(obs, activeScene, swapScene)
		} else {
			swapScene = brbScene
			activeScene = handleNoFace(obs, activeScene, swapScene)
		}
	}

}

// Takes the currently active scene, the scene to swap to, and returns the new
// active scene, i.e. the swap scene.
func handleFace(obs *OBS, activeScene, swapScene string) string {
	msg := "Detected a face"

	if activeScene == swapScene {
		msg += fmt.Sprintf(", but we're already in the %q scene so no point in setting it again", swapScene)
		log.Println(msg)
		return activeScene
	}

	log.Println(msg)

	if _, err := obs.Client.Scenes.SetCurrentScene(&scenes.SetCurrentSceneParams{SceneName: swapScene}); err != nil {
		log.Fatal(err)
	}

	if _, err := obs.Client.SceneItems.DeleteSceneItem(&sceneitems.DeleteSceneItemParams{
		Scene: activeScene,
		Item:  &typedefs.Item{Name: swapScene},
	}); err != nil {
		log.Fatal(err)
	}

	return swapScene
}

// Takes the currently active scene, the scene to swap to, and returns the new
// active scene, i.e. the swap scene.
func handleNoFace(obs *OBS, activeScene, swapScene string) string {
	msg := "No faces"

	if activeScene == swapScene {
		msg += fmt.Sprintf(", but we're already in the %q scene so no point in setting it again", swapScene)
		log.Println(msg)
		return activeScene
	}

	log.Println(msg)
	log.Printf("Setting the %q scene", swapScene)

	if _, err := obs.Client.Scenes.SetCurrentScene(&scenes.SetCurrentSceneParams{SceneName: swapScene}); err != nil {
		log.Fatal(err)
	}

	if _, err := obs.Client.SceneItems.AddSceneItem(&sceneitems.AddSceneItemParams{
		SceneName:  swapScene,
		SourceName: activeScene,
		SetVisible: true,
	}); err != nil {
		log.Fatal(err)
	}

	return swapScene
}

func getenv(envvar string, d string) string {
	val := os.Getenv(envvar)
	if val == "" {
		return d
	}

	return val
}
