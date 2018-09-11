package main

import(
   "github.com/sirupsen/logrus"
   prefixed "github.com/x-cray/logrus-prefixed-formatter"
   "fmt"
   "time"
   "log"
   "strconv"
)

var(
   Cfr = new(prefixed.TextFormatter)       // Formated
)

// ----------------------------------
// Informer to command line
// ----------------------------------
func Inf(Proc, Info, Tp string){

	Cfr.TimestampFormat = "2006-01-02 15:04:05"
	Cfr.FullTimestamp   = true
	logrus.SetFormatter(Cfr)
	log := logrus.WithField("prefix", Proc)

     switch Tp {
     case "i": log.Info(Info)
     case "w": log.Warn(Info)
     case "d": log.Debug(Info)
     case "e": log.Error(Info)
     case "p": log.Panic(Info)
     default:  log.Fatal(Info)
     }
}

// -------------------------------------
// Informer to command line
// -------------------------------------
func Prn(Info string){
	 fmt.Println(Info)
}


/***************************************************************
  Company     Essentia
  Description Error
  Time        11-12-2017
  Title       Check errors
 ****************************************************************/
func Err(Er error, Txt string) {
	if Er != nil {
	 	 // log.Println("ERROR : " + Txt, "Description : " + Er.Error())
     log.Println("ERROR : " + Txt)
		 return
	}
}

/************************************************************************
 * String -> int64
 ************************************************************************/
func STI(StrParameter string) int64 {
	k, err := strconv.ParseInt(StrParameter, 10, 64)
	if err != nil {
		k = 0
	}
	return k
}

/*******************************************************************
 * String -> int
 ******************************************************************/
func Sti(StrParameter string) int {
	k, err := strconv.Atoi(StrParameter)
	if err != nil {
		k = 0
	}
	return k
}

/*****************************************************************
 * String -> float64
 ****************************************************************/
func Stf(StrParameter string) float64 {
	k, err := strconv.ParseFloat(StrParameter, 2)
	if err != nil {
		k = 0
	}
	return k
}

/****************************************************
 * Конвертация Int to Str
 ****************************************************/
func InttoStr(Ints int) string {
	//str := strconv.FormatInt(Intt64, 10)      // Выдает конвертацию 2000-wqut
	//str := strconv.Itoa64(Int64)              // use base 10 for sanity purpose
	str := fmt.Sprintf("%d", Ints)
	return str
}

/*****************************************************
 * Конвертация Int64 to Str
 *****************************************************/
func Int64toStr(Int64 int64) string {
	//str := strconv.FormatInt(Intt64, 10)      // Выдает конвертацию 2000-wqut
	//str := strconv.Itoa64(Int64)              // use base 10 for sanity purpose
	str := fmt.Sprintf("%d", Int64)
	return str
}

/*****************************************************
 *   Title       : Current time in format (YYYY-MM-DD HH:MM:SS)
 * 	 Date        : 2015-12-14
 *****************************************************/
func CTM() string {
	 return time.Now().Format("2006-01-02 15:04:05")
}

/*****************************************************
// Logurs
 *****************************************************/
func Logus() {
	 l:=logrus.Fields{"animal": "walrus", "Id":"Id001"}
     logrus.WithFields(l).Info("A walrus appears")	

	 l=logrus.Fields{"animal": "walrus", "Id":"Id001", "Idw":"Id001"}
     logrus.WithFields(l).Info("A walrus appears")	

     requestLogger := logrus.WithFields(logrus.Fields{"request_id": "Req-001", "user_ip": "Userip"})
     requestLogger.Info("something happened on that request") 
     requestLogger.Warn("something not great happened")
}
