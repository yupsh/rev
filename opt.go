package command

type SeparateFlag bool

const (
	Separate   SeparateFlag = true
	NoSeparate SeparateFlag = false
)

type flags struct {
	Separate SeparateFlag
}

func (s SeparateFlag) Configure(flags *flags) { flags.Separate = s }
