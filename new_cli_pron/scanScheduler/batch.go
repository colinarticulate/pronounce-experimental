package scanScheduler

import (
  "bufio"
  // "fmt"
  // "log"
  "os"
  // "os/exec"
  // "path"
  "strings"
)

/*
func (s Scheduler) newBatchScan(frate string) batchScan {
  if b, ok := s.batchResults[[]string{frate}]; ok {
    // We've already got a batchScan for this frate so just return it
    return b
  }
  
  return batchScan{
    running,
    path.Join(s.outfolder, "batch" + frate + ".log"),
    "",
    []chan bool{},
  }
}
*/
func (b batchScan) getCmnVec(logfile string) []string {
  vec := []string{}
  f, err := os.Open(logfile)
  if err != nil {
    return []string{}
  }
  defer f.Close()
  
  s := bufio.NewScanner(f)
  for s.Scan() {
    l := s.Text()
    tokens := strings.Fields(l)
    if len(tokens) < 4 {
      // We're looking for something like
      // INFO: cmn.c(133): CMN: 49.60 ...
      // so any line with less than 4 tokens is NOT what we're after
      continue
    }
    // This is horribly specific but not sure there's a better way 
    if tokens[0] == "INFO:" && tokens[2] == "CMN:" {
      vec = tokens[3:]
    }
  }
  return vec
}

/*
func (b *batchScan) runBatchScan(frate string) {
  args := []string{
    "-adcin", "yes", "-cepdir", cepdir, "-cepext", ".wav", "-ctl", ctlfile, "-logfn", batchfile,
  }
  _, err := exec.Command("pocketsphinx_batch", args...).Output()
  if err != nil {
    fmt.Println("Oops, check pocketsphinx settings? args are...", args)
  }  
  b.cmnVec = strings.Join(b.getCmnVec(), ",")
}
*/