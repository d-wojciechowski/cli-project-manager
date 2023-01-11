/*
Copyright © 2023 Dominik Wojciechowski

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the project/projects",
	Run: func(cmd *cobra.Command, args []string) {
		flag := cmd.Flag("projectName")
		if !flag.Changed {
			traverseAll(cmd.Flag("workDir").Value.String())
		}
	},
}

func traverseAll(workDir string) {
	file, err := os.Open(workDir)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	dirEntries, err := file.ReadDir(-1)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	for _, dir := range dirEntries {
		if dir.IsDir() {
			dirPath := workDir + string(os.PathSeparator) + dir.Name() + string(os.PathSeparator) + ".git"
			_, err := os.Open(dirPath)
			if err != nil {
				fmt.Errorf(err.Error())
				return
			}
			commandString := fmt.Sprintf(`git --git-dir=%s pull`, dirPath)
			cmd := execCommand(commandString)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				//logger.Errorf("Command %s finished with error %s", command.Command, err.Error())
				//s.Error(err)
				fmt.Errorf(err.Error())
				return
			}
			errorPiper, err := cmd.StderrPipe()
			if err != nil {
				//logger.Errorf("Command %s finished with error %s", command.Command, err.Error())
				//s.Error(err)
				fmt.Errorf(err.Error())
				return
			}
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				//logger.Infof("Command %s standard out opened", command.Command)
				pipeReader(stdout)
				wg.Done()
				//logger.Infof("Command %s standard out closed", command.Command)
			}()

			go func() {
				//logger.Infof("Command %s error out opened", command.Command)
				pipeReader(errorPiper)
				wg.Done()
				//logger.Infof("Command %s error out closed", command.Command)
			}()

			_ = cmd.Start()
			wg.Wait()

			if err != nil {
				//logger.Errorf("Command %s finished with error %s", command.Command, err.Error())
				//s.Error(err)
				fmt.Errorf(err.Error())
				return
			}
		}
	}
}

func pipeReader(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}

func execCommand(cmd string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/U", "/c", cmd)
	} else {
		return exec.Command("sh", "-c", cmd)
	}
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	updateCmd.Flags().StringP("workDir", "d", ".", "Specifies the workdir for application. Default current place")
	updateCmd.Flags().StringP("projectName", "p", "all", "Specifies the project to be updated. Default all")
	updateCmd.Flags().BoolP("sequential", "s", false, "Specifies if updates are sequential or parallel. Default parallel")
	viper.BindPFlag("projectName", rootCmd.Flags().Lookup("projectName"))
	viper.BindPFlag("sequential", rootCmd.Flags().Lookup("sequential"))
	viper.BindPFlag("workDir", rootCmd.Flags().Lookup("workDir"))

}
