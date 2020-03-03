package app

type Log struct {
    File string `yaml:",omitempty"`
    Size string `yaml:",omitempty"`
    Num  int    `yaml:",omitempty"`
}

type Cmd struct {
    Name string
    Argv []string `yaml:",omitempty"`
    Env []string `yaml:",omitempty"`
    Dir string `yaml:",omitempty"`
}

type Entry struct {
    Exp string
    Lock bool
    Cmd Cmd
    Log Log `yaml:",omitempty"`
    Stdout Log `yaml:",omitempty"`
    Stderr Log `yaml:",omitempty"`
}
