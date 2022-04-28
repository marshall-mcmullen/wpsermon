package pkg

import (
	// STDLIB
	"fmt"
	"sync"

	// External
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*
for video in "${videos[@]}"; do
    increment total_frames $(ffprobe -v error -select_streams v:0 -count_packets -show_entries stream=nb_read_packets -of csv=p=0 "${video}")
done

rm -f "${output}.mp4"
ffmpeg -v error                         \
    -vsync 0                            \
    -i "${TOPDIR}/assets/intro.mp4"     \
    -i "${tmpdir}/trimmed.mp4"          \
    -i "${TOPDIR}/assets/ending.mp4"    \
    -filter_complex "[0:v:0][0:a:0][1:v:0][1:a:0][2:v:0][2:a:0]concat=n=3:v=1:a=1[outv][outa]" \
    -map "[outv]" -map "[outa]"         \
    -progress "${tmpdir}/progress.txt"  \
    "${output}.mp4" &
pid=$!


    | ffmpeg -v error         \
      -ss "${start_time}"     \
      -i ${HOME}/audio.mp4"   \
      -ss "${start_time}"     \
      -i "${HOME}/video.mp4"  \
      -c copy                 \
      "${tmpdir}/trimmed.mp4" \
*/

func Trim(data *Data) {

	var waitGroup sync.WaitGroup

	progress := widget.NewProgressBar()

	window := fyne.CurrentApp().NewWindow("Trimming")
	window.SetContent(container.NewVBox(progress))
	window.Resize(fyne.NewSize(200, 300))
	window.CenterOnScreen()
	window.Show()
	window.SetCloseIntercept(func() { CheckError(fmt.Errorf("Trimming cancelled")) })

	waitGroup.Add(1)
	go func() {




		waitGroup.Done()
	}()

	waitGroup.Wait()
	window.Close()
}

