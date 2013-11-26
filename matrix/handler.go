package matrix

type outputmode int

const (
	OM_Replace outputmode = 0 // replace input extension
	OM_Prepend outputmode = 1 // prepend to input extension
	OM_Append  outputmode = 2 // append to input extension
	OM_Discard outputmode = 3 // there is no output
)

// handler is nil when canHandle is false
type HandlerConstructor func(*File) (handler Handler, canHandle bool)

type FileExt struct {
	// May be set to "*" to catch all or a string such as "js" to match all files with the "js" file extention
	// Is set to an empty string if the output does not directly correlate to an input
	Input string
	// Is true if handler takes all matching inputs at once
	InputMulti bool
	// May be set to a string such as "js"
	// Is set to an empty string if the input does not directly correlate to an output
	Output string
	// Specifies how the output suffix should be placed relative to the input suffix and if there is an output
	OutputMode outputmode
}

type Handler interface {
	// File extensions supported for input/output
	FileExts() []*FileExt
	GenerateOutput(*File) []*File
	InputFile() *File
}
