package opt

// Boolean flag types with constants
type SeparateFlag bool
const (
	Separate   SeparateFlag = true
	NoSeparate SeparateFlag = false
)

// Flags represents the configuration options for the rev command
type Flags struct {
	Separate SeparateFlag // Reverse each word separately instead of whole line
}

// Configure methods for the opt system
func (s SeparateFlag) Configure(flags *Flags) { flags.Separate = s }
