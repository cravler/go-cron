package main

import (
    "io"
    "os"
    "log"
    "path"
    "strings"
    "syscall"
    "os/exec"
    "os/signal"

    "github.com/spf13/cobra"
    "github.com/robfig/cron/v3"

    "github.com/cravler/go-cron/internal/app"
)

var version = "0.x-dev"

var logger = log.New(os.Stderr, "", log.LstdFlags)

func main() {
    rootCmdName := path.Base(os.Args[0])
    rootCmd := app.NewRootCmd(rootCmdName, version, func(c *cobra.Command, args []string) error {
        workdir, _ := c.Flags().GetString("workdir")
        verbose, _ := c.Flags().GetBool("verbose")

        runApp(args, workdir, verbose)

        return nil
    })

    rootCmd.Flags().StringP("workdir", "w", "", "Working directory")
    rootCmd.Flags().BoolP("verbose", "v", false, "Verbose output")

    if err := rootCmd.Execute(); err != nil {
        rootCmd.Println(err)
        os.Exit(1)
    }
}

func runApp(args []string, workdir string, verbose bool) {
    l := cron.PrintfLogger(logger)
    p := cron.NewParser(cron.SecondOptional |cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
    c := cron.New(cron.WithParser(p), cron.WithLogger(l))

    entries := []app.Entry{}
    for _, crontab := range args {
        if err := app.GetEntries(crontab, &entries); err != nil {
            logger.Printf("%s\n", err.Error())
        }
    }

    for _, raw := range entries {
        entry := raw

        tz, exp, err := app.Parse(entry.Exp)
        if err != nil {
            logger.Printf("%s\n", err.Error())
            continue
        }

        if strings.TrimSpace(tz + " " + exp) != strings.TrimSpace(entry.Exp) {
            entry.Exp = entry.Exp + " (" + exp + ")"
        }

        if err := app.PrepareLogs(&entry); err != nil {
            logger.Printf("%s\n", err.Error())
        }

        if verbose {
            d, _ := app.DumpYaml("Add", entry)
            logger.Printf("%s\n", d)
        }

        schedule, err := p.Parse(strings.TrimSpace(tz + " " + exp))
        if err != nil {
            logger.Printf("%s\n", err.Error())
            continue
        }

        fn := func() {
            if verbose {
                d, _ := app.DumpYaml("Execute", entry)
                logger.Printf("%s\n", d)
            }

            cmd := exec.Command(entry.Cmd.Name, entry.Cmd.Argv...)
            cmd.Env = os.Environ()
            if entry.Cmd.Env != nil {
                cmd.Env = append(cmd.Env, entry.Cmd.Env...)
            }
            if entry.Cmd.Dir != "" {
                cmd.Dir = entry.Cmd.Dir
            } else if workdir != "" {
                cmd.Dir = workdir
            }

            var outFile, errFile io.WriteCloser

            if entry.Log.File != "" {
                file, err := app.OpenLogFile(entry.Log, "", 0)
                if err != nil {
                    logger.Printf("%s\n", err.Error())
                } else {
                    defer file.Close()
                    outFile = file
                    errFile = file
                }
            } else {
                if entry.Stdout.File != "" {
                    file, err := app.OpenLogFile(entry.Stdout, entry.Log.Size, entry.Log.Num)
                    if err != nil {
                        logger.Printf("%s\n", err.Error())
                    } else {
                        defer file.Close()
                        outFile = file
                    }
                }
                if entry.Stderr.File != "" {
                    file, err := app.OpenLogFile(entry.Stderr, entry.Log.Size, entry.Log.Num)
                    if err != nil {
                        logger.Printf("%s\n", err.Error())
                    } else {
                        defer file.Close()
                        errFile = file
                    }
                }
            }

            if outFile != nil {
                cmd.Stdout = outFile
            }
            if errFile != nil {
                cmd.Stderr = errFile
            }

            if err := cmd.Start(); err != nil {
                logger.Printf("%s\n", err.Error())
            }
            cmd.Wait()
        }

        var job cron.Job
        job = cron.FuncJob(fn)
        if entry.Lock {
            job = cron.NewChain(cron.SkipIfStillRunning(l)).Then(job)
        }

        c.Schedule(schedule, job)
    }

    s := make(chan os.Signal, 1)
    signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
    defer signal.Stop(s)

    c.Start()
    <-s
    <-c.Stop().Done()
}
