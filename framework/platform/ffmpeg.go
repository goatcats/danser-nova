package platform

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/wieku/danser-go/framework/files"
)

const errMsg = "ffmpeg not found! Please make sure it's installed in danser directory or in PATH. Follow download instructions at https://github.com/Wieku/danser-go/wiki/FFmpeg"

var ffmpegInit bool
var ffPath string

func PrepareFFMpeg(cmdName string, args ...string) (*exec.Cmd, error) {
	if !ffmpegInit {
		ffmpegInit = true

		ffmpegExec, err := files.GetCommandExec("ffmpeg", "ffmpeg")
		if err != nil {
			return nil, fmt.Errorf(errMsg)
		}

		ffPath = filepath.Dir(ffmpegExec)
		log.Println("FFmpeg exec location:", ffmpegExec)
	} else if ffPath == "" {
		return nil, fmt.Errorf(errMsg)
	}

	cmd := exec.Command(filepath.Join(ffPath, cmdName), args...)
	cmd.Dir = ffPath

	return cmd, nil
}
