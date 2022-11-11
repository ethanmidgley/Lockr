package main

import (
	"Lockr/pkg/archive"
	"Lockr/pkg/encrypt"
	"Lockr/pkg/progress"
	"Lockr/pkg/validation"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/manifoldco/promptui"
)

func main() {
	fmt.Println("\033[32m", `
@@@        @@@@@@    @@@@@@@  @@@  @@@  @@@@@@@   
@@@       @@@@@@@@  @@@@@@@@  @@@  @@@  @@@@@@@@  
@@!       @@!  @@@  !@@       @@!  !@@  @@!  @@@  
!@!       !@!  @!@  !@!       !@!  @!!  !@!  @!@  
@!!       @!@  !@!  !@!       @!@@!@!   @!@!!@!   
!!!       !@!  !!!  !!!       !!@!!!    !!@!@!    
!!:       !!:  !!!  :!!       !!: :!!   !!: :!!   
 :!:      :!:  !:!  :!:       :!:  !:!  :!:  !:!  
 :: ::::  ::::: ::   ::: :::   ::  :::  ::   :::  
: :: : :   : :  :    :: :: :   :   :::   :   : :`)
	fmt.Println("")
	fmt.Println("")
	fmt.Printf(`
--------What would you like to do today?--------
	`)
	fmt.Printf("\033[0m")

	prompt := promptui.Select{
		Label: "Select Lockr action",
		Items: []string{"Encrypt", "Decrypt"},
	}

	_, action, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch action {
	case "Encrypt":

		encryptMenu := promptui.Select{
			Label: "What would you like to encrypt",
			Items: []string{"Single file", "Folder"},
		}

		_, encryptAction, err := encryptMenu.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		tempFile, err := os.CreateTemp("", "temp.zip")
		if err != nil {
			fmt.Println("Failed to create temp file. This is a serious problem ")
			return
		}
		switch encryptAction {
		case "Single file":
			fileLocationPrompt := promptui.Prompt{Label: "Enter file location and name", Validate: validation.NotEmpty, Default: ""}
			filename, err := fileLocationPrompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
			if !validation.FileExist(filename) {
				fmt.Println("File does not exist to decrypt")
				return
			}

			err = archive.SingleToZip(filename, tempFile.Name())
			if err != nil {
				fmt.Println("Failed to archive file inputed")
				return
			}

			break
		case "Folder":
			folderLocationPrompt := promptui.Prompt{Label: "Enter folder location and name", Validate: validation.NotEmpty, Default: ""}
			foldername, err := folderLocationPrompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
			if !validation.FileExist(foldername) {
				fmt.Println("File does not exist to decrypt")
				return
			}
			if foldername[len(foldername)-1] != '/' {
				foldername += "/"
			}

			err = archive.RecursiveZip(foldername, tempFile.Name())
			if err != nil {
				log.Panic(err)
			}
			break
		}

		passwordPrompt := promptui.Prompt{Label: "Enter password to encrypt data", Mask: '*'}
		password, err := passwordPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		newFilenamePrompt := promptui.Prompt{Label: "Enter filename for the encrypted data"}
		newFilename, err := newFilenamePrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		dst, err := os.Create(newFilename + ".Lockr")
		if err != nil {
			fmt.Printf("Failed to create new destination file %v\n", err)
			return
		}

		encStream, err := encrypt.NewEncryptReader(tempFile, []byte(password))

		pg := &progress.Monitor{Message: "Encrypting file", CompleteMessage: "Encryption complete", Iteration: 1}
		_, err = io.Copy(dst, io.TeeReader(encStream, pg))
		if err != nil {
			fmt.Printf("Failed to stream encrypted file to output %v\n", err)
			return
		}
		pg.Complete()

		tempFile.Close()
		os.Remove(tempFile.Name())

		fmt.Println("Successfully encrypted file")

		break

	case "Decrypt":

		filenamePrompt := promptui.Prompt{Label: "Enter filename of data to decrypt"}
		filename, err := filenamePrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		passwordPrompt := promptui.Prompt{Label: "Enter password to decrypt data", Mask: '*'}
		password, err := passwordPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if !validation.FileExist(filename) {
			fmt.Println("File does not exist to decrypt")
			return
		}

		encStream, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Failed to read file %v\n", err)
			return
		}

		dstStream, err := os.CreateTemp("", "temp.zip")
		if err != nil {
			fmt.Printf("Failed to create temp zip file %v\n", err)
			return
		}

		src, err := encrypt.NewDecryptReader(encStream, []byte(password))
		if err != nil {
			fmt.Println("Uhoh error probably wrong password")
			return
		}

		pg := &progress.Monitor{Message: "Decrypting file", CompleteMessage: "Decryption complete", Iteration: 1}
		_, err = io.Copy(dstStream, io.TeeReader(src, pg))
		if err != nil {
			fmt.Printf("Failed to decrypt stream %v\n", err)
			return
		}
		pg.Complete()

		err = archive.Unzip(dstStream.Name(), filename+" unlocked")
		if err != nil {
			fmt.Printf("Failed to unzip data %v\n", err)
			return
		}

		dstStream.Close()

		fmt.Println("Successfully decrypted")

		break
	}

}
