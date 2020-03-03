package app

import (
    "io"
    "os"
    "fmt"
    "time"
    "regexp"
    "reflect"
    "strconv"
    "strings"

    "github.com/go-pkg-hub/logrotate"
)

func PrepareLogs(entry *Entry) error {
    PrepareLog(&entry.Stdout)
    PrepareLog(&entry.Stderr)

    if entry.Log.File != "" {
        PrepareLog(&entry.Log)
        entry.Stdout = Log{}
        entry.Stderr = Log{}
    } else if entry.Stdout.File == entry.Stderr.File {
        l, err := MergeLogs(entry.Stderr, entry.Stdout)
        if err == nil {
            if reflect.DeepEqual(entry.Log, Log{}) {
                entry.Log = l
            } else {
                entry.Log.File = l.File
                if entry.Log.Size == "" {
                    entry.Log.Size = l.Size
                }
                if entry.Log.Num == 0 {
                    entry.Log.Num = l.Num
                }
            }
            PrepareLog(&entry.Log)
        }

        entry.Stdout = Log{}
        entry.Stderr = Log{}

        if err != nil {
            return err
        }
    }

    return nil
}

func PrepareLog(l *Log) {
    switch l.File {
    case "/dev/stdin", "/dev/stdout", "/dev/stderr":
        l.Size = ""
        l.Num = 0
    }
}

func OpenLogFile(l Log, defaultSize string, defaultNum int) (io.WriteCloser, error) {
    switch l.File {
    case "/dev/stdin", "/dev/stdout", "/dev/stderr":
        return os.OpenFile(l.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    }

    num := defaultNum
    if l.Num > 0 {
        num = l.Num
    }

    size := defaultSize
    if l.Size != "" {
        size = l.Size
    }

    opts := []logrotate.Option{
        logrotate.WithMaxSize(logrotate.StringToSize(size)),
        logrotate.WithMaxFiles(num),
    }

    return logrotate.New(l.File, opts...)
}

func GetExpression(descriptor string) (string, error) {
    const every = "@every "
    if strings.HasPrefix(descriptor, every) {
        duration, err := time.ParseDuration(descriptor[len(every):])
        if err != nil {
            return "", fmt.Errorf("failed to parse duration %s: %s", descriptor, err)
        }

        re := regexp.MustCompile(`(?P<val>\d+)(?P<unit>\w{1,1})`)

        d := make(map[string]int)
        matches := re.FindAllStringSubmatch(fmt.Sprint(duration), -1)
        for i := 0; i < len(matches); i++ {
            val, _ := strconv.Atoi(matches[i][1])
            unit := matches[i][2]
            d[unit] = val
        }

        exp := ""
        if duration > 59 * time.Minute {
            exp = fmt.Sprintf("%d %d */%d * * ?", d["s"], d["m"], d["h"])
        } else if duration > 59 * time.Second {
            exp = fmt.Sprintf("%d */%d * * * ?", d["s"], d["m"])
        } else {
            exp = fmt.Sprintf("*/%d * * * * ?", d["s"])
        }

        return exp, nil
    }

    return descriptor, nil
}

func Parse(spec string) (tz, exp string, err error) {
    tz = ""
    if strings.HasPrefix(spec, "TZ=") || strings.HasPrefix(spec, "CRON_TZ=") {
        i := strings.Index(spec, " ")
        tz = spec[:i]
        spec = strings.TrimSpace(spec[i:])
    }

    exp, err = GetExpression(spec)

    return tz, exp, err
}


