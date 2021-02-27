package main

import (
    "fmt"
    "os"
    "github.com/AdamSLevy/jsonrpc2/v14"
    "log"
    "github.com/fsnotify/fsnotify"
    "bufio"
)

func addReservedPeer(peer string) {
    fmt.Println("Adding peer", peer )
    var c jsonrpc2.Client
    params := []string{peer}
    var result int
    fmt.Println("sending request" )
    err := c.Request(nil, "http://localhost:9933", "system_addReservedPeer", params, &result)
    if _, ok := err.(jsonrpc2.Error); ok {
        fmt.Println("request error" )
        log.Fatal(err)
        os.Exit(1)
    }
    if err != nil {
        // The json unmarshaler cannot handle a null response object, which is the response we expect
        if err.Error() != "unexpected end of JSON input" {
            fmt.Println("json marshaling or network error" )
            log.Fatal(err)
            os.Exit(1)
        }
    }
}

func loadPeers(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        addReservedPeer(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func watchChanges(filename string) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    done := make(chan bool)
    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                if event.Op&fsnotify.Write == fsnotify.Write {
                    loadPeers(filename)
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("error:", err)
            }
        }
    }()

    err = watcher.Add(filename)
    if err != nil {
        log.Fatal(err)
    }
    <-done
}

func printHelp(progname string) {
    fmt.Printf(`Use:
%v [filename]
`, progname)
}

func main() {
    args := os.Args
    if len(args) != 2 {
        printHelp(args[0])
        os.Exit(1)
    }
    loadPeers(args[1])
    watchChanges(args[1])
}


