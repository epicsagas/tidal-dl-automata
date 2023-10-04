package main

import (
    "bufio"
    "log"
    "net/url"
    "os"
    "os/exec"
    "strings"
)

func main() {
        queueFilePath := "queue.txt"
	    queueFile, queueFileErr := os.OpenFile(queueFilePath, os.O_RDWR, 0644)
        if queueFileErr != nil {
            log.Fatal(queueFileErr)
            return
        }
        defer queueFile.Close()

        tempFilePath := "temp.txt"
        tempFile, tempFileErr := os.Create(tempFilePath)
        if tempFileErr != nil {
            log.Fatal(tempFileErr)
            return
        }
        defer tempFile.Close()

        scanner := bufio.NewScanner(queueFile)
        for scanner.Scan() {
            link := strings.TrimSpace(scanner.Text())
            _, err := url.ParseRequestURI(link)
            if err != nil {
                log.Printf("You have wrong url. err=%+v url=%+v\n", err, link)
                continue
            } else {
                cmd := exec.Command("tidal-dl", "-l", link)
                cmd.Stdout = os.Stdout
                cmd.Stderr = os.Stderr

                if cmdErr := cmd.Run(); err != nil {
                    log.Println(cmdErr)
                }
            }
        }

        if err := scanner.Err(); err != nil {
            log.Println(err)
            return
        }

        queueFile.Close()
        if err := os.Remove(queueFilePath); err != nil {
            log.Println("Error removing the queue file:", err)
            return
        }

        if err := os.Rename(tempFilePath, queueFilePath); err != nil {
            log.Println("Error renaming the temporary file:", err)
            return
        }

        log.Println("Download completed!")
}
