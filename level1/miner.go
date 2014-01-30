package main

import "fmt"

import "strconv"

import "strings"

import (
  "crypto/sha1"
  "io"
)

import "encoding/hex"

import "bytes"

import "os/exec"

import "runtime"

func RunCmd(cmd string) (string) {
  split := strings.Fields(cmd)
  output, _ := exec.Command(split[0], split[1:]...).Output()
  return strings.TrimSpace(string(output[:]))
}

func main() {
  runtime.GOMAXPROCS(8)

  for {
    fmt.Printf("Obtaining new copy of ledger...\n\n")
    RunCmd("./clean")

    fmt.Printf("Fetching parameters...\n\n")

    difficulty := RunCmd("cat level1/difficulty.txt")
    tree := RunCmd("git -C level1 write-tree")
    parent := RunCmd("git -C level1 rev-parse HEAD")
    time := RunCmd("date +%s")

    fmt.Printf("\tDifficulty: %s\n\tTree: %s\n\tParent: %s\n\tTime: %s\n\n", difficulty, tree, parent, time)

    fmt.Printf("Mining for Gitcoins...\n\n")

    counter := make(chan int)
    finish := make(chan int)
    go func() {
      for {
        problem := strconv.Itoa(<-counter)

        body := fmt.Sprintf("tree %v\nparent %v\nauthor CTF user <me@example.com> %v +0000\ncommitter CTF user <me@example.com> %v +0000\n\nGive me a Gitcoin\n\n%x\n", tree, parent, time, time, problem)
        count := len(body)
        commit := fmt.Sprintf("commit %v\x00%v", count, body)

        hasher := sha1.New()
        io.WriteString(hasher, commit)
        sha1 := fmt.Sprintf("%x", hasher.Sum(nil))

        sha1_value, _ := hex.DecodeString(sha1[:len(difficulty)])
        difficulty_value, _ := hex.DecodeString(fmt.Sprintf("%s", difficulty))
        if bytes.Compare(sha1_value, difficulty_value) <= 0 {
          fmt.Printf("Found a Gitcoin\n\n\tValue: %x\n\tSHA1: %s\n\n", problem, sha1)
          commitCmd := fmt.Sprintf("./commit %s %s %s %s", problem, tree, parent, time)
          _ = RunCmd(commitCmd)
          finish <- 1
        }
      }
    }()

    go func() {
      i := 0
      for {
        counter <- i
        i = i+1
      }
    }()

    _ = <- finish
    fmt.Printf("\n\n\n")
  }
}