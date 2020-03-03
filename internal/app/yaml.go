package app

import (
    "io"
    "os"
    "bufio"
    "regexp"
    "path/filepath"
    "gopkg.in/yaml.v3"
)

func DumpYaml(prefix string, i interface{}) (out []byte, err error) {
    j := make(map[interface{}]interface{})
    j[prefix] = i
    return yaml.Marshal(&j)
}

func MergeLogs(l1 Log, l2 Log) (Log, error) {
    l := l1

    d, err := yaml.Marshal(&l2)
    if err != nil {
        return Log{}, err
    }

    err = yaml.Unmarshal(d, &l)
    if err != nil {
        return Log{}, err
    }

    return l, nil
}

func ParseEntries(r io.Reader, entries *[]Entry) (err error) {
    decoder := yaml.NewDecoder(r)

    for err == nil {
        entry := Entry{}
        err = decoder.Decode(&entry)
        if nil != err {
            if err == io.EOF {
                err = nil
            }
            return
        }
        *entries = append(*entries, entry)
    }

    return
}

func LoadEntries(s string, entries *[]Entry) (error) {
    f, err := os.Open(s)

    if nil != err {
        return err
    }

    defer f.Close()

    r := bufio.NewReader(f)

    return ParseEntries(r, entries)
}

func GetEntries(s string, entries *[]Entry) (error) {
    f, err := os.Stat(s)
    if os.IsNotExist(err) {
        return err
    }

    if !f.IsDir() {
        return LoadEntries(s, entries)
    }

    return filepath.Walk(s, func(p string, f os.FileInfo, _ error) error {
        if !f.IsDir() {
            r, err := regexp.MatchString(".yaml", f.Name())
            if err == nil && r {
                err = LoadEntries(p, entries)
                if err != nil {
                    return err
                }
            }
        }
        return nil
    })
}
