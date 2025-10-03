package dove

import "io"

type Renderer interface {
	Render(w io.Writer) error
}
