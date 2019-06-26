package wizard

//	"github.com/AlecAivazis/survey"

// TODO - DO NOT CHECK IN
// This is parked as the survey package is not playing nice with glide.

func requestFile(prompt string, defaultValue string) string {
	file := defaultValue
	//	prompt := &survey.Input{Message: prompt}
	//	survey.AskOne(prompt, &file)

	return file
}
