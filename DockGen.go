package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorMagenta = "\033[35m"
	textBold     = "\033[1m"
)

func main() {
	fmt.Println(colorGreen + textBold + "Welcome! I will create your Dockerfile in seconds 🎉✨" + colorReset)
	fmt.Println()
	fmt.Println(colorRed + textBold + "Made with ❤️ by Shivam" + colorReset)
	languages := []string{"node", "python", "rust", "go", "java", "ubuntu"}
	languagesString := strings.Join(languages, ", ")
	if _, err := os.Stat("Dockerfile"); err == nil {
		fmt.Println()
		fmt.Print(colorMagenta + textBold + "A Dockerfile already present. Do you want to overwrite it? (y/n): " + colorReset)
		overwriteInput, err := getUserInput()
		if err != nil {
			fmt.Println()
			fmt.Println("Error reading input:", err)
			return
		}

		if strings.ToLower(overwriteInput) != "y" {
			fmt.Println(" ")
			fmt.Println(colorMagenta + textBold + "Cool choice! Keeping the existing Dockerfile safely with you👋." + colorReset)
			return
		}
	}

	fmt.Println()
	fmt.Println(colorMagenta + textBold + "Available application platforms are: " + languagesString + colorReset)
	for {
		fmt.Println()
		fmt.Print(colorMagenta + textBold + "Enter your application platform used by your project: " + colorReset)
		userInputLanguage, err := getUserInput()
		if err != nil {
			fmt.Println()
			fmt.Println("Error reading input:", err)
			return
		}

		if isValidLanguage(userInputLanguage, languages) {
			fmt.Println()
			fmt.Print(colorMagenta + textBold + "Enter the version of your application that you want (e.g., 1.0): " + colorReset)
			userInputVersion, err := getUserInput()
			if err != nil {
				fmt.Println()
				fmt.Println("Error reading input:", err)
				return
			}
			fmt.Println()

			fmt.Println(colorRed + "This step will be repeated, if you wish to escape this step type 'done' " + colorReset)
			userLabelCommands := getUserMultiLineCommands(colorMagenta + textBold + "Enter the Label (e.g. description=, maintainer=)" + colorReset)
			fmt.Println()

			fmt.Println(colorRed + "This step will be repeated, if you wish to escape this step type 'done' " + colorReset)
			userRunCommands := getUserMultiLineCommands(colorMagenta + textBold + "Enter the command that installs everything required for your application (e.g., 'apt-get' install -y vim or npm install)" + colorReset)
			fmt.Println()

			fmt.Println(colorRed + "This step will be repeated, if you wish to escape this step type 'done' " + colorReset)
			userExposeCommands := getUserMultiLineCommands(colorMagenta + textBold + "Enter the port on which your app will listen (type 'done' to escape this step)" + colorReset)
			userCmdCommands := getUserCommand(colorMagenta + textBold + "Enter the final Command that starts your application " + colorReset)

			err = generateDockerfile(userInputLanguage, userInputVersion, userRunCommands, userLabelCommands, userExposeCommands, userCmdCommands)
			if err != nil {
				fmt.Println()
				fmt.Println(colorRed + textBold + "Error generating Dockerfile:", err)
			} else {
				fmt.Println()
				fmt.Println(colorGreen + textBold + "Dockerfile generated successfully!" + colorReset)
			}
			break
		} else {
			fmt.Println()
			fmt.Println(colorMagenta + textBold + "Please try again. Valid application platforms are:", languagesString + colorReset)
		}
	}
}

func getUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	userInput = strings.TrimSpace(userInput)
	return userInput, nil
}

func isValidLanguage(userInput string, languages []string) bool {
	for _, lang := range languages {
		if userInput == strings.ToLower(lang) {
			return true
		}
	}
	return false
}

func getUserMultiLineCommands(prompt string) []string {
	var commands []string
	for {
		fmt.Println()
		fmt.Printf("%s: ", prompt)
		userInput, err := getUserInput()
		if err != nil {
			fmt.Println()
			fmt.Println("Error reading input:", err)
			return nil
		}

		if userInput == "done" {
			break
		}

		commands = append(commands, userInput)
	}

	return commands
}

func getUserCommand(commandType string) string {
	fmt.Println()
	fmt.Printf("%s: ", commandType)
	userInput, err := getUserInput()
	if err != nil {
		fmt.Println()
		fmt.Println("Error reading input:", err)
		return ""
	}
	return userInput
}

func generateDockerfile(selectedLanguage, selectedVersion string, userRunCommands, userLabelCommands, userExposeCommands []string, userCmdCommands string) error {
	dockerfileContent := fmt.Sprintf(`FROM %s:%s
	
WORKDIR /app

COPY . /app

`, selectedLanguage, selectedVersion)

	for i, labelCommand := range userLabelCommands {
		if i > 0 {
			dockerfileContent += "\n\n"
		}
		dockerfileContent += fmt.Sprintf("LABEL %s", labelCommand)
	}
	dockerfileContent += "\n\n" // Add a newline here

	for i, runCommand := range userRunCommands {
		if i > 0 {
			dockerfileContent += "\n"
		}
		dockerfileContent += fmt.Sprintf("RUN %s", runCommand)
	}

	if len(userExposeCommands) > 0 {
		dockerfileContent += "\n\n"
		for i, exposeCommand := range userExposeCommands {
			if i > 0 {
				dockerfileContent += "\n"
			}
			dockerfileContent += fmt.Sprintf("EXPOSE %s", exposeCommand)
		}
	}

	if userCmdCommands != "" {
		cmdArgs := strings.Fields(userCmdCommands)
		cmdJSONArray := fmt.Sprintf(`CMD ["%s"]`, strings.Join(cmdArgs, `", "`))
		dockerfileContent += "\n\n"
		dockerfileContent += fmt.Sprintf("%s", cmdJSONArray)
	}

	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(dockerfileContent)
	if err != nil {
		return err
	}

	return nil
}
