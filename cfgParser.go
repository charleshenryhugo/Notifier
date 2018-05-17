package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

/*------structs corresponding to the config file for different notifiers------*/

//NotifConfig is the most initial struct(class) corresponding to config-file for notifiers
type NotifConfig struct {
	Notifiers Notifiers `yaml:"notifiers"`
}

//Notifiers is the struct that contains all the possible notifiers parsed from config file
//After modifying the notifier-config-file, also please add the new notifier here
//e.g. AWSNotifier AWSNotifier `yaml:"awsnotifier"`
type Notifiers struct {
	SMTPEmailNotifier SmtpEmailNotifier `yaml:"smtpemailnotifier"`
	SlackNotifier     SlackNotifier     `yaml:"slacknotifier"`
}

//SmtpEmailNotifier is the struct corresponding to the yaml:smtpemailnotifier in the config file
type SmtpEmailNotifier struct {
	Type     string `yaml:"type"`
	State    bool   `yaml:"state"`
	Account  string `yaml:"account"`
	Pwd      string `yaml:"pwd"`
	SMTPHost string `yaml:"SMTPHost"`
	SMTPPort string `yaml:"SMTPPort"`
}

//SlackNotifier is the struct corresponding to the yaml:slacknotifier in the config file
type SlackNotifier struct {
	Type        string   `yaml:"type"`
	State       bool     `yaml:"state"`
	Token       string   `yaml:"token"`
	AsUser      bool     `yaml:"asUser"`
	UserName    string   `yaml:"userName"`
	IconEmoji   string   `yaml:"iconEmoji"`
	WebhookURLs []string `yaml:"WebhookURLs"`
}

//Add new Notifier struct here:
//e.g. type AWSNotifier struct {}

/*------please add new Notifiers above this line------*/

//initViper initializes a viper for yaml parsing
//specify the target *.yaml file to "file" parameter
//return a viper instance(reference type)
func initViper(file string) *viper.Viper {
	//initialize an viper for notifiers
	nviper := viper.New()
	//name of config file(without extension)
	nviper.SetConfigName(file)
	//paths to look for notifyrc file in
	nviper.AddConfigPath(".")
	nviper.AddConfigPath("$HOME")
	//tell the viper instance to watchConfig
	nviper.WatchConfig()
	//provide a function for viper to run each time a config change occurs
	nviper.OnConfigChange(
		func(e fsnotify.Event) {
			log.Println("config file changed:", e.Name)
		})
	return nviper
}

//parse notifiers objects from the *.yaml file specified by "file"
//using viper
//return a Notifiers struct
func parseNotifiers(file string) (Notifiers, ERR) {
	//initialize viper to parse notifyrcFile
	nviper := initViper(notifyrcFile)
	//find and read notifyrc file
	err := nviper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return Notifiers{}, NOTIFRC_PARSE_ERR
	}

	//initialize Notifiers
	var ntfCfg NotifConfig
	//unmarshall into Notifiers
	if err := nviper.Unmarshal(&ntfCfg); err != nil {
		log.Println(err)
		return Notifiers{}, NOTIFRC_PARSE_ERR
	}
	return ntfCfg.Notifiers, SUCCESS
}

//DfltConfig is the most initial struct(class) corresponding to config-file for default settings
type DfltConfig struct {
	Defaults Defaults `yaml:"defaults"`
}

//Defaults contains all the default settings stored in the defaultsFile
//If you modify the defaultsFile, please also modify this struct correspondingly
type Defaults struct {
	EmailListFile string `yaml:"emailListFile"`
	SlackListFile string `yaml:"slackListFile"`
	Subject       string `yaml:"subject"`
	Message       string `yaml:"message"`
	MessageFile   string `yaml:"messageFile"`
}

//parse the Defaults object from *.yaml file
//return a Defaults struct
func parseDefaults(file string) (Defaults, ERR) {
	dviper := initViper(defaultsFile)

	//find and read defaults file
	err := dviper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return Defaults{}, DFLTS_PARSE_ERR
	}

	//initialize Defaults object
	var dfltCfg DfltConfig
	//unmarshall into Defaults
	if err := dviper.Unmarshal(&dfltCfg); err != nil {
		log.Println(err)
		return Defaults{}, DFLTS_PARSE_ERR
	}

	return dfltCfg.Defaults, SUCCESS
}

/*------the methods of the Defaults struct will only be called when the corresponding argument is blank------*/
/*------e.g. GetDfltmsg will be called only if the input message is "" ------*/

//GetDfltSbjt returns the default subject set by the defaultsFile
func (dflt *Defaults) GetDfltSbjt() string {
	return dflt.Subject
}

//GetDfltmsg returns the default message set by the defaultsFile
func (dflt *Defaults) GetDfltmsg() string {
	//get message from the default message file, only if the file is available
	if fileBytes, err := ioutil.ReadFile(dflt.MessageFile); err == nil {
		//only when the msg file is read successfully, we rewrite the msg
		return string(fileBytes)
	}
	//if file reading is failed, then return default message directly
	return dflt.Message
}

//GetDfltSlackList returns default slack IDs stored in the "slackListFile" which is set by defaultsFile
func (dflt *Defaults) GetDfltSlackList() []string {
	//return those slack user IDs stored in the default file, only if the file is available
	if fileBytes, err := ioutil.ReadFile(dflt.SlackListFile); err == nil {
		return strings.Fields(string(fileBytes))
	}
	return []string{}
}

//GetDfltEmailList returns default email addrs stored in the "EmailListFile" which is set by defaultsFile
func (dflt *Defaults) GetDfltEmailList() []string {
	//return those email addrs stored in the default file, only if the file is available
	if fileBytes, err := ioutil.ReadFile(dflt.EmailListFile); err == nil {
		return strings.Fields(string(fileBytes))
	}
	return []string{}
}

/*------ these methods of Defaults struct above will be called only if the input message is "" ------*/

//cfgRead gets value of a specific item in cfgFile
//and return it as an interface{}
func cfgRead(item, cfgFile string) (interface{}, error) {
	viper.SetConfigName(cfgFile)
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		return nil, err
	}
	return viper.Get(item), nil
}

//cfgWrite overwrites a specified item to newVal in cfgFile
func cfgWrite(item string, newVal interface{}, cfgFile string) error {
	viper.SetConfigName(cfgFile)
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		return err
	}
	viper.Set(item, newVal)
	return viper.WriteConfig()
}

//CfgDflt overwrites a specified item to newVal in defaultsFile
func CfgDflt(item string, newVal interface{}) error {
	return cfgWrite(item, newVal, defaultsFile)
}

//CfgNtfyrc overwrites a specified item to newVal in notifyrc.yml
func CfgNtfyrc(item string, newVal interface{}) error {
	return cfgWrite(item, newVal, notifyrcFile)
}

//CfgDfltSbjt overwrites default subject in defaultsFile
func CfgDfltSbjt(newSbjt string) error {
	err := CfgDflt("defaults.subject", newSbjt)
	if err == nil {
		log.Println("default subject/title reset as:", newSbjt)
	}
	return err
}

//CfgDfltMsg overwrites default message in defaultsFile
func CfgDfltMsg(newMsg string) error {
	err := CfgDflt("defaults.message", newMsg)
	if err == nil {
		log.Println("default message reset as:", newMsg)
	}
	return err
}

//CfgDfltMsgFile overwrites default message-file in defaultsFile
func CfgDfltMsgFile(newMsgFile string) error {
	err := CfgDflt("defaults.messageFile", newMsgFile)
	if err == nil {
		log.Println("default message file reset as:", newMsgFile)
	}
	return err
}

//CfgDfltSlackListFile overwrites default SlackListFile in defaultsFile
func CfgDfltSlackListFile(newFile string) error {
	err := CfgDflt("defaults.slackListFile", newFile)
	if err == nil {
		log.Println("default slack list file reset as:", newFile)
	}
	return err
}

//CfgDfltEmailListFile overwrites default EmailListFile in defaultsFile
func CfgDfltEmailListFile(newFile string) error {
	err := CfgDflt("defaults.emailListFile", newFile)
	if err == nil {
		log.Println("default email list file reset as:", newFile)
	}
	return err
}

//CfgToggStat toggles state between on and off for a specific notifier type
//it will modify notifyrc.yml
func CfgToggStat(ntfName string) error {
	item := "notifiers." + ntfName + ".state"
	state, err := cfgRead(item, notifyrcFile)
	if err != nil {
		log.Println(err)
		return err
	}
	err = CfgNtfyrc(item, !state.(bool))
	if err == nil {
		log.Println("toggle state of ", ntfName, " to ", !state.(bool))
	}
	return err
}
