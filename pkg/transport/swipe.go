package transport

import (
	"github.com/fesyunoff/api/pkg/controller"
	. "github.com/swipe-io/swipe/v2"
)

func Swipe() {
	Build(
		Service(
			// ReadmeEnable(),
			// ReadmeOutput("./"),

			HTTPServer(),
			Interface((*controller.Service)(nil), ""),

			JSONRPCEnable(),
			// JSONRPCDocEnable(),
			// JSONRPCDocOutput("./docs"),
			JSONRPCPath("/{method:.*}"),

			// OpenapiEnable(),
			// OpenapiOutput("./docs"),

			MethodDefaultOptions(
				Logging(true),
				Instrumenting(true),
			),
		),
		// ),
	)
}
