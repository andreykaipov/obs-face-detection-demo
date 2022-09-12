package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	sceneitems "github.com/andreykaipov/goobs/api/requests/scene_items"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/typedefs"
	"github.com/armon/circbuf"
)

var (
	obsFaceSource         string
	obsFaceCheckInterval  int
	obsFaceCheckThreshold int

	sceneOG     string
	sceneActive string
	sceneBRB    string
)

func init() {
	obsFaceSource = getenv("OBS_FACE_SOURCE", "Webcam")
	sceneBRB = getenv("OBS_SCENE_BRB", "")

	obsFaceCheckInterval, _ = strconv.Atoi(getenv("OBS_FACE_CHECK_INTERVAL", "15"))
	obsFaceCheckThreshold, _ = strconv.Atoi(getenv("OBS_FACE_CHECK_THRESHOLD", "3"))
}

func main() {
	obs, err := NewOBS(os.Getenv("OBS_HOST"), os.Getenv("OBS_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}

	scene, err := obs.Client.Scenes.GetCurrentScene()
	if err != nil {
		log.Fatal(err)
	}

	sceneOG = scene.Name
	sceneActive = sceneOG

	log.Printf("Current scene is %q. Using that as our original and active scene.", scene.Name)

	if sceneActive == sceneBRB {
		log.Fatalf("Active scene %q is set to brb scene %q already.", sceneActive, sceneBRB)
	}

	buf, _ := circbuf.NewBuffer(int64(obsFaceCheckThreshold))

	for range time.Tick(time.Duration(obsFaceCheckInterval) * time.Second) {
		face, err := obs.DetectFace(obsFaceSource)
		if err != nil {
			log.Fatal(err)
		}

		if face {
			buf.WriteByte('1')
		} else {
			buf.WriteByte('0')
		}

		fmt.Println(buf)

		if strings.Count(buf.String(), "0") == int(buf.Size()) {
			handleNoFace(obs)
		} else {
			handleFace(obs)
		}
	}

}

func handleFace(obs *OBS) {
	msg := "Detected a face"

	if sceneActive == sceneOG {
		msg += fmt.Sprintf(", but we're already in the %q scene so no point in setting it again", sceneOG)
		log.Println(msg)
		return
	}

	log.Println(msg)
	log.Printf("Setting the %q scene", sceneOG)

	if _, err := obs.Client.Scenes.SetCurrentScene(&scenes.SetCurrentSceneParams{SceneName: sceneOG}); err != nil {
		log.Fatal(err)
	}

	if _, err := obs.Client.SceneItems.DeleteSceneItem(&sceneitems.DeleteSceneItemParams{
		Scene: sceneBRB + " active blurred",
		Item:  &typedefs.Item{Name: sceneOG},
	}); err != nil {
		log.Fatal(err)
	}

	sceneActive = sceneOG
}

func handleNoFace(obs *OBS) {
	msg := "No faces"

	if sceneActive == sceneBRB {
		msg += fmt.Sprintf(", but we're already in the %q scene so no point in setting it again", sceneBRB)
		log.Println(msg)
		return
	}

	log.Println(msg)
	log.Printf("Setting the %q scene", sceneBRB)

	if _, err := obs.Client.Scenes.SetCurrentScene(&scenes.SetCurrentSceneParams{SceneName: sceneBRB}); err != nil {
		log.Fatal(err)
	}

	if _, err := obs.Client.SceneItems.AddSceneItem(&sceneitems.AddSceneItemParams{
		SceneName:  sceneBRB + " active blurred",
		SourceName: sceneOG,
		SetVisible: true,
	}); err != nil {
		log.Fatal(err)
	}

	sceneActive = sceneBRB
}

func getenv(envvar string, d string) string {
	val := os.Getenv(envvar)
	if val == "" {
		return d
	}

	return val
}
