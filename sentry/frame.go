package sentry

// NewFrame assembles a stacktrace frame out of `runtime.Frame`.
func NewFrame(f runtime.Frame) Frame {
	abspath := f.File
	filename := f.File
	function := f.Function
	var module string

	if filename != "" {
		filename = extractFilename(filename)
	} else {
		filename = unknown
	}

	if abspath == "" {
		abspath = unknown
	}

	if function != "" {
		module, function = deconstructFunctionName(function)
	}

	frame := Frame{
		AbsPath:  abspath,
		Filename: filename,
		Lineno:   f.Line,
		Module:   module,
		Function: function,
	}

	frame.InApp = isInAppFrame(frame)

	return frame
}

type Frame struct {
	Function    string                 `json:"function,omitempty"`
	Symbol      string                 `json:"symbol,omitempty"`
	Module      string                 `json:"module,omitempty"`
	Package     string                 `json:"package,omitempty"`
	Filename    string                 `json:"filename,omitempty"`
	AbsPath     string                 `json:"abs_path,omitempty"`
	Lineno      int                    `json:"lineno,omitempty"`
	Colno       int                    `json:"colno,omitempty"`
	PreContext  []string               `json:"pre_context,omitempty"`
	ContextLine string                 `json:"context_line,omitempty"`
	PostContext []string               `json:"post_context,omitempty"`
	InApp       bool                   `json:"in_app,omitempty"`
	Vars        map[string]interface{} `json:"vars,omitempty"`
}

func extractFrames(pcs []uintptr) []Frame {
	var frames []Frame
	callersFrames := runtime.CallersFrames(pcs)

	for {
		callerFrame, more := callersFrames.Next()

		frames = append([]Frame{
			NewFrame(callerFrame),
		}, frames...)

		if !more {
			break
		}
	}

	return frames
}

func filterFrames(frames []Frame) []Frame {
	isTestFileRegexp := regexp.MustCompile(`getsentry/sentry-go/.+_test.go`)
	isExampleFileRegexp := regexp.MustCompile(`getsentry/sentry-go/example/`)
	filteredFrames := make([]Frame, 0, len(frames))

	for _, frame := range frames {
		// go runtime frames
		if frame.Module == "runtime" || frame.Module == "testing" {
			continue
		}
		// sentry internal frames
		isTestFile := isTestFileRegexp.MatchString(frame.AbsPath)
		isExampleFile := isExampleFileRegexp.MatchString(frame.AbsPath)
		if strings.Contains(frame.AbsPath, "github.com/getsentry/sentry-go") &&
			!isTestFile &&
			!isExampleFile {
			continue
		}
		filteredFrames = append(filteredFrames, frame)
	}

	return filteredFrames
}

func extractFilename(path string) string {
	_, file := filepath.Split(path)
	return file
}

func isInAppFrame(frame Frame) bool {
	if strings.HasPrefix(frame.AbsPath, build.Default.GOROOT) ||
		strings.Contains(frame.Module, "vendor") ||
		strings.Contains(frame.Module, "third_party") {
		return false
	}

	return true
}

// Transform `runtime/debug.*T·ptrmethod` into `{ module: runtime/debug, function: *T.ptrmethod }`
func deconstructFunctionName(name string) (module string, function string) {
	if idx := strings.LastIndex(name, "."); idx != -1 {
		module = name[:idx]
		function = name[idx+1:]
	}
	function = strings.Replace(function, "·", ".", -1)
	return module, function
}

func callerFunctionName() string {
	pcs := make([]uintptr, 1)
	runtime.Callers(3, pcs)
	callersFrames := runtime.CallersFrames(pcs)
	callerFrame, _ := callersFrames.Next()
	_, function := deconstructFunctionName(callerFrame.Function)
	return function
}
