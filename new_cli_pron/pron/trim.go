package pron

import (
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cryptix/wav"
	"github.com/maxhawkins/go-webrtcvad"
)

type voicedInterval struct {
	start, duration int
}

//===================================================================================================
//  __      __   _    ___ _____ ___    __   ___   ___              _   _   _
//  \ \    / /__| |__| _ \_   _/ __|   \ \ / /_\ |   \     ___ ___| |_| |_(_)_ _  __ _ ___
//   \ \/\/ / -_) '_ \   / | || (__     \ V / _ \| |) |   (_-</ -_)  _|  _| | ' \/ _` (_-<
//    \_/\_/\___|_.__/_|_\ |_| \___|     \_/_/ \_\___/    /__/\___|\__|\__|_|_||_\__, /__/
//                                                                               |___/
//===================================================================================================
// The number used with WebRtc sets the VAD operating mode. A more aggressive (higher mode) VAD is more
// restrictive in reporting speech. Put in other words the probability of being speech when the VAD
// returns 1 is increased with increasing mode. As a consequence also the missed detection rate goes up.
//
// Aggressiveness mode (0, 1, 2, or 3).
/*
VAD.Mode.NORMAL   ... 0?
Constant for normal voice detection mode. Suitable for high bitrate, low-noise data. May classify noise as voice, too. The default value if mode is omitted in the constructor.

VAD.Mode.LOW_BITRATE   ... 1?
Detection mode optimised for low-bitrate audio.

VAD.Mode.AGGRESSIVE   ..... 2?
Detection mode best suited for somewhat noisy, lower quality audio.

VAD.Mode.VERY_AGGRESSIVE  ..... 3?
Detection mode with lowest miss-rate. Works well for most inputs.


The WebRTC VAD only accepts 16-bit mono PCM audio, sampled at 8000, 16000, 32000 or 48000 Hz. A frame must be either 10, 20, or 30 ms in duration:
Optionally, set its aggressiveness mode, which is an integer between 0 and 3. 0 is the least aggressive about filtering out non-speech, 3 is the most aggressive.
*/

/*
Addtionaly
==========
The VAD engine requires mono, 16-bit PCM audio with a sample rate of 8, 16, 32 or 48 KHz as input. The input should be an audio segment of 10, 20 or 30 milliseconds.
When the audio input is 16 Khz, the input array should thus be either of length 160, 320 or 480.
https://github.com/jitsi/jitsi-webrtc-vad-wrapper

For example, if your sample rate is 16000 Hz, then the only allowed frame/chunk sizes are 16000 * ({10,20,30} / 1000) = 160, 320 or 480 samples.
Since each sample is 2 bytes (16 bits), the only allowed frame/chunk sizes are 320, 640, or 960 bytes.
https://github.com/wiseman/py-webrtcvad/issues/30
*/

func fetchTrimBounds(audiofile string, phons []phoneme) (float64, float64) {
	// Original we set this quite aggressively, but may?? be causing issues when background is noisy, needs more investigation - PE August 2019
	//start2, duration2 := webRtcBounds(audiofile, 2)
	//start3, _ := webRtcBounds(audiofile, 3)
	dir, file := filepath.Split(audiofile)
	ext := filepath.Ext(file)
	sincAudiofile := filepath.Join(dir, file[:len(file)-len(ext)]+"_lowpass"+ext)

	out, err := exec.Command("soxi", "-D", audiofile).Output()
	if err != nil {
		log.Panic(err)
	}
	// Yuk! out contains a float terminated by '\n' so strip the '\n'. There must
	// be a better way...
	//
	length_audio, err := strconv.ParseFloat(string(out[:len(out)-1]), 64)
	if err != nil {
		log.Panic(err)
	}

	// //sox
	// _, err = exec.Command("sox", audiofile, sincAudiofile, "sinc", "5000-500").Output()
	// debug("\n length of audio = ", length_audio)

	// //_, err := exec.Command("sox", audiofile, sincAudiofile).Output()
	// if err != nil {
	// 	debug("Call to sox failed with err, ", err)
	// }
	//Meena
	_, err = exec.Command("mt_cv_mvns", audiofile, sincAudiofile, "1", "16000", "25").Output()
	debug("\n length of audio = ", length_audio)
	if err != nil {
		debug("Call to mt_cv_mvns failed with err, ", err)
	}

	start2 := 0.0
	duration2 := 0.0
	start3 := 0.0
	duration3 := 0.0

	start2, duration2 = webRtcBounds(sincAudiofile, 0)
	start3, duration3 = webRtcBounds(sincAudiofile, 3)

	//start0, duration0 := webRtcBounds(sincAudiofile, 1)
	//start1, duration1 := webRtcBounds(sincAudiofile, 3)

	start := 0.0
	end := 0.0
	duration := 0.0

	end2 := start2 + duration2
	end3 := start3 + duration3
	end = math.Max(end2, end3)

	/*
	  if math.Min(start2, start3) == 0 {
	     if start2 == start3 {  // they are both 0
	      start = 0
	     } else if start2 == 0 {
	       start = start3
	       } else if start3 == 0 {
	       start = start2}
	  }else{
	  start = math.Min(start2, start3)
	  }
	*/

	start = math.Min(start2, start3)

	// if start-0.3 <= 0 {
	// 	start = 0
	// } else {
	// 	start = start - 0.25
	// }

	// Use the settings below for production - PE 2 June 2022 -- the ones above are for trimming scraped audio for training
	if start-0.2 <= 0 {
		start = 0
	} else {
		start = start - 0.10
	}

	/*
	   if start2 <= 0.2 {
	     if start3 > 0.5 {
	        start = start3 -0.3   /// this is because start3 is generally the more aggressive
	        } else {
	        start = start3
	        }

	   }

	   if start3 <= 0.2 {
	   start = start2
	   }

	   if start <= 0.2 {
	   start = 0.2
	   }
	*/

	debug("\n start2, duration2, end2 = ", start2, duration2, end2)
	debug("\n start3, duration3, end3 = ", start3, duration3, end3)
	//debug("\n start0, duration0 = ", start0, duration0)
	//debug("\n start1, duration1 = ", start1, duration1)

	duration = end - start

	debug("\n Initial start, duration, end = ", start, duration, end)

	//Adjustments
	//start = start - 0.2333  //*********************** adjusted to try and fix tomo replied where the voiceless plosive "p" goes to zero and VAD thinks it's a silence gap
	// 0.4 is added to cover contraceptives1_tomo  .... and potentially others than have voiceless plosives towards the end ... though big silent gaps could be an issue

	//
	//if start > 0.5{
	//start = start - 0.2       // This is required for the S words where VAD starts late, example seat1_tressa, the s goes missing without an earlier 0 period  // 0.2
	//duration = duration + 0.2
	// //start = start
	// //duration = duration
	//}

	duration = duration + 0.45 // 0.4 is required to get contraceptives1_tomo to work
	//duration = duration

	//check if start + duration is greater than the end of the audio file
	if (start + duration) >= length_audio {
		duration = length_audio - start - 0.02
	}

	guard_end := 0.0
	guard_end = start + duration

	debug("\n Cut at: start, duration = ", start, duration, "    guard_end = ", guard_end, "\n")
	return start, duration

}

func webRtcBounds(audiofile string, mode int) (float64, float64) {
	info, err := os.Stat(audiofile)
	if err != nil {
		debug("webRtcBounds: call to os.Stat failed. err =", err)
		log.Panic(err)
	}

	file, err := os.Open(audiofile)
	if err != nil {
		debug("webRtcBounds: failed to open file. err =", err)
		log.Panic(err)
	}
	defer file.Close()

	wavReader, err := wav.NewReader(file, info.Size())
	if err != nil {
		debug("webRtcBounds: call to wav.NewReader failed. err =", err)
		log.Panic(err)
	}

	reader, err := wavReader.GetDumbReader()
	if err != nil {
		debug("webRtcBounds: call to wav.GetDumbReader failed. err =", err)
		log.Panic(err)
	}

	wavInfo := wavReader.GetFile()
	rate := int(wavInfo.SampleRate)
	if wavInfo.Channels != 1 {
		debug("webRtcBounds: expected mono file")
		log.Panic("expected mono file")
	}
	if rate != 16000 {
		debug("webRtcBounds: expected 16kHz file")
		log.Panic("expected 16kHz file")
	}

	vad, err := webrtcvad.New()
	// vad, err := artVad.New()
	if err != nil {
		debug("webRtcBounds: call to webrtcvad.New failed. err =", err)
		log.Panic(err)
	}

	if err := vad.SetMode(mode); err != nil {
		debug("webRtcBounds: call to vad.SetMode failed. err =", err)
		log.Panic(err)
	}

	//frame := make([]byte, 160*2)               //  160/16kHz=10ms but each sample is 2 bytes.  Therefore, 320=10ms, 640=20ms, 960=30ms.
	frame := make([]byte, 640) //640

	debug("\n\n\nlen(frame)=", len(frame))
	debug("rate =", rate)
	//debug("\nThe frame is = ", frame)
	//if ok := vad.ValidRateAndFrameLength(rate, len(frame)); !ok {
	//  debug("\nwebRtcBounds: Valid frame rate & length checker failed")
	//  log.Panic("webRTC: invalid rate or frame length")
	//}

	var isActive bool
	var offset int

	var start, duration int
	var thisStart int

	report := func() {
		_ = time.Duration(offset) * time.Second / time.Duration(rate) / 2
	}

	for {
		_, err := io.ReadFull(reader, frame)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			debug("webRtcBounds: call to io.ReadFull failed. err =", err)
			log.Panic(err)
		}

		frameActive, err := vad.Process(rate, frame)
		if err != nil {
			debug("webRtcBounds: call to vad.Process failed. err =", err)
			log.Panic(err)
		}

		if isActive != frameActive || offset == 0 {
			// isActive = frameActive
			report()
		}
		// Find word on the fly
		// Declare start, duration, both initialised to 0
		// On isActive set potStart, on !isactive check to see if duration > last
		// saved duration. If it is start = potStart, duration == 'new duration'
		if isActive != frameActive || offset == 0 {
			if isActive {
				if offset-thisStart > duration {
					start = thisStart
					duration = offset - start
				}
			} else {
				thisStart = offset
			}
			isActive = frameActive
		}
		offset += len(frame)
	}
	report()
	if isActive {
		if offset-thisStart > duration {
			start = thisStart
			duration = offset - start
		}
	}
	startSecs := float64(start) / float64(rate*2)
	if startSecs > 0.1 {
		startSecs -= 0.1
	}
	durationSecs := (float64(duration) / float64(rate*2)) + 0.1
	return startSecs, durationSecs
}

func trimAudio(audiofile string, phons []phoneme) string {
	start, duration := fetchTrimBounds(audiofile, phons)

	dir, file := filepath.Split(audiofile)
	ext := filepath.Ext(file)
	trimmedfile := filepath.Join(dir, file[:len(file)-len(ext)]+"_trimmed"+ext)

	cmd := exec.Command("sox", audiofile, trimmedfile, "trim", strconv.FormatFloat(start, 'f', -1, 64), strconv.FormatFloat(duration, 'f', -1, 64))

	err := cmd.Run()
	if err != nil {
		debug("err =", err)
	}
	return trimmedfile
}
