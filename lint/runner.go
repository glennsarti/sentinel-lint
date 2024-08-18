package lint

type Runner interface {
	Config() Config

	Run() (Issues, error)
}
